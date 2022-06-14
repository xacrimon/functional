package functional

// TODO:
// - OptionOkOr*
// - OptionIter*
// - Result
// - Iterator
// - IterAdvanceBy
// - Iter collect to result

type Option[T any] struct {
	value T
	some  bool
}

func OptionSome[T any](value T) Option[T] {
	return Option[T]{value, true}
}

func OptionNone[T any]() Option[T] {
	return Option[T]{some: false}
}

func OptionIsSome[T any](opt Option[T]) bool {
	return opt.some
}

func OptionIsNone[T any](opt Option[T]) bool {
	return !opt.some
}

func OptionValue[T any](opt Option[T]) T {
	return opt.value
}

func AsRef[T any](opt *Option[T]) Option[*T] {
	if !opt.some {
		return OptionNone[*T]()
	}

	return OptionSome(&opt.value)
}

func OptionUnwrap[T any](opt Option[T]) T {
	if !opt.some {
		panic("unwrap on none")
	}

	return opt.value
}

func OptionExpect[T any](opt Option[T], message string) T {
	if !opt.some {
		panic(message)
	}

	return opt.value
}

func OptionUnwrapOr[T any](opt Option[T], def T) T {
	if !opt.some {
		return def
	}

	return opt.value
}

func OptionUnwrapOrElse[T any](opt Option[T], def func() T) T {
	if !opt.some {
		return def()
	}

	return opt.value
}

func OptionMap[T, U any](opt Option[T], transform func(T) U) Option[U] {
	if !opt.some {
		return OptionNone[U]()
	}

	return OptionSome(transform(opt.value))
}

func OptionMapOr[T, U any](opt Option[T], transform func(T) U, def U) U {
	if !opt.some {
		return def
	}

	return transform(opt.value)
}

func OptionMapOrElse[T, U any](opt Option[T], transform func(T) U, def func() U) U {
	if !opt.some {
		return def()
	}

	return transform(opt.value)
}

func OptionAnd[T any](optA Option[T], optB Option[T]) Option[T] {
	if !optA.some {
		return optA
	}

	return optB
}

func OptionAndThen[T, U any](opt Option[T], then func(T) Option[U]) Option[U] {
	if !opt.some {
		return OptionNone[U]()
	}

	return then(opt.value)
}

func OptionFilter[T any](opt Option[T], predicate func(T) bool) Option[T] {
	if !opt.some || !predicate(opt.value) {
		return OptionNone[T]()
	}

	return opt
}

func OptionOr[T any](optA Option[T], optB Option[T]) Option[T] {
	if optA.some {
		return optA
	}

	return optB
}

func OptionElse[T any](optA Option[T], optB func() Option[T]) Option[T] {
	if optA.some {
		return optA
	}

	return optB()
}

func OptionXor[T any](optA Option[T], optB Option[T]) Option[T] {
	switch {
	case optA.some && !optB.some:
		return optA
	case !optA.some && optB.some:
		return optB
	default:
		return OptionNone[T]()
	}
}

func OptionInsert[T any](opt *Option[T], value T) *T {
	opt.value = value
	opt.some = true
	return &opt.value
}

func OptionGetOrInsert[T any](opt *Option[T], value T) *T {
	if !opt.some {
		opt.value = value
		opt.some = true
	}

	return &opt.value
}

func OptionGetOrInsertWith[T any](opt *Option[T], value func() T) *T {
	if !opt.some {
		opt.value = value()
		opt.some = true
	}

	return &opt.value
}

func OptionTake[T any](opt *Option[T]) Option[T] {
	if !opt.some {
		return OptionNone[T]()
	}

	opt.some = false
	return OptionSome(opt.value)
}

func OptionReplace[T any](opt *Option[T], value T) Option[T] {
	old := *opt
	*opt = OptionSome(value)
	return old
}

func OptionContains[T comparable](opt Option[T], value T) bool {
	return opt.some && opt.value == value
}

func OptionZip[T1, T2 any](optA Option[T1], optB Option[T2]) Option[Cons[T1, T2]] {
	if !optA.some || !optB.some {
		return OptionNone[Cons[T1, T2]]()
	}

	return OptionSome(Cons[T1, T2]{optA.value, optB.value})
}

func OptionUnzip[T1, T2 any](opt Option[Cons[T1, T2]]) (Option[T1], Option[T2]) {
	if !opt.some {
		return OptionNone[T1](), OptionNone[T2]()
	}

	return OptionSome(opt.value.A), OptionSome(opt.value.B)
}

func OptionZipWith[T1, T2, U any](optA Option[T1], optB Option[T2], f func(T1, T2) U) Option[U] {
	if !optA.some || !optB.some {
		return OptionNone[U]()
	}

	return OptionSome(f(optA.value, optB.value))
}

func OptionFlatten[T any](opt Option[Option[T]]) Option[T] {
	if !opt.some {
		return OptionNone[T]()
	}

	return opt.value
}
