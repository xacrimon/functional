package functional

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterSlice(t *testing.T) {
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
	got := IterFold(iter, IterIntoSlice[int])
	require.Equal(t, expected, got)
}

func TestIterate(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := SliceIter(slice)

	i := 0
	iter2, cancel := Iterate[int](iter)
	defer cancel()
	for got := range iter2 {
		require.Equal(t, slice[i], got)
		i++
	}
}
