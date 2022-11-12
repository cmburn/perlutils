package cpanmeta

import (
	"encoding/json"

	// local
	"github.com/cmburn/perlutils/internal"
	"github.com/cmburn/perlutils/version"
)

type Spec struct {
	// Author is a list of the distribution's authors/maintainers.
	//
	// This takes the form of a list of strings in the form "NAME <EMAIL>".
	Author []string `json:"author"`

	// DynamicConfig is a boolean indicating whether the distribution
	// requires dynamic configuration through a Build.PL or Makefile.PL
	// script.
	DynamicConfig bool `json:"dynamic_config"`

	// GeneratedBy is the name of the program that generated the META.json
	// file.
	//
	// It takes the form of "NAME version VERSION"
	GeneratedBy string `json:"generated_by"`

	// License is a list of the distribution's licenses
	License []License `json:"license"`

	// MetaSpec is the version of the META spec that the META.json or
	// META.yml file conforms to.
	MetaSpec APIVersion `json:"meta-spec"`

	// Name is the name of the distribution.
	Name string `json:"name"`

	// ReleaseStatus is the release status of the distribution.
	ReleaseStatus ReleaseStatus `json:"release_status"`

	// Version is the version of the distribution.
	Version version.JSON `json:"version"`

	// Description is a short description of the distribution.
	Description string `json:"description"`

	// Keywords is a list of keywords that describe the distribution.
	Keywords []string `json:"keywords"`

	// NoIndex contains lists of directories and files that should not be
	// indexed by CPAN clients.
	NoIndex NoIndex `json:"no_index"`

	// OptionalFeatures is a map of optional features that the
	// distribution provides.
	OptionalFeatures map[string]OptionalFeature `json:"optional_features"`

	// Prereqs contains the prerequisites for the distribution.
	Prereqs Prereqs `json:"prereqs"`

	// Provides is a map of files that the distribution provides.
	Provides map[string]File `json:"provides"`

	// Resources contains info on resources related to the distribution.
	Resources Resources `json:"resources"`
}

func (s *Spec) UnmarshalJSON(data []byte) error {
	var v struct {
		Author           []string                   `json:"author"`
		DynamicConfig    interface{}                `json:"dynamic_config"`
		GeneratedBy      string                     `json:"generated_by"`
		License          []License                  `json:"license"`
		MetaSpec         APIVersion                 `json:"meta-spec"`
		Name             string                     `json:"name"`
		ReleaseStatus    ReleaseStatus              `json:"release_status"`
		Version          version.JSON               `json:"version"`
		Description      string                     `json:"description"`
		Keywords         []string                   `json:"keywords"`
		NoIndex          NoIndex                    `json:"no_index"`
		OptionalFeatures map[string]OptionalFeature `json:"optional_features"`
		Prereqs          Prereqs                    `json:"prereqs"`
		Provides         map[string]File            `json:"provides"`
		Resources        Resources                  `json:"resources"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if err := internal.CoerceBool(&s.DynamicConfig,
		v.DynamicConfig); err != nil {
		return err
	}
	s.Author = v.Author
	s.GeneratedBy = v.GeneratedBy
	s.License = v.License
	s.MetaSpec = v.MetaSpec
	s.Name = v.Name
	s.ReleaseStatus = v.ReleaseStatus
	s.Version = v.Version
	s.Description = v.Description
	s.Keywords = v.Keywords
	s.NoIndex = v.NoIndex
	s.OptionalFeatures = v.OptionalFeatures
	s.Prereqs = v.Prereqs
	s.Provides = v.Provides
	s.Resources = v.Resources
	return nil
}
