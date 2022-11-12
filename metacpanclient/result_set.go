package metacpanclient

import (
	"encoding/json"
)

type ResultSet[T result] interface {
	HasScroller() bool
	Scroller() *Scroll[T]
	Items() []T
	Aggregations() map[string]interface{}
	Next() (T, error)
	Total() int

	_type() Type
	setClient(*Client)
}

type resultSetFetch[T result] struct {
	itemIndex int
	items     []T
	total     int
	client    *Client
}

func (r *resultSetFetch[T]) setClient(mc *Client) {
	for _, item := range r.items {
		item.setClient(mc)
	}
	r.client = mc
}

func (r *resultSetFetch[T]) Aggregations() map[string]interface{} {
	return nil
}

func (r *resultSetFetch[T]) HasScroller() bool {
	return false
}

func (r *resultSetFetch[T]) Scroller() *Scroll[T] {
	return nil
}

func (r *resultSetFetch[T]) _type() Type {
	return getType[T]()
}

func (r *resultSetFetch[T]) Items() []T {
	return r.items[r.itemIndex:]
}

func (r *resultSetFetch[T]) Total() int {
	return r.total
}

func (r *resultSetFetch[T]) Next() (T, error) {
	var null T
	if r.itemIndex >= len(r.items) {
		return null, nil
	}
	v := r.items[r.itemIndex]
	r.itemIndex++
	return v, nil
}

func newResultSetFetch[T result](items []T, total int) *resultSetFetch[T] {
	r := &resultSetFetch[T]{
		items: items,
	}
	r.total = total
	return r
}

type resultSetScroll[T result] struct {
	scroller *Scroll[T]
}

func (r *resultSetScroll[T]) Aggregations() map[string]interface{} {
	return r.scroller.aggregations
}

func (r *resultSetScroll[T]) Next() (T, error) {
	v, err := r.scroller.Next()
	if err != nil {
		var blank T
		return blank, err
	}
	return v, nil
}

func (r *resultSetScroll[T]) Items() []T {
	return nil
}

func (r *resultSetScroll[T]) HasScroller() bool {
	return true
}

func (r *resultSetScroll[T]) Scroller() *Scroll[T] {
	return r.scroller
}

func (r *resultSetScroll[T]) _type() Type {
	return getType[T]()
}

func (r *resultSetScroll[T]) setClient(mc *Client) {
	r.scroller.mc = mc
}

func (r *resultSetScroll[T]) Total() int {
	return r.scroller.total
}

func (r *resultSetScroll[T]) UnmarshalJSON(data []byte) error {
	r.scroller = NewScroll[T]("", 0, 0)
	return json.Unmarshal(data, &r.scroller)
}

func newResultSetScroll[T result](s *Scroll[T]) *resultSetScroll[T] {
	return &resultSetScroll[T]{s}
}
