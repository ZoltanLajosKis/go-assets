package assets

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func retrieve(loc string) (*file, error) {
	if strings.HasPrefix(loc, "http://") || strings.HasPrefix(loc, "https://") {
		return retrieveHTTP(loc)
	}
	return retrieveFS(loc)
}

func retrieveHTTP(url string) (*file, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http status: " + string(resp.StatusCode))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	modTime, err := http.ParseTime(resp.Header.Get("Last-Modified"))
	if err != nil {
		modTime = time.Now()
	}

	return &file{url, data, modTime}, nil
}

func retrieveFS(loc string) (*file, error) {
	f, err := os.Open(loc)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	modTime := time.Now()
	info, err := f.Stat()
	if err == nil {
		modTime = info.ModTime()
	}

	return &file{loc, data, modTime}, nil
}
