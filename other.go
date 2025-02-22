package genh

import (
	"reflect"

	"github.com/alpineiq/genh/internal"
)

type (
	Signed      = internal.Signed
	Unsigned    = internal.Unsigned
	Integer     = internal.Integer
	Float       = internal.Float
	Complex     = internal.Complex
	Ordered     = internal.Ordered
	EncoderType = internal.EncoderType
	DecoderType = internal.DecoderType
)

// ValuesToPtrs converts a slice of values to a slice of pointers
// optionally copying the values instead of pointing to them in the original slice.
func ValuesToPtrs[T any](vals []T, copy bool) []*T {
	out := make([]*T, 0, len(vals))
	for i := range vals {
		var v *T
		if copy {
			cp := vals[i]
			v = &cp
		} else {
			v = &vals[i]
		}
		out = append(out, v)
	}
	return out
}

func PtrsToValues[T any](vals []*T) []T {
	out := make([]T, 0, len(vals))
	for i := range vals {
		v := vals[i]
		out = append(out, *v)
	}
	return out
}

func Iff[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func IffFn[T any](cond bool, a, b func() T) T {
	if cond {
		return a()
	}
	return b()
}

func FirstNonZero[T any](args ...T) T {
	for _, arg := range args {
		if v := reflect.ValueOf(arg); v.IsValid() && !v.IsZero() {
			return arg
		}
	}
	return args[0]
}

func FirstNonZeroPtr[T any](args ...*T) *T {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}
	return nil
}

func FirstNonZeroCmp[T comparable](args ...T) T {
	var zero T
	for _, arg := range args {
		if arg != zero {
			return arg
		}
	}
	return zero
}
