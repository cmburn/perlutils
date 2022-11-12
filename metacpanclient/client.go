// Package metacpanclient is a client for the MetaCPAN API, equivalent to
// MetaCPAN::Client. It is written to mirror the API and behavior as closely as
// possible, although the implementation varies in some places.
package metacpanclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	// third party
	"github.com/elastic/go-elasticsearch/v7"

	// local
	pui "github.com/cmburn/perlutils/internal"
	"github.com/cmburn/perlutils/version"
)

/*
Client is the main interface to the MetaCPAN API. It is equivalent to
[MetaCPAN::Client] in Perl. It is thread-safe.

[MetaCPAN::Client]: https://metacpan.org/pod/MetaCPAN::Client
*/
type Client struct {
	auto         *transport
	domains      []string
	es           *elasticsearch.Client
	post         *transport
	scrollCh     chan string
	scrolls      map[string]mortal
	scrollsMu    sync.Mutex
	stopCh       chan struct{}
	stopScrollCh chan struct{}
	timeout      time.Duration
	userAgent    string
	debug        bool
	size         uint16
}

func (mc *Client) AllAuthors() (ResultSet[*Author], error) {
	return rsSearch[*Author](mc, map[string]interface{}{
		"match_all": map[string]interface{}{},
	})
}

func (mc *Client) AllDistributions() (ResultSet[*Distribution], error) {
	return rsSearch[*Distribution](mc, map[string]interface{}{
		"match_all": map[string]interface{}{},
	})
}

func (mc *Client) AllFavorites() (ResultSet[*Favorite], error) {
	return rsSearch[*Favorite](mc, map[string]interface{}{
		"match_all": map[string]interface{}{},
	})
}

func (mc *Client) AllModules() (ResultSet[*Module], error) {
	return rsSearch[*Module](mc, map[string]interface{}{
		"match_all": map[string]interface{}{},
	})
}

func (mc *Client) AllReleases() (ResultSet[*Release], error) {
	return rsSearch[*Release](mc, map[string]interface{}{
		"match_all": map[string]interface{}{},
	})
}

func (mc *Client) Author(s string) (*Author, error) {
	return getResult[*Author](mc, "/author", s)
}

func (mc *Client) AuthorSearch(args map[string]interface{}) (
	ResultSet[*Author], error) {
	return rsSearch[*Author](mc, args)
}

func (mc *Client) Autocomplete(s string) ([]*File, error) {
	sb := strings.Builder{}
	sb.WriteString("q=")
	sb.WriteString(s)
	req, err := query[struct {
		Hits struct {
			Hits []struct {
				Type   Type  `json:"_type"`
				Fields *File `json:"fields"`
			} `json:"hits"`
		}
	}](mc, "/search/autocomplete", "", sb.String())
	if err != nil {
		return nil, err
	}

	files := make([]*File, 0, len(req.Hits.Hits))
	for _, hit := range req.Hits.Hits {
		if hit.Type != TypeFile {
			panic("unexpected type: " + hit.Type.String())
		}
		files = append(files, hit.Fields)
	}

	return files, nil
}

func (mc *Client) AutocompleteSuggest(s string) ([]*File, error) {
	sb := strings.Builder{}
	sb.WriteString("q=")
	sb.WriteString(s)
	req, err := query[struct {
		Suggestions []*File `json:"suggestions"`
	}](mc, "/search/autocomplete/suggest", "", sb.String())
	if err != nil {
		return nil, err
	}
	return req.Suggestions, nil
}

func (mc *Client) Close() error {
	mc.stopCh <- struct{}{}
	return nil
}

func (mc *Client) Cover(s string) (*Cover, error) {
	return get[*Cover](mc, "/cover", s)
}

func (mc *Client) Distribution(s string) (*Distribution, error) {
	return getResult[*Distribution](mc, "/distribution", s)
}

func (mc *Client) DistributionSearch(args map[string]interface{}) (
	ResultSet[*Distribution], error) {
	return rsSearch[*Distribution](mc, args)
}

func (mc *Client) DownloadURL(release string, r *version.Range,
	dev bool) (*DownloadURL, error) {
	qb := strings.Builder{}
	haveVersion := r != nil
	if haveVersion {
		qb.WriteString("version=")
		qb.WriteString(r.URLString())
		if dev {
			qb.WriteString("&dev=1")
		}
	} else if dev {
		qb.WriteString("dev=1")
	}
	return query[*DownloadURL](mc, "/download_url", release, qb.String())
}

