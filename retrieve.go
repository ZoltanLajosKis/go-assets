package assets

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// RetrieveError is returned when there is a problem retrieving an asset source
type RetrieveError struct {
	Location string
	Err      error
}

func (e *RetrieveError) Error() string {
	return e.Location + ": " + e.Err.Error()
}

func retrieve(loc string) ([]byte, time.Time, error) {
	if strings.HasPrefix(loc, "http://") || strings.HasPrefix(loc, "https://") {
		return retrieveHTTP(loc)
	}
	return retrieveFS(loc)
}

func retrieveHTTP(loc string) ([]byte, time.Time, error) {
	resp, err := http.Get(loc)
	if err != nil {
		return nil, time.Time{}, &RetrieveError{loc, err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, time.Time{}, &RetrieveError{loc, errors.New(string(resp.StatusCode))}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, time.Time{}, &RetrieveError{loc, err}
	}

	modTime, err := http.ParseTime(resp.Header.Get("Last-Modified"))
	if err != nil {
		modTime = time.Now()
	}

	return data, modTime, nil
}

func retrieveFS(loc string) ([]byte, time.Time, error) {
	f, err := os.Open(loc)
	if err != nil {
		return nil, time.Time{}, &RetrieveError{loc, err}
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, time.Time{}, &RetrieveError{loc, err}
	}

	modTime := time.Now()
	info, err := f.Stat()
	if err == nil {
		modTime = info.ModTime()
	}

	return data, modTime, nil
}
