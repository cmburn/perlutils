//nolint

package version

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func getFractionValue(s string, t *testing.T) []int64 {
	if s == "" {
		// should only happen in lax decimal shenanigans
		return nil
	}
	s = strings.TrimPrefix(s, ".")
	expectedValues := (len(s) / 3) + 1
	if (len(s) % 3) == 0 {
		expectedValues--
	}
	stringValues := make([]string, expectedValues)
	currentString := ""
	for i := range s {
		if i%3 == 0 && i != 0 {
			stringValues[i/3-1] = currentString
			currentString = ""
		}
		currentString += string(s[i])
	}
	// have to pad the last value with zeros until three wide
	for i := len(currentString); i < 3; i++ {
		currentString += "0"
	}
	stringValues[len(stringValues)-1] = currentString
	values := make([]int64, len(stringValues))
	for i, s := range stringValues {
		val, err := strconv.Atoi(s)
		if err != nil {
			// if we make it this far and fail, something's
			// fundamentally wrong with the code.
			t.Fatalf("getFractionValue(%q) failed to convert "+
				"%q to int: %v", s, s, err)
		}
		values[i] = int64(val)
	}

	return values
}

func TestVersion_Equal(t *testing.T) {
	t.Parallel()
	undef, err := Parse("undef")
	if err != nil {
		t.Fatal(err)
	}
	standard, err := Parse("v1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	zero, err := Parse("0")
	if err != nil {
		t.Fatal(err)
	}
	alpha, err := Parse("v1.2.3_0")
	if err != nil {
		t.Fatal(err)
	}
	if !alpha.Equal(&alpha) {
		t.Errorf("alpha.Equal(&alpha) should be true")
	}
	if !standard.Equal(&standard) {
		t.Errorf("standard.Equal(&standard) should be true")
	}
	if !undef.Equal(&undef) {
		t.Errorf("undef.Equal(&undef) should be true")
	}
	if !undef.Equal(&zero) {
		t.Errorf("undef.Equal(&zero) should be true")
	}
	if !zero.Equal(&undef) {
		t.Errorf("zero.Equal(&undef) should be true")
	}
	if !zero.Equal(&zero) {
		t.Errorf("zero.Equal(&zero) should be true")
	}
	if alpha.Equal(&standard) {
		t.Errorf("alpha.Equal(&standard) should be false")
	}
	if alpha.Equal(&undef) {
		t.Errorf("alpha.Equal(&undef) should be false")
	}
	if alpha.Equal(&zero) {
		t.Errorf("alpha.Equal(&zero) should be false")
	}
	if standard.Equal(&alpha) {
		t.Errorf("standard.Equal(&alpha) should be false")
	}
	if standard.Equal(&undef) {
		t.Errorf("standard.Equal(&undef) should be false")
	}
	if standard.Equal(&zero) {
		t.Errorf("standard.Equal(&zero) should be false")
	}
	if undef.Equal(&alpha) {
		t.Errorf("undef.Equal(&alpha) should be false")
	}
	if undef.Equal(&standard) {
		t.Errorf("undef.Equal(&standard) should be false")
	}
	if zero.Equal(&alpha) {
		t.Errorf("zero.Equal(&alpha) should be false")
	}
	if zero.Equal(&standard) {
		t.Errorf("zero.Equal(&standard) should be false")
	}
}

func TestVersion_GreaterThan(t *testing.T) {
	t.Parallel()
	undef, err := Parse("undef")
	if err != nil {
		t.Fatal(err)
	}
	standard, err := Parse("v1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	zero, err := Parse("0")
	if err != nil {
		t.Fatal(err)
	}
	alpha, err := Parse("v1.2.3_0")
	if err != nil {
		t.Fatal(err)
	}
	if !alpha.GreaterThan(&standard) {
		t.Errorf("alpha.GreaterThan(&standard) should be true")
	}
	if !alpha.GreaterThan(&undef) {
		t.Errorf("alpha.GreaterThan(&undef) should be true")
	}
	if !alpha.GreaterThan(&zero) {
		t.Errorf("alpha.GreaterThan(&zero) should be true")
	}
	if !standard.GreaterThan(&undef) {
		t.Errorf("standard.GreaterThan(&undef) should be true")
	}
	if !standard.GreaterThan(&zero) {
		t.Errorf("standard.GreaterThan(&zero) should be true")
	}
	if alpha.GreaterThan(&alpha) {
		t.Errorf("alpha.GreaterThan(&alpha) should be false")
	}
	if standard.GreaterThan(&alpha) {
		t.Errorf("standard.GreaterThan(&alpha) should be false")
	}
	if standard.GreaterThan(&standard) {
		t.Errorf("standard.GreaterThan(&standard) should be false")
	}
	if undef.GreaterThan(&alpha) {
		t.Errorf("undef.GreaterThan(&alpha) should be false")
	}
	if undef.GreaterThan(&standard) {
		t.Errorf("undef.GreaterThan(&standard) should be false")
	}
	if undef.GreaterThan(&undef) {
		t.Errorf("undef.GreaterThan(&undef) should be false")
	}
	if undef.GreaterThan(&zero) {
		t.Errorf("undef.GreaterThan(&zero) should be false")
	}
	if zero.GreaterThan(&alpha) {
		t.Errorf("zero.GreaterThan(&alpha) should be false")
	}
	if zero.GreaterThan(&standard) {
		t.Errorf("zero.GreaterThan(&standard) should be false")
	}
	if zero.GreaterThan(&undef) {
		t.Errorf("zero.GreaterThan(&undef) should be false")
	}
	if zero.GreaterThan(&zero) {
		t.Errorf("zero.GreaterThan(&zero) should be false")
	}
}

func TestVersion_GreaterThanOrEqual(t *testing.T) {
	t.Parallel()
	undef, err := Parse("undef")
	if err != nil {
		t.Fatal(err)
	}
	standard, err := Parse("v1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	zero, err := Parse("0")
	if err != nil {
		t.Fatal(err)
	}
	alpha, err := Parse("v1.2.3_0")
	if err != nil {
		t.Fatal(err)
	}
	if !alpha.GreaterThanOrEqual(&alpha) {
		t.Errorf("alpha.GreaterThanOrEqual(&alpha) should be true")
	}
	if !alpha.GreaterThanOrEqual(&standard) {
		t.Errorf("alpha.GreaterThanOrEqual(&standard) should be true")
	}
	if !alpha.GreaterThanOrEqual(&undef) {
		t.Errorf("alpha.GreaterThanOrEqual(&undef) should be true")
	}
	if !alpha.GreaterThanOrEqual(&zero) {
		t.Errorf("alpha.GreaterThanOrEqual(&zero) should be true")
	}
	if !standard.GreaterThanOrEqual(&standard) {
		t.Errorf("standard.GreaterThanOrEqual(&standard) should be true")
	}
	if !standard.GreaterThanOrEqual(&undef) {
		t.Errorf("standard.GreaterThanOrEqual(&undef) should be true")
	}
	if !standard.GreaterThanOrEqual(&zero) {
		t.Errorf("standard.GreaterThanOrEqual(&zero) should be true")
	}
	if !undef.GreaterThanOrEqual(&undef) {
		t.Errorf("undef.GreaterThanOrEqual(&undef) should be true")
	}
	if !undef.GreaterThanOrEqual(&zero) {
		t.Errorf("undef.GreaterThanOrEqual(&zero) should be true")
	}
	if !zero.GreaterThanOrEqual(&undef) {
		t.Errorf("zero.GreaterThanOrEqual(&undef) should be true")
	}
	if !zero.GreaterThanOrEqual(&zero) {
		t.Errorf("zero.GreaterThanOrEqual(&zero) should be true")
	}
	if standard.GreaterThanOrEqual(&alpha) {
		t.Errorf("standard.GreaterThanOrEqual(&alpha) should be false")
	}
	if undef.GreaterThanOrEqual(&alpha) {
		t.Errorf("undef.GreaterThanOrEqual(&alpha) should be false")
	}
	if undef.GreaterThanOrEqual(&standard) {
		t.Errorf("undef.GreaterThanOrEqual(&standard) should be false")
	}
	if zero.GreaterThanOrEqual(&alpha) {
		t.Errorf("zero.GreaterThanOrEqual(&alpha) should be false")
	}
	if zero.GreaterThanOrEqual(&standard) {
		t.Errorf("zero.GreaterThanOrEqual(&standard) should be false")
	}
}

func TestVersion_LessThan(t *testing.T) {
	t.Parallel()
	undef, err := Parse("undef")
	if err != nil {
		t.Fatal(err)
	}
	standard, err := Parse("v1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	zero, err := Parse("0")
	if err != nil {
		t.Fatal(err)
	}
	alpha, err := Parse("v1.2.3_0")
	if err != nil {
		t.Fatal(err)
	}
	if !standard.LessThan(&alpha) {
		t.Errorf("standard.LessThan(&alpha) should be true")
	}
	if !undef.LessThan(&alpha) {
		t.Errorf("undef.LessThan(&alpha) should be true")
	}
	if !undef.LessThan(&standard) {
		t.Errorf("undef.LessThan(&standard) should be true")
	}
	if !zero.LessThan(&alpha) {
		t.Errorf("zero.LessThan(&alpha) should be true")
	}
	if !zero.LessThan(&standard) {
		t.Errorf("zero.LessThan(&standard) should be true")
	}
	if alpha.LessThan(&alpha) {
		t.Errorf("alpha.LessThan(&alpha) should be false")
	}
	if alpha.LessThan(&standard) {
		t.Errorf("alpha.LessThan(&standard) should be false")
	}
	if alpha.LessThan(&undef) {
		t.Errorf("alpha.LessThan(&undef) should be false")
	}
	if alpha.LessThan(&zero) {
		t.Errorf("alpha.LessThan(&zero) should be false")
	}
	if standard.LessThan(&standard) {
		t.Errorf("standard.LessThan(&standard) should be false")
	}
	if standard.LessThan(&undef) {
		t.Errorf("standard.LessThan(&undef) should be false")
	}
	if standard.LessThan(&zero) {
		t.Errorf("standard.LessThan(&zero) should be false")
	}
	if undef.LessThan(&undef) {
		t.Errorf("undef.LessThan(&undef) should be false")
	}
	if undef.LessThan(&zero) {
		t.Errorf("undef.LessThan(&zero) should be false")
	}
	if zero.LessThan(&undef) {
		t.Errorf("zero.LessThan(&undef) should be false")
	}
	if zero.LessThan(&zero) {
		t.Errorf("zero.LessThan(&zero) should be false")
	}
}

func TestVersion_LessThanOrEqual(t *testing.T) {
	t.Parallel()
	undef, err := Parse("undef")
	if err != nil {
		t.Fatal(err)
	}
	standard, err := Parse("v1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	zero, err := Parse("0")
	if err != nil {
		t.Fatal(err)
	}
	alpha, err := Parse("v1.2.3_0")
	if err != nil {
		t.Fatal(err)
	}
	if !alpha.LessThanOrEqual(&alpha) {
		t.Errorf("alpha.LessThanOrEqual(&alpha) should be true")
	}
	if !standard.LessThanOrEqual(&alpha) {
		t.Errorf("standard.LessThanOrEqual(&alpha) should be true")
	}
	if !standard.LessThanOrEqual(&standard) {
		t.Errorf("standard.LessThanOrEqual(&standard) should be true")
	}
	if !undef.LessThanOrEqual(&alpha) {
		t.Errorf("undef.LessThanOrEqual(&alpha) should be true")
	}
	if !undef.LessThanOrEqual(&standard) {
		t.Errorf("undef.LessThanOrEqual(&standard) should be true")
	}
	if !undef.LessThanOrEqual(&undef) {
		t.Errorf("undef.LessThanOrEqual(&undef) should be true")
	}
	if !undef.LessThanOrEqual(&zero) {
		t.Errorf("undef.LessThanOrEqual(&zero) should be true")
	}
	if !zero.LessThanOrEqual(&alpha) {
		t.Errorf("zero.LessThanOrEqual(&alpha) should be true")
	}
	if !zero.LessThanOrEqual(&standard) {
		t.Errorf("zero.LessThanOrEqual(&standard) should be true")
	}
	if !zero.LessThanOrEqual(&undef) {
		t.Errorf("zero.LessThanOrEqual(&undef) should be true")
	}
	if !zero.LessThanOrEqual(&zero) {
		t.Errorf("zero.LessThanOrEqual(&zero) should be true")
	}
	if alpha.LessThanOrEqual(&standard) {
		t.Errorf("alpha.LessThanOrEqual(&standard) should be false")
	}
	if alpha.LessThanOrEqual(&undef) {
		t.Errorf("alpha.LessThanOrEqual(&undef) should be false")
	}
	if alpha.LessThanOrEqual(&zero) {
		t.Errorf("alpha.LessThanOrEqual(&zero) should be false")
	}
	if standard.LessThanOrEqual(&undef) {
		t.Errorf("standard.LessThanOrEqual(&undef) should be false")
	}
	if standard.LessThanOrEqual(&zero) {
		t.Errorf("standard.LessThanOrEqual(&zero) should be false")
	}
}

func TestVersion_NotEqual(t *testing.T) {
	t.Parallel()
	undef, err := Parse("undef")
	if err != nil {
		t.Fatal(err)
	}
	standard, err := Parse("v1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	zero, err := Parse("0")
	if err != nil {
		t.Fatal(err)
	}
	alpha, err := Parse("v1.2.3_0")
	if err != nil {
		t.Fatal(err)
	}
	if !alpha.NotEqual(&standard) {
		t.Errorf("alpha.NotEqual(&standard) should be true")
	}
	if !alpha.NotEqual(&undef) {
		t.Errorf("alpha.NotEqual(&undef) should be true")
	}
	if !alpha.NotEqual(&zero) {
		t.Errorf("alpha.NotEqual(&zero) should be true")
	}
	if !standard.NotEqual(&alpha) {
		t.Errorf("standard.NotEqual(&alpha) should be true")
	}
	if !standard.NotEqual(&undef) {
		t.Errorf("standard.NotEqual(&undef) should be true")
	}
	if !standard.NotEqual(&zero) {
		t.Errorf("standard.NotEqual(&zero) should be true")
	}
	if !undef.NotEqual(&alpha) {
		t.Errorf("undef.NotEqual(&alpha) should be true")
	}
	if !undef.NotEqual(&standard) {
		t.Errorf("undef.NotEqual(&standard) should be true")
	}
	if !zero.NotEqual(&alpha) {
		t.Errorf("zero.NotEqual(&alpha) should be true")
	}
	if !zero.NotEqual(&standard) {
		t.Errorf("zero.NotEqual(&standard) should be true")
	}
	if alpha.NotEqual(&alpha) {
		t.Errorf("alpha.NotEqual(&alpha) should be false")
	}
	if standard.NotEqual(&standard) {
		t.Errorf("standard.NotEqual(&standard) should be false")
	}
	if undef.NotEqual(&undef) {
		t.Errorf("undef.NotEqual(&undef) should be false")
	}
	if undef.NotEqual(&zero) {
		t.Errorf("undef.NotEqual(&zero) should be false")
	}
	if zero.NotEqual(&undef) {
		t.Errorf("zero.NotEqual(&undef) should be false")
	}
	if zero.NotEqual(&zero) {
		t.Errorf("zero.NotEqual(&zero) should be false")
	}
}

func TestGetFractionValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input  string
		output []int64
	}{
		{"2", []int64{200}},
		{"02", []int64{20}},
		{"002", []int64{2}},
		{"0023", []int64{2, 300}},
		{"00203", []int64{2, 30}},
		{"002003", []int64{2, 3}},
		{"0020034", []int64{2, 3, 400}},
		{"00200304", []int64{2, 3, 40}},
		{"002003004", []int64{2, 3, 4}},
	}
	for _, test := range tests {
		values := getFractionValue(test.input, t)
		if len(values) != len(test.output) {
			t.Errorf("getFractionValue(%q) => %d, expected %d",
				test.input, len(values), len(test.output))
		}
		for i, v := range values {
			if v != test.output[i] {
				t.Errorf("getFractionValue(%q)[%d] => %d, expected %d",
					test.input, i, v, test.output[i])
			}
		}
	}
}

func TestNewPerlVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version  string
		expected Version
	}{
		{".1", Version{
			original: ".1",
			alpha:    false,
			version:  []int64{0, 100},
		}},
		{".1.2", Version{
			original: ".1.2",
			alpha:    false,
			version:  []int64{0, 1, 2},
		}},
		{"0", Version{
			original: "0",
			alpha:    false,
			version:  []int64{0},
		}},
		{"0.0", Version{
			original: "0.0",
			alpha:    false,
			version:  []int64{0, 0},
		}},
		{"0.123", Version{
			original: "0.123",
			alpha:    false,
			version:  []int64{0, 123},
		}},
		{"01", Version{
			original: "01",
			alpha:    false,
			version:  []int64{1},
		}},
		{"01.0203", Version{
			original: "01.0203",
			alpha:    false,
			version:  []int64{1, 20, 300},
		}},
		{"1.", Version{
			original: "1.",
			alpha:    false,
			version:  []int64{1, 0},
		}},
		{"1.00", Version{
			original: "1.00",
			alpha:    false,
			version:  []int64{1, 0},
		}},
		{"1.00001", Version{
			original: "1.00001",
			alpha:    false,
			version:  []int64{1, 0, 10},
		}},
		{"1.002", Version{
			original: "1.002",
			alpha:    false,
			version:  []int64{1, 2},
		}},
		{"1.002003", Version{
			original: "1.002003",
			alpha:    false,
			version:  []int64{1, 2, 3},
		}},
		{"1.00203", Version{
			original: "1.00203",
			alpha:    false,
			version:  []int64{1, 2, 30},
		}},
		{"1.0023", Version{
			original: "1.0023",
			alpha:    false,
			version:  []int64{1, 2, 300},
		}},
		{"1.02", Version{
			original: "1.02",
			alpha:    false,
			version:  []int64{1, 20},
		}},
		{"1.0203", Version{
			original: "1.0203",
			alpha:    false,
			version:  []int64{1, 20, 300},
		}},
		{"1.02_03", Version{
			original: "1.02_03",
			alpha:    true,
			version:  []int64{1, 20, 300},
		}},
		{"1.2", Version{
			original: "1.2",
			alpha:    false,
			version:  []int64{1, 200},
		}},
		{"1.2.3", Version{
			original: "1.2.3",
			alpha:    false,
			version:  []int64{1, 2, 3},
		}},
		{"1.2345_01", Version{
			original: "1.2345_01",
			alpha:    true,
			version:  []int64{1, 234, 501},
		}},
		{"12.345", Version{
			original: "12.345",
			alpha:    false,
			version:  []int64{12, 345},
		}},
		{"42", Version{
			original: "42",
			alpha:    false,
			version:  []int64{42},
		}},
		{"undef", Version{
			original: "undef",
			alpha:    false,
			version:  []int64{0},
		}},
		{"v0", Version{
			original: "v0",
			alpha:    false,
			version:  []int64{0, 0, 0},
		}},
		{"v0.0.0", Version{
			original: "v0.0.0",
			alpha:    false,
			version:  []int64{0, 0, 0},
		}},
		{"v0.1.2", Version{
			original: "v0.1.2",
			alpha:    false,
			version:  []int64{0, 1, 2},
		}},
		{"v01", Version{
			original: "v01",
			alpha:    false,
			version:  []int64{1, 0, 0},
		}},
		{"v01.02.03", Version{
			original: "v01.02.03",
			alpha:    false,
			version:  []int64{1, 2, 3},
		}},
		{"v1", Version{
			original: "v1",
			alpha:    false,
			version:  []int64{1, 0, 0},
		}},
		{"v1.02_03", Version{
			original: "v1.02_03",
			alpha:    true,
			version:  []int64{1, 203, 0},
		}},
		{"v1.2", Version{
			original: "v1.2",
			alpha:    false,
			version:  []int64{1, 2, 0},
		}},
		{"v1.2.3", Version{
			original: "v1.2.3",
			alpha:    false,
			version:  []int64{1, 2, 3},
		}},
		{"v1.2.3.4", Version{
			original: "v1.2.3.4",
			alpha:    false,
			version:  []int64{1, 2, 3, 4},
		}},
		{"v1.2.30", Version{
			original: "v1.2.30",
			alpha:    false,
			version:  []int64{1, 2, 30},
		}},
		{"v1.2.3_0", Version{
			original: "v1.2.3_0",
			alpha:    true,
			version:  []int64{1, 2, 30},
		}},
		{"v1.2345.6", Version{
			original: "v1.2345.6",
			alpha:    false,
			version:  []int64{1, 2345, 6},
		}},
		{"v1.2_3", Version{
			original: "v1.2_3",
			alpha:    true,
			version:  []int64{1, 23, 0},
		}},
	}
	for _, test := range tests {
		pv, err := Parse(test.version)
		if err != nil {
			t.Fatalf("NewPerlVersion(%q) returned error: %v", test.version, err)
		}
		if pv.original != test.expected.original {
			t.Errorf("NewPerlVersion(%q).original => %q, expected %q",
				test.version, pv.original,
				test.expected.original)
		}
		if pv.alpha != test.expected.alpha {
			t.Errorf("NewPerlVersion(%q).alpha => %t, expected %t",
				test.version, pv.alpha,
				test.expected.alpha)
		}
		if len(pv.version) != len(test.expected.version) {
			t.Errorf("len(NewPerlVersion(%q).Version) => %d, "+
				"expected %d", test.version, len(pv.version),
				len(test.expected.version))
			continue // prevent segfault
		}
		for i, v := range pv.version {
			if v != (test.expected.version)[i] {
				t.Errorf("NewPerlVersion(%q).Version[%d] "+
					"=> %d, expected %d", test.version, i,
					v, (test.expected.version)[i])
			}
		}
	}
}

func TestPerlVersion_IsAlpha(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version  string
		expected bool
	}{
		{".1", false},
		{".1.2", false},
		{"0", false},
		{"0.0", false},
		{"0.123", false},
		{"01", false},
		{"01.0203", false},
		{"1.", false},
		{"1.00", false},
		{"1.00001", false},
		{"1.002", false},
		{"1.002003", false},
		{"1.00203", false},
		{"1.0023", false},
		{"1.02", false},
		{"1.0203", false},
		{"1.02_03", true},
		{"1.2", false},
		{"1.2.3", false},
		{"1.2345_01", true},
		{"12.345", false},
		{"42", false},
		{"undef", false},
		{"v0", false},
		{"v0.0.0", false},
		{"v0.1.2", false},
		{"v01", false},
		{"v01.02.03", false},
		{"v1", false},
		{"v1.02_03", true},
		{"v1.2", false},
		{"v1.2.3", false},
		{"v1.2.3.4", false},
		{"v1.2.30", false},
		{"v1.2.3_0", true},
		{"v1.2345.6", false},
		{"v1.2_3", true},
		{"1.11111111111", false},
		{"2147483647.000", false},
	}
	for _, test := range tests {
		pv, err := Parse(test.version)
		if err != nil {
			t.Fatalf("NewPerlVersion(%q) returned error: %v",
				test.version, err)
		}
		if pv.IsAlpha() != test.expected {
			t.Errorf("NewPerlVersion(%q).IsAlpha() => %t,"+
				" expected %t", test.version, pv.IsAlpha(),
				test.expected)
		}
	}
}

func TestPerlVersion_IsQv(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version  string
		expected bool
	}{
		{".1", false},
		{".1.2", true},
		{"0", false},
		{"0.0", false},
		{"0.123", false},
		{"01", false},
		{"01.0203", false},
		{"1.", false},
		{"1.00", false},
		{"1.00001", false},
		{"1.002", false},
		{"1.002003", false},
		{"1.00203", false},
		{"1.0023", false},
		{"1.02", false},
		{"1.0203", false},
		{"1.02_03", false},
		{"1.2", false},
		{"1.2.3", true},
		{"1.2345_01", false},
		{"12.345", false},
		{"42", false},
		{"undef", false},
		{"v0", true},
		{"v0.0.0", true},
		{"v0.1.2", true},
		{"v01", true},
		{"v01.02.03", true},
		{"v1", true},
		{"v1.02_03", true},
		{"v1.2", true},
		{"v1.2.3", true},
		{"v1.2.3.4", true},
		{"v1.2.30", true},
		{"v1.2.3_0", true},
		{"v1.2345.6", true},
		{"v1.2_3", true},
		{"1.11111111111", false},
		{"2147483647.000", false},
	}
	for _, test := range tests {
		pv, err := Parse(test.version)
		if err != nil {
			t.Fatalf("NewPerlVersion(%q) returned error: %v",
				test.version, err)
		}
		if pv.IsQv() != test.expected {
			t.Errorf("NewPerlVersion(%q).IsQv() => %t,"+
				" expected %t", test.version, pv.IsQv(),
				test.expected)
		}
	}
}

