package internal

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

func coerce[T any](into *T, v *interface{}) bool {
	if *v == nil {
		var null T
		*into = null
		return true
	}
	switch v2 := (*v).(type) {
	case json.RawMessage:
		*v = string(v2)
	case []byte:
		*v = string(v2)
	}
	return false
}

func CoerceBool(into *bool, v interface{}) error {
	if coerce[bool](into, &v) {
		return nil
	}
	switch v := v.(type) {
	case bool:
		*into = v
	case int:
		*into = v != 0
	case float64:
		*into = v != 0
	case string:
		v = strings.Trim(v, "\"")
		switch v {
		case "true":
			*into = true
		case "false":
			*into = false
		default:
			i := 0
			err := CoerceInt(&i, v)
			if err == nil {
				*into = i != 0
				return nil
			}
			return ErrInvalidType
		}
	default:
		return ErrInvalidType
	}

	return nil
}

func CoerceFloat(into *float64, v interface{}) error {
	if coerce[float64](into, &v) {
		return nil
	}

	switch v := v.(type) {
	case float64:
		*into = v
	case int:
		*into = float64(v)
	case string:
		v = strings.Trim(v, "\"")
		if len(v) == 0 {
			*into = 0
			break
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		*into = f
	case bool:
		if v {
			*into = 1
		} else {
			*into = 0
		}
	default:
		return ErrInvalidType
	}
	return nil
}

func CoerceInt(into *int, v interface{}) error {
	if coerce[int](into, &v) {
		return nil
	}

	switch v := v.(type) {
	case int:
		*into = v
	case int64:
		*into = int(v)
	case float64:
		*into = int(v)
	case string:
		v = strings.Trim(v, "\"")
		if len(v) == 0 {
			*into = 0
			break
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		*into = i
	case bool:
		if v {
			*into = 1
		} else {
			*into = 0
		}
	default:
		return ErrInvalidType
	}

	return nil
}

func CoerceInt64(into *int64, v interface{}) error {
	if coerce[int64](into, &v) {
		return nil
	}

	switch v := v.(type) {
	case int64:
		*into = v
	case int:
		*into = int64(v)
	case float64:
		*into = int64(v)
	case string:
		v = strings.Trim(v, "\"")
		if len(v) == 0 {
			*into = 0
			break
		}
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		*into = i
	case bool:
		if v {
			*into = 1
		} else {
			*into = 0
		}
	default:
		return ErrInvalidType
	}
	return nil
}

func CoerceStringArray(into *[]string, v interface{}) error {
	if v == nil {
		*into = nil
		return nil
	}
	switch v := v.(type) {
	case []interface{}:
		for _, i := range v {
			switch i := i.(type) {
			case string:
				*into = append(*into, i)
			default:
				return ErrInvalidType
			}
		}
	case []byte:
		var vi interface{}
		if err := json.Unmarshal(v, &vi); err != nil {
			return err
		}
		return CoerceStringArray(into, vi)
	case string:
		str := strings.Trim(v, "\"")
		*into = []string{str}
	default:
		return ErrInvalidType
	}
	return nil
}

func CoerceUInt32(into *uint32, v interface{}) error {
	if coerce[uint32](into, &v) {
		return nil
	}

	switch v := v.(type) {
	case uint32:
		*into = v
	case int:
		*into = uint32(v)
	case float64:
		*into = uint32(v)
	case string:
		v = strings.Trim(v, "\"")
		if len(v) == 0 {
			*into = 0
			break
		}
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return err
		}
		*into = uint32(i)
	case bool:
		if v {
			*into = 1
		} else {
			*into = 0
		}
	default:
		return ErrInvalidType
	}
	return nil
}

type CoercedBool struct{ Value bool }

func (c *CoercedBool) UnmarshalJSON(data []byte) error {
	return CoerceBool(&c.Value, data)
}

type CoercedInt struct{ Value int }

func (c *CoercedInt) UnmarshalJSON(data []byte) error {
	return CoerceInt(&c.Value, data)
}

type CoercedInt64 struct{ Value int64 }

func (c *CoercedInt64) UnmarshalJSON(data []byte) error {
	return CoerceInt64(&c.Value, data)
}

type CoercedUInt32 struct{ Value uint32 }

func (c *CoercedUInt32) UnmarshalJSON(data []byte) error {
	return CoerceUInt32(&c.Value, data)
}

type CoercedFloat struct{ Value float64 }

func (c *CoercedFloat) UnmarshalJSON(data []byte) error {
	return CoerceFloat(&c.Value, data)
}

type CoercedStringArray struct{ Value []string }

func (c *CoercedStringArray) UnmarshalJSON(data []byte) error {
	return CoerceStringArray(&c.Value, data)
}

var (
	ErrInvalidType = errors.New("invalid type")
)
