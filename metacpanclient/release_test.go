package metacpanclient

import (
	"strings"
	"testing"
)

func TestRelease_Changes(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	r, err := mc.Release("Mojolicious")
	if err != nil {
		t.Fatal(err)
	}
	c, err := r.Changes()
	if err != nil {
		t.Fatal(err)
	}
	if c == "" {
		t.Fatal("expected changes")
	}
	r.mc = nil
	_, err = r.Changes()
	if err != nil {
		t.Error(err)
	}
	r.changes = ""
	_, err = r.Changes()
	if err != ErrNilClient {
		t.Errorf("expected error")
	}
}

func TestRelease_MetaCPANURL(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	r, err := mc.Release("Mojolicious")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(r.MetaCPANURL(), "https://metacpan.org/release/"+
		"SRI/Mojolicious-") {
		t.Errorf("unexpected metacpan url")
	}
}
