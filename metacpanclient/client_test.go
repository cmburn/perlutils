package metacpanclient

import (
	"os"
	"strings"
	"testing"
	"time"

	// local
	"github.com/cmburn/perlutils/version"

	// external
	"github.com/relvacode/iso8601"
)

// The All* tests run against *every* record on CPAN, so they are slow and
// should only be run when preparing a release. It also puts a strain on
// MetaCPAN itself, and we don't want to be jerks.
//
// AllAuthors ~ 30s
// AllDistributions ~ 1m30s
// AllFavorites ~ 1m30s
// AllModules ~ 10m
// AllReleases ~ 5m50s

func TestClient_AllAuthors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestClient_AllAuthors")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	rs, err := mc.AllAuthors()
	if err != nil {
		t.Fatal(err)
	}
	if rs.Total() < 10000 {
		t.Errorf("not enough authors")
	}
	for a, err := rs.Next(); a != nil || err != nil; a, err = rs.Next() {
		if err != nil {
			t.Errorf("error: %s", err)
			if a == nil {
				break
			}
			continue
		}
	}
}

func TestClient_AllDistributions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestClient_AllDistributions")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	rs, err := mc.AllDistributions()
	if err != nil {
		t.Fatal(err)
	}
	if rs.Total() < 10000 {
		t.Errorf("not enough distributions")
	}
	for d, err := rs.Next(); d != nil || err != nil; d, err = rs.Next() {
		if err != nil {
			t.Errorf("error: %s", err)
			if d == nil {
				break
			}
			continue
		}
	}
}

func TestClient_AllFavorites(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestClient_AllFavorites")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	rs, err := mc.AllFavorites()
	if err != nil {
		t.Fatal(err)
	}
	if rs.Total() < 10000 {
		t.Errorf("not enough favorites")
	}
	for f, err := rs.Next(); f != nil || err != nil; f, err = rs.Next() {
		if err != nil {
			t.Errorf("error: %s", err)
			if f == nil {
				break
			}
			continue
		}
	}
}

func TestClient_AllModules(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestClient_AllModules")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	rs, err := mc.AllModules()
	if err != nil {
		t.Fatal(err)
	}
	if rs.Total() < 10000000 {
		t.Errorf("not enough modules")
	}
	for m, err := rs.Next(); m != nil || err != nil; m, err = rs.Next() {
		if err != nil {
			t.Errorf("error: %s", err)
			if m == nil {
				break
			}
			continue
		}
	}
}

func TestClient_AllReleases(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestClient_AllReleases")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	rs, err := mc.AllReleases()
	if err != nil {
		t.Fatal(err)
	}
	if rs.Total() < 100000 {
		t.Errorf("not enough releases")
	}
	for r, err := rs.Next(); r != nil || err != nil; r, err = rs.Next() {
		if err != nil {
			t.Errorf("error: %s", err)
			if r == nil {
				break
			}
			continue
		}
	}
}

func TestClient_Author(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	a, err := mc.Author("SRI")
	if err != nil {
		t.Fatal(err)
	}
	if a.Updated.After(sriUpdated) {
		t.Fatal("SRI updated since test was written")
	}
	if len(a.Blog) != 1 {
		t.Errorf("unexpected blog count")
	}
	if a.Blog[0].URL != "http://blog.kraih.com" {
		t.Errorf("unexpected blog URL")
	}
	if a.Blog[0].Feed != "http://feeds.feedburner.com/kraih" {
		t.Errorf("unexpected blog feed")
	}
	if len(a.Donation) != 3 {
		t.Errorf("unexpected donation count")
	}
	if a.Donation[0].Name != "paypal" {
		t.Errorf("unexpected donation name")
	}
	if a.Donation[0].ID != "kraihx@gmail.com" {
		t.Errorf("unexpected donation ID")
	}
	if a.Donation[1].Name != "wishlist" {
		t.Errorf("unexpected donation name")
	}
	if a.Donation[1].ID != "http://www.amazon.de/registry/wishlist/"+
		"28ZQM8GX6WJNF" {
		t.Errorf("unexpected donation ID")
	}
	if a.Donation[2].Name != "flattr" {
		t.Errorf("unexpected donation name")
	}
	if a.Donation[2].ID != "kraih" {
		t.Errorf("unexpected donation ID")
	}
	if a.Country != "DE" {
		t.Errorf("unexpected country")
	}
	if a.Name != "Sebastian Riedel" {
		t.Errorf("unexpected name")
	}
	if len(a.Website) != 1 {
		t.Errorf("unexpected website count")
	}
	if a.Website[0] != "http://mojolicio.us" {
		t.Errorf("unexpected website")
	}
	if a.GravatarURL != "https://secure.gravatar.com/avatar/4a49eb49e0b98"+
		"ed1a1fb30b7d39baac3?s=130&d=http%3A%2F%2Fwww.gravatar.com%2F"+
		"avatar%2Fbfa97d786f12ee3381f97bc909b88e11%3Fs%3D130%26d%3Did"+
		"enticon" {
		t.Errorf("unexpected gravatar URL")
	}
	if a.User != "6CqBa8BRRoSa-Bs2KlrqAQ" {
		t.Errorf("unexpected user ID")
	}
	for i, r := range []struct {
		name string
		id   string
	}{
		{name: "coderwall", id: "kraih"},
		{name: "github", id: "kraih"},
		{name: "gittip", id: "kraih"},
		{name: "googleplus", id: "107078354112715770721"},
		{name: "perlmonks", id: "sri"},
		{name: "tumblr", id: "kraih"},
		{name: "twitter", id: "kraih"},
	} {
		if a.Profile[i].Name != r.name {
			t.Errorf("unexpected profile name")
		}
		if a.Profile[i].ID != r.id {
			t.Errorf("unexpected profile ID")
		}
	}
}

