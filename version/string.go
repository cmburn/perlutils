package version

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	// local
	"github.com/cmburn/perlutils/internal"
)

// JSON is a wrapper around Version that implements json.Marshaler and
// json.Unmarshaler as a string.
type JSON struct {
	Version
}

func (j *JSON) MarshalJSON() ([]byte, error) {
	str := j.Raw()
	return json.Marshal(str)
}

func (j *JSON) UnmarshalJSON(b []byte) error {
	// unmarshal as a string
	var str string
	if b[0] != '"' {
		str = string(b)
	} else if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	str = strings.Trim(str, `"`)
	version, err := Parse(str)
	if err != nil {
		return err
	}
	*j = JSON{version}
	return nil
}

func (j *JSON) String() string {
	return j.Raw()
}

type JSONNoFail struct{ Version }

func (j *JSONNoFail) MarshalJSON() ([]byte, error) {
	str := j.Raw()
	return json.Marshal(str)
}

func (j *JSONNoFail) UnmarshalJSON(b []byte) error {
	// unmarshal as a string
	var str string
	if b[0] != '"' {
		str = string(b)
	} else if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	str = strings.Trim(str, `"`)
	version, err := Parse(str)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error parsing version %q: %s",
			str, err)
		*j = JSONNoFail{
			Version: Version{
				original: str,
				version:  []int64{0},
			},
		}
		return nil
	}
	*j = JSONNoFail{version}
	return nil
}

func (j *JSONNoFail) String() string {
	return j.Raw()
}

// RangeJSON is a wrapper around Range that implements json.Marshaler and
// json.Unmarshaler as a string.
type RangeJSON struct {
	Range
}

func (rs *RangeJSON) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	r, err := ParseRange(s)
	if err != nil {
		return err
	}
	rs.Range = *r
	return nil
}

func (rs *RangeJSON) MarshalJSON() ([]byte, error) {
	return internal.WrapEnumTypeJSON(&rs.Range)
}

var _ json.Marshaler = (*JSON)(nil)
var _ json.Unmarshaler = (*JSON)(nil)
var _ json.Marshaler = (*RangeJSON)(nil)
var _ json.Unmarshaler = (*RangeJSON)(nil)
