# Go Package: ctxkey

  [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/cstockton/pkg/ctxkey)
  [![Go Report Card](https://goreportcard.com/badge/github.com/cstockton/pkg/ctxkey?style=flat-square)](https://goreportcard.com/report/github.com/cstockton/pkg/ctxkey)
  [![Coverage Status](https://img.shields.io/codecov/c/github/cstockton/pkg/master.svg?style=flat-square)](https://codecov.io/gh/cstockton/pkg/src/master/ctxkey/ctxkey.go)

  > Get:
  > ```bash
  > go get -u github.com/cstockton/pkg/ctxkey
  > ```
  >
  > Example:
  > ```Go
  > var (
  >   AppKey  = ctxkey.New(`mypkg.App`)
  >   LogKey  = ctxkey.New(`mypkg.App.Log`)
  >   AuthKey = ctxkey.New(`mypkg.App.Auth`)
  >   UserKey = ctxkey.New(`mypkg.App.User`)
  > )
  > fmt.Println(AppKey)
  > fmt.Println(` ->`, LogKey)
  > fmt.Println(` ->`, AuthKey)
  > fmt.Println(` ->`, UserKey)
  > ```
  >
  > Output:
  > ```Go
  > Key(mypkg.App)
  >   -> Key(mypkg.App.Log)
  >   -> Key(mypkg.App.Auth)
  >   -> Key(mypkg.App.User)
  > ```


## Summary

Package ctxkey is like errors.New() for context keys. I created it to prevent
code duplication and provide clear intent for public context keys. It is very
small and has zero dependencies, so feel free to copy and paste it locally if
you would like.


### Uniqueness

  All context keys created with New() will be globally unique during program
  execution for the lifespan of the Key. Uniqueness is satisfied by returning
  pointers to a private key struct, forcing comparison operations in context
  value retrieval to use pointer equality within the current Go's specification
  for [comparison operators](https://golang.org/ref/spec#Comparison_operators).

  > Example:
  > ```Go
  > k1, k2 := ctxkey.New(`mypkg.App`), ctxkey.New(`mypkg.App`)
  > fmt.Print(k1 != k2)
  > ```
  >
  > Output:
  > ```Go
  > true
  > ```


### Allocations

  All context keys created with New() will not cause additional
  allocations when they are associated with a context value as seen with the
  common idiom of `type ctxKey int`. This is because Keys issued by New() fit
  within an interface{} value.

  > Example:
  > ```Go
  > type IntKey int
  > intKey, ctxKey := IntKey(0), ctxkey.New(`MyKey`)
  >
  > ctx := context.Background()
  > fmt.Println(`IntKey allocs:`, testing.AllocsPerRun(1000, func() {
  >   _ = context.WithValue(ctx, intKey, nil)
  > }))
  > fmt.Println(`ctxkey allocs:`, testing.AllocsPerRun(1000, func() {
  > _ = context.WithValue(ctx, ctxKey, nil)
  > }))
  > ```
  >
  > Output:
  > ```Go
  > Output:
  > IntKey allocs: 2
  > ctxkey allocs: 1
  > ```
