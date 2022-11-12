package metacpanclient

import (
	"strings"
	"testing"
)

func TestNewPod(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	const prefix = "http://search.cpan.org/perldoc?"
	pod, err := NewPod("Mojolicious", prefix, mc)
	if err != nil {
		t.Fatal(err)
	}
	if pod.mc != mc {
		t.Error("unexpected NewPod")
	}
	if pod.Name != mojo {
		t.Error("unexpected NewPod")
	}
	if pod.URLPrefix != "http://search.cpan.org/perldoc?" {
		t.Error("unexpected NewPod")
	}
	if !strings.Contains(pod.HTML, prefix+"Mojolicious") {
		t.Error("unexpected NewPod")
	}
	_, err = NewPod("not a module", prefix, mc)
	if err == nil {
		t.Error("unexpected NewPod")
	}
	_, err = NewPod("Mojolicious", "", nil)
	if err != ErrNilClient {
		t.Error("unexpected NewPod")
	}
}
