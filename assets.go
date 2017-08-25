package assets

import (
	"fmt"
	"log"
	"net/http"
	"time"

	mfs "github.com/ZoltanLajosKis/go-mapfs"
	"github.com/shurcooL/vfsgen"
	"golang.org/x/tools/godoc/vfs/httpfs"
)

// Source describes an asset source to be retrieved and processed.
type Source struct {
	Path     string
	Location string
	Checksum *Checksum
	Archive  *Archive
}

// Opts provides optional parameters to the Compile function.
type Opts struct {
	// BuildTags are the build tags in the generated source code.
	// Defaults to no tags.
	BuildTags string

	// VariableComment is the comment of the variable in the generated source code.
	// Defaults to "<VariableName> implements a http.FileSystem.".
	VariableComment string
}

type file struct {
	path    string
	data    []byte
	modTime time.Time
}

// Retrieve retrieves and processes the specified asset sources, and returns
// them using a http.FileSystem interface.
func Retrieve(sources []*Source) (http.FileSystem, error) {
	files := make(mfs.Files)

	for i, source := range sources {
		log.Printf("Processing asset source (%d/%d): %s ...", i+1, len(sources), source.Location)

		file, err := retrieve(source.Location)
		if err != nil {
			return nil, &RetrieveError{source.Location, err}
		}

		if source.Checksum != nil {
			err = verifyChecksum(source.Checksum, file.data)
			if err != nil {
				return nil, &ChecksumError{source.Location, err}
			}
		}

		if source.Archive == nil {
			log.Printf("Created asset: %s ...", source.Path)
			files[source.Path] = &mfs.File{file.data, file.modTime}
			continue
		}

		archFiles, err := processArchive(source.Archive, file.data)
		if err != nil {
			return nil, &ArchiveError{source.Location, err}
		}

		for _, file := range archFiles {
			log.Printf("Created asset: %s ...", file.path)
			files[file.path] = &mfs.File{file.data, file.modTime}
		}

	}

	fs, err := mfs.New(files)
	if err != nil {
		return nil, err
	}

	return httpfs.New(fs), nil
}

// Compile retrieves and processes the specified asset sources, and
// compiles them to the specified variable in the source file.
func Compile(sources []*Source, filePath string, pkgName string, varName string, opts *Opts) error {

	fs, err := Retrieve(sources)
	if err != nil {
		return err
	}

	if opts == nil {
		opts = &Opts{}
	}

	if opts.VariableComment == "" {
		opts.VariableComment = fmt.Sprintf("%s implements a http.FileSystem.", varName)
	}

	err = vfsgen.Generate(fs, vfsgen.Options{
		Filename:        filePath,
		PackageName:     pkgName,
		BuildTags:       opts.BuildTags,
		VariableName:    varName,
		VariableComment: opts.VariableComment,
	})
	if err != nil {
		return err
	}

	return nil
}

// RetrieveError is returned when there is a problem retrieving an asset source
type RetrieveError struct {
	Location string
	Err      error
}

func (e *RetrieveError) Error() string {
	return e.Location + ": " + e.Err.Error()
}

// ChecksumError is returned when there is a checksum problem with an asset source
type ChecksumError struct {
	Location string
	Err      error
}

func (e *ChecksumError) Error() string {
	return e.Location + ": " + e.Err.Error()
}

// ArchiveError is returned when there is a problem processing the archive
type ArchiveError struct {
	Path string
	Err  error
}

func (e *ArchiveError) Error() string {
	return e.Path + ": " + e.Err.Error()
}
