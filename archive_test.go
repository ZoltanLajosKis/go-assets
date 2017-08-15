package assets

import (
	"archive/zip"
	"bytes"
	"os"
	"testing"
	"time"

	mfs "github.com/ZoltanLajosKis/go-mapfs"
)

func TestArchiveNil(t *testing.T) {
	fs := make(mfs.Files)
	mt := time.Unix(1500000000, 0)

	err := processArchive(nil, "test/file1.txt", []byte("Test"), mt, fs)
	assertEqual(t, err, nil)

	f1, _ := fs["test/file1.txt"]
	assertEqual(t, f1, &mfs.File{[]byte("Test"), mt})

	_, ok := fs["test/file2.txt"]
	assertEqual(t, ok, false)
}

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

	err := processArchive(&Archive{ArchiveZip, nil}, "arch.zip", buf.Bytes(), time.Now(), fs)
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

	err := processArchive(&Archive{ArchiveZip, mapper}, "arch.zip", buf.Bytes(), time.Now(), fs)
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

	err := processArchive(&Archive{ArchiveZip, nil}, "arch.zip", []byte("1234"), time.Now(), fs)
	assertEqual(t, err, &ArchiveError{"arch.zip", zip.ErrFormat})
}

func TestArchiveUnknown(t *testing.T) {
	fs := make(mfs.Files)
	mt := time.Unix(1500000000, 0)

	err := processArchive(&Archive{-1, nil}, "test/file1.txt", []byte("Test"), mt, fs)
	assertEqual(t, err, ErrArchiveUnknown)
}
