# go-assets

[![Build Status](https://travis-ci.org/ZoltanLajosKis/go-assets.svg?branch=master)](https://travis-ci.org/ZoltanLajosKis/go-assets)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZoltanLajosKis/go-assets)](https://goreportcard.com/report/github.com/ZoltanLajosKis/go-assets)
[![Coverage Status](https://coveralls.io/repos/github/ZoltanLajosKis/go-assets/badge.svg?branch=master)](https://coveralls.io/github/ZoltanLajosKis/go-assets?branch=master)
[![GoDoc](https://godoc.org/github.com/ZoltanLajosKis/go-assets?status.svg)](https://godoc.org/github.com/ZoltanLajosKis/go-assets)

Go-assets collects assets from multiple sources. Assets can be stored in
memory as a [http.FileSystem][httpfs] and provided for your application. They
can also be combined into a go source code using [vfsgen][vfsgen]. This is
best used with `go generate` to embed the assets into binaries.

Usage
-----
The `Retrieve` function collects assets into memory and returns them as a
[http.FileSystem][httpfs]. The `Compile` function creates a go source file
from the collected assets. This source code contains a single variable that
points to the file system.
```go
func Retrieve(sources []*Source) (http.FileSystem, error)

func Compile(sources []*Source, filePath string, pkgName string, varName string, opts *Opts) error
```
With `Compile`, the `filePath` argument specifies the location of the asset
source. `pkgName` and `varName` specifies the package and variable name to use.
The optional `opts` parameter can specify build tags to be added to the source
file (`BuildTags`) and a custom comment text for the variable
(`VariableComment`).

Each asset source is described with the below structure.
```go
type Source struct {
  Path     string
  Location string
  Checksum *Checksum
  Archive  *Archive
}
```
Here `Path` tells the path of the resulting asset(s) in the output file system
(details below), while `Location`, `Checksum` and `Archive` each correspond to
one of the threep rocessing steps below.

Path strings in `Path` and other field should only use forward slashes (`/`)
for separator.

### 1. Retrieval
The asset source file is retrieved from the specified `Location`.

If `Location` starts with `http://` or `https://`, it is considered a URL, and
the file is downloaded.  
If `Location` contains a [glob pattern][globpattern], the pattern is applied
to the local file system, and all matching files are retrieved.  
Otherwise, `Location` is assumed to contain a file path, and that file is
retrieved from the local file system.

If multiple files were retrieved (this can only happen when using a
[glob pattern][globpattern]), processing stops here. The path for these assets
is calculated as follows.

1. take the longest path prefix in `Location` that does not contain glob
   patterns
2. trim this prefix from the file paths (as matched by the pattern)
3. prefix the previous value with the value in `Path`

As an example, if `Path` is set to `"assets"`, and `Location` is set to
`"input/data/data[0-9].txt"`, the resulting files will be located at
`"assets/data0.txt"`, `"assets/data1.txt"`, and so on.

If only a single file was retrieved, processing continues. If the `Checksum`
field is not `nil`, the file checksum is verified and processing halts with an
error on mismatch. Then, if `Archive` is not `nil`, the file is processed
as an archive.

If `Archive` is `nil`, processing stops here, and thefile is stored at the
path specified by `Path`.


### 2. Checksum
Checksum verification can be requested with the following structure.
```go
type Checksum struct {
	Algo  ChecksumAlgo
	Value string
}
```
Here `Algo` specifies the checksum algorithm to use, and `Value` is the
precalculated checksum value. If a mismatch is found during verification,
an error is returned.
The currently supported checksum algorithms are: `MD5`, `SHA1`, `SHA256` and
`SHA512`.


### 3. Archive extraction
Archive extraction can be requested with the following structure.
```go
type Archive struct {
	Format     ArchiveFormat
	PathMapper PathMapper
}
```
Here `Format` specifies the type of archive. Currently `Zip` and `TarGz` are
supported.

The `PathMapper` function can be used to filter files from the archive and
specify custom paths for them. If it is set to `nil`, all files are kept and
are stored at the same path they were found in the archive.
```go
type PathMapper func(string) string
```
The `PathMapper` function is invoked for each file path in the archive. If it
returns `""` the file is dropped. Otherwise the file is kept and stored at the
path returned by the function.

For common cases the `ReMap` function can be used for generating a `PathMapper`
function.
```go
func ReMap(pattern string, replacement string) PathMapper
```
The `pattern` argument takes a [regular expression][re]. This pattern is
matched against all file paths in the archive and only matching files are kept.
For these files the `replacement` string specifies the storage path.
In order to ensure unique paths for files, the pattern should contain
capturing groups and the replacement string should contain backreferences.


Example
-------
See a usage example in the [gmdd][gmdd] project.



[vfsgen]: https://github.com/shurcooL/vfsgen
[httpfs]: https://golang.org/pkg/net/http/#FileSystem
[globpattern]: https://golang.org/pkg/path/filepath/#Match
[re]: https://github.com/google/re2/wiki/Syntax
[gmdd]: https://github.com/ZoltanLajosKis/gmdd/blob/master/generate/assets.go#L12
