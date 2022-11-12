// Package version is a Go implementation of Perl's version.pm. It's
// written for the purpose of working with Perl packages in the context of a
// larger, multi-language monorepo. It's written to be a bug-for-bug compatible
// implementation of Perl's Version.pm.
package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	// local
	"github.com/cmburn/perlutils/version/internal"
)

const (
	LaxVersionRegex    = internal.LaxVersionRegex
	StrictVersionRegex = internal.StrictVersionRegex
)

var (
	laxRegexp    = regexp.MustCompile(LaxVersionRegex)
	strictRegexp = regexp.MustCompile(StrictVersionRegex)
)

// Version is a direct, mapping of Perl's Version::Internal, methods and
// all. It's meant to be opaque, as the internal representation might change
// if the need arises.
type Version struct {
	original string
	version  []int64
	alpha    bool
	qv       bool
}

///////////////////////////////////////////////////////////////////////////////
// Direct methods from Perl                                                  //
///////////////////////////////////////////////////////////////////////////////

// These methods are directly copied from Perl's Version.pm.

// IsAlpha checks whether a Version is an alpha Version. This is implied by the
// presence of an underscore in the Version. For example, "1.2_3" is an alpha
// Version (equating to "v1.203.0", but that's due to Perl's bizarre Version
// semantics).
func (v *Version) IsAlpha() bool {
	return v.alpha
}

// IsQv checks whether a Version is a qv Version. This is indicated by a 'v' at
// the beginning of the Version. For example, "v1.2.3" is a qv Version, while
// "1.2.3" is not. The versions are not equal either- "1.2.3" is represented
// as "v1.200.300", and "v1.2.3" is still just "v1.2.3".
func (v *Version) IsQv() bool {
	return v.qv
}

// Normal is a convenience function for normalizing a Version string. It
// returns it in standardized qv form, with at least three subversions.
func (v *Version) Normal() string {
	num := len(v.version)
	if num < 3 {
		num = 3
	}
	fixed := make([]int64, num)
	copy(fixed, v.version)
	asStrings := make([]string, num)
	for i, v := range fixed {
		asStrings[i] = fmt.Sprintf("%d", v)
	}
	return "v" + strings.Join(asStrings, ".")
}

// Numify returns the numeric Version of a Version string. For example,
// "v1.2.3" would return 1.002003. This is useful for quick comparisons, and
// embedding in maps, though if you have a Version with many subversions, it's
// probably better to use the relevant comparison methods (which are probably
// faster regardless).
func (v *Version) Numify() float64 {
	if len(v.version) == 1 {
		return float64(v.version[0])
	}
	asStrings := make([]string, len(v.version)-1)
	for i, v := range v.version[1:] {
		asStrings[i] = strconv.FormatInt(v, 10)
		// pad with zeros
		for len(asStrings[i]) < 3 {
			asStrings[i] = "0" + asStrings[i]
		}
	}
	tail := strings.Join(asStrings, "")
	str := fmt.Sprintf("%d.%s", v.version[0], tail)
	out, _ := strconv.ParseFloat(str, 64)
	return out
}

// Stringify matches its Perl equivalent- functionally it acts the same as Raw,
// however if the Version is undefined, it returns "0".
func (v *Version) Stringify() string {
	if v.original == "undef" {
		return "0"
	}
	return v.original
}

// Raw returns the original representation of the Version.
func (v *Version) Raw() string {
	return v.original
}

func (v *Version) InRange(r *Range) bool {
	return r.Contains(v)
}

// MarshalJSON implements the json.Marshaler interface. This allows for caching
// of the Version.
func (v *Version) MarshalJSON() ([]byte, error) {
	data := struct {
		Original string  `json:"original"`
		Alpha    bool    `json:"alpha"`
		Qv       bool    `json:"qv"`
		Version  []int64 `json:"Version"`
	}{
		Original: v.original,
		Alpha:    v.alpha,
		Qv:       v.qv,
		Version:  v.version,
	}
	return json.Marshal(&data)
}

// Version returns the Version as a slice of integers.
func (v *Version) Version() []int64 {
	// return duplicate
	return append([]int64{}, v.version...)
}

