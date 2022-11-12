package version

import "strings"

type RangeSpecifier struct {
	Condition RangeCondition `json:"range_type"`
	Version   Version        `json:"version"`
}

func (rc *RangeSpecifier) contains(v *Version) bool {
	switch rc.Condition {
	case RangeConditionGreaterThan:
		return v.GreaterThan(&rc.Version)
	case RangeConditionGreaterThanOrEqual:
		return v.GreaterThanOrEqual(&rc.Version)
	case RangeConditionLessThan:
		return v.LessThan(&rc.Version)
	case RangeConditionLessThanOrEqual:
		return v.LessThanOrEqual(&rc.Version)
	case RangeConditionNotEqual:
		return v.NotEqual(&rc.Version)
	case RangeConditionEqual:
		fallthrough
	case rangeConditionNone, rangeConditionUndef:
		return v.Equal(&rc.Version)
	}
	return false
}

func (rc *RangeSpecifier) String() string {
	if rc.Condition == rangeConditionUndef {
		// not possible
		panic("invalid RangeSpecifier")
	}
	if rc.Condition == rangeConditionNone {
		return rc.Version.Raw()
	}
	return rc.Condition.String() + rc.Version.Raw()
}

func (rc *RangeSpecifier) URLString() string {
	return rc.Condition.URLString() + rc.Version.Raw()
}

func parseCondition(rc *RangeSpecifier, s string) error {
	s = strings.TrimSpace(s)
	symbol := string(s[0])
	foundSymbol := true
	switch s[0] {
	case '>':
		fallthrough
	case '<':
		fallthrough
	case '=':
		fallthrough
	case '!':
		if s[1] == '=' {
			symbol += "="
		}
	default:
		foundSymbol = false
	}
	if foundSymbol {
		rt, err := rangeTypeFromString(symbol)
		if err != nil {
			rc.Condition = rangeConditionUndef
			return err
		}
		rc.Condition = rt
		s = strings.TrimSpace(s[len(symbol):])
	}
	v, err := Parse(s)
	if err != nil {
		rc.Condition = rangeConditionUndef
		return err
	}
	if rc.Condition == rangeConditionUndef {
		rc.Condition = rangeConditionNone
	}
	rc.Version = v
	return nil
}
