package a_test

import (
	"testing"
	//"time"

	a "github.com/bluemon0919/go-training/paralleltest"
)

func TestIsOdd(t *testing.T) {
	t.Parallel()

	// time.Sleepはゴールーチンが切り替わるから-parall=1でも効果が出る
	// busyloopだとゴールーチンが切り替わらない
	//time.Sleep(1 * time.Second)

	busyloop()
	got, want := a.IsOdd(5), true
	if got != want {
		t.Errorf("want %v but got %v", want, got)
	}
}

func TestIsEven(t *testing.T) {
	t.Parallel()

	//time.Sleep(1 * time.Second)

	busyloop()
	got, want := a.IsEven(5), false
	if got != want {
		t.Errorf("want %v but got %v", want, got)
	}
}

func busyloop() {
	var n int
	for i := 0; i <= 10000; i++ {
		for j := 0; j <= 100000; j++ {n++}
	}
}
