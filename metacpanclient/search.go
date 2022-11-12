package metacpanclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	// local
	pui "github.com/cmburn/perlutils/internal"

	// external
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type search[T result] struct {
	fields        []string
	source        []string
	Query         map[string]interface{} `json:"query,omitempty"`
	filter        map[string]interface{}
	sort          []string
	haveWildcards bool
	scrollerSize  uint16
	scrollerTime  time.Duration
	mc            *Client
}

type searchConfig struct {
	fields []string
	source []string
	query  map[string]interface{}
	filter map[string]interface{}
	sort   []map[string]map[string]string
	time   time.Duration
	size   uint16
	mc     *Client
}

func newSearch[T result](config *searchConfig) (*search[T], error) {
	// fail early if we need to
	if config.query == nil {
		return nil, errNilQuery
	}
	if config.time > 24*time.Hour {
		return nil, errScrollTimeMax
	}
	if config.size > 10000 {
		return nil, errScrollSizeMax
	}
	wc, err := adjust(config.query)
	if err != nil {
		return nil, err
	}
	if config.time == 0 {
		config.time = config.mc.timeout
	}
	if config.size == 0 {
		config.size = config.mc.size
	}
	sorts := make([]string, 0, len(config.sort))
	for _, s := range config.sort {
		if len(s) != 1 {
			return nil, errInvalidSort
		}
		sb := strings.Builder{}
		for k, v := range s {
			if o, ok := v["order"]; ok {
				sb.WriteString(k)
				sb.WriteString(":")
				sb.WriteString(o)
			} else {
				return nil, errInvalidSort
			}
		}
		sorts = append(sorts, sb.String())
	}
	s := &search[T]{
		Query:         config.query,
		filter:        config.filter,
		fields:        config.fields,
		source:        config.source,
		scrollerTime:  config.time,
		scrollerSize:  config.size,
		sort:          sorts,
		haveWildcards: wc,
		mc:            config.mc,
	}
	return s, nil
}

func (s *search[T]) do() (*resultSetScroll[T], error) {
	opts, err := s.build(s.mc.es)
	if err != nil {
		return nil, err
	}
	res, err := s.mc.es.Search(opts...)
	if res != nil {
		defer pui.CloseBody(res.Body)
	}
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, errors.New(res.Status())
	}
	var rs resultSetScroll[T]
	if err := json.NewDecoder(res.Body).Decode(&rs); err != nil {
		return nil, err
	}
	rs.scroller.timeout = s.scrollerTime
	rs.scroller.mc = s.mc
	rs.scroller.baseURL = s.mc.domains[0]
	return &rs, nil
}

func (s *search[T]) build(es *elasticsearch.Client) (
	[]func(*esapi.SearchRequest), error) {
	nOpts := 2
	for _, cond := range []bool{s.fields != nil, s.scrollerSize > 0,
		s.scrollerTime > 0, s.sort != nil, s.source != nil,
		s.haveWildcards} {
		if cond {
			nOpts++
		}
	}
	opts := make([]func(*esapi.SearchRequest), 0, nOpts)
	body, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(body)
	opts = append(opts, es.Search.WithBody(reader))
	t := getType[T]()
	opts = append(opts, es.Search.WithIndex(t.String()))
	if s.haveWildcards {
		opts = append(opts, es.Search.WithAllowNoIndices(true))
	}
	if s.scrollerTime > 0 {
		opts = append(opts, es.Search.WithScroll(s.scrollerTime))
	}
	if s.scrollerSize > 0 {
		opts = append(opts, es.Search.WithSize(int(s.scrollerSize)))
	}
	if s.sort != nil {
		opts = append(opts, es.Search.WithSort(s.sort...))
	}
	if s.fields != nil {
		opts = append(opts, es.Search.WithStoredFields(s.fields...))
	}
	if s.source != nil {
		opts = append(opts, es.Search.WithSource(s.source...))
	}
	return opts, nil
}

func adjustCondition(m map[string]interface{},
	key, newKey string) (bool, error) {
	v := m[key]
	delete(m, key)
	wc := false
	if vm, ok := v.(map[string]interface{}); ok {
		arr := make([]interface{}, 0, len(vm))
		for k, v := range vm {
			arr = append(arr, map[string]interface{}{k: v})
		}
		v = arr
	}
	_, ok := v.([]interface{})
	if !ok {
		return false, fmt.Errorf("invalid %s condition", key)
	}
	for _, i := range v.([]interface{}) {
		haveWC, err := adjust(i)
		if err != nil {
			return false, err
		}
		wc = wc || haveWC
	}
	oldVal := m["bool"]
	if oldVal == nil {
		m["bool"] = map[string]interface{}{newKey: v}
	} else {
		oldVal.(map[string]interface{})[newKey] = v
	}
	return wc, nil
}

func isValue(i interface{}) bool {
	switch i.(type) {
	case float32, float64, int, int16, int32, int64, int8, string, uint,
		uint16, uint32, uint64, uint8, uintptr:
		return true
	default:
		return false
	}
}

func adjustMap(m map[string]interface{}) (bool, error) {
	isBasic := true
	minimumShouldMatch := false
	wc := false
	for _, k := range []string{"either", "all", "not"} {
		if _, ok := m[k]; ok {
			newKey := ""
			switch k {
			case "either":
				newKey = "should"
				minimumShouldMatch = true
			case "all":
				newKey = "must"
			case "not":
				newKey = "must_not"
			}
			isBasic = false
			haveWC, err := adjustCondition(m, k, newKey)
			if err != nil {
				return false, err
			}
			wc = wc || haveWC
		}
	}
	if minimumShouldMatch {
		m["bool"].(map[string]interface{})["minimum_should_match"] = 1
	}
	if isBasic {
		if len(m) != 1 {
			return false, errInvalidSearchEntry
		}
		var key string
		var value interface{}
		for l, w := range m {
			key = l
			value = w
		}
		if key == "match_all" {
			v, ok := value.(map[string]interface{})
			if !ok || len(v) != 0 {
				return false, errInvalidMatchAllEntry
			}
			// do nothing
			return wc, nil
		}
		if !isValue(value) {
			return false, errInvalidSearchEntry
		}
		queryType := "term"
		if str, ok := value.(string); ok {
			if strings.ContainsRune(str, '*') {
				queryType = "wildcard"
				wc = true
			}
		}
		m[queryType] = map[string]interface{}{
			key: value,
		}
		delete(m, key)
	}
	return wc, nil
}

func adjust(i interface{}) (bool, error) {
	switch v := i.(type) {
	case map[string]interface{}:
		m := v
		wc, err := adjustMap(m)
		if err != nil {
			return false, err
		}
		return wc, nil
	default:
		return false, errInvalidSearchEntry
	}
}

var (
	errNilQuery             = errors.New("query cannot be nil")
	errInvalidSort          = errors.New("invalid sort")
	errInvalidSearchEntry   = errors.New("invalid search entry")
	errInvalidMatchAllEntry = errors.New("invalid match_all entry")
	errScrollTimeMax        = errors.New("scroll time cannot be greater " +
		"than 24 hours")
	errScrollSizeMax = errors.New("scroll size cannot be greater " +
		"than 10000")
)