func TestClient_AuthorSearch(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	args := map[string]interface{}{
		"name": "Sebastian Riedel",
	}
	rs, err := mc.AuthorSearch(args)
	if err != nil {
		t.Fatal(err)
	}
	if rs.Total() < 1 {
		t.Errorf("unexpected total")
	}
	foundSRI := false
	for a, err := rs.Next(); a != nil || err != nil; a, err = rs.Next() {
		if err != nil {
			t.Errorf("error: %s", err)
			if a == nil {
				break
			}
			continue
		}
		if a.User == "6CqBa8BRRoSa-Bs2KlrqAQ" {
			foundSRI = true
		}
	}
	if !foundSRI {
		t.Errorf("SRI not found")
	}
}

func TestClient_Autocomplete(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	acs, err := mc.Autocomplete(moose)
	if err != nil {
		t.Fatal(err)
	}
	if len(acs) < 1 {
		t.Errorf("expected at least one result")
	}
	foundMoose := false
	for _, ac := range acs {
		if ac.Distribution == moose {
			foundMoose = true
			break
		}
	}
	if !foundMoose {
		t.Errorf("expected to find Moose.pm")
	}
}

func TestClient_AutocompleteSuggest(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	acs, err := mc.AutocompleteSuggest(moose)
	if err != nil {
		t.Fatal(err)
	}
	if len(acs) < 1 {
		t.Errorf("expected at least one result")
	}
	foundMoose := false
	for _, ac := range acs {
		if ac.Distribution == moose {
			foundMoose = true
			break
		}
	}
	if !foundMoose {
		t.Errorf("expected to find Moose.pm")
	}
}

func TestClient_Cover(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	c, err := mc.Cover("Mojolicious-9.28")
	if err != nil {
		t.Fatal(err)
	}
	if c.Distribution != mojo {
		t.Errorf("unexpected distribution")
	}
	if !strings.HasPrefix(c.Release, "Mojolicious-") {
		t.Errorf("unexpected release")
	}
}

func TestClient_Distribution(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	d, err := mc.Distribution(mojo)
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != mojo {
		t.Errorf("unexpected name")
	}
}

func TestClient_DistributionSearch(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	args := map[string]interface{}{
		"name": mojo,
	}
	rs, err := mc.DistributionSearch(args)
	if err != nil {
		t.Fatal(err)
	}
	if rs.Total() != 1 {
		t.Errorf("unexpected total")
	}
	d, err := rs.Next()
	if err != nil {
		t.Fatal(err)
	}
	if d == nil {
		t.Fatal("no distribution")
	}
	if d.Name != mojo {
		t.Errorf("unexpected name")
	}
}

