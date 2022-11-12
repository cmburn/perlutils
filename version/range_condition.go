package version

import (
	"errors"
	"strings"

	// local
	"github.com/cmburn/perlutils/internal"
)

type RangeCondition int32

const (
	rangeConditionUndef RangeCondition = iota // unexported
	RangeConditionGreaterThan
	RangeConditionGreaterThanOrEqual
	RangeConditionLessThan
	RangeConditionLessThanOrEqual
	RangeConditionNotEqual
	RangeConditionEqual
	rangeConditionNone // implicit equal, unexported
)

func (c *RangeCondition) String() string {
	if c == nil {
		// not possible unless someone is messing with the internals
		panic("nil RangeCondition")
	}
	switch *c {
	case RangeConditionGreaterThan:
		return ">"
	case RangeConditionGreaterThanOrEqual:
		return ">="
	case RangeConditionLessThan:
		return "<"
	case RangeConditionLessThanOrEqual:
		return "<="
	case RangeConditionNotEqual:
		return "!="
	case RangeConditionEqual:
		return "=="
	case rangeConditionNone, rangeConditionUndef:
		return ""
	default:
		// also not possible
		panic("invalid RangeCondition")
	}
}

func (c *RangeCondition) URLString() string {
	if *c == rangeConditionNone {
		return "=="
	}
	return c.String()
}

func rangeTypeFromString(s string) (RangeCondition, error) {
	switch s {
	case ">":
		return RangeConditionGreaterThan, nil
	case ">=":
		return RangeConditionGreaterThanOrEqual, nil
	case "<":
		return RangeConditionLessThan, nil
	case "<=":
		return RangeConditionLessThanOrEqual, nil
	case "!=":
		return RangeConditionNotEqual, nil
	case "=":
		return RangeConditionEqual, nil
	case "":
		return rangeConditionNone, nil
	}
	return rangeConditionUndef, ErrInvalidRangeCondition
}

func (c *RangeCondition) MarshalJSON() ([]byte, error) {
	return internal.WrapEnumTypeJSON(c)
}

func (c *RangeCondition) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	r, err := rangeTypeFromString(s)
	if err != nil {
		return err
	}
	*c = r
	return nil
}

var (
	ErrInvalidRangeCondition = errors.New("invalid range condition type")
)
