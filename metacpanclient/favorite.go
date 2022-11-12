package metacpanclient

import (
	"github.com/relvacode/iso8601"
)

// Favorite is an object specifying info about a favorite on MetaCPAN.
type Favorite struct {
	// Author is the PAUSE ID of the author.
	Author string `json:"author"`

	// Date is the date the user favorited the distribution.
	Date iso8601.Time `json:"date"`

	// Distribution is the name of the distribution.
	Distribution string `json:"distribution"`

	// ID is the ID of the favorite.
	ID string `json:"id"`

	// Release is the name of the release.
	Release string `json:"release"`

	// User is the UUID of the user who favorited the distribution.
	User string `json:"user"`

	hasUA
}

func (f *Favorite) _type() Type {
	return TypeFavorite
}
