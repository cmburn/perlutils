package cpanmeta

import "github.com/cmburn/perlutils/version"

type Phase struct {
	// Conflicts is a list of modules that conflict with this phase.
	Conflicts map[string]version.JSON `json:"conflicts"`

	// Recommends is a list of modules that are recommended for this phase.
	Recommends map[string]version.JSON `json:"recommends"`

	// Requires is a list of modules that are required for this phase.
	Requires map[string]version.JSON `json:"requires"`

	// Suggests is a list of modules that are suggested for this phase.
	Suggests map[string]version.JSON `json:"suggests"`
}
