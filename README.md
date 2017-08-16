# go-assets

[![Build Status](https://travis-ci.org/ZoltanLajosKis/go-assets.svg?branch=master)](https://travis-ci.org/ZoltanLajosKis/go-assets)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZoltanLajosKis/go-assets)](https://goreportcard.com/report/github.com/ZoltanLajosKis/go-assets)
[![Coverage Status](https://coveralls.io/repos/github/ZoltanLajosKis/go-assets/badge.svg?branch=master)](https://coveralls.io/github/ZoltanLajosKis/go-assets?branch=master)
[![GoDoc](https://godoc.org/github.com/ZoltanLajosKis/go-assets?status.svg)](https://godoc.org/github.com/ZoltanLajosKis/go-assets)

Go-assets collects assets from different sources and transforms them into Go
code using [vfsgen][vfsgen]. The generated code implements a
[http.FileSystem][httpfs] that contains the assets. This can be used by HTTP
servers for serving those assets. Go-assets is best used with `go generate`
to generate source files before compiling the project.

Go-assets has the following features:
- retrieve assets via http(s) or from the local filesystem
- verify the checksum of assets
- extract and selectively filter files from archives


Usage
-----
Generate source code from assets with `go generate` before compiling.
```go
// +build ignore

package main

import (
  "io"
  "log"
  "strings"

  as "github.com/ZoltanLajosKis/go-assets"
)

var (
  sources = []*as.AssetSource{
    // Retrieve "assets/text.txt" from the local file system. Store it as
    // "docs/text1.txt".
    {"docs/text1.txt",
      "assets/text.txt", nil, nil},
    // Retrieve "remote.txt" from www.example.com. Verify the MD5 checksum.
    // Store it as "docs/text2.txt".
    {"docs/text2.txt",
      "https://www.example.com/remote.txt",
      &as.Checksum{as.ChecksumMD5, "1234567890abcdef1234567890abcdef"}, nil},
    // Retrieve "images.zip" from www.example.com. Verify the MD5 checksum.
    // Extract files from the archive, only keep files in the "img" directory
    // and store them in the "images" directory (see imagesMapper function).
    {"images.zip",
      "https://www.example.com/images/images.zip",
      &as.Checksum{as.ChecksumMD5, "1234567890abcdef1234567890abcdef"},
      &as.Archive{as.ArchiveZip, imagesMapper}},
  }

  imagesMapper = func(path string) string {
    if strings.HasPrefix(path, "img/") {
      return strings.Replace(path, "img/", "images/", 1)
    }
    return ""
  }
)

func main() {
  if err := as.Compile(sources,
    "assets/assets.go",   // output source file
    "assets",             // package name in the source file
    "FS",                 // variable name for the file system
    nil,
  ); err != nil {
    log.Panic(err)
  }
}
```

Example of using the generated source code to serve the files.
```go
package main

import (
  "net/http"
)

func main() {
  http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assets.FS)))
  http.ListenAndServe(":3000", nil)
}
```


[vfsgen]: https://github.com/shurcooL/vfsgen
[httpfs]: https://golang.org/pkg/net/http/#FileSystem
