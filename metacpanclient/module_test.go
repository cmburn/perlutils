package metacpanclient

import "testing"

func TestModule_MetaCPANURL(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	m, err := mc.Module("Mojolicious")
	if err != nil {
		t.Fatal(err)
	}
	if m.MetaCPANURL() != "https://fastapi.metacpan.org/v1/pod/release/"+
		"SRI/Mojolicious-9.29/lib/Mojolicious.pm" {
		t.Error("unexpected MetaCPANURL")
	}
}

func TestModule_Package(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	m, err := mc.Module("Mojolicious")
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := m.Package()
	if err != nil {
		t.Fatal(err)
	}
	if pkg == nil {
		t.Error("unexpected Package")
	}
	m.mc = nil
	pkg, err = m.Package()
	if err != nil {
		t.Fatal(err)
	}
	if pkg == nil {
		t.Error("unexpected Package")
	}
	m = &Module{}
	_, err = m.Package()
	if err != ErrNilClient {
		t.Error("unexpected Package")
	}
}

func TestModule_Permission(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	m, err := mc.Module("Mojolicious")
	if err != nil {
		t.Fatal(err)
	}
	perm, err := m.Permission()
	if err != nil {
		t.Fatal(err)
	}
	if perm == nil {
		t.Error("unexpected Permission")
	}
	m.mc = nil
	pkg, err := m.Permission()
	// ensure we've cached the permission
	if err != nil {
		t.Fatal(err)
	}
	if pkg == nil {
		t.Error("unexpected Permission")
	}
	m = &Module{}
	_, err = m.Permission()
	if err != ErrNilClient {
		t.Error("unexpected Permission")
	}
}

func TestModule_type(t *testing.T) {
	t.Parallel()
	m := &Module{}
	if m._type() != TypeModule {
		t.Error("unexpected _type")
	}
}

func TestModule_setClient(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	m := &Module{}
	m.setClient(mc)
	if m.mc != mc {
		t.Error("unexpected setClient")
	}
}
