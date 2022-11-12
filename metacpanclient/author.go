package metacpanclient

import (
	"encoding/json"
	"strings"

	// local
	pui "github.com/cmburn/perlutils/internal"

	// external
	"github.com/relvacode/iso8601"
)

// Author contains all available information about a specific author on CPAN.
type Author struct {
	// Blog is a list of blog URLs and feeds.
	Blog []struct {
		// URL is the URL of the blog.
		URL string `json:"url"`
		// Feed is the URL of the blog's RSS feed.
		Feed string `json:"feed"`
	} `json:"blog"`

	// Donation is a list of donation URLs.
	Donation []struct {
		// Name is the name of the donation service.
		Name string `json:"name"`

		// ID is the ID of the author on the donation service.
		ID string `json:"id"`
	} `json:"donation"`

	// Links is a list of links to the author's CPAN pages.
	Links struct {
		// BackCPANDirectory is a link to the author's BackCPAN
		// directory.
		BackCPANDirectory string `json:"backpan_directory"`

		// CPANDirectory is a link to the author's CPAN directory.
		CPANDirectory string `json:"cpan_directory"`

		// CPANTS is a link to the author's CPANTS page.
		CPANTS string `json:"cpants"`

		// CPANTestersMatrix is a link to the author's CPAN Testers
		// matrix page.
		CPANTestersMatrix string `json:"cpantesters_matrix"`

		// CPANTestersReports is a link to the author's CPAN Testers
		// reports page.
		CPANTestersReports string `json:"cpantesters_reports"`

		// MetaCPANExplorer is a link to the author's MetaCPAN
		// explorer page.
		MetaCPANExplorer string `json:"metacpan_explorer"`
	} `json:"links"`

	// PerlMongers is a set of Perl Mongers groups that the author
	// belongs to.
	PerlMongers []struct {
		// Name is the name of the perl monger's site.
		// Usually {city}.pm (i.e. Chicago.pm).
		Name string `json:"name"`
		// URL is the site's URL.
		URL string `json:"url"`
	} `json:"perlmongers"`

	// Profile contains a list of profiles the author has on various sites.
	Profile []struct {
		// Name is the name of the service the profile is on.
		Name string `json:"name"`
		// ID is the author's ID on the service.
		ID string `json:"id"`
	} `json:"profile"`

	// ReleaseCount is the number of releases the author has made.
	ReleaseCount struct {
		// BackpanOnly is the number of releases the author has
		// available on BackPAN.
		BackpanOnly int `json:"backpan-only"`

		// CPAN is the number of releases the author has available
		// on CPAN.
		CPAN int `json:"cpan"`

		// Latest is the number of unique releases the author has
		// made on CPAN.
		Latest int `json:"latest"`
	} `json:"release_count"`

	// ASCIIName is the author's name in ASCII, if available.
	ASCIIName string `json:"ascii_name"`

	// City is the author's city, if provided.
	City string `json:"city"`

	// Country is the author's country, if provided.
	Country string `json:"country"`

	// Directory is the author's directory on CPAN.
	Directory string `json:"dir"`

	// Email is a list of email addresses the author has provided.
	Email []string `json:"email"`

	// Extra is a list of extra information the author has provided.
	Extra map[string]string `json:"extra"`

	// GravatarURL is the URL of the author's gravatar.
	GravatarURL string `json:"gravatar_url"`

	// Name is the author's name.
	Name string `json:"name"`

	// PauseID is the author's PAUSE ID.
	PauseID string `json:"pauseid"`

	// RealName is the author's real name, if provided.
	Region string `json:"region"`

	// Updated is the last time the author's information was updated, in
	// ISO8601 format.
	Updated iso8601.Time `json:"updated"`

	// User is the author's unique user ID on MetaCPAN.
	User string `json:"user"`

	// Website is a list of URLs the author has provided.
	Website []string `json:"website"`

	hasUA
	releases ResultSet[*Release]
}

// MetaCPANURL returns the URL of the author's MetaCPAN page.
func (a *Author) MetaCPANURL() string {
	const baseURL = MetaCPANURL + "/author/"
	sb := strings.Builder{}
	sb.WriteString(baseURL)
	sb.WriteString(a.PauseID)
	return sb.String()
}

// Releases returns a ResultSet[*Release] of the author's releases.
func (a *Author) Releases() (ResultSet[*Release], error) {
	if a.releases != nil {
		return a.releases, nil
	}
	if a.mc == nil {
		return nil, ErrNilClient
	}
	r, err := a.mc.ReleaseSearch(map[string]interface{}{
		"author": a.PauseID,
	})
	if err != nil {
		return nil, err
	}
	a.releases = r
	return r, nil
}

func (a *Author) UnmarshalJSON(data []byte) error {
	var v struct {
		Blog []struct {
			URL  string `json:"url"`
			Feed string `json:"feed"`
		} `json:"blog"`
		Donation []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"donation"`
		Links struct {
			BackCPANDirectory  string `json:"backpan_directory"`
			CPANDirectory      string `json:"cpan_directory"`
			CPANTS             string `json:"cpants"`
			CPANTestersMatrix  string `json:"cpantesters_matrix"`
			CPANTestersReports string `json:"cpantesters_reports"`
			MetaCPANExplorer   string `json:"metacpan_explorer"`
		} `json:"links"`
		PerlMongers json.RawMessage
		Profile     []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"profile"`
		ASCIIName    string                 `json:"ascii_name"`
		City         string                 `json:"city"`
		Country      string                 `json:"country"`
		Directory    string                 `json:"dir"`
		Email        pui.CoercedStringArray `json:"email"`
		Extra        map[string]string      `json:"extra"`
		GravatarURL  string                 `json:"gravatar_url"`
		Name         string                 `json:"name"`
		PauseID      string                 `json:"pauseid"`
		Region       string                 `json:"region"`
		ReleaseCount struct {
			BackpanOnly int `json:"backpan-only"`
			CPAN        int `json:"cpan"`
			Latest      int `json:"latest"`
		} `json:"release_count"`
		Updated iso8601.Time           `json:"updated"`
		User    string                 `json:"user"`
		Website pui.CoercedStringArray `json:"website"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*a = Author{
		Blog:         v.Blog,
		Donation:     v.Donation,
		Links:        v.Links,
		Profile:      v.Profile,
		ASCIIName:    v.ASCIIName,
		City:         v.City,
		Country:      v.Country,
		Directory:    v.Directory,
		Extra:        v.Extra,
		GravatarURL:  v.GravatarURL,
		Name:         v.Name,
		PauseID:      v.PauseID,
		Region:       v.Region,
		ReleaseCount: v.ReleaseCount,
		Updated:      v.Updated,
		User:         v.User,
		Email:        v.Email.Value,
		Website:      v.Website.Value,
	}
	type pm = struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	if len(v.PerlMongers) > 0 {
		if v.PerlMongers[0] == '[' {
			var pms []pm
			if err := json.Unmarshal(v.PerlMongers,
				&pms); err != nil {
				return err
			}
			a.PerlMongers = make([]pm, len(pms))
			copy(a.PerlMongers, pms)
		} else {
			var p pm
			if err := json.Unmarshal(v.PerlMongers,
				&p); err != nil {
				return err
			}
			a.PerlMongers = []pm{p}
		}
	}
	return nil
}

func (a *Author) _type() Type {
	return TypeAuthor
}
