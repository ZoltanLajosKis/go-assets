package assets

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRetrieveHttp(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/assets.txt" {
			fmt.Fprintf(w, "Assets.")
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer svr.Close()
	file, err := retrieve(strings.Join([]string{svr.URL, "/assets.txt"}, ""))
	assertEqual(t, err, nil)
	assertEqual(t, file.data, []byte("Assets."))

	_, err = retrieve(strings.Join([]string{svr.URL, "/xxxx"}, ""))
	assertNotEqual(t, err, nil)

	_, err = retrieve("http://invalid.u.r.l")
	assertNotEqual(t, err, nil)
	assertEqual(t, strings.Contains(err.Error(), "http://invalid.u.r.l"), true)
}

func TestRetrieveFile(t *testing.T) {
	file, err := retrieve("retrieve_test.go")
	assertEqual(t, err, nil)
	assertEqual(t, len(file.data) > 0, true)
}

func TestRetrieveFileNotExist(t *testing.T) {
	_, err := retrieve("xxxx")
	assertNotEqual(t, err, nil)
}
