package internal

// This file holds the implementation of the parser for perl's lax versioning
// spec. It's largely deprecated by the Perl project, but there's still a good
// deal of modules on CPAN that use it.

import (
	"regexp"
	"strings"
)

var (
	laxRegexp = regexp.MustCompile(LaxVersionRegex)
)

type laxDotted struct {
	integer           string // Version A
	dottedGroup       string
	alpha             string
	secondInteger     string // Version B
	secondDottedGroup string
	secondAlpha       string
}

type laxDecimal struct {
	integer        string // Version A
	fraction       string
	alpha          string
	secondFraction string // Version B
	secondAlpha    string
}

type lax struct {
	original       string
	undef          string
	dotted         string
	dottedMatches  laxDotted
	decimal        string
	decimalMatches laxDecimal
}

func (d laxDotted) toPerlVersionA(original string) VersionData {

	isAlpha := d.alpha != ""
	sb := strings.Builder{}
	sb.WriteString(d.dottedGroup)
	if isAlpha {
		sb.WriteString(strings.TrimPrefix(d.alpha, "_"))
	}
	var minors []int64
	if sb.Len() > 0 {
		minors = dottedToMinors(sb.String())
	}
	numValues := len(minors)
	if numValues < minValues {
		// implied zeroes in v-qualified lax Version
		numValues = minValues
	}
	values := make([]int64, numValues)
	values[0] = mustParseInt64(d.integer)
	if minors != nil {
		copy(values[1:], minors)
	}
	return VersionData{
		Original: original,
		Alpha:    isAlpha,
		QV:       true,
		Version:  values,
	}
}

func (d laxDotted) toPerlVersionB(original string) VersionData {
	// This particular case is a bit tricky. If there's three values,
	// *implied* zeroes included, it counts as a quoted lax Version.
	sb := strings.Builder{}
	sb.WriteString(d.secondDottedGroup)
	isAlpha := d.secondAlpha != ""
	if isAlpha {
		sb.WriteString(strings.TrimPrefix(d.secondAlpha, "_"))
	}
	minors := dottedToMinors(sb.String())
	numValues := len(minors)
	impliedZero := d.secondDottedGroup[0] == '.' && d.secondInteger == ""
	if d.secondInteger != "" || impliedZero {
		numValues++
	}
	values := make([]int64, numValues)
	if d.secondInteger != "" {
		values[0] = mustParseInt64(d.secondInteger)
	} else if impliedZero {
		values[0] = 0
	}
	if minors != nil {
		copy(values[1:], minors)
	}
	return VersionData{
		Original: original,
		Alpha:    d.secondAlpha != "",
		QV:       numValues == minValues,
		Version:  values,
	}
}

func (d laxDotted) toPerlVersion(original string) VersionData {
	switch {
	case d.integer != "":
		return d.toPerlVersionA(original)
	case d.secondDottedGroup != "":
		return d.toPerlVersionB(original)
	default:
		panic("unreachable")
	}
}

func (d laxDecimal) toPerlVersionA(original string) (VersionData, error) {
	// Due to, what I can tell, is a runtime check, this is the only
	// subset of versioning that can error out. It happens when there's an
	// alpha part but no fractional part. It would make sense to change the
	// regex, but I'm hesitant to deviate from the Perl versioning spec.
	// Example: "1_0"
	sb := strings.Builder{}
	sb.WriteString(d.fraction)
	isAlpha := d.alpha != ""
	if isAlpha {
		if d.fraction == "" {
			return VersionData{}, errAlphaWithoutDecimal
		}
		sb.WriteString(strings.TrimPrefix(d.alpha, "_"))
	}
	fractions := getFractionValue(sb.String())
	numValues := len(fractions) + 1
	impliedZeroEnd := original[len(original)-1] == '.' && d.fraction == ""
	if impliedZeroEnd {
		numValues++
	}
	values := make([]int64, numValues)
	values[0] = mustParseInt64(d.integer)
	if fractions != nil {
		copy(values[1:], fractions)
	}
	if impliedZeroEnd {
		values[len(values)-1] = 0
	}
	return VersionData{
		Original: original,
		Alpha:    d.alpha != "",
		QV:       false,
		Version:  values,
	}, nil
}

func (d laxDecimal) toPerlVersionB(original string) VersionData {
	sb := strings.Builder{}
	sb.WriteString(d.secondFraction)
	isAlpha := d.secondAlpha != ""
	if isAlpha {
		sb.WriteString(strings.TrimPrefix(d.secondAlpha, "_"))
	}
	fractions := getFractionValue(sb.String())
	if fractions == nil {
		panic("unreachable")
	}
	values := make([]int64, len(fractions)+1) // implied zero
	values[0] = 0
	copy(values[1:], fractions)
	return VersionData{
		Original: original,
		Alpha:    d.secondAlpha != "",
		QV:       false,
		Version:  values,
	}
}

func (d laxDecimal) toPerlVersion(original string) (VersionData, error) {
	switch {
	case d.integer != "":
		return d.toPerlVersionA(original)
	case d.secondFraction != "":
		return d.toPerlVersionB(original), nil
	default:
		panic("unreachable")
	}
}

func (d lax) toPerlVersion() (VersionData, error) {
	switch {
	case d.undef != "":
		return VersionData{
			Original: d.original,
			Alpha:    false,
			QV:       false,
			Version:  []int64{0},
		}, nil
	case d.dotted != "":
		return d.dottedMatches.toPerlVersion(d.original), nil
	case d.decimal != "":
		return d.decimalMatches.toPerlVersion(d.original)
	default:
		panic("unreachable")
	}
}

func LaxVersion(matches []string) (VersionData, error) {
	return lax{
		original: matches[0],
		undef:    matches[1],
		dotted:   matches[2],
		dottedMatches: laxDotted{
			integer:           matches[3],
			dottedGroup:       matches[4],
			alpha:             matches[5],
			secondInteger:     matches[6],
			secondDottedGroup: matches[7],
			secondAlpha:       matches[8],
		},
		decimal: matches[9],
		decimalMatches: laxDecimal{
			integer:        matches[10],
			fraction:       matches[11],
			alpha:          matches[12],
			secondFraction: matches[13],
			secondAlpha:    matches[14],
		},
	}.toPerlVersion()
}

const (
	minValues = 3
)
