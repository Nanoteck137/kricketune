package types

type WorkDir string

func (d WorkDir) String() string {
	return string(d)
}

type Change[T any] struct {
	Value   T
	Changed bool
}
