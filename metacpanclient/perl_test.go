// This particular file is licensed under the same terms as Perl itself.

package metacpanclient

// I've attempted to keep the directly adapted tests in a separate file. These
// do not cover all cases; rather, they are a direct translation of the
// original tests. A full test suite is in client_test.go.

import (
	"os"
	"regexp"
	"strings"
	"testing"

	// local
	"github.com/cmburn/perlutils/version"
)

func TestPerl_MetaCPAN_Client_Author(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	author, err := mc.Author("XSAWYERX")
	if err != nil {
		t.Fatal(err)
	}
	if author.PauseID != "XSAWYERX" {
		t.Errorf("incorrect author")
	}
	args := map[string]interface{}{
		"either": []interface{}{
			map[string]interface{}{
				"name": "Dave *",
			},
			map[string]interface{}{
				"name": "David *",
			},
		},
	}
	daves, err := mc.AuthorSearch(args)
	if err != nil {
		t.Fatal(err)
	}
	mostDaves := daves.Total()
	if mostDaves < 200 {
		t.Errorf("not a lot of Daves")
	}
	args = map[string]interface{}{
		"either": []interface{}{
			map[string]interface{}{
				"name": "Dave *",
			},
			map[string]interface{}{
				"name": "David *",
			},
		},
		"not": []interface{}{
			map[string]interface{}{
				"name": "Dave S*",
			},
			map[string]interface{}{
				"name": "David S*",
			},
		},
	}
	daves, err = mc.AuthorSearch(args)
	if err != nil {
		t.Fatal(err)
	}
	if daves.Total() >= mostDaves {
		t.Errorf("too many Daves")
	}
	args = map[string]interface{}{
		"either": []interface{}{
			map[string]interface{}{
				"all": []interface{}{
					map[string]interface{}{
						"name": "Dave *",
					},
					map[string]interface{}{
						"email": "*gmail.com",
					},
				},
			},
			map[string]interface{}{
				"all": []interface{}{
					map[string]interface{}{
						"name": "David *",
					},
					map[string]interface{}{
						"email": "*gmail.com",
					},
				},
			},
		},
	}
	daves, err = mc.AuthorSearch(args)
	if err != nil {
		t.Fatal(err)
	}
	if daves.Total() > mostDaves {
		t.Errorf("too many Daves")
	}
	for dave, err := daves.Next(); dave != nil; dave, err = daves.Next() {
		if err != nil {
			t.Fatal(err)
		}
		foundGmail := false
		for _, email := range dave.Email {
			if strings.HasSuffix(email, "@gmail.com") {
				foundGmail = true
				break
			}
		}
		if !foundGmail {
			t.Errorf("Dave doesn't have a gmail address")
		}
	}
	daves = nil
	args = map[string]interface{}{
		"all": []interface{}{
			map[string]interface{}{
				"name": "John *",
			},
			map[string]interface{}{
				"email": "*gmail.com",
			},
		},
	}
	johns, err := mc.AuthorSearch(args)
	if err != nil {
		t.Fatal(err)
	}
	if johns.Total() <= 0 {
		t.Errorf("expected some Johns, got %d", johns.Total())
	}
	for john, err := johns.Next(); john != nil; john, err = johns.Next() {
		if err != nil {
			t.Fatal(err)
		}
		foundGmail := false
		for _, email := range john.Email {
			if strings.HasSuffix(email, "@gmail.com") {
				foundGmail = true
				break
			}
		}
		if !foundGmail {
			t.Errorf("John doesn't have a gmail address")
		}
	}
}

func TestPerl_MetaCPAN_Client_Cover(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	cover, err := mc.Cover("Moose-2.2007")
	if err != nil {
		t.Fatal(err)
	}
	if cover.Release != "Moose-2.2007" {
		t.Errorf("expected Moose-2.2007, got %s", cover.Release)
	}
}

func TestPerl_MetaCPAN_Client_Distribution(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	dist, err := mc.Distribution("Business-ISBN")
	if err != nil {
		t.Fatal(err)
	}
	if dist._type() != TypeDistribution {
		_t := dist._type()
		t.Fatalf("expected a distribution, got %s",
			_t.String())
	}
	if dist.Name != "Business-ISBN" {
		t.Errorf("expected Business-ISBN, got %s", dist.Name)
	}
	rt := dist.RT()
	if rt.Closed == 0 {
		t.Errorf("expected some closed RT tickets")
	}
	gh := dist.Github()
	if gh.Closed == 0 {
		t.Errorf("expected some closed GH tickets")
	}

}

