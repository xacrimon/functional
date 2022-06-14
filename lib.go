package goiter

type Iterator[T any] interface {
	next()
}