func TestPerlVersion_Normal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version  string
		expected string
	}{
		{".1", "v0.100.0"},
		{".1.2", "v0.1.2"},
		{"0", "v0.0.0"},
		{"0.0", "v0.0.0"},
		{"0.123", "v0.123.0"},
		{"01", "v1.0.0"},
		{"01.0203", "v1.20.300"},
		{"1.", "v1.0.0"},
		{"1.00", "v1.0.0"},
		{"1.00001", "v1.0.10"},
		{"1.002", "v1.2.0"},
		{"1.002003", "v1.2.3"},
		{"1.00203", "v1.2.30"},
		{"1.0023", "v1.2.300"},
		{"1.02", "v1.20.0"},
		{"1.0203", "v1.20.300"},
		{"1.02_03", "v1.20.300"},
		{"1.2", "v1.200.0"},
		{"1.2.3", "v1.2.3"},
		{"1.2345_01", "v1.234.501"},
		{"12.345", "v12.345.0"},
		{"42", "v42.0.0"},
		{"undef", "v0.0.0"},
		{"v0", "v0.0.0"},
		{"v0.0.0", "v0.0.0"},
		{"v0.1.2", "v0.1.2"},
		{"v01", "v1.0.0"},
		{"v01.02.03", "v1.2.3"},
		{"v1", "v1.0.0"},
		{"v1.02_03", "v1.203.0"},
		{"v1.2", "v1.2.0"},
		{"v1.2.3", "v1.2.3"},
		{"v1.2.3.4", "v1.2.3.4"},
		{"v1.2.30", "v1.2.30"},
		{"v1.2.3_0", "v1.2.30"},
		{"v1.2345.6", "v1.2345.6"},
		{"v1.2_3", "v1.23.0"},
		{"1.11111111111", "v1.111.111.111.110"},
		{"2147483647.000", "v2147483647.0.0"},
	}
	for _, test := range tests {
		pv, err := Parse(test.version)
		if err != nil {
			t.Fatalf("NewPerlVersion(%q) returned error: %v",
				test.version, err)
		}
		if pv.Normal() != test.expected {
			t.Errorf("NewPerlVersion(%q).Normal() => %q,"+
				" expected %q", test.version, pv.Normal(),
				test.expected)
		}
	}
}

