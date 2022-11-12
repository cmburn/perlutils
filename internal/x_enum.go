package internal

import (
	"fmt"
	"strings"
)

type XEnumKind[T ~int] interface {
	~int
	fmt.Stringer
	Parse(string) (T, error)
}

type XEnum[T XEnumKind[T]] struct {
	Kind   T
	XValue string
}

func (f *XEnum[T]) MarshalJSON() ([]byte, error) {
	return WrapEnumTypeJSON(f.Kind)
}

func (f *XEnum[T]) String() string {
	if f.XValue != "" {
		return f.XValue
	}
	return f.Kind.String()
}

func (f *XEnum[T]) UnmarshalJSON(b []byte) error {
	s, err := UnwrapJSONString(b)
	if err != nil {
		return err
	}
	if strings.HasPrefix(s, "x_") {
		f.XValue = s
		return nil
	}
	var v T
	v, err = v.Parse(s)
	if err != nil {
		return err
	}
	f.Kind = v
	return nil
}
