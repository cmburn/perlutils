package metacpanclient

import "testing"

func TestDistribution_Github(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	d, err := mc.Distribution("Mojolicious")
	if err != nil {
		t.Fatal(err)
	}
	gh := d.Github()
	if gh.Source == "" {
		t.Error("expected source")
	}
	if gh.Closed == 0 {
		t.Error("expected closed")
	}
}

func TestDistribution_RT(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	d, err := mc.Distribution("Mojolicious")
	if err != nil {
		t.Fatal(err)
	}
	rt := d.RT()
	if rt.Source == "" {
		t.Error("expected source")
	}
	if rt.Closed == 0 {
		t.Error("expected closed")
	}
}

func TestDistribution_type(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	d, err := mc.Distribution("Mojolicious")
	if err != nil {
		t.Fatal(err)
	}
	if d._type() != TypeDistribution {
		t.Errorf("unexpected type")
	}
}
