// Package ctxkey is like errors.New() for context keys.
//
// Invariants
//
// All context keys created with New() will be globally unique during program
// execution for the lifespan of the Key and will not cause additional
// allocations when they are associated with a context value. Uniqueness is
// satisfied by returning pointers to a private key struct, forcing comparison
// operations in context value retrieval to use pointer equality within the
// current Go's specification:
//
//   https://golang.org/ref/spec#Comparison_operators
//
// Additional allocations means that you will not allocate for the context key
// as seen with this common idiom:
//
//   type key int
//   var myKey key = 0
//   // 3 allocations, context struct, myKey and someValue
//   context.WithValue(ctx, myKey, SomeValue)
//
// Keys issued by New() fit within an interface{} value, i.e.:
//
//   var myKey = ctxkey.New(`MyKey`)
//   // 2 allocations, context struct and someValue
//   context.WithValue(ctx, myKey, SomeValue)
//
package ctxkey

// A Key holds a globally unique named value to associate with context Values.
//
// Name will return a non-unique string for documentation and troubleshooting
// purposes. It is not likely to be unique as your app crosses package
// boundaries.
type Key interface {
	Name() string
}

// New returns a globally unique context Key that will not cause allocations
// associated with the given name.
func New(name string) Key {
	return &key{name}
}

type key struct{ name string }

func (k *key) Name() string {
	return k.name
}

const (
	lTag = `Key(`
	rTag = `)`
)

func (k key) String() string {
	return lTag + k.Name() + rTag
}
