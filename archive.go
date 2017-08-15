package assets

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"time"

	mfs "github.com/ZoltanLajosKis/go-mapfs"
)

// ArchiveFormat enumerates archive formats.
type ArchiveFormat int

const (
	// ArchiveZip is for the zip file format.
	ArchiveZip = iota
)

var (
	// ErrArchiveUnknown is returned when an invalid archive format is specified
	ErrArchiveUnknown = errors.New("unknown archive format")
)

// ArchiveError is returned when there is a problem processing the archive
type ArchiveError struct {
	Name string
	Err  error
}

func (e *ArchiveError) Error() string {
	return e.Name + ": " + e.Err.Error()
}

// PathMapper specifies a function that is executed on all files in the archive.
// The archived file will be stored under the returned path (ignoring the asset
// source name). If "" is returned, the file is ignored.
type PathMapper func(string) string

// Archive describes an archive format for the asset source.
type Archive struct {
	Format     ArchiveFormat
	PathMapper PathMapper
}

func processArchive(arch *Archive, name string, data []byte, modTime time.Time, files mfs.Files) error {
	if arch == nil {
		log.Printf("Created asset: %s ...", name)
		files[name] = &mfs.File{data, modTime}
		return nil
	}

	switch arch.Format {
	case ArchiveZip:
		return processZip(arch, name, data, files)
	default:
		return ErrArchiveUnknown
	}
}

func processZip(arch *Archive, name string, data []byte, files mfs.Files) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return &ArchiveError{name, err}
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
			return &ArchiveError{name, err}
		}

		fdata, err := ioutil.ReadAll(fr)
		if err != nil {
			return &ArchiveError{name, err}
		}

		log.Printf("Created asset: %s ...", fp)
		files[fp] = &mfs.File{fdata, fh.ModTime()}
	}

	return nil
}

func mapPath(mapper PathMapper, path string) string {
	if mapper == nil {
		return path
	}

	return mapper(path)
}
