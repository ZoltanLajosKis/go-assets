package assets

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"time"

	mfs "github.com/ZoltanLajosKis/go-mapfs"
)

// ArchiveFormat enumerates archive formats.
type ArchiveFormat int

const (
	// Zip is the zip file format.
	Zip = iota
	// TarGz is the tar.gz file format.
	TarGz
)

var (
	// ErrArchiveUnknown is returned when an invalid archive format is specified
	ErrArchiveUnknown = errors.New("unknown archive format")
)

// ArchiveError is returned when there is a problem processing the archive
type ArchiveError struct {
	Path string
	Err  error
}

func (e *ArchiveError) Error() string {
	return e.Path + ": " + e.Err.Error()
}

// PathMapper specifies a function that is executed on all files in the archive.
// The mapper receives the full path to each file in the archive and returns
// the path to use in the asset file system. If "" is returned, the file is
// dropped.
type PathMapper func(string) string

// ReMap returns a PathMapper that compares file paths to the input pattern
// and maps matches to the replacement string (see Regexp.ReplaceAllString).
func ReMap(pattern string, replacement string) PathMapper {
	re := regexp.MustCompile(pattern)
	return func(filePath string) string {
		if !re.MatchString(filePath) {
			return ""
		}
		return re.ReplaceAllString(filePath, replacement)
	}
}

// Archive describes an archive format for the asset source.
type Archive struct {
	Format     ArchiveFormat
	PathMapper PathMapper
}

func processArchive(arch *Archive, path string, data []byte, modTime time.Time, files mfs.Files) error {
	switch arch.Format {
	case Zip:
		return processZip(arch, path, data, files)
	case TarGz:
		return processTarGz(arch, path, data, files)
	default:
		return ErrArchiveUnknown
	}
}

func processZip(arch *Archive, path string, data []byte, files mfs.Files) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return &ArchiveError{path, err}
	}

	for _, fh := range r.File {
		if fh.FileInfo().IsDir() {
			continue
		}

		fp := mapPath(arch.PathMapper, fh.Name)
		if fp == "" {
			continue
		}

		fr, err := fh.Open()
		if err != nil {
			return &ArchiveError{path, err}
		}

		fdata, err := ioutil.ReadAll(fr)
		if err != nil {
			return &ArchiveError{path, err}
		}

		log.Printf("Created asset: %s ...", fp)
		files[fp] = &mfs.File{fdata, fh.ModTime()}
	}

	return nil
}

func processTarGz(arch *Archive, path string, data []byte, files mfs.Files) error {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return &ArchiveError{path, err}
	}

	r := tar.NewReader(zr)

	for {
		h, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &ArchiveError{path, err}
		}

		if h.Typeflag != tar.TypeReg && h.Typeflag != tar.TypeRegA {
			continue
		}

		fp := mapPath(arch.PathMapper, h.Name)
		if fp == "" {
			continue
		}

		fdata, err := ioutil.ReadAll(r)
		if err != nil {
			return &ArchiveError{path, err}
		}

		log.Printf("Created asset: %s ...", fp)
		files[fp] = &mfs.File{fdata, h.ModTime}
	}

	return nil
}

func mapPath(mapper PathMapper, path string) string {
	if mapper == nil {
		return path
	}

	return mapper(path)
}
