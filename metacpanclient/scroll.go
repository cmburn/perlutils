package metacpanclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	// local
	"github.com/cmburn/perlutils/internal"
)

type Scroll[T result] struct {
	aggregations map[string]interface{}
	baseURL      string
	buffer       []T
	bufferIndex  int
	mc           *Client
	scrollID     string
	timeout      time.Duration
	lastUpdate   time.Time
	total        int
	size         uint16
	registered   bool
	mu           sync.Mutex
}

func (s *Scroll[T]) kill(force bool) (bool, error) {
	if force {
		s.mu.Lock()
	} else if !s.mu.TryLock() {
		return false, nil
	}
	defer s.mu.Unlock()
	sb := strings.Builder{}
	sb.WriteString(s.baseURL)
	sb.WriteRune('/')
	sb.WriteString("_search/scroll?scroll=")
	sb.WriteString(adjustTimeString(s.timeout))
	req, err := http.NewRequest(http.MethodDelete, sb.String(),
		strings.NewReader(s.scrollID))
	if err != nil {
		return true, fmt.Errorf("error creating DELETE request for "+
			"scroll '%s': %s", s.scrollID, err)
	}
	switch resp, err := s.mc.auto.RoundTrip(req); {
	case err != nil:
		return true, fmt.Errorf("error sending DELETE request for "+
			"scroll '%s': %s", s.scrollID, err)
	case resp.StatusCode != http.StatusOK:
		return true, fmt.Errorf("error sending DELETE request for "+
			"scroll '%s': %s", s.scrollID, resp.Status)
	default:
		internal.CloseBody(resp.Body)
		return true, nil
	}
}

func (s *Scroll[T]) time() time.Duration {
	return s.timeout
}

func (s *Scroll[T]) id() string {
	return s.scrollID
}

func (s *Scroll[T]) timedOut() bool {
	return time.Since(s.lastUpdate) > s.timeout
}

// mortal is a type that can be killed. Really only intended for *Scroll[T],
// but since we use a generic type we needed to make it an interface
type mortal interface {
	kill(bool) (bool, error)
	time() time.Duration
	timedOut() bool
	id() string
}

func (s *Scroll[T]) Next() (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var null T
	if s.timedOut() {
		return null, errors.New("scroll timed out")
	}
	if !s.registered {
		s.mc.registerScroll(s)
		s.registered = true
	}
	if s.bufferIndex == len(s.buffer) {
		if err := s.fetchNext(); err != nil {
			return null, err
		}
	}
	if len(s.buffer) == 0 {
		return null, nil
	}
	s.lastUpdate = time.Now()
	v := s.buffer[s.bufferIndex]
	s.buffer[s.bufferIndex] = null
	s.bufferIndex++
	return v, nil
}

func (s *Scroll[T]) Total() int {
	return s.total
}

func (s *Scroll[T]) Type() Type {
	return getType[T]()
}

func (s *Scroll[T]) UnmarshalJSON(data []byte) error {
	var v struct {
		ScrollID string `json:"_scroll_id"`
		Hits     struct {
			Total int      `json:"total"`
			Hits  []hit[T] `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	items := make([]T, len(v.Hits.Hits))
	for i, hit := range v.Hits.Hits {
		items[i] = hit.Source
	}
	s.total = v.Hits.Total
	s.buffer = items
	s.scrollID = v.ScrollID
	s.lastUpdate = time.Now()
	return nil
}

func NewScroll[T result](baseURL string, dur time.Duration,
	size uint16) *Scroll[T] {
	return &Scroll[T]{
		timeout:    dur,
		size:       size,
		baseURL:    baseURL,
		lastUpdate: time.Now(),
	}
}

func (s *Scroll[T]) fetchNext() error {
	// assumes we've already locked the buffer in Next()
	if s.scrollID == "" {
		return errors.New("no scroll id")
	}
	if s.mc == nil {
		return errors.New("no client")
	}
	// have to do this directly, or we hit an instantiation error, because
	// we'd have to instantiate a Scroll[T] even though it would never be
	// used
	req := &request[*wrapper[struct {
		Hits struct {
			Hits []T `json:"hits"`
		} `json:"hits"`
	}]]{
		Domain:  "",
		BaseURL: "",
		debug:   s.mc.Debug(),
		mc:      s.mc,
	}
	scrollTime := s.timeout
	if scrollTime == 0 {
		scrollTime = s.mc.timeout
	}
	scrollSize := s.size
	if scrollSize == 0 {
		scrollSize = s.mc.size
	}
	path := fmt.Sprintf("/_search/scroll/%s?scroll=%s&size=%d",
		s.scrollID, adjustTimeString(scrollTime), scrollSize)
	set, err := req.fetch(path, nil)
	s.bufferIndex = 0
	if err != nil {

		return err
	}
	items := set.Result.Hits.Hits
	s.total = len(items)
	s.buffer = items
	return nil
}
