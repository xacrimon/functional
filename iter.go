package functional

import "reflect"

type Iter[T any] interface {
	Next() Option[T]
}

type sliceIter[T any] struct {
	slice []T
	index int
}

func (iter *sliceIter[T]) Next() Option[T] {
	if iter.index >= len(iter.slice) {
		return OptionNone[T]()
	}

	value := iter.slice[iter.index]
	iter.index++
	return OptionSome(value)
}

func SliceIter[T any](slice []T) Iter[T] {
	return &sliceIter[T]{slice, 0}
}

type mapIter[K, V any] struct {
	inner *reflect.MapIter
	index int
}

func (iter *mapIter[K, V]) Next() Option[Cons[K, V]] {
	if !iter.inner.Next() {
		return OptionNone[Cons[K, V]]()
	}

	key := iter.inner.Key().Interface().(K)
	value := iter.inner.Value().Interface().(V)
	return OptionSome(Cons[K, V]{key, value})
}

func MapIter[K comparable, V any](m map[K]V) Iter[Cons[K, V]] {
	return &mapIter[K, V]{reflect.ValueOf(m).MapRange(), 0}
}

func Count[T, I Iter[T]](iter I) int {
	var i int
	for OptionIsSome(iter.Next()) {
		i++
	}

	return i
}
