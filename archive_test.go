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
)

func TestArchiveZip(t *testing.T) {
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

	files, err := processArchive(&Archive{Zip, nil}, buf.Bytes())
	assertEqual(t, err, nil)

	assertEqual(t, len(files), 3)
	assertEqual(t, files[0], &file{"test/file1.txt", []byte("File 1"), mt1.UTC()})
	assertEqual(t, files[1], &file{"test/file2.txt", []byte("File 2"), mt2.UTC()})
	assertEqual(t, files[2], &file{"test/file3.txt", []byte("File 3"), mt3.UTC()})
}

func TestArchiveZipFilter(t *testing.T) {
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

	files, err := processArchive(&Archive{Zip, mapper}, buf.Bytes())
	assertEqual(t, err, nil)

	assertEqual(t, len(files), 2)
	assertEqual(t, files[0], &file{"test/file1.txt", []byte("File 1"), mt1.UTC()})
	assertEqual(t, files[1], &file{"test/file2.txt", []byte("File 2"), mt2.UTC()})
}

func TestArchiveZipReMap(t *testing.T) {
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

	files, err := processArchive(&Archive{Zip, ReMap("(test/file[12].txt)", "${1}")}, buf.Bytes())
	assertEqual(t, err, nil)

	assertEqual(t, len(files), 2)
	assertEqual(t, files[0], &file{"test/file1.txt", []byte("File 1"), mt1.UTC()})
	assertEqual(t, files[1], &file{"test/file2.txt", []byte("File 2"), mt2.UTC()})
}

func TestArchiveZipInvalid(t *testing.T) {
	_, err := processArchive(&Archive{Zip, nil}, []byte("1234"))
	assertEqual(t, err, zip.ErrFormat)
}

func TestArchiveTarGz(t *testing.T) {
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

	files, err := processArchive(&Archive{TarGz, nil}, buf.Bytes())
	assertEqual(t, err, nil)

	assertEqual(t, len(files), 3)
	assertEqual(t, files[0], &file{"test/file1.txt", []byte("File 1"), mt1})
	assertEqual(t, files[1], &file{"test/file2.txt", []byte("File 2"), mt2})
	assertEqual(t, files[2], &file{"test/file3.txt", []byte("File 3"), mt3})
}

func TestArchiveTarGzFilter(t *testing.T) {
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

	files, err := processArchive(&Archive{TarGz, mapper}, buf.Bytes())
	assertEqual(t, err, nil)

	assertEqual(t, len(files), 2)
	assertEqual(t, files[0], &file{"test/file1.txt", []byte("File 1"), mt1})
	assertEqual(t, files[1], &file{"test/file2.txt", []byte("File 2"), mt2})
}

func TestArchiveTarGzInvalid(t *testing.T) {
	_, err := processArchive(&Archive{TarGz, nil}, []byte("1234"))
	assertEqual(t, err, io.ErrUnexpectedEOF)
}

func TestArchiveUnknown(t *testing.T) {
	_, err := processArchive(&Archive{-1, nil}, []byte("Test"))
	assertEqual(t, err, ErrArchiveUnknown)
}
