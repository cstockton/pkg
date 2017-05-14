package ctxkey_test

import (
	"context"
	"fmt"
	"testing"

	ctxkey "github.com/cstockton/pkg/ctxkey"
)

func Example() {
	var (
		AppKey  = ctxkey.New(`mypkg.App`)
		LogKey  = ctxkey.New(`mypkg.App.Log`)
		AuthKey = ctxkey.New(`mypkg.App.Auth`)
		UserKey = ctxkey.New(`mypkg.App.User`)
	)

	fmt.Println(AppKey)
	fmt.Println(` ->`, LogKey)
	fmt.Println(` ->`, AuthKey)
	fmt.Println(` ->`, UserKey)

	// Output:
	// Key(mypkg.App)
	//  -> Key(mypkg.App.Log)
	//  -> Key(mypkg.App.Auth)
	//  -> Key(mypkg.App.User)
}

func Example_uniqueness() {
	k1, k2 := ctxkey.New(`mypkg.App`), ctxkey.New(`mypkg.App`)
	fmt.Print(k1 != k2)

	// Output:
	// true
}

func Example_allocations() {
	type IntKey int
	intKey, ctxKey := IntKey(0), ctxkey.New(`MyKey`)

	ctx := context.Background()
	fmt.Println(`IntKey allocs:`, testing.AllocsPerRun(1000, func() {
		_ = context.WithValue(ctx, intKey, nil)
	}))
	fmt.Println(`ctxkey allocs:`, testing.AllocsPerRun(1000, func() {
		_ = context.WithValue(ctx, ctxKey, nil)
	}))

	// Output:
	// IntKey allocs: 2
	// ctxkey allocs: 1
}
