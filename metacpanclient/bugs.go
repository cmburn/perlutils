package metacpanclient

import (
	"encoding/json"

	// local
	pui "github.com/cmburn/perlutils/internal"
)

// Bugs contains information about bugs within a distribution
type Bugs struct {
	// Open is the number of open bugs on a given repo
	Open int `json:"open"`
	// Closed is the number of closed bugs on a given repo
	Closed int `json:"closed"`
	// Active is the number of active bugs on a given repo
	Active int `json:"active"`
	// Source is the URL of the repo
	Source string `json:"source"`
}

func (b *Bugs) UnmarshalJSON(data []byte) error {
	var v struct {
		Open   pui.CoercedInt `json:"open"`
		Closed pui.CoercedInt `json:"closed"`
		Active pui.CoercedInt `json:"active"`
		Source string         `json:"source"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*b = Bugs{
		Open:   v.Open.Value,
		Closed: v.Closed.Value,
		Active: v.Active.Value,
		Source: v.Source,
	}
	return nil
}
