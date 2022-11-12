package metacpanclient

import (
	"fmt"

	// local
	pui "github.com/cmburn/perlutils/internal"
)

type Type int

const (
	typeUndef Type = iota
	TypeAuthor
	TypeDistribution
	TypeFavorite
	TypeFile
	TypeModule
	TypeRating
	TypeRelease
)

func (r *Type) String() string {
	switch *r {
	case TypeAuthor:
		return "author"
	case TypeDistribution:
		return "distribution"
	case TypeFavorite:
		return "favorite"
	case TypeFile:
		return "file"
	case TypeModule:
		return "module"
	case TypeRating:
		return "rating"
	case TypeRelease:
		return "release"
	case typeUndef:
		return pui.Undef
	default:
		return ""
	}
}

func (r *Type) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"author"`:
		*r = TypeAuthor
	case `"distribution"`:
		*r = TypeDistribution
	case `"favorite"`:
		*r = TypeFavorite
	case `"file"`:
		*r = TypeFile
	case `"module"`:
		*r = TypeModule
	case `"rating"`:
		*r = TypeRating
	case `"release"`:
		*r = TypeRelease
	default:
		return fmt.Errorf("invalid result type: %s", string(data))
	}
	return nil
}

func (r *Type) MarshalJSON() ([]byte, error) {
	return pui.WrapEnumTypeJSON(r)
}