func TestClient_DownloadURL(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	d, err := mc.DownloadURL(mojo, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	const mojoPrefix = "https://cpan.metacpan.org/authors/id/S/SR/SRI/" +
		"Mojolicious-"
	if !strings.HasPrefix(d.DownloadURL, mojoPrefix) {
		t.Errorf("unexpected download URL")
	}
	vr := version.MustParseRange("9.27")
	// Test with a specific version
	d, err = mc.DownloadURL(mojo, vr, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(d.DownloadURL, mojoPrefix+"9.27.tar.gz") {
		t.Errorf("unexpected download URL")
	}

	// try with dev version
	_, err = mc.DownloadURL(mojo, vr, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Mirror(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	// this is currently the only mirror listed
	m, err := mc.Mirror("www.cpan.org")
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "www.cpan.org" {
		t.Errorf("expected www.cpan.org, got %s", m.Name)
	}
}

func TestClient_Module(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	m, err := mc.Module(mojo)

	if err != nil {
		t.Fatal(err)
	}

	if m.Abstract != "Real-time web framework" {
		t.Errorf("unexpected abstract")
	}

	if m.Author != "SRI" {
		t.Error("unexpected author")
	}

	if !m.Authorized {
		t.Error("unexpected authorized")
	}

	if m.Name != "Mojolicious.pm" {
		t.Error("unexpected name")
	}

	if m.Path != "lib/Mojolicious.pm" {
		t.Error("unexpected path")
	}

	if !strings.HasPrefix(m.Description, "An amazing real-time web") {
		t.Error("unexpected description")
	}
}

func TestClient_ModuleSearch(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	args := map[string]interface{}{
		"path": "lib/Mojolicious.pm",
	}
	rs, err := mc.ModuleSearch(args)
	if err != nil {
		t.Fatal(err)
	}
	if rs.Total() < 1 {
		t.Errorf("unexpected total")
	}
	foundMojo := false
	for m, err := rs.Next(); err == nil; m, err = rs.Next() {
		if err != nil {
			t.Fatal(err)
		}
		if m.Name == "Mojolicious.pm" {
			foundMojo = true
			break
		}
	}
	if !foundMojo {
		t.Errorf("expected to find Mojolicious.pm")
	}
}

func TestClient_Package(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	p, err := mc.Package(mojo)
	if err != nil {
		t.Fatal(err)
	}
	if p.ModuleName != mojo {
		t.Errorf("unexpected name")
	}
}

func TestClient_Permission(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	p, err := mc.Permission(mojo)
	if err != nil {
		t.Fatal(err)
	}
	if p.ModuleName != mojo {
		t.Errorf("unexpected name")
	}
}

func TestClient_Pod(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	p, err := mc.Pod(mojo)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(p.Plain, "NAME\n    Mojolicious - ") {
		t.Errorf("expected pod")
	}
}

func TestClient_Rating(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	r, err := mc.Rating(map[string]interface{}{
		"distribution": mojo,
	})
	if err != nil {
		t.Fatal(err)
	}
	rating, err := r.Next()
	if err != nil {
		t.Fatal(err)
	}
	if rating.Distribution != mojo {
		t.Errorf("unexpected distribution")
	}
}

func TestClient_Recent(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	r, err := mc.Recent(100)
	if err != nil {
		t.Fatal(err)
	}
	if r.Total() != 100 {
		t.Errorf("unexpected total")
	}
}

func TestClient_Release(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	r, err := mc.Release(mojo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Distribution != mojo {
		t.Errorf("unexpected distribution")
	}
}

func TestClient_ReleasedToday(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	r, err := mc.ReleasedToday()
	if err != nil {
		t.Fatal(err)
	}
	if r.Total() < 1 {
		t.Errorf("unexpected total")
	}
}

func TestClient_ReverseDependencies(t *testing.T) {
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	r, err := mc.ReverseDependencies(mojo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Total() < 1 {
		t.Errorf("expected at least one result")
	}

}

// Utility functions

func tCloseClient(mc *Client, t *testing.T) {
	err := mc.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func tNewClient(t *testing.T) *Client {
	mc, err := NewClient(debugClient, 0, 0, "", tDomains...)
	if err != nil {
		t.Fatal(err)
	}
	return mc
}

func init() {
	sri, err := iso8601.ParseString("2016-03-27T18:27:34")
	if err != nil {
		panic(err) // should never happen
	}
	sriUpdated = sri
	if v, ok := os.LookupEnv("METACPAN_DOMAINS"); ok {
		tDomains = strings.Split(v, ",")
	}
}

var (
	sriUpdated time.Time
	tDomains   []string
)

const (
	debugClient = false
	moose       = "Moose"
	mojo        = "Mojolicious"
)
