package internal

import "testing"

type tEnum int

const (
	tEnumUndef tEnum = iota
	tEnumOne
)

func (t *tEnum) String() string {
	switch *t {
	case tEnumOne:
		return "one"
	case tEnumUndef:
		fallthrough
	default:
		return "undef"
	}
}

func TestWrapEnumTypeJSON(t *testing.T) {
	t.Parallel()
	e := tEnumUndef
	_, err := WrapEnumTypeJSON(&e)
	if err != ErrInvalidEnumValue {
		t.Error("expected error")
	}
	e = tEnumOne
	b, err := WrapEnumTypeJSON(&e)
	if err != nil {
		t.Errorf("WrapEnumTypeJSON: %v", err)
	}
	if string(b) != `"one"` {
		t.Errorf("expected `one`, got %s", string(b))
	}
}
