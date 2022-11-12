package metacpanclient

import "testing"

func TestType_MarshalJSON(t *testing.T) {
	t.Parallel()
	for _, s := range []Type{
		TypeAuthor,
		TypeDistribution,
		TypeFavorite,
		TypeFile,
		TypeModule,
		TypeRating,
		TypeRelease,
	} {
		b, err := s.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != `"`+s.String()+`"` {
			t.Fatal("wrong value")
		}
	}
}

func TestType_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	for _, s := range []string{
		"author",
		"distribution",
		"favorite",
		"file",
		"module",
		"rating",
		"release",
	} {
		var r Type
		if err := r.UnmarshalJSON([]byte(`"` + s + `"`)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestType_String(t *testing.T) {
	t.Parallel()
	for _, s := range []struct {
		v Type
		s string
	}{
		{TypeAuthor, "author"},
		{TypeDistribution, "distribution"},
		{TypeFavorite, "favorite"},
		{TypeFile, "file"},
		{TypeModule, "module"},
		{TypeRating, "rating"},
		{TypeRelease, "release"},
		{Type(-1), ""},
	} {
		if s.v.String() != s.s {
			t.Fatal("wrong value")
		}
	}
}
