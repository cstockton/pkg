# Go Package: frameof

  [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/cstockton/pkg/frameof)
  [![Go Report Card](https://goreportcard.com/badge/github.com/cstockton/pkg/frameof?style=flat-square)](https://goreportcard.com/report/github.com/cstockton/pkg/frameof)
  [![Coverage Status](https://img.shields.io/codecov/c/github/cstockton/pkg/frameof/master.svg?style=flat-square)](https://codecov.io/github/cstockton/pkg/frameof?branch=master)

  > Get:
  > ```bash
  > go get -u github.com/cstockton/pkg/frameof
  > ```
  >
  > Example:
  > ```Go
  > // Basic types
  > fr := frameof.Call()
  >
  > fmt.Printf("Frame in %v() at %v:%d for package %v.\n",
  >     fr.Name(), filepath.Base(fr.File), fr.Line, fr.Pkg())
  > ```
  >
  > Output:
  > ```Go
  > Frame in Example() at example_test.go:12 for package frameof_test.
  > ```


## Summary

Package frameof provides a few simple helpers to access call frames from the
perspective of this packages caller.
