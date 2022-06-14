package functional

import (
	"reflect"
	"sync"
)

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

func (iter *mapIter[K, V]) Next() Option[Tuple2[K, V]] {
	if !iter.inner.Next() {
		return OptionNone[Tuple2[K, V]]()
	}

	key := iter.inner.Key().Interface().(K)
	value := iter.inner.Value().Interface().(V)
	return OptionSome(Tuple2[K, V]{key, value})
}

func MapIter[K comparable, V any](m map[K]V) Iter[Tuple2[K, V]] {
	return &mapIter[K, V]{reflect.ValueOf(m).MapRange(), 0}
}

func Count[T, I Iter[T]](iter I) int {
	var i int
	for OptionIsSome(iter.Next()) {
		i++
	}

	return i
}

type iterMap[T, U any, I Iter[T]] struct {
	inner I
	f     func(T) U
}

func (iter *iterMap[T, U, I]) Next() Option[U] {
	item := iter.inner.Next()
	return OptionMap(item, iter.f)
}

func IterMap[T, U any, I Iter[T]](iter I, f func(T) U) Iter[U] {
	return &iterMap[T, U, I]{iter, f}
}

func ForEach[T any, I Iter[T]](iter I, f func(T)) {
	for {
		item := iter.Next()
		if OptionIsNone(item) {
			break
		}

		f(OptionValue(item))
	}
}

type CollectFunction[Acc, T any] func(Acc, T) Acc

func IterIntoSlice[T any](acc []T, item T) []T {
	return append(acc, item)
}

func IterIntoMap[K comparable, V any](acc map[K]V, item Tuple2[K, V]) map[K]V {
	if acc == nil {
		acc = make(map[K]V)
	}

	acc[item.A] = item.B
	return acc
}

func IterCollect[Acc, T any, I Iter[T]](iter I, f CollectFunction[Acc, T]) Acc {
	var acc Acc
	ForEach(iter, func(item T) { acc = f(acc, item) })
	return acc
}

func Iterate[T any, I Iter[T]](iter I) (<-chan T, func()) {
	ch := make(chan T, 1)
	once := &sync.Once{}
	close := func() { once.Do(func() { close(ch) }) }

	go func() {
		for {
			item := iter.Next()
			if OptionIsNone(item) {
				close()
				break
			}

			ch <- OptionValue(item)
		}
	}()

	return ch, close
}
