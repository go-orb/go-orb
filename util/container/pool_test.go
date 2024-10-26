package container

import (
	"testing"
)

func TestPool(t *testing.T) {
	type Foo struct{ n int }
	n := 42
	p := NewPool(func() Foo { return Foo{n} })
	v1 := p.Get()

	if got, want := v1.n, 42; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
	v1.n = 43
	p.Put(v1)

	n = 1234

	v2 := p.Get()
	if got, want := v2.n, 43; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	v3 := p.Get()
	if got, want := v3.n, 1234; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
}