func TestPerl_MetaCPAN_Client_DownloadURL(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	du, err := mc.DownloadURL("Moose", nil, false)
	if err != nil {
		t.Fatal(err) // probably won't pass anything else, bail
	}
	if !isSHA256Regex.MatchString(du.ChecksumSHA256) {
		t.Errorf("expected a checksum, got empty string")
	}
	vr, err := version.ParseRange("1.01")
	if err != nil {
		t.Error(err)
	}
	du, err = mc.DownloadURL("Moose", vr, false)
	if err != nil {
		t.Error(err)
	}
	if du.Version.Raw() != "1.01" {
		t.Error("incorrect version")
	}
	if du.DownloadURL != "https://cpan.metacpan.org/authors/id/F/FL/"+
		"FLORA/Moose-1.01.tar.gz" {
		t.Error("incorrect DownloadURL for Moose-1.01")
	}
	if du.ChecksumSHA256 != "f4424f4d709907dea8bc9de2a37b9d3f"+
		"ef4f87775a8c102f432c48a1fdf8067b" {
		t.Error("incorrect SHA256 for Moose-1.01.tar.gz")
	}
	if du.ChecksumMD5 != "f13f9c203d099f5dc6117f59bda96340" {
		t.Error("incorrect MD5 for Moose-1.01.tar.gz")
	}
	vr, err = version.ParseRange(">1.01,<=2.00")
	if err != nil {
		t.Fatal(err)
	}
	du, err = mc.DownloadURL("Moose", vr, false)
	if err != nil {
		t.Fatal(err)
	}
	if du.Version.Raw() != "1.07" {
		t.Error("incorrect version")
	}
	if du.DownloadURL != "https://cpan.metacpan.org/authors/id/F/FL/"+
		"FLORA/Moose-1.07.tar.gz" {
		t.Error("incorrect DownloadURL for Moose-1.07")
	}
	du, err = mc.DownloadURL("Moose", nil, true)
	if err != nil {
		t.Error(err)
	}
	if du.Status.Kind != ReleaseStatusKindLatest {
		t.Error("incorrect status")
	}
	vr, err = version.ParseRange(">0.21,<0.27")
	if err != nil {
		t.Fatal(err)
	}
	du, err = mc.DownloadURL("Try::Tiny", vr, true)
	if err != nil {
		t.Fatal(err)
	}
	if du.Version.Raw() != "0.22" {
		t.Error("incorrect version")
	}

}

func TestPerl_MetaCPAN_Client_Favorite(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	args := map[string]interface{}{
		"author": "XSAWYERX",
	}
	favs, err := mc.Favorite(args)
	if err != nil {
		t.Fatal(err)
	}
	if favs._type() != TypeFavorite {
		ty := favs._type()
		t.Errorf("expected type favorite, got %s", ty.String())
	}
	if !favs.HasScroller() {
		t.Errorf("expected scrolled")
	}
}

func TestPerl_MetaCPAN_Client_File(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	f, err := mc.File("DOY/Moose-2.0001/lib/Moose.pm")
	if err != nil {
		t.Fatal(err)
	}
	if f.Author != "DOY" {
		t.Errorf("expected DOY, got %s", f.Author)
	}
	if f.Distribution != moose {
		t.Errorf("expected Moose, got %s", f.Distribution)
	}
	if f.Name != "Moose.pm" {
		t.Errorf("expected Moose.pm, got %s", f.Name)
	}
	if f.Path != "lib/Moose.pm" {
		t.Errorf("expected lib/Moose.pm, got %s", f.Path)
	}
	if f.Release != "Moose-2.0001" {
		t.Errorf("expected Moose-2.2001, got %s", f.Release)
	}
	if f.Version.Raw() != "2.0001" {
		t.Errorf("expected 2.2001, got %s", f.Version.Raw())
	}
	const expected = MetaCPANURL + "/source/DOY/Moose-2.0001/lib/Moose.pm"
	if f.MetaCPANURL() != expected {
		t.Errorf("expected %s, got %s", expected, f.MetaCPANURL())
	}
	src, err := f.Source()
	if err != nil {
		t.Error(err)
	} else if !strings.Contains(src, "package Moose;") {
		t.Errorf("expected Moose.pm source, got %s", src)
	}
	pod, err := f.Pod("plain")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(pod, "Moose - A postmodern object system for Perl") {
		t.Errorf("doesn't seem to be Moose.pm POD, got %s", pod)
	}
}