func TestPerlVersion_Numify(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version  string
		expected float64
	}{
		{".1", 0.100},
		{".1.2", 0.001002},
		{"0", 0.000},
		{"0.0", 0.000},
		{"0.123", 0.123},
		{"01", 1.000},
		{"01.0203", 1.020300},
		{"1.", 1.000},
		{"1.00", 1.000},
		{"1.00001", 1.000010},
		{"1.002", 1.002},
		{"1.002003", 1.002003},
		{"1.00203", 1.002030},
		{"1.0023", 1.002300},
		{"1.02", 1.020},
		{"1.0203", 1.020300},
		{"1.02_03", 1.020300},
		{"1.2", 1.200},
		{"1.2.3", 1.002003},
		{"1.2345_01", 1.234501},
		{"12.345", 12.345},
		{"42", 42.000},
		{"undef", 0.000},
		{"v0", 0.000000},
		{"v0.0.0", 0.000000},
		{"v0.1.2", 0.001002},
		{"v01", 1.000000},
		{"v01.02.03", 1.002003},
		{"v1", 1.000000},
		{"v1.02_03", 1.203000},
		{"v1.2", 1.002000},
		{"v1.2.3", 1.002003},
		{"v1.2.3.4", 1.002003004},
		{"v1.2.30", 1.002030},
		{"v1.2.3_0", 1.002030},
		{"v1.2345.6", 1.2345006},
		{"v1.2_3", 1.023000},
		{"1.11111111111", 1.111111111110},
		{"2147483647.000", 2147483647.000},
	}
	for _, test := range tests {
		pv, err := Parse(test.version)
		if err != nil {
			t.Fatalf("NewPerlVersion(%q) returned error: %v",
				test.version, err)
		}
		if pv.Numify() != test.expected {
			t.Errorf("NewPerlVersion(%q).Numify() => %f,"+
				" expected %f", test.version, pv.Numify(),
				test.expected)
		}
	}
}

