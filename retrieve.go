package assets

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	// ErrNoMatch is returned when the glob does not match any file.
	ErrNoMatch = errors.New("no match")
)

func retrieve(loc string) ([]*file, error) {
	if strings.HasPrefix(loc, "http://") || strings.HasPrefix(loc, "https://") {
		return retrieveHTTP(loc)
	}

	if hasMeta(loc) {
		return retrieveGlob(loc)
	}

	return retrieveFile(loc)
}

func retrieveHTTP(url string) ([]*file, error) {
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

	return []*file{&file{url, data, modTime}}, nil
}

func retrieveFile(loc string) ([]*file, error) {
	f, err := os.Open(filepath.FromSlash(loc))
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

	return []*file{&file{loc, data, modTime}}, nil
}

func retrieveGlob(loc string) ([]*file, error) {
	// find longest prefix not containing globs
	dirs := strings.Split(loc, "/")
	i := 0

	for ; i < len(dirs); i++ {
		if hasMeta(dirs[i]) {
			break
		}
	}

	root := strings.Join(dirs[:i], "/") + "/"

	matches, err := filepath.Glob(filepath.FromSlash(loc))
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, ErrNoMatch
	}

	files := []*file{}

	for _, match := range matches {
		path := strings.TrimPrefix(filepath.ToSlash(match), root)

		f, err := os.Open(match)
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			f.Close()
			return nil, err
		}

		modTime := time.Now()
		info, err := f.Stat()
		if err == nil {
			modTime = info.ModTime()
		}

		f.Close()

		files = append(files, &file{path, data, modTime})
	}

	return files, nil
}

func hasMeta(path string) bool {
	return strings.ContainsAny(path, "*?[")
}
