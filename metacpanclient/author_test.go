package metacpanclient

import (
	"encoding/json"
	"testing"
)

func TestAuthor_MetaCPANURL(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	a, err := mc.Author("SRI")
	if err != nil {
		t.Fatal(err)
	}
	if a.MetaCPANURL() != "https://fastapi.metacpan.org/v1/author/SRI" {
		t.Errorf("unexpected metacpan url")
	}
}

func TestAuthor_Releases(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	a, err := mc.Author("SRI")
	if err != nil {
		t.Fatal(err)
	}
	r, err := a.Releases()
	if err != nil {
		t.Fatal(err)
	}
	if r.Total() < 1 {
		t.Errorf("expected at least one result")
	}
	a.mc = nil
	_, err = a.Releases()
	if err != nil {
		t.Fatal(err)
	}
	a.releases = nil
	_, err = a.Releases()
	if err != ErrNilClient {
		t.Errorf("expected error")
	}
}

func TestAuthor_type(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	a, err := mc.Author("SRI")
	if err != nil {
		t.Fatal(err)
	}
	if a._type() != TypeAuthor {
		t.Errorf("unexpected type")
	}
}

func TestAuthor_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	var a Author
	err := json.Unmarshal([]byte(`{"name":{}}`), &a)
	if err == nil {
		t.Error("expected error")
	}
	err = json.Unmarshal([]byte(`{"email": {}}`), &a)
	if err == nil {
		t.Error("expected error")
	}
	err = json.Unmarshal([]byte(`{"website": {}}`), &a)
	if err == nil {
		t.Error("expected error")
	}
}