func TestPerlVersion_Original(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version  string
		expected string
	}{
		{".1", ".1"},
		{".1.2", ".1.2"},
		{"0", "0"},
		{"0.0", "0.0"},
		{"0.123", "0.123"},
		{"01", "01"},
		{"01.0203", "01.0203"},
		{"1.", "1."},
		{"1.00", "1.00"},
		{"1.00001", "1.00001"},
		{"1.002", "1.002"},
		{"1.002003", "1.002003"},
		{"1.00203", "1.00203"},
		{"1.0023", "1.0023"},
		{"1.02", "1.02"},
		{"1.0203", "1.0203"},
		{"1.02_03", "1.02_03"},
		{"1.2", "1.2"},
		{"1.2.3", "1.2.3"},
		{"1.2345_01", "1.2345_01"},
		{"12.345", "12.345"},
		{"42", "42"},
		{"undef", "undef"},
		{"v0", "v0"},
		{"v0.0.0", "v0.0.0"},
		{"v0.1.2", "v0.1.2"},
		{"v01", "v01"},
		{"v01.02.03", "v01.02.03"},
		{"v1", "v1"},
		{"v1.02_03", "v1.02_03"},
		{"v1.2", "v1.2"},
		{"v1.2.3", "v1.2.3"},
		{"v1.2.3.4", "v1.2.3.4"},
		{"v1.2.30", "v1.2.30"},
		{"v1.2.3_0", "v1.2.3_0"},
		{"v1.2345.6", "v1.2345.6"},
		{"v1.2_3", "v1.2_3"},
		{"1.11111111111", "1.11111111111"},
		{"2147483647.000", "2147483647.000"},
	}
	for _, test := range tests {
		pv, err := Parse(test.version)
		if err != nil {
			t.Fatalf("NewPerlVersion(%q) returned error: %v",
				test.version, err)
		}
		if pv.Raw() != test.expected {
			t.Errorf("NewPerlVersion(%q).Stringify() => %q,"+
				" expected %q", test.version, pv.Raw(),
				test.expected)
		}
	}
}

