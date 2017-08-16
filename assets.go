package assets

import (
	"fmt"
	"log"

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

// Compile retrieves and processes the specified asset sources, and
// compiles them to source code in the specified file
func Compile(assets []*Source, filePath string, pkgName string, varName string, opts *Opts) error {
	files := make(mfs.Files)

	for i, asset := range assets {
		log.Printf("Processing asset source (%d/%d): %s ...", i+1, len(assets), asset.Location)

		data, modTime, err := retrieve(asset.Location)
		if err != nil {
			return err
		}

		err = verifyChecksum(asset.Checksum, data)
		if err != nil {
			return err
		}

		err = processArchive(asset.Archive, asset.Path, data, modTime, files)
		if err != nil {
			return err
		}
	}

	fs, err := mfs.New(files)
	if err != nil {
		return err
	}

	if opts == nil {
		opts = &Opts{}
	}

	if opts.VariableComment == "" {
		opts.VariableComment = fmt.Sprintf("%s implements a http.FileSystem.", varName)
	}

	err = vfsgen.Generate(httpfs.New(fs), vfsgen.Options{
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
