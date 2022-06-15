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

func IterCollectSlice[T any](acc []T, item T) []T {
	return append(acc, item)
}

func IterCollectMap[K comparable, V any](acc map[K]V, item Tuple2[K, V]) map[K]V {
	if acc == nil {
		acc = make(map[K]V)
	}

	acc[item.A] = item.B
	return acc
}

func IterFold[Acc, T any, I Iter[T]](iter I, f func(Acc, T) Acc) Acc {
	var acc Acc
	ForEach(iter, func(item T) { acc = f(acc, item) })
	return acc
}

func Iterate[T any, I Iter[T]](iter I) (<-chan T, func()) {
	ch := make(chan T)
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

type iterFunc[T any] struct {
	f func() Option[T]
}

func (iter *iterFunc[T]) Next() Option[T] {
	return iter.f()
}

func IterFromFunc[T any](f func() Option[T]) Iter[T] {
	return &iterFunc[T]{f}
}

type iterFilter[T any] struct {
	inner     Iter[T]
	predicate func(*T) bool
}

func (iter *iterFilter[T]) Next() Option[T] {
	for {
		item := iter.inner.Next()
		if OptionIsNone(item) || iter.predicate(&item.value) {
			return item
		}
	}
}

func IterFilter[T any](iter Iter[T], predicate func(*T) bool) Iter[T] {
	return &iterFilter[T]{iter, predicate}
}

type iterTake[T any] struct {
	inner Iter[T]
	left  int
}

func (iter *iterTake[T]) Next() Option[T] {
	if iter.left > 0 {
		iter.left--
		return iter.inner.Next()
	}

	return OptionNone[T]()
}

func IterTake[T any](iter Iter[T], amount int) Iter[T] {
	return &iterTake[T]{iter, amount}
}

type iterFilterMap[T, U any] struct {
	inner Iter[T]
	f     func(T) Option[U]
}

func (iter *iterFilterMap[T, U]) Next() Option[U] {
	for {
		item := iter.inner.Next()
		if OptionIsNone(item) {
			return OptionNone[U]()
		}

		mapped := iter.f(item.value)
		if OptionIsSome(mapped) {
			return mapped
		}
	}
}

func IterFilterMap[T, U any](iter Iter[T], f func(T) Option[U]) Iter[U] {
	return &iterFilterMap[T, U]{iter, f}
}

type iterSkip[T any] struct {
	inner Iter[T]
	left  int
}

func (iter *iterSkip[T]) Next() Option[T] {
	for iter.left > 0 {
		iter.left--
		iter.inner.Next()
	}

	return iter.inner.Next()
}

func IterSkip[T any](iter Iter[T], amount int) Iter[T] {
	return &iterSkip[T]{iter, amount}
}

type iterSkipWhile[T any] struct {
	inner     Iter[T]
	predicate func(*T) bool
}

func (iter *iterSkipWhile[T]) Next() Option[T] {
	for iter.predicate != nil {
		item := iter.inner.Next()
		if OptionIsNone(item) {
			return item
		}

		if !iter.predicate(&item.value) {
			iter.predicate = nil
			return item
		}
	}

	return iter.inner.Next()
}

func IterSkipWhile[T any](iter Iter[T], predicate func(*T) bool) Iter[T] {
	return &iterSkipWhile[T]{iter, predicate}
}

type iterZip[T1, T2 any] struct {
	iter1 Iter[T1]
	iter2 Iter[T2]
}

func (iter *iterZip[T1, T2]) Next() Option[Tuple2[T1, T2]] {
	item1 := iter.iter1.Next()
	item2 := iter.iter2.Next()

	if OptionIsNone(item1) || OptionIsNone(item2) {
		return OptionNone[Tuple2[T1, T2]]()
	}

	return OptionSome(Tuple2[T1, T2]{item1.value, item2.value})
}

func IterZip[T1, T2 any](iter1 Iter[T1], iter2 Iter[T2]) Iter[Tuple2[T1, T2]] {
	return &iterZip[T1, T2]{iter1, iter2}
}
