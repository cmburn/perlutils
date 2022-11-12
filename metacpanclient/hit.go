package metacpanclient

import (
	"encoding/json"
	"fmt"
)

type hit[T result] struct {
	Source T    `json:"_source"`
	Type   Type `json:"_type"`
}

func (h *hit[T]) UnmarshalJSON(data []byte) error {
	// can't alias in generic functions yet
	var v struct {
		Source T    `json:"_source"`
		Type   Type `json:"_type"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t := getType[T]()
	if v.Type == TypeFile && t == TypeModule || v.Type == t {
		*h = v
		return nil
	}
	return fmt.Errorf("expected type %v, got %v",
		getType[T](), v.Type)
}
