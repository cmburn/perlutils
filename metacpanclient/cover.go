package metacpanclient

import (
	"encoding/json"

	// local
	pui "github.com/cmburn/perlutils/internal"
	"github.com/cmburn/perlutils/version"
)

// Cover contains information about a distribution's test coverage.
type Cover struct {
	// Criteria is the coverage criteria for the distribution.
	Criteria struct {
		// Branch contains the coverage statistic by code branch.
		Branch float64 `json:"branch"`

		// Condition is the coverage statistic by conditional.
		Condition float64 `json:"condition"`

		// Statement is the coverage statistic by statement.
		Statement float64 `json:"statement"`

		// Subroutine is the coverage statistic by subroutine.
		Subroutine float64 `json:"subroutine"`

		// Total is the overall coverage.
		Total float64 `json:"total"`
	} `json:"criteria"`

	// Distribution is the name of the distribution this Cover pertains to.
	Distribution string `json:"Distribution"`

	// Release is the name of the release this Cover pertains to.
	Release string `json:"release"`

	// Version is the version of the release this Cover pertains to.
	Version version.JSON `json:"version"`
}

func (c *Cover) UnmarshalJSON(data []byte) error {
	var v struct {
		Criteria struct {
			Branch     pui.CoercedFloat `json:"branch"`
			Condition  pui.CoercedFloat `json:"condition"`
			Statement  pui.CoercedFloat `json:"statement"`
			Subroutine pui.CoercedFloat `json:"subroutine"`
			Total      pui.CoercedFloat `json:"total"`
		} `json:"criteria"`
		Distribution string       `json:"Distribution"`
		Release      string       `json:"release"`
		Version      version.JSON `json:"version"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*c = Cover{
		Criteria: struct {
			// B
			Branch     float64 `json:"branch"`
			Condition  float64 `json:"condition"`
			Statement  float64 `json:"statement"`
			Subroutine float64 `json:"subroutine"`
			Total      float64 `json:"total"`
		}{
			Branch:     v.Criteria.Branch.Value,
			Condition:  v.Criteria.Condition.Value,
			Statement:  v.Criteria.Statement.Value,
			Subroutine: v.Criteria.Subroutine.Value,
			Total:      v.Criteria.Total.Value,
		},
		Distribution: v.Distribution,
		Release:      v.Release,
		Version:      v.Version,
	}
	return nil
}