func TestPerl_MetaCPAN_Client_Module(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	m, err := mc.Module("MetaCPAN::API")
	if err != nil {
		t.Fatal(err)
	}
	if m.Distribution != "MetaCPAN-API" {
		t.Errorf("expected MetaCPAN-API, got %s", m.Distribution)
	}
	if m.Name != "API.pm" {
		t.Errorf("expected API.pm, got %s", m.Name)
	}
	if m.Path != "lib/MetaCPAN/API.pm" {
		t.Errorf("expected lib/MetaCPAN/API.pm, got %s", m.Path)
	}

	rsm, err := mc.ModuleSearch(map[string]interface{}{
		"path": "lib/MetaCPAN/API.pm",
	})
	if err != nil {
		t.Fatal(err)
	}
	if rsm.Total() <= 0 {
		t.Errorf("expected at least one result")
	}
}

func TestPerl_MetaCPAN_Client_Package(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	pack, err := mc.Package("Moose")
	if err != nil {
		t.Fatal(err)
	}
	if pack.ModuleName != "Moose" {
		t.Errorf("expected Moose, got %s", pack.ModuleName)
	}
}

func TestPerl_MetaCPAN_Client_Permission(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	perm, err := mc.Permission("MooseX::Types")
	if err != nil {
		t.Fatal(err)
	}
	if perm.ModuleName != "MooseX::Types" {
		t.Errorf("expected MooseX::Types, got %s", perm.ModuleName)
	}
}

func TestPerl_MetaCPAN_Client_Pod(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	pod, err := mc.Pod("MetaCPAN::API")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(pod.Plain, "=head1") {
		t.Errorf("expected pod, got %s", pod.Plain)
	}
}

func TestPerl_MetaCPAN_Client_Rating(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	args := map[string]interface{}{
		"distribution": "Moose",
	}
	r, err := mc.Rating(args)
	if err != nil {
		t.Fatal(err)
	}
	rt, err := r.Next()
	if err != nil {
		t.Fatal(err)
	}
	if rt.Distribution != "Moose" {
		t.Errorf("expected Moose, got %s", rt.Distribution)
	}
}

func TestPerl_MetaCPAN_Client_Release(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	rel, err := mc.Release("MetaCPAN-API")
	if err != nil {
		t.Fatal(err)
	}
	if rel.Distribution != "MetaCPAN-API" {
		t.Errorf("expected MetaCPAN-API, got %s", rel.Distribution)
	}
	if !isSHA256Regex.MatchString(rel.ChecksumSHA256) {
		t.Errorf("expected SHA256, got %s", rel.ChecksumSHA256)
	}
	if !isMD5Regex.MatchString(rel.ChecksumMD5) {
		t.Errorf("expected MD5, got %s", rel.ChecksumMD5)
	}
}

func TestPerl_MetaCPAN_Client_ReverseDependencies(t *testing.T) {
	if !runPerlTests {
		t.Skip("skipping conformance test.")
	}
	t.Parallel()
	mc := tNewClient(t)
	defer tCloseClient(mc, t)
	rd, err := mc.ReverseDependencies("MetaCPAN::Client")
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for dep, err := rd.Next(); dep != nil; dep, err = rd.Next() {
		if err != nil {
			t.Error(err)
			continue
		}
		count++
	}
}

func init() {
	if s, b := os.LookupEnv("RUN_PERL_TESTS"); b && s == "1" {
		runPerlTests = true
	}
}

var (
	isMD5Regex    = regexp.MustCompile(`^[a-f0-9]{32}$`)
	isSHA256Regex = regexp.MustCompile(`^[a-f0-9]{64}$`)
	runPerlTests  = false
)
