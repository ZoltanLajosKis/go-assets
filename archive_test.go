package assets

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"testing"
	"time"

	mfs "github.com/ZoltanLajosKis/go-mapfs"
)

func TestArchiveZip(t *testing.T) {
	fs := make(mfs.Files)

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	fh1 := &zip.FileHeader{Name: "test/file1.txt"}
	mt1 := time.Unix(1300000000, 0)
	fh1.SetModTime(mt1)
	f1, _ := w.CreateHeader(fh1)
	f1.Write([]byte("File 1"))

	fh2 := &zip.FileHeader{Name: "test/file2.txt"}
	mt2 := time.Unix(1400000000, 0)
	fh2.SetModTime(mt2)
	f2, _ := w.CreateHeader(fh2)
	f2.Write([]byte("File 2"))

	fh3 := &zip.FileHeader{Name: "test/file3.txt"}
	mt3 := time.Unix(1500000000, 0)
	fh3.SetModTime(mt3)
	f3, _ := w.CreateHeader(fh3)
	f3.Write([]byte("File 3"))

	dh := &zip.FileHeader{Name: "test/dir"}
	dh.SetMode(0755 | os.ModeDir)
	w.CreateHeader(dh)

	w.Close()

	err := processArchive(&Archive{Zip, nil}, buf.Bytes(), fs)
	assertEqual(t, err, nil)

	file1, _ := fs["test/file1.txt"]
	assertEqual(t, file1, &mfs.File{[]byte("File 1"), mt1.UTC()})
	file2, _ := fs["test/file2.txt"]
	assertEqual(t, file2, &mfs.File{[]byte("File 2"), mt2.UTC()})
	file3, _ := fs["test/file3.txt"]
	assertEqual(t, file3, &mfs.File{[]byte("File 3"), mt3.UTC()})
	_, ok := fs["test/dir"]
	assertEqual(t, ok, false)
}

func TestArchiveZipFilter(t *testing.T) {
	fs := make(mfs.Files)

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	fh1 := &zip.FileHeader{Name: "test/file1.txt"}
	mt1 := time.Unix(1300000000, 0)
	fh1.SetModTime(mt1)
	f1, _ := w.CreateHeader(fh1)
	f1.Write([]byte("File 1"))

	fh2 := &zip.FileHeader{Name: "test/file2.txt"}
	mt2 := time.Unix(1400000000, 0)
	fh2.SetModTime(mt2)
	f2, _ := w.CreateHeader(fh2)
	f2.Write([]byte("File 2"))

	fh3 := &zip.FileHeader{Name: "test/file3.txt"}
	mt3 := time.Unix(1500000000, 0)
	fh3.SetModTime(mt3)
	f3, _ := w.CreateHeader(fh3)
	f3.Write([]byte("File 3"))

	dh := &zip.FileHeader{Name: "test/dir"}
	dh.SetMode(0755 | os.ModeDir)
	w.CreateHeader(dh)

	w.Close()

	mapper := func(s string) string {
		switch s {
		case "test/file1.txt", "test/file2.txt":
			return s
		default:
			return ""
		}
	}

	err := processArchive(&Archive{Zip, mapper}, buf.Bytes(), fs)
	assertEqual(t, err, nil)

	file1, _ := fs["test/file1.txt"]
	assertEqual(t, file1, &mfs.File{[]byte("File 1"), mt1.UTC()})
	file2, _ := fs["test/file2.txt"]
	assertEqual(t, file2, &mfs.File{[]byte("File 2"), mt2.UTC()})
	_, ok := fs["test/file3.txt"]
	assertEqual(t, ok, false)
	_, ok = fs["test/dir"]
	assertEqual(t, ok, false)
}

func TestArchiveZipReMap(t *testing.T) {
	fs := make(mfs.Files)

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	fh1 := &zip.FileHeader{Name: "test/file1.txt"}
	mt1 := time.Unix(1300000000, 0)
	fh1.SetModTime(mt1)
	f1, _ := w.CreateHeader(fh1)
	f1.Write([]byte("File 1"))

	fh2 := &zip.FileHeader{Name: "test/file2.txt"}
	mt2 := time.Unix(1400000000, 0)
	fh2.SetModTime(mt2)
	f2, _ := w.CreateHeader(fh2)
	f2.Write([]byte("File 2"))

	fh3 := &zip.FileHeader{Name: "test/file3.txt"}
	mt3 := time.Unix(1500000000, 0)
	fh3.SetModTime(mt3)
	f3, _ := w.CreateHeader(fh3)
	f3.Write([]byte("File 3"))

	dh := &zip.FileHeader{Name: "test/dir"}
	dh.SetMode(0755 | os.ModeDir)
	w.CreateHeader(dh)

	w.Close()

	err := processArchive(&Archive{Zip, ReMap("(test/file[12].txt)", "${1}")}, buf.Bytes(), fs)
	assertEqual(t, err, nil)

	file1, _ := fs["test/file1.txt"]
	assertEqual(t, file1, &mfs.File{[]byte("File 1"), mt1.UTC()})
	file2, _ := fs["test/file2.txt"]
	assertEqual(t, file2, &mfs.File{[]byte("File 2"), mt2.UTC()})
	_, ok := fs["test/file3.txt"]
	assertEqual(t, ok, false)
	_, ok = fs["test/dir"]
	assertEqual(t, ok, false)
}