// UnmarshalJSON implements the json.Unmarshaler interface. This allows for
// extracting the Version from a cached Version.
func (v *Version) UnmarshalJSON(data []byte) error {
	var obj struct {
		Original string  `json:"original"`
		Alpha    bool    `json:"alpha"`
		Qv       bool    `json:"qv"`
		Version  []int64 `json:"Version"`
	}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	v.original = obj.Original
	v.alpha = obj.Alpha
	v.qv = obj.Qv
	v.version = obj.Version
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Comparisons                                                               //
///////////////////////////////////////////////////////////////////////////////

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// LessThan checks whether a Version is older than another.
func (v *Version) LessThan(other *Version) bool {
	length := min(len(v.version), len(other.version))
	for i := 0; i < length; i++ {
		if v.version[i] < other.version[i] {
			return true
		}
		if v.version[i] > other.version[i] {
			return false
		}
	}
	return false
}

// GreaterThan checks whether a Version is newer than another.
func (v *Version) GreaterThan(other *Version) bool {
	length := min(len(v.version), len(other.version))
	for i := 0; i < length; i++ {
		if v.version[i] > other.version[i] {
			return true
		}
		if v.version[i] < other.version[i] {
			return false
		}
	}
	return false
}

// Equal checks whether two versions are the same. This doesn't strictly
// mean they're identical, it means, for example, "v5.34" counts as the same as
// "v5.34.0" *or* "v5.34.1".
func (v *Version) Equal(other *Version) bool {
	return !(v.LessThan(other) || v.GreaterThan(other))
}

// LessThanOrEqual checks whether a Version is older or equivalent to
// another. Same as (LessThan || Equal).
func (v *Version) LessThanOrEqual(other *Version) bool {
	return v.Equal(other) || v.LessThan(other)
}

// GreaterThanOrEqual checks whether a Version is newer or equivalent to
// another. Same as (GreaterThan || Equal).
func (v *Version) GreaterThanOrEqual(other *Version) bool {
	return v.Equal(other) || v.GreaterThan(other)
}

// NotEqual checks whether two versions are not the same. Same as
// !(Equal).
func (v *Version) NotEqual(other *Version) bool {
	return !v.Equal(other)
}

// Compare compares two versions. It returns -1 if the receiver is older,
// 0 if they're equivalent, and 1 if the receiver is newer.
func (v *Version) Compare(other *Version) int {
	if v.LessThan(other) {
		return -1
	}
	if v.GreaterThan(other) {
		return 1
	}
	return 0
}

func internalToVersion(dst *Version, src *internal.VersionData) {
	dst.original = src.Original
	dst.alpha = src.Alpha
	dst.qv = src.QV
	dst.version = src.Version
}

// Parse parses a string into a Version. The string can be either a lax or
// strict versioning scheme, as defined in Version::Internals.
func Parse(version string) (Version, error) {
	version = strings.TrimSpace(version)
	endsWithAlpha := strings.HasSuffix(version, "_")
	laxMatch := laxRegexp.FindStringSubmatch(version)
	strictMatch := strictRegexp.FindStringSubmatch(version)
	tryLax := laxMatch != nil
	if !tryLax && endsWithAlpha {
		version = strings.TrimSuffix(version, "_")
		laxMatch = laxRegexp.FindStringSubmatch(version)
		tryLax = laxMatch != nil
	}
	tryStrict := strictMatch != nil
	preferLax := (tryLax && tryStrict) &&
		(len(laxMatch[0]) > len(strictMatch[0]))
	parsed := false
	v := Version{}
	if tryLax {
		data, err := internal.LaxVersion(laxMatch)
		if err == nil {
			if endsWithAlpha {
				data.Alpha = true
				data.Original = fmt.Sprintf("%s_",
					data.Original)
			}
			internalToVersion(&v, &data)
			parsed = true
		}
	}
	if tryStrict && (!preferLax || !parsed) {
		data := internal.StrictVersion(strictMatch)
		internalToVersion(&v, &data)
		parsed = true
	}

	if parsed {
		return v, nil
	}

	return Version{}, errors.New("invalid Version string: " + version)
}

func parseMulti(a, b string) (Version, Version, error) {
	aPv, err := Parse(a)
	if err != nil {
		return Version{}, Version{}, err
	}
	bPv, err := Parse(b)
	if err != nil {
		return Version{}, Version{}, err
	}
	return aPv, bPv, nil
}

// Undef returns a new, undefined Version.
func Undef() Version {
	return Version{
		original: "undef",
		alpha:    false,
		qv:       false,
		version:  []int64{0},
	}
}

// CompatibleWith checks if candidate is compatible with target. Panics if
// it can't parse either of the two; this is meant for already-validated
// versions. If you need to check, use the respective Version method
// Version.GreaterThanOrEqual.
func CompatibleWith(candidate, target string) bool {
	a, b, err := parseMulti(candidate, target)
	if err != nil {
		panic(err)
	}
	return a.GreaterThanOrEqual(&b)
}

// MustParse is for parsing a Version string that must be valid. It panics
// if it can't parse the string. You probably want Parse(), unless you're
// dealing with an internal cache.
func MustParse(version string) Version {
	v, err := Parse(version)
	if err != nil {
		panic(err)
	}
	return v
}

// LooksValid returns true if the Version looks valid. It can sometimes return
// false positives in edge cases that can't be easily checked, but if it says
// it's not valid, it's not valid.
func LooksValid(version string) bool {
	return laxRegexp.MatchString(version) ||
		strictRegexp.MatchString(version)
}
