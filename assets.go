package assets

import (
	"fmt"
	"log"
	"net/http"

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

// Retrieve retrieves and processes the specified asset sources, and returns
// them using a http.FileSystem interface.
func Retrieve(sources []*Source) (http.FileSystem, error) {
	files := make(mfs.Files)

	for i, source := range sources {
		log.Printf("Processing asset source (%d/%d): %s ...", i+1, len(sources), source.Location)

		data, modTime, err := retrieve(source.Location)
		if err != nil {
			return nil, err
		}

		if source.Checksum != nil {
			err = verifyChecksum(source.Checksum, data)
			if err != nil {
				return nil, err
			}
		}

		if source.Archive == nil {
			log.Printf("Created asset: %s ...", source.Path)
			files[source.Path] = &mfs.File{data, modTime}
			continue
		}

		err = processArchive(source.Archive, source.Path, data, modTime, files)
		if err != nil {
			return nil, err
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
