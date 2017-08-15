package assets

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
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

	assets := []*AssetSource{
		{"assets.txt",
			strings.Join([]string{svr.URL, "/assets.txt"}, ""), nil, nil},
		{"retrieve_test.go",
			"retrieve_test.go", nil, nil},
	}

	err = Compile(assets, filepath.Join(dir, "assets.go"), "assets", "fs")
	assertEqual(t, err, nil)
}

func TestGenerateRetrieveError(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-assets")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	assets := []*AssetSource{
		{"xxxx",
			"xxxx", nil, nil},
	}

	err = Compile(assets, filepath.Join(dir, "assets.go"), "assets", "fs")
	assertNotEqual(t, err, nil)
}
func TestGenerateChecksumError(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-assets")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	assets := []*AssetSource{
		{"retrieve_test.go",
			"retrieve_test.go", &Checksum{ChecksumMD5, "1234"}, nil},
	}

	err = Compile(assets, filepath.Join(dir, "assets.go"), "assets", "fs")
	assertNotEqual(t, err, nil)
}

func TestGenerateArchiveError(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-assets")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	assets := []*AssetSource{
		{"retrieve_test.go",
			"retrieve_test.go", nil, &Archive{ArchiveZip, nil}},
	}

	err = Compile(assets, filepath.Join(dir, "assets.go"), "assets", "fs")
	assertNotEqual(t, err, nil)
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