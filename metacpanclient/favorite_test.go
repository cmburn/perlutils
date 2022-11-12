package metacpanclient

import "testing"

func TestFavorite_type(t *testing.T) {
	t.Parallel()
	f := &Favorite{}
	if f._type() != TypeFavorite {
		t.Errorf("unexpected type")
	}
}
