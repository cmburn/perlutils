package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const Undef = "undef"

func CloseBody(body io.Closer) {
	if body == nil {
		return
	}
	err := body.Close()
	if err != nil {
		// no point in dealing with this error, if we can't reach
		// stderr, something's too broken to fix anyway
		_, _ = fmt.Fprint(os.Stderr, err.Error())
	}
}

func UnwrapJSONString(b []byte) (string, error) {
	if len(b) < 2 {
		return "", ErrInvalidJSONString
	}
	if b[0] != '"' || b[len(b)-1] != '"' {
		return "", ErrInvalidJSONString
	}
	return string(b[1 : len(b)-1]), nil
}

func WrapEnumTypeJSON(s fmt.Stringer) ([]byte, error) {
	str := s.String()
	if str == "undef" {
		return nil, ErrInvalidEnumValue
	}
	sb := strings.Builder{}
	sb.WriteRune('"')
	sb.WriteString(str)
	sb.WriteRune('"')
	return []byte(sb.String()), nil
}

var (
	ErrInvalidEnumValue  = errors.New("invalid enum value")
	ErrInvalidJSONString = errors.New("invalid JSON string")
)

const (
	PackageVersion = "1.0.0"
)
