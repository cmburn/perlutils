package metacpanclient

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	// local
	pui "github.com/cmburn/perlutils/internal"
)

type request[T result] struct {
	Domain  string
	BaseURL string
	debug   bool
	mc      *Client
}

func (r *request[T]) fetch(s string, m map[string]interface{}) (T, error) {
	if r.mc == nil {
		var null T
		return null, ErrNilClient
	}
	var b []byte
	if m != nil {
		var err error
		b, err = json.Marshal(m)
		if err != nil {
			var null T
			return null, err
		}
	}
	var rc io.ReadCloser
	var err error
	rc, err = r.mc.request(r.Domain, s, b)
	if err != nil {
		var null T
		return null, err
	}
	defer pui.CloseBody(rc)
	var v T
	err = json.NewDecoder(rc).Decode(&v)
	if err != nil {
		var null T
		return null, err
	}
	return v, nil
}

type Request[T result] struct{ request[T] }

func (r *Request[T]) IsDebug() bool {
	return r.debug
}

func (r *Request[T]) EnableDebug() {
	r.debug = true
}

func (r *Request[T]) DisableDebug() {
	r.debug = false
}

func (r *Request[T]) SetClient(mc *Client) error {
	r.mc = mc
	return nil
}

func (r *Request[T]) Fetch(s string, m map[string]interface{}) (T, error) {
	if r.mc == nil {
		var null T
		return null, ErrNilClient
	}
	return r.request.fetch(s, m)
}

func (r *Request[T]) SSearch(params map[string]interface{}) (*Scroll[T],
	error) {
	if r.mc == nil {
		return nil, ErrNilClient
	}
	var query, filter map[string]interface{}
	var scrollerSize uint16
	var scrollerTime time.Duration
	var fields, source []string
	var sort []map[string]map[string]string

	if v, ok := params["filter"]; ok {
		if filter, ok = v.(map[string]interface{}); !ok {
			return nil, errFilterType
		}
		delete(params, "filter")
	}
	if v, ok := params["size"]; ok {
		v2, ok := v.(int)
		if !ok {
			return nil, errScrollerSizeType
		}
		scrollerSize = uint16(v2)
		delete(params, "size")
	}
	if v, ok := params["time"]; ok {
		if scrollerTime, ok = v.(time.Duration); !ok {
			return nil, errScrollerTimeType
		}
		delete(params, "time")
	}
	if v, ok := params["fields"]; ok {
		if fields, ok = v.([]string); !ok {
			return nil, errFieldsType
		}
		delete(params, "fields")
	}
	if v, ok := params["source"]; ok {
		if source, ok = v.([]string); !ok {
			return nil, errSourceType
		}
		delete(params, "source")
	}
	if v, ok := params["sort"]; ok {
		if sort, ok = v.([]map[string]map[string]string); !ok {
			return nil, errSortType
		}
		delete(params, "sort")
	}
	if v, ok := params["query"]; ok {
		if query, ok = v.(map[string]interface{}); !ok {
			return nil, errQueryType
		}
		delete(params, "query")
		if len(params) != 0 {
			return nil, errors.New("unknown parameters")
		}
	} else {
		query = params
	}
	sr, err := newSearch[T](&searchConfig{
		query:  query,
		filter: filter,
		size:   scrollerSize,
		time:   scrollerTime,
		fields: fields,
		source: source,
		sort:   sort,
		mc:     r.mc,
	})
	if err != nil {
		return nil, err
	}
	rs, err := sr.do()
	if err != nil {
		return nil, err
	}
	s := rs.Scroller()

	return s, err
}

func NewRequest[T result](domain, baseURL string, debug bool,
	mc *Client) *Request[T] {
	return &Request[T]{
		request[T]{
			Domain:  domain,
			BaseURL: baseURL,
			debug:   debug,
			mc:      mc,
		},
	}
}

func newRequest[T result](mc *Client) *Request[T] {
	return NewRequest[T]("", "", mc.Debug(), mc)
}

var (
	errFilterType = errors.New("filter must be of type " +
		"map[string]interface{}")
	errQueryType = errors.New("query must be of type " +
		"map[string]interface{}")
	errScrollerSizeType = errors.New("scroller_size must be of type int")
	errScrollerTimeType = errors.New("scroller_time must be of type " +
		"time.Duration")
	errFieldsType = errors.New("fields must be of type []string")
	errSourceType = errors.New("source must be of type []string")
	errSortType   = errors.New("sort must be of type " +
		"[]map[string]map[string]string")
)