func TestArchiveZipInvalid(t *testing.T) {
	fs := make(mfs.Files)

	err := processArchive(&Archive{Zip, nil}, []byte("1234"), fs)
	assertEqual(t, err, zip.ErrFormat)
}

func TestArchiveTarGz(t *testing.T) {
	fs := make(mfs.Files)
	buf := new(bytes.Buffer)
	zw := gzip.NewWriter(buf)
	w := tar.NewWriter(zw)

	mt1 := time.Unix(1300000000, 0)
	fh1 := &tar.Header{Name: "test/file1.txt", Size: int64(6), ModTime: mt1}
	w.WriteHeader(fh1)
	w.Write([]byte("File 1"))

	mt2 := time.Unix(1400000000, 0)
	fh2 := &tar.Header{Name: "test/file2.txt", Size: int64(6), ModTime: mt2}
	w.WriteHeader(fh2)
	w.Write([]byte("File 2"))

	mt3 := time.Unix(1500000000, 0)
	fh3 := &tar.Header{Name: "test/file3.txt", Size: int64(6), ModTime: mt3}
	w.WriteHeader(fh3)
	w.Write([]byte("File 3"))

	dh := &tar.Header{Name: "test/dir", Typeflag: tar.TypeDir}
	w.WriteHeader(dh)

	w.Close()
	zw.Close()

	err := processArchive(&Archive{TarGz, nil}, buf.Bytes(), fs)
	assertEqual(t, err, nil)

	file1, _ := fs["test/file1.txt"]
	assertEqual(t, file1, &mfs.File{[]byte("File 1"), mt1})
	file2, _ := fs["test/file2.txt"]
	assertEqual(t, file2, &mfs.File{[]byte("File 2"), mt2})
	file3, _ := fs["test/file3.txt"]
	assertEqual(t, file3, &mfs.File{[]byte("File 3"), mt3})
	_, ok := fs["test/dir"]
	assertEqual(t, ok, false)
}

func TestArchiveTarGzFilter(t *testing.T) {
	fs := make(mfs.Files)
	buf := new(bytes.Buffer)
	zw := gzip.NewWriter(buf)
	w := tar.NewWriter(zw)

	mt1 := time.Unix(1300000000, 0)
	fh1 := &tar.Header{Name: "test/file1.txt", Size: int64(6), ModTime: mt1}
	w.WriteHeader(fh1)
	w.Write([]byte("File 1"))

	mt2 := time.Unix(1400000000, 0)
	fh2 := &tar.Header{Name: "test/file2.txt", Size: int64(6), ModTime: mt2}
	w.WriteHeader(fh2)
	w.Write([]byte("File 2"))

	mt3 := time.Unix(1500000000, 0)
	fh3 := &tar.Header{Name: "test/file3.txt", Size: int64(6), ModTime: mt3}
	w.WriteHeader(fh3)
	w.Write([]byte("File 3"))

	dh := &tar.Header{Name: "test/dir", Typeflag: tar.TypeDir}
	w.WriteHeader(dh)

	w.Close()
	zw.Close()

	mapper := func(s string) string {
		switch s {
		case "test/file1.txt", "test/file2.txt":
			return s
		default:
			return ""
		}
	}

	err := processArchive(&Archive{TarGz, mapper}, buf.Bytes(), fs)
	assertEqual(t, err, nil)

	file1, _ := fs["test/file1.txt"]
	assertEqual(t, file1, &mfs.File{[]byte("File 1"), mt1})
	file2, _ := fs["test/file2.txt"]
	assertEqual(t, file2, &mfs.File{[]byte("File 2"), mt2})
	_, ok := fs["test/file3.txt"]
	assertEqual(t, ok, false)
	_, ok = fs["test/dir"]
	assertEqual(t, ok, false)
}

func TestArchiveTarGzInvalid(t *testing.T) {
	fs := make(mfs.Files)

	err := processArchive(&Archive{TarGz, nil}, []byte("1234"), fs)
	assertEqual(t, err, io.ErrUnexpectedEOF)
}

func TestArchiveUnknown(t *testing.T) {
	fs := make(mfs.Files)

	err := processArchive(&Archive{-1, nil}, []byte("Test"), fs)
	assertEqual(t, err, ErrArchiveUnknown)
}
