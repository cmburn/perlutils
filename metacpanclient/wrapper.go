package metacpanclient

import (
	"encoding/json"
)

type wrapper[T any] struct {
	Result T
}

func (w *wrapper[T]) _type() Type {
	return typeUndef
}

func (w *wrapper[T]) setClient(mc *Client) { /* no-op */ }

func (w *wrapper[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &w.Result)
}