func TestPerlVersion_Stringify(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version  string
		expected string
	}{
		{".1", ".1"},
		{".1.2", ".1.2"},
		{"0", "0"},
		{"0.0", "0.0"},
		{"0.123", "0.123"},
		{"01", "01"},
		{"01.0203", "01.0203"},
		{"1.", "1."},
		{"1.00", "1.00"},
		{"1.00001", "1.00001"},
		{"1.002", "1.002"},
		{"1.002003", "1.002003"},
		{"1.00203", "1.00203"},
		{"1.0023", "1.0023"},
		{"1.02", "1.02"},
		{"1.0203", "1.0203"},
		{"1.02_03", "1.02_03"},
		{"1.2", "1.2"},
		{"1.2.3", "1.2.3"},
		{"1.2345_01", "1.2345_01"},
		{"12.345", "12.345"},
		{"42", "42"},
		{"undef", "0"},
		{"v0", "v0"},
		{"v0.0.0", "v0.0.0"},
		{"v0.1.2", "v0.1.2"},
		{"v01", "v01"},
		{"v01.02.03", "v01.02.03"},
		{"v1", "v1"},
		{"v1.02_03", "v1.02_03"},
		{"v1.2", "v1.2"},
		{"v1.2.3", "v1.2.3"},
		{"v1.2.3.4", "v1.2.3.4"},
		{"v1.2.30", "v1.2.30"},
		{"v1.2.3_0", "v1.2.3_0"},
		{"v1.2345.6", "v1.2345.6"},
		{"v1.2_3", "v1.2_3"},
		{"1.11111111111", "1.11111111111"},
		{"2147483647.000", "2147483647.000"},
	}
	for _, test := range tests {
		pv, err := Parse(test.version)
		if err != nil {
			t.Fatalf("NewPerlVersion(%q) returned error: %v",
				test.version, err)
		}
		if pv.Stringify() != test.expected {
			t.Errorf("NewPerlVersion(%q).Stringify() => %q,"+
				" expected %q", test.version, pv.Stringify(),
				test.expected)
		}
	}
}

