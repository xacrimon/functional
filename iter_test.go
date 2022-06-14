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
