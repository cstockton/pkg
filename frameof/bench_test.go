package frameof

import (
	"runtime"
	"testing"
)

func BenchmarkPC(b *testing.B) {
	b.Run(`Call`, func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fr := Call()
			if exp, got := 13, fr.Line; exp != got {
				b.Fatalf(`exp line %v; got %v`, exp, got)
			}
		}
	})
	b.Run(`Skip`, func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fr := Skip(0)
			if exp, got := 23, fr.Line; exp != got {
				b.Fatalf(`exp line %v; got %v`, exp, got)
			}
		}
	})
	b.Run(`runtime.Caller`, func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, Line, _ := runtime.Caller(0)
			if got := Line; got == 0 {
				b.Fatalf(`exp non-zero line num; got %v`, got)
			}
		}
	})
	b.Run(`Caller`, func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fr := Caller()
			if got := fr.Line; got == 0 {
				b.Fatalf(`exp non-zero line num; got %v`, got)
			}
		}
	})
}
