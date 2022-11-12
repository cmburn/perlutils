package metacpanclient

import (
	"strings"
	"testing"
)

func TestFile_MetaCPANURL(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	f, err := mc.File("SRI/Mojolicious-9.28/lib/Mojolicious.pm")
	if err != nil {
		t.Fatal(err)
	}
	if f.MetaCPANURL() != "https://fastapi.metacpan.org/v1/source/SRI/Mojolicious-9.28/lib/Mojolicious.pm" {
		t.Error("unexpected MetaCPANURL")
	}
}

func TestFile_Pod(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	f, err := mc.File("SRI/Mojolicious-9.28/lib/Mojolicious.pm")
	if err != nil {
		t.Fatal(err)
	}
	// ensure that we return an error if the kind is not supported
	_, err = f.Pod("foo")
	if err != ErrInvalidPodKind {
		t.Error("unexpected Pod")
	}
	for _, kind := range []PodKind{
		PodKindHTML,
		PodKindPlain,
		PodKindXMarkdown,
		PodKindXPod,
	} {
		pod, err := f.Pod(kind)
		if err != nil {
			t.Error(err)
		}
		if pod == "" {
			t.Error("unexpected Pod")
		}
	}
	f = &File{}
	_, err = f.Pod(PodKindHTML)
	if err != ErrNilClient {
		t.Error("unexpected Pod")
	}
}

func TestFile_Source(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	f, err := mc.File("SRI/Mojolicious-9.28/lib/Mojolicious.pm")
	if err != nil {
		t.Fatal(err)
	}
	src, err := f.Source()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(src, "package Mojolicious;") {
		t.Error("unexpected source")
	}
	f.src = ""
	f.Author = "///"
	_, err = f.Source()
	if err == nil {
		t.Error("unexpected source")
	}
	f = &File{}
	_, err = f.Source()
	if err != ErrNilClient {
		t.Error("unexpected source")
	}
}

func TestFile_type(t *testing.T) {
	t.Parallel()
	f := &File{}
	if f._type() != TypeFile {
		t.Error("unexpected _type")
	}
}
