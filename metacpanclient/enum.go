package metacpanclient

import (
	"errors"
	"fmt"

	// local
	"github.com/cmburn/perlutils/internal"
)

type Maturity = internal.XEnum[MaturityKind]
type Phase = internal.XEnum[PhaseKind]
type Relationship = internal.XEnum[RelationshipKind]
type ReleaseStatus = internal.XEnum[ReleaseStatusKind]

type MaturityKind int
type PhaseKind int
type RelationshipKind int
type ReleaseStatusKind int

const (
	maturityKindUndef MaturityKind = iota
	MaturityKindReleased
	MaturityKindDeveloper
)

const (
	strMaturityKindReleased  = "released"
	strMaturityKindDeveloper = "developer"
)

func (m MaturityKind) String() string {
	switch m {
	case MaturityKindReleased:
		return strMaturityKindReleased
	case MaturityKindDeveloper:
		return strMaturityKindDeveloper
	case maturityKindUndef:
		fallthrough
	default:
		return internal.Undef
	}
}

func (m MaturityKind) Parse(s string) (MaturityKind, error) {
	switch s {
	case "released":
		return MaturityKindReleased, nil
	case "developer":
		return MaturityKindDeveloper, nil
	default:
		return maturityKindUndef, errors.New("invalid maturity kind")
	}
}

const (
	phaseKindUndef PhaseKind = iota
	PhaseKindBuild
	PhaseKindConfigure
	PhaseKindDevelop
	PhaseKindRuntime
	PhaseKindTest
)

const (
	strPhaseKindBuild     = "build"
	strPhaseKindConfigure = "configure"
	strPhaseKindDevelop   = "develop"
	strPhaseKindRuntime   = "runtime"
	strPhaseKindTest      = "test"
)

func (p PhaseKind) String() string {
	switch p {
	case PhaseKindConfigure:
		return strPhaseKindConfigure
	case PhaseKindBuild:
		return strPhaseKindBuild
	case PhaseKindRuntime:
		return strPhaseKindRuntime
	case PhaseKindTest:
		return strPhaseKindTest
	case PhaseKindDevelop:
		return strPhaseKindDevelop
	case phaseKindUndef:
		fallthrough
	default:
		return internal.Undef
	}
}

func (p PhaseKind) Parse(s string) (PhaseKind, error) {
	switch s {
	case strPhaseKindConfigure:
		return PhaseKindConfigure, nil
	case strPhaseKindBuild:
		return PhaseKindBuild, nil
	case strPhaseKindRuntime:
		return PhaseKindRuntime, nil
	case strPhaseKindTest:
		return PhaseKindTest, nil
	case strPhaseKindDevelop:
		return PhaseKindDevelop, nil
	default:
		return phaseKindUndef, errors.New("invalid phase kind")
	}
}

const (
	relationshipKindUndef RelationshipKind = iota
	RelationshipKindRequires
	RelationshipKindRecommends
	RelationshipKindSuggests
	RelationshipKindConflicts
)

const (
	strRelationshipKindRequires   = "requires"
	strRelationshipKindRecommends = "recommends"
	strRelationshipKindSuggests   = "suggests"
	strRelationshipKindConflicts  = "conflicts"
)

func (r RelationshipKind) String() string {
	switch r {
	case RelationshipKindRequires:
		return strRelationshipKindRequires
	case RelationshipKindRecommends:
		return strRelationshipKindRecommends
	case RelationshipKindSuggests:
		return strRelationshipKindSuggests
	case RelationshipKindConflicts:
		return strRelationshipKindConflicts
	case relationshipKindUndef:
		fallthrough
	default:
		return internal.Undef
	}
}

func (r RelationshipKind) Parse(s string) (RelationshipKind, error) {
	switch s {
	case strRelationshipKindRequires:
		return RelationshipKindRequires, nil
	case strRelationshipKindRecommends:
		return RelationshipKindRecommends, nil
	case strRelationshipKindSuggests:
		return RelationshipKindSuggests, nil
	case strRelationshipKindConflicts:
		return RelationshipKindConflicts, nil
	default:
		return relationshipKindUndef,
			errors.New("invalid relationship kind")
	}
}

const (
	releaseStatusKindUndef ReleaseStatusKind = iota
	ReleaseStatusKindLatest
	ReleaseStatusKindCPAN
	ReleaseStatusKindBackpan
)

const (
	strReleaseStatusKindLatest  = "latest"
	strReleaseStatusKindCPAN    = "cpan"
	strReleaseStatusKindBackpan = "backpan"
)

func (s ReleaseStatusKind) String() string {
	switch s {
	case ReleaseStatusKindLatest:
		return strReleaseStatusKindLatest
	case ReleaseStatusKindCPAN:
		return strReleaseStatusKindCPAN
	case ReleaseStatusKindBackpan:
		return strReleaseStatusKindBackpan
	case releaseStatusKindUndef:
		fallthrough
	default:
		return internal.Undef
	}
}

func (s ReleaseStatusKind) Parse(str string) (ReleaseStatusKind, error) {
	switch str {
	case "latest":
		return ReleaseStatusKindLatest, nil
	case "cpan":
		return ReleaseStatusKindCPAN, nil
	case "backpan":
		return ReleaseStatusKindBackpan, nil
	default:
		return releaseStatusKindUndef, fmt.Errorf("unknown status: %s", str)
	}
}
