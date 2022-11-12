package metacpanclient

import "github.com/cmburn/perlutils/version"

type Package struct {
	// Author is the PAUSE ID of the author
	Author string `json:"author"`

	// ModuleName is the name of the module of the package.
	ModuleName string `json:"module_name"`

	// Distribution is the distribution in which this module exists.
	Distribution string `json:"distribution"`

	// DistVersion is the latest version of the distribution.
	DistVersion version.JSON `json:"dist_version"`

	// Version is the latest version of the module
	Version version.JSON `json:"version"`
}
