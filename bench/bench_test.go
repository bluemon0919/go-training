package a

import "testing"

func BenchmarkWithDefer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithDefer()
	}
}

func BenchmarkWithoutDefer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithoutDefer()
	}
}
