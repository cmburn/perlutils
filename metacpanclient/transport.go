package metacpanclient

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type transport struct {
	userAgent string
	post      bool
	debug     bool
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.post {
		req.Method = "POST"
	}
	req.Header.Set("User-Agent", t.userAgent)
	if t.debug {
		// ok if not covered by tests
		dump, _ := httputil.DumpRequestOut(req, true)
		fmt.Printf("%s", dump)
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if t.debug && resp != nil {
		dump, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("%s", dump)
	}
	return resp, nil
}

func newTransport(userAgent string, post, debug bool) *transport {
	return &transport{userAgent: userAgent, post: post, debug: debug}
}
