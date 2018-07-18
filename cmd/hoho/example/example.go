package example

import "context"

type Example interface {
	Foo(ctx context.Context, a int) error
	Bar(i int, j string) (int, error)
	Baz(i, k, j int) (x, y string)
	BarBaz() int
}
