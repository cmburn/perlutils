package internal

type VersionData struct {
	Original string
	Alpha    bool
	QV       bool
	Version  []int64
}

func init() {
	strictRegexp.Longest()
	laxRegexp.Longest()
}

// suffix notation:
//     R: regex constant
//     P: repeating/plus group
//     ${N}P: repeating N or more times
//     Nc: Non capturing group

// shared between lax and strict:
const (
	fractionR = `(\.[0-9]+)`
)

// strict regexes:
const (
	strictIntR         = `(0|[1-9][0-9]*)`
	strictDottedNcR    = `(?:\.[0-9]{1,3})`
	strictDotted2PR    = `(` + strictDottedNcR + `{2,})`
	strictDecimalFormR = `(` + strictIntR + fractionR + `?)`
	strictDottedFormR  = `(v` + strictIntR + strictDotted2PR + `)`
)

// lax regexes:
const (
	laxIntR         = `([0-9]+)`
	laxDottedNcR    = `(?:\.[0-9]+)`
	laxDotted2PR    = `(` + laxDottedNcR + `{2,})`
	laxDottedPR     = `(` + laxDottedNcR + `+)`
	laxAlphaR       = `(_[0-9]+)`
	laxUndefR       = `(undef)`
	laxDecimalFormR = `(` + laxIntR + `(?:` + fractionR + `|\.)?` +
		laxAlphaR + `?|` + fractionR + laxAlphaR + `?)`
	laxDottedFormR = `(v` + laxIntR + `(?:` + laxDottedPR + laxAlphaR +
		`?)?|` + laxIntR + `?` + laxDotted2PR + laxAlphaR + `?)`
)

// LaxVersionRegex is a regular expression that matches a Perl Version string,
// under the documented rules under Version::regexp. It is a direct adaptation
// to Go's regex-engine. The Lax Version has a few interesting edge cases, but
// so there's actually four different forms it has to cover.
const LaxVersionRegex = `(?:` + laxUndefR + `|` + laxDottedFormR + `|` +
	laxDecimalFormR + `)$`

// StrictVersionRegex is a regular expression that matches a Perl Version
// string,under the documented rules under Version::regexp. Strict versioning
// is highly recommended, both by the Perl project and someone who's just
// had to write a parser for the lax Version.
const StrictVersionRegex = `(?:` + strictDecimalFormR + `|` +
	strictDottedFormR + `)$`
