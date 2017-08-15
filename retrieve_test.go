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
	data, _, err := retrieve(strings.Join([]string{svr.URL, "/assets.txt"}, ""))
	assertEqual(t, err, nil)
	assertEqual(t, data, []byte("Assets."))

	_, _, err = retrieve(strings.Join([]string{svr.URL, "/xxxx"}, ""))
	assertNotEqual(t, err, nil)

	_, _, err = retrieve("http://invalid.u.r.l")
	assertNotEqual(t, err, nil)
	assertEqual(t, strings.Contains(err.Error(), "http://invalid.u.r.l"), true)
}

func TestRetrieveFile(t *testing.T) {
	data, _, err := retrieve("retrieve_test.go")
	assertEqual(t, err, nil)
	assertEqual(t, len(data) > 0, true)
}

func TestRetrieveFileNotExist(t *testing.T) {
	_, _, err := retrieve("xxxx")
	assertNotEqual(t, err, nil)
}