func (mc *Client) Favorite(args map[string]interface{}) (ResultSet[*Favorite],
	error) {
	return rsSearch[*Favorite](mc, args)
}

func (mc *Client) File(s string) (*File, error) {
	return getResult[*File](mc, "/file", s)
}

func (mc *Client) Mirror(s string) (*Mirror, error) {
	return get[*Mirror](mc, "/mirror", s)
}

func (mc *Client) Module(s string) (*Module, error) {
	return getResult[*Module](mc, "/module", s)
}

func (mc *Client) ModuleSearch(args map[string]interface{}) (
	ResultSet[*Module], error) {
	return rsSearch[*Module](mc, args)
}

func (mc *Client) Package(s string) (*Package, error) {
	return get[*Package](mc, "/package", s)
}

func (mc *Client) Permission(s string) (*Permission, error) {
	return get[*Permission](mc, "/permission", s)
}

func (mc *Client) Pod(s string) (*Pod, error) {
	pod, err := NewPod(s, "", mc)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (mc *Client) Rating(args map[string]interface{}) (ResultSet[*Rating],
	error) {
	return rsSearch[*Rating](mc, args)
}

func (mc *Client) Recent(count uint16) (ResultSet[*Release], error) {
	return mc.recent(&searchConfig{
		query: map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		sort: []map[string]map[string]string{
			{
				"date": map[string]string{
					"order": "desc",
				},
			},
		},
		size: count,
	})
}

func (mc *Client) Release(s string) (*Release, error) {
	return getResult[*Release](mc, "/release", s)
}

func (mc *Client) ReleaseSearch(args map[string]interface{}) (
	ResultSet[*Release], error) {
	return rsSearch[*Release](mc, args)
}

func (mc *Client) ReleasedToday() (ResultSet[*Release], error) {
	return mc.recent(&searchConfig{
		query: map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		sort: []map[string]map[string]string{
			{
				"date": map[string]string{
					"order": "desc",
				},
			},
		},
		filter: map[string]interface{}{
			"range": map[string]interface{}{
				"date": map[string]interface{}{
					"from": "now-1d+0h",
				},
			},
		},
	})
}

func (mc *Client) ReverseDependencies(s string) (ResultSet[*Release], error) {
	params := map[string]interface{}{
		"size": mc.size,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"filter": map[string]interface{}{
			"and": []map[string]interface{}{
				{
					"term": map[string]interface{}{
						"status": "latest",
					},
				},
				{
					"term": map[string]interface{}{
						"authorized": true,
					},
				},
			},
		},
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	s = strings.ReplaceAll(s, "::", "-")
	url := buildRequestURL("/reverse_dependencies/dist", s, "")
	rc, err := mc.request("", url, body)
	defer pui.CloseBody(rc)
	if err != nil {
		return nil, err
	}
	var v struct {
		Data []*Release `json:"data"`
	}
	if err = json.NewDecoder(rc).Decode(&v); err != nil {
		return nil, err
	}
	return newResultSetFetch[*Release](v.Data, len(v.Data)), nil
}

// private methods

func (mc *Client) doRequest(domain, path string, body []byte,
	rc *io.ReadCloser, eb *strings.Builder) bool {
	var err error
	var resp *http.Response
	defer func() {
		if err != nil {
			eb.WriteString("request to domain '")
			eb.WriteString(domain)
			eb.WriteString("' failed: ")
			eb.WriteString(err.Error())
			eb.WriteRune('\n')
			if resp != nil {
				pui.CloseBody(resp.Body)
			}
		}
	}()
	b := (io.Reader)(nil)
	if body != nil {
		b = bytes.NewReader(body)
	}
	ub := strings.Builder{}
	ub.WriteString(domain)
	if !strings.HasSuffix(domain, "/") {
		ub.WriteRune('/')
	}
	ub.WriteString(path)
	var req *http.Request
	if b == nil {
		req, err = http.NewRequest(http.MethodGet, ub.String(), b)
	} else {
		req, err = http.NewRequest(http.MethodPost, ub.String(), b)
	}
	if err != nil {
		return false
	}
	if resp, err = mc.auto.RoundTrip(req); err != nil {
		return false
	}
	if resp.StatusCode < 300 {
		*rc = resp.Body
		return true
	}
	if err == nil {
		err = fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	return false
}

func (mc *Client) killScroll(id string, early bool) {
	mc.scrollsMu.Lock()
	defer mc.scrollsMu.Unlock()
	if scroll, ok := mc.scrolls[id]; ok {
		switch dead, err := scroll.kill(early); {
		case err != nil:
			log.Println(err)
			fallthrough
		case dead:
			delete(mc.scrolls, id)
		default:
			// scroll is still alive, so we'll try again later
			if !early {
				mc.registerScroll(scroll)
			}
		}
	}
}

func (mc *Client) recent(sc *searchConfig) (*resultSetFetch[*Release], error) {
	sc.mc = mc
	s, err := newSearch[*Release](sc)
	if err != nil {
		return nil, err
	}
	rs, err := s.do()
	if err != nil {
		return nil, err
	}
	rsf := &resultSetFetch[*Release]{
		total:  len(rs.scroller.buffer),
		items:  make([]*Release, len(rs.scroller.buffer)),
		client: mc,
	}
	copy(rsf.items, rs.scroller.buffer)
	return rsf, nil

}

func (mc *Client) registerScroll(scroll mortal) {
	mc.scrollsMu.Lock()
	defer mc.scrollsMu.Unlock()
	if mc.scrolls == nil {
		mc.scrolls = make(map[string]mortal)
	}
	id := scroll.id()
	t := scroll.time()
	mc.scrolls[id] = scroll
	go func() {
		time.Sleep(t)
		mc.scrollCh <- id
	}()
}

func (mc *Client) request(domain, path string, body []byte) (io.ReadCloser,
	error) {
	eb := strings.Builder{}
	var rc io.ReadCloser
	if domain != "" && mc.doRequest(domain, path, body, &rc, &eb) {
		return rc, nil
	}
	for _, domain := range mc.domains {
		if mc.doRequest(domain, path, body, &rc, &eb) {
			return rc, nil
		}
	}
	return nil, errors.New(strings.Trim(eb.String(), "\n"))
}

func (mc *Client) scrollKiller() {
	for {
		select {
		case <-mc.stopScrollCh:
			return
		case id := <-mc.scrollCh:
			mc.killScroll(id, false)
		}
	}
}

func (mc *Client) stopper() {
	<-mc.stopCh
	mc.stopScrollCh <- struct{}{}
	mc.scrollsMu.Lock()
	defer mc.scrollsMu.Unlock()
	for _, v := range mc.scrolls {
		// goroutine so if we have something else with a mutex lock on
		// us, we don't get stuck
		go func(v mortal) {
			_, err := v.kill(true)
			if err != nil {
				log.Println(err)
			}
		}(v)
	}
}

// Attributes

func (mc *Client) Domains() []string {
	dup := make([]string, len(mc.domains))
	copy(dup, mc.domains)
	return dup
}

func (mc *Client) Debug() bool {
	return mc.debug
}

func (mc *Client) UserAgent() string {
	return mc.userAgent
}

// NewClient returns a new MetaCPAN client. If domain is empty, the default
// domain is used ("https://fastapi.metacpan.org").
func NewClient(debug bool, scrollSize uint16, scrollTime time.Duration,
	userAgent string, domains ...string) (*Client, error) {
	const (
		defaultUserAgent = "perl_utils.metacpan.client/" +
			pui.PackageVersion
		defaultScrollSize = 100
		defaultScrollTime = 5 * time.Minute
	)
	if scrollSize == 0 {
		scrollSize = defaultScrollSize
	}
	if scrollTime == 0 {
		scrollTime = defaultScrollTime
	}
	if userAgent == "" {
		userAgent = defaultUserAgent
	}
	if domains == nil {
		domains = []string{"https://fastapi.metacpan.org/v1"}
	}
	for i, domain := range domains {
		d := strings.TrimSuffix(domain, "/")
		// d = strings.TrimSuffix(d, APIVersion)
		domains[i] = strings.TrimSuffix(d, "/")
	}
	auto := newTransport(userAgent, false, debug)
	post := newTransport(userAgent, true, debug)
	esConf := elasticsearch.Config{Addresses: domains, Transport: post}
	es, err := elasticsearch.NewClient(esConf)
	if err != nil {
		return nil, err
	}
	killCh := make(chan struct{})
	killScrollCh := make(chan struct{})
	scrollCh := make(chan string)
	mc := &Client{
		domains:      domains,
		debug:        debug,
		size:         scrollSize,
		timeout:      scrollTime,
		es:           es,
		auto:         auto,
		post:         post,
		stopCh:       killCh,
		stopScrollCh: killScrollCh,
		scrollCh:     scrollCh,
	}
	go mc.scrollKiller()
	go mc.stopper()
	return mc, nil
}

const (
	MetaCPANURL = "https://fastapi.metacpan.org/v1"
)