func TestVersion_Version(t *testing.T) {
	t.Parallel()
	tests := []struct {
		version  string
		expected []int64
	}{
		{".1", []int64{0, 100}},
		{".1.2", []int64{0, 1, 2}},
		{"0", []int64{0}},
		{"0.0", []int64{0, 0}},
		{"0.123", []int64{0, 123}},
		{"01", []int64{1}},
		{"01.0203", []int64{1, 20, 300}},
		{"1.", []int64{1, 0}},
		{"1.00", []int64{1, 0}},
		{"1.00001", []int64{1, 0, 10}},
		{"1.002", []int64{1, 2}},
		{"1.002003", []int64{1, 2, 3}},
		{"1.00203", []int64{1, 2, 30}},
		{"1.0023", []int64{1, 2, 300}},
		{"1.02", []int64{1, 20}},
		{"1.0203", []int64{1, 20, 300}},
		{"1.02_03", []int64{1, 20, 300}},
		{"1.2", []int64{1, 200}},
		{"1.2.3", []int64{1, 2, 3}},
		{"1.2345_01", []int64{1, 234, 501}},
		{"12.345", []int64{12, 345}},
		{"42", []int64{42}},
		{"undef", []int64{0}},
		{"v0", []int64{0, 0, 0}},
		{"v0.0.0", []int64{0, 0, 0}},
		{"v0.1.2", []int64{0, 1, 2}},
		{"v01", []int64{1, 0, 0}},
		{"v01.02.03", []int64{1, 2, 3}},
		{"v1", []int64{1, 0, 0}},
		{"v1.02_03", []int64{1, 203, 0}},
		{"v1.2", []int64{1, 2, 0}},
		{"v1.2.3", []int64{1, 2, 3}},
		{"v1.2.3.4", []int64{1, 2, 3, 4}},
		{"v1.2.30", []int64{1, 2, 30}},
		{"v1.2.3_0", []int64{1, 2, 30}},
		{"v1.2345.6", []int64{1, 2345, 6}},
		{"v1.2_3", []int64{1, 23, 0}},
	}
	for _, test := range tests {
		pv, err := Parse(test.version)
		if err != nil {
			t.Fatalf("NewPerlVersion(%q) returned error: %v",
				test.version, err)
		}
		v := pv.Version()
		if !reflect.DeepEqual(v, test.expected) {
			t.Errorf("NewPerlVersion(%q).Version() => %v,"+
				" expected %v", test.version, v, test.expected)
		}
		// make sure it's a copy, not a reference
		v[0] = -1
		if pv.Version()[0] == -1 {
			t.Errorf("NewPerlVersion(%q).Version() returned a "+
				"reference to the internal Version slice",
				test.version)
		}
	}
}

