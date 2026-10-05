package pkg


type Result[T any] struct {
	Value T
	Err   error
}

func (r Result[T]) Map[K any](fn func (v T) K) Result[K] {
	if r.Err != nil {
		return Result[K]{Err: r.Err}
	}
	return Result[K]{Value: fn(r.Value)}
}

func (r Result[T]) Apply[K any](fn func (v T) (K, error)) Result[K] {
	if r.Err != nil {
		return Result[K]{Err: r.Err}
	}
	return FromTuple(fn(r.Value))
}

func FromTuple[T any](v T, err error) Result[T] {
	return Result[T]{Value: v, Err: err}
}
