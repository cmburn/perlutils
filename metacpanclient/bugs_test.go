package metacpanclient

import (
	"encoding/json"
	"testing"
)

func TestBugs_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	var b Bugs
	err := json.Unmarshal([]byte(`{"open":1,"closed":1,"active":1,"source":{}}`), &b)
	if err == nil {
		t.Error("expected error")
	}
	err = json.Unmarshal([]byte(`{"open":{},"closed":1,"active":1,"source":"foo"}`), &b)
	if err == nil {
		t.Error("expected error")
	}
	err = json.Unmarshal([]byte(`{"open":1,"closed":{},"active":0,"source":"foo"}`), &b)
	if err == nil {
		t.Error("expected error")
	}
	err = json.Unmarshal([]byte(`{"open":1,"closed":1,"active":{},"source":"foo"}`), &b)
	if err == nil {
		t.Error("expected error")
	}
}
