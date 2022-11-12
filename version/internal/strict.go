package internal

import (
	"regexp"
	"strings"
)

var (
	strictRegexp = regexp.MustCompile(StrictVersionRegex)
)

type strictDecimalForm struct {
	integerPart  string
	fractionPart string
}

type strictDottedForm struct {
	integerPart string
	dottedGroup string
}

type strict struct {
	original       string
	decimal        string
	decimalMatches strictDecimalForm
	dotted         string
	dottedMatches  strictDottedForm
}

func (d strictDecimalForm) toPerlVersion(original string) VersionData {
	pv := VersionData{
		Original: original,
		Alpha:    false,
		QV:       false,
	}
	trimmed := strings.TrimPrefix(d.fractionPart, ".")
	fracValues := getFractionValue(trimmed)
	pv.Version = make([]int64, len(fracValues)+1)
	pv.Version[0] = mustParseInt64(d.integerPart)
	copy(pv.Version[1:], fracValues)
	return pv
}

func (d strictDottedForm) toPerlVersion(original string) VersionData {
	pv := VersionData{
		Original: original,
		Alpha:    false,
		QV:       true,
	}
	trimmed := strings.TrimPrefix(d.dottedGroup, ".")
	minors := strings.Split(trimmed, ".")
	pv.Version = make([]int64, len(minors)+1)
	pv.Version[0] = mustParseInt64(d.integerPart)
	for i, part := range minors {
		pv.Version[i+1] = mustParseInt64(part)
	}
	return pv
}

func (d strict) toPerlVersion() VersionData {
	switch {
	case d.decimal != "":
		return d.decimalMatches.toPerlVersion(d.original)
	case d.dotted != "":
		return d.dottedMatches.toPerlVersion(
			d.original)
	default:
		panic("logic error: strictRegexp matched but no Version found")
	}
}

func StrictVersion(matches []string) VersionData {
	return strict{
		original: matches[0],
		decimal:  matches[1],
		decimalMatches: strictDecimalForm{
			integerPart:  matches[2],
			fractionPart: matches[3],
		},
		dotted: matches[4],
		dottedMatches: strictDottedForm{
			integerPart: matches[5],
			dottedGroup: matches[6],
		},
	}.toPerlVersion()
}
