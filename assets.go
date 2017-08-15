package assets

import (
	"log"

	mfs "github.com/ZoltanLajosKis/go-mapfs"
	"github.com/shurcooL/vfsgen"
	"golang.org/x/tools/godoc/vfs/httpfs"
)

// AssetSource describe an asset source to be retrieved and processed.
type AssetSource struct {
	Name     string
	Location string
	Checksum *Checksum
	Archive  *Archive
}

// Compile retrieves and processes the specified asset sources, and
// compiles them to source code in the specified file
func Compile(assets []*AssetSource, filePath string, pkgName string, varName string) error {
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

		err = processArchive(asset.Archive, asset.Name, data, modTime, files)
		if err != nil {
			return err
		}
	}

	fs, err := mfs.New(files)
	if err != nil {
		return err
	}

	err = vfsgen.Generate(httpfs.New(fs), vfsgen.Options{
		Filename:     filePath,
		PackageName:  pkgName,
		VariableName: varName,
	})
	if err != nil {
		return err
	}

	return nil
}
