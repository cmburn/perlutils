package metacpanclient

import (
	"testing"

	// local
	"github.com/cmburn/perlutils/internal"
)

func runEnumParseTest[T internal.XEnumKind[T]](t *testing.T,
	outputs map[string]T,
	errVal T) {
	var value T
	for k, v := range outputs {
		out, err := value.Parse(k)
		if err != nil {
			t.Fatal(err)
		}
		if out != v {
			t.Error("wrong value")
		}
	}
	out, err := value.Parse("")
	if err == nil {
		t.Fatal("expected error")
	}
	if out != errVal {
		t.Error("wrong value")
	}
}

func runEnumStringTest[T internal.XEnumKind[T]](t *testing.T,
	inputs map[T]string) {
	for k, v := range inputs {
		if k.String() != v {
			t.Error("wrong value")
		}
	}
	if T(-1).String() != internal.Undef {
		t.Fatal("wrong value")
	}
}

func TestMaturityKind_Parse(t *testing.T) {
	t.Parallel()
	runEnumParseTest[MaturityKind](t, map[string]MaturityKind{
		strMaturityKindReleased:  MaturityKindReleased,
		strMaturityKindDeveloper: MaturityKindDeveloper,
	}, maturityKindUndef)
}

func TestMaturityKind_String(t *testing.T) {
	t.Parallel()
	runEnumStringTest[MaturityKind](t, map[MaturityKind]string{
		MaturityKindReleased:  strMaturityKindReleased,
		MaturityKindDeveloper: strMaturityKindDeveloper,
	})
}

func TestPhaseKind_Parse(t *testing.T) {
	t.Parallel()
	runEnumParseTest[PhaseKind](t, map[string]PhaseKind{
		strPhaseKindConfigure: PhaseKindConfigure,
		strPhaseKindBuild:     PhaseKindBuild,
		strPhaseKindRuntime:   PhaseKindRuntime,
		strPhaseKindTest:      PhaseKindTest,
		strPhaseKindDevelop:   PhaseKindDevelop,
	}, phaseKindUndef)
}

func TestPhaseKind_String(t *testing.T) {
	t.Parallel()
	runEnumStringTest[PhaseKind](t, map[PhaseKind]string{
		PhaseKindConfigure: strPhaseKindConfigure,
		PhaseKindBuild:     strPhaseKindBuild,
		PhaseKindRuntime:   strPhaseKindRuntime,
		PhaseKindTest:      strPhaseKindTest,
		PhaseKindDevelop:   strPhaseKindDevelop,
	})
}

func TestRelationshipKind_Parse(t *testing.T) {
	t.Parallel()
	runEnumParseTest[RelationshipKind](t, map[string]RelationshipKind{
		strRelationshipKindRequires:   RelationshipKindRequires,
		strRelationshipKindRecommends: RelationshipKindRecommends,
		strRelationshipKindSuggests:   RelationshipKindSuggests,
		strRelationshipKindConflicts:  RelationshipKindConflicts,
	}, relationshipKindUndef)
}

func TestRelationshipKind_String(t *testing.T) {
	t.Parallel()
	runEnumStringTest[RelationshipKind](t, map[RelationshipKind]string{
		RelationshipKindRequires:   strRelationshipKindRequires,
		RelationshipKindRecommends: strRelationshipKindRecommends,
		RelationshipKindSuggests:   strRelationshipKindSuggests,
		RelationshipKindConflicts:  strRelationshipKindConflicts,
	})
}

func TestReleaseStatusKind_Parse(t *testing.T) {
	t.Parallel()
	runEnumParseTest[ReleaseStatusKind](t, map[string]ReleaseStatusKind{
		strReleaseStatusKindLatest:  ReleaseStatusKindLatest,
		strReleaseStatusKindCPAN:    ReleaseStatusKindCPAN,
		strReleaseStatusKindBackpan: ReleaseStatusKindBackpan,
	}, releaseStatusKindUndef)
}

func TestReleaseStatusKind_String(t *testing.T) {
	t.Parallel()
	runEnumStringTest[ReleaseStatusKind](t, map[ReleaseStatusKind]string{
		ReleaseStatusKindLatest:  strReleaseStatusKindLatest,
		ReleaseStatusKindCPAN:    strReleaseStatusKindCPAN,
		ReleaseStatusKindBackpan: strReleaseStatusKindBackpan,
	})
}