func TestVersion_MarshalJSON(t *testing.T) {
	t.Parallel()
	input := Version{
		original: "v1.2.3",
		alpha:    false,
		qv:       true,
		version:  []int64{1, 2, 3},
	}
	data, err := json.Marshal(&input)
	if err != nil {
		t.Errorf("Version.MarshalJSON() returned error: %v", err)
	}
	var actual Version
	if err := json.Unmarshal(data, &actual); err != nil {
		t.Errorf("Version.UnmarshalJSON() returned error: %v", err)
	}
	if !reflect.DeepEqual(actual, input) {
		t.Errorf("Version.UnmarshalJSON() => %+v, expected %+v",
			actual, input)
	}
}

type rangeTestStruct struct {
	input    string
	expected Range
}

func TestParseRange(t *testing.T) {
	t.Parallel()
	// the version parsing itself is already tested in version_test.go,
	// so we're not throwing any edge cases at it here
	tests := []rangeTestStruct{
		{
			input: "v1",
			expected: Range{
				conditions: []RangeSpecifier{
					{
						Condition: rangeConditionNone,
						Version:   MustParse("v1"),
					},
				},
			},
		},
		{
			input: "v1, v2",
			expected: Range{
				conditions: []RangeSpecifier{
					{
						Condition: rangeConditionNone,
						Version:   MustParse("v1"),
					},
					{
						Condition: rangeConditionNone,
						Version:   MustParse("v2"),
					},
				},
			},
		},
		{
			input: "!=v1, > v2, <=v3.1.5",
			expected: Range{
				conditions: []RangeSpecifier{
					{
						Condition: RangeConditionNotEqual,
						Version:   MustParse("v1"),
					},
					{
						Condition: RangeConditionGreaterThan,
						Version:   MustParse("v2"),
					},
					{
						Condition: RangeConditionLessThanOrEqual,
						Version:   MustParse("v3.1.5"),
					},
				},
			},
		},
	}
	for _, test := range tests {
		r, err := ParseRange(test.input)
		if err != nil {
			t.Errorf("ParseRange(%q) returned error: %v",
				test.input, err)
		}
		if !reflect.DeepEqual(*r, test.expected) {
			t.Errorf("ParseRange(%q) returned %v, expected %v",
				test.input, *r, test.expected)
		}
	}

}
