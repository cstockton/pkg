package ctxkey

import (
	"context"
	"reflect"
	"testing"
)

func set(ctx context.Context, k Key, val interface{}) context.Context {
	return context.WithValue(ctx, k, val)
}

func get(ctx context.Context, k Key) interface{} {
	return ctx.Value(k)
}

func has(ctx context.Context, k Key) bool {
	return get(ctx, k) != nil
}

func TestDuplicates(t *testing.T) {
	ctx := context.Background()
	k1, k2 := New(`Dup`), New(`Dup`)
	T1, T2 := k1.(*key), k2.(*key) // fixes go-vet false positives
	v1, v2 := `v1`, `v2`

	// k1 must not equal k2
	if k1 == k2 {
		t.Fatalf(`must not have equality between %v (%[1]p) and %v (%[2]p)`, T1, T2)
	}

	ctx = set(ctx, k1, v1)
	if !has(ctx, k1) {
		t.Fatalf(`must find %v with %v (%[2]p)`, v1, T1)
	}
	if has(ctx, k2) {
		t.Fatalf(`must not find %v for %v (%[2]p) within %v (%[3]p)`, v1, T1, T2)
	}
	if !reflect.DeepEqual(get(ctx, k1), v1) {
		t.Fatalf(`must find %v in %v (%[2]p)`, v1, T1)
	}

	ctx = set(ctx, k2, v1)
	if !has(ctx, k1) {
		t.Fatalf(`must find %v with %v (%[2]p)`, v1, T1)
	}
	if !has(ctx, k2) {
		t.Fatalf(`must find %v with %v (%[2]p)`, v2, T2)
	}
	if !reflect.DeepEqual(get(ctx, k1), get(ctx, T2)) {
		t.Fatalf(`must not have equality between %v (%[1]p) and %v (%[2]p)`, T1, T2)
	}
}

func TestKeyAllocation(t *testing.T) {
	bg := context.Background()

	// common go idiom
	t.Run(`IntKey`, func(t *testing.T) {
		type MyKey int
		exp, k := 2, MyKey(0)

		var ctx context.Context
		allocs := int(testing.AllocsPerRun(1000, func() {
			ctx = context.WithValue(bg, k, nil)
		}))
		_ = ctx
		if exp != allocs {
			t.Fatalf(`exp %d allocs; got %d`, exp, allocs)
		}
	})

	// using key by val
	t.Run(`KeyByVal`, func(t *testing.T) {
		exp, k := 2, key{`TestKey`}

		var ctx context.Context
		allocs := int(testing.AllocsPerRun(1000, func() {
			ctx = context.WithValue(bg, k, nil)
		}))
		_ = ctx
		if exp != allocs {
			t.Fatalf(`exp %d allocs; got %d`, exp, allocs)
		}
	})

	// using pointer to key should not alloc
	t.Run(`KeyByPtr`, func(t *testing.T) {
		exp, k := 1, &key{`TestKey`}

		var ctx context.Context
		allocs := int(testing.AllocsPerRun(1000, func() {
			ctx = context.WithValue(bg, k, nil)
		}))
		_ = ctx
		if exp != allocs {
			t.Fatalf(`exp %d allocs; got %d`, exp, allocs)
		}
	})

	// above identical to
	t.Run(`IfcByPtr`, func(t *testing.T) {
		var k Key = &key{`TestKey`}
		exp := 1

		var ctx context.Context
		allocs := int(testing.AllocsPerRun(1000, func() {
			ctx = context.WithValue(bg, k, nil)
		}))
		_ = ctx
		if exp != allocs {
			t.Fatalf(`exp %d allocs; got %d`, exp, allocs)
		}
	})

	// verify with New()
	t.Run(`KeyByNew`, func(t *testing.T) {
		exp, k := 1, New(`TestKey`)

		var ctx context.Context
		allocs := int(testing.AllocsPerRun(1000, func() {
			ctx = context.WithValue(bg, k, nil)
		}))
		_ = ctx
		if exp != allocs {
			t.Fatalf(`exp %d allocs; got %d`, exp, allocs)
		}
	})
}
