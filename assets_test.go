package assets

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCompile(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-assets")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/assets.txt" {
			fmt.Fprintf(w, "Assets.")
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer svr.Close()

	sources := []*Source{
		{"assets.txt",
			strings.Join([]string{svr.URL, "/assets.txt"}, ""), nil, nil},
		{"retrieve_test.go",
			"retrieve_test.go", nil, nil},
	}

	err = Compile(sources, filepath.Join(dir, "assets.go"), "assets", "fs", nil)
	assertEqual(t, err, nil)
}

func TestRetrieveGlob(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-assets")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = os.MkdirAll(dir+"/test/t1", 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(dir+"/test/t2", 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(dir+"/test/t1/file1.txt", []byte("File 1"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(dir+"/test/t1/file2.txt", []byte("File 2"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(dir+"/test/t2/file3.txt", []byte("File 3"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(dir+"/test/t2/file4.txt", []byte("File 4"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	sources := []*Source{
		{"newdir",
			dir + "/test/t[12]/file[123].txt", nil, nil},
	}

	fs, err := Retrieve(sources)
	assertEqual(t, err, nil)

	file1, err := fs.Open("newdir/t1/file1.txt")
	assertEqual(t, err, nil)
	fdata1, err := ioutil.ReadAll(file1)
	assertEqual(t, err, nil)
	assertEqual(t, fdata1, []byte("File 1"))

	file2, err := fs.Open("newdir/t1/file2.txt")
	assertEqual(t, err, nil)
	fdata2, err := ioutil.ReadAll(file2)
	assertEqual(t, err, nil)
	assertEqual(t, fdata2, []byte("File 2"))

	file3, err := fs.Open("newdir/t2/file3.txt")
	assertEqual(t, err, nil)
	fdata3, err := ioutil.ReadAll(file3)
	assertEqual(t, err, nil)
	assertEqual(t, fdata3, []byte("File 3"))

	_, err = fs.Open("newdir/t2/file4.txt")
	assertNotEqual(t, err, nil)
}

func TestCompileArchive(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-assets")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

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

	err = ioutil.WriteFile(path.Join(dir, "arch.zip"), buf.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}

	sources := []*Source{
		{"arch.zip",
			path.Join(dir, "arch.zip"), nil, &Archive{Zip, nil}},
	}

	fs, err := Retrieve(sources)
	assertEqual(t, err, nil)

	file1, err := fs.Open("test/file1.txt")
	assertEqual(t, err, nil)
	fdata1, err := ioutil.ReadAll(file1)
	assertEqual(t, err, nil)
	fstat1, err := file1.Stat()
	assertEqual(t, err, nil)
	assertEqual(t, fdata1, []byte("File 1"))
	assertEqual(t, fstat1.ModTime(), mt1.UTC())

	file2, err := fs.Open("test/file2.txt")
	assertEqual(t, err, nil)
	fdata2, err := ioutil.ReadAll(file2)
	assertEqual(t, err, nil)
	fstat2, err := file2.Stat()
	assertEqual(t, err, nil)
	assertEqual(t, fdata2, []byte("File 2"))
	assertEqual(t, fstat2.ModTime(), mt2.UTC())

	file3, err := fs.Open("test/file3.txt")
	assertEqual(t, err, nil)
	fdata3, err := ioutil.ReadAll(file3)
	assertEqual(t, err, nil)
	fstat3, err := file3.Stat()
	assertEqual(t, err, nil)
	assertEqual(t, fdata3, []byte("File 3"))
	assertEqual(t, fstat3.ModTime(), mt3.UTC())

	_, err = fs.Open("test/dir")
	assertNotEqual(t, err, nil)
}

func TestCompileRetrieveError(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-assets")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	sources := []*Source{
		{"xxxx",
			"xxxx", nil, nil},
	}

	err = Compile(sources, filepath.Join(dir, "assets.go"), "assets", "fs", nil)
	assertEqual(t, reflect.TypeOf(err).String(), "*assets.RetrieveError")
	assertEqual(t, err.Error(), "xxxx: open xxxx: no such file or directory")
}

func TestCompileChecksumError(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-assets")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	sources := []*Source{
		{"retrieve_test.go",
			"retrieve_test.go", &Checksum{MD5, "1234"}, nil},
	}

	err = Compile(sources, filepath.Join(dir, "assets.go"), "assets", "fs", nil)
	assertEqual(t, err, &ChecksumError{"retrieve_test.go", ErrChecksumMismatch})
	assertEqual(t, err.Error(), "retrieve_test.go: checksum mismatch")
}

func TestCompileArchiveError(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-assets")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	sources := []*Source{
		{"retrieve_test.go",
			"retrieve_test.go", nil, &Archive{Zip, nil}},
	}

	err = Compile(sources, filepath.Join(dir, "assets.go"), "assets", "fs", nil)
	assertEqual(t, err, &ArchiveError{"retrieve_test.go", zip.ErrFormat})
	assertEqual(t, err.Error(), "retrieve_test.go: zip: not a valid zip file")
}

func assertEqual(t *testing.T, exp interface{}, act interface{}) {
	if reflect.DeepEqual(exp, act) {
		return
	}
	t.Fatal(fmt.Sprintf("%v != %v", exp, act))
}

func assertNotEqual(t *testing.T, exp interface{}, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		return
	}
	t.Fatal(fmt.Sprintf("%v == %v", exp, act))
}
