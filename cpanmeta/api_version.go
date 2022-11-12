package cpanmeta

import "github.com/cmburn/perlutils/version"

type APIVersion struct {
	// Version is the version of CPAN::Meta::Spec that the metadata
	// conforms to.
	Version version.JSON `json:"version"`

	// URL is the URL of the specification document.
	URL string `json:"url"`
}
