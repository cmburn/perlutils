package cpanmeta

type Resources struct {
	// License is a list of URLs to the license for the distribution.
	License []string `json:"license"`

	// Homepage is a URL to the homepage for the distribution.
	Homepage string `json:"homepage"`

	// BugTracker contains information about the bug tracker for the
	// distribution.
	BugTracker struct {
		// Web is a URL to the web interface for the bug tracker.
		Web string `json:"web"`

		// Mailto is an email address to which bug reports should be
		// sent.
		MailTo string `json:"mailto"`
	} `json:"bugtracker"`

	// Repository is a URL to the repository for the distribution.
	Repository struct {
		// Type is the type of repository.
		Type string `json:"type"`

		// URL is the URL of the repository.
		URL string `json:"url"`

		// Web is a URL to the web interface for the repository.
		Web string `json:"web"`
	}
}
