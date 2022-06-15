package functional

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceIter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := SliceIter(slice)

	for _, expected := range slice {
		got := iter.Next()
		require.Equal(t, OptionSome(expected), got)
	}
}

func TestIterMap(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	expected := []int{4, 8, 12, 16, 20}
	f := func(i int) int { return i * 2 }
	iter := SliceIter(slice)
	iter = IterMap(iter, f)
	iter = IterMap(iter, f)

	for _, expected := range expected {
		got := iter.Next()
		require.Equal(t, OptionSome(expected), got)
	}
}

func TestIterSliceCollect(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4, 6, 8, 10}
	iter := SliceIter(slice)
	iter = IterMap(iter, func(i int) int { return i * 2 })
	got := IterFold(iter, nil, IterCollectSlice[int])
	require.Equal(t, expected, got)
}

func TestIterate(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := SliceIter(slice)

	i := 0
	iter2, cancel := Iterate(iter)
	defer cancel()
	for got := range iter2 {
		require.Equal(t, slice[i], got)
		i++
	}
}

func TestFuncAndFilterAndTake(t *testing.T) {
	x := 0
	f := func() Option[int] {
		x++
		return OptionSome(x)
	}

	numbers := IterFromFunc(f)
	even := IterFilter(numbers, func(i *int) bool { return *i%2 == 0 })
	firstFive := IterTake(even, 5)
	expected := []int{2, 4, 6, 8, 10}

	for _, expected := range expected {
		got := firstFive.Next()
		require.Equal(t, OptionSome(expected), got)
	}

	require.True(t, OptionIsNone(firstFive.Next()))
}

func TestFuncAndFilterMap(t *testing.T) {
	x := 0
	f := func() Option[int] {
		x++
		return OptionSome(x)
	}

	numbers := IterFromFunc(f)
	evenPlusOne := IterFilterMap(numbers, func(i int) Option[int] {
		if i%2 == 0 {
			return OptionSome(i + 1)
		}

		return OptionNone[int]()
	})

	expected := []int{3, 5, 7, 9, 11}

	for _, expected := range expected {
		got := evenPlusOne.Next()
		require.Equal(t, OptionSome(expected), got)
	}
}

func TestSkip(t *testing.T) {
	x := 0
	f := func() Option[int] {
		x++
		return OptionSome(x)
	}

	numbers := IterFromFunc(f)
	skipped := IterSkip(numbers, 5)
	expected := []int{6, 7, 8, 9, 10}

	for _, expected := range expected {
		got := skipped.Next()
		require.Equal(t, OptionSome(expected), got)
	}
}

func TestSkipWhile(t *testing.T) {
	x := 0
	f := func() Option[int] {
		x++
		return OptionSome(x)
	}

	numbers := IterFromFunc(f)
	skipped := IterSkipWhile(numbers, func(i *int) bool { return *i < 10 })
	expected := []int{10, 11, 12, 13, 14}

	for _, expected := range expected {
		got := skipped.Next()
		require.Equal(t, OptionSome(expected), got)
	}
}

func TestZip(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{6, 7, 8, 9, 10}
	expected := []int{7, 9, 11, 13, 15}

	iterA := SliceIter(a)
	iterB := SliceIter(b)
	iter := IterZip(iterA, iterB)

	for _, expected := range expected {
		got := OptionMap(iter.Next(), func(t Tuple2[int, int]) int { return t.A + t.B })
		require.Equal(t, OptionSome(expected), got)
	}
}
