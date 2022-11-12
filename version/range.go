package version

import (
	"encoding/json"
	"strings"
)

// Range is a set of conditions that must all be true for a version to be
// considered to be in the range. The conditions are ANDed together.
type Range struct {
	conditions []RangeSpecifier
}

func (r *Range) Contains(v *Version) bool {
	for _, rc := range r.conditions {
		if !rc.contains(v) {
			return false
		}
	}
	return true
}

func (r *Range) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.conditions)
}

func (r *Range) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.conditions)
}

func (r *Range) String() string {
	if len(r.conditions) == 1 &&
		r.conditions[0].Condition == RangeConditionEqual {
		return r.conditions[0].Version.Raw()
	}
	sb := strings.Builder{}
	for i, rc := range r.conditions {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(rc.String())
	}

	return sb.String()
}

func (r *Range) URLString() string {
	sb := strings.Builder{}
	for i, rc := range r.conditions {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(rc.URLString())
	}
	return sb.String()
}

func NewRange(conditions []RangeSpecifier) (*Range, error) {
	// make sure there's no undefined conditions
	for _, rc := range conditions {
		switch rc.Condition {
		case RangeConditionEqual, RangeConditionGreaterThan,
			RangeConditionGreaterThanOrEqual,
			RangeConditionLessThan, RangeConditionLessThanOrEqual,
			RangeConditionNotEqual, rangeConditionNone:
			break
		case rangeConditionUndef:
			fallthrough
		default:
			return nil, ErrInvalidRangeCondition
		}
	}
	r := &Range{}
	r.conditions = make([]RangeSpecifier, 0, len(conditions))
	copy(r.conditions, conditions)
	return r, nil
}

// NewRangeNoFail is like NewRange, but panics if there's an error
func NewRangeNoFail(conditions []RangeSpecifier) *Range {
	r, err := NewRange(conditions)
	if err != nil {
		panic(err)
	}
	return r
}

// ParseRange parses a string into a Range, in accordance with
// CPAN::Meta::Spec.
func ParseRange(s string) (*Range, error) {
	// leading spaces, trailing spaces, and a comma at the end are all
	// allowed
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, ",")
	s = strings.TrimSpace(s) // catch any space before the comma
	condStrings := strings.Split(s, ",")
	r := &Range{}
	for _, condString := range condStrings {
		rc := RangeSpecifier{}
		if err := parseCondition(&rc, condString); err != nil {
			return nil, err
		}
		r.conditions = append(r.conditions, rc)
	}
	return r, nil
}

// MustParseRange is like ParseRange, but panics if there's an error
func MustParseRange(s string) *Range {
	r, err := ParseRange(s)
	if err != nil {
		panic(err)
	}
	return r
}
