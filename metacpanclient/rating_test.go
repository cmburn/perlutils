package metacpanclient

import "testing"

func TestRating_type(t *testing.T) {
	t.Parallel()
	r := &Rating{}
	if r._type() != TypeRating {
		t.Error("unexpected Rating._type")
	}
}
