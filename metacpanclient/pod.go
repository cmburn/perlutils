package metacpanclient

import (
	"errors"
	"io"
	"strings"

	// local
	pui "github.com/cmburn/perlutils/internal"
)

type PodKind string

const (
	PodKindHTML      PodKind = "html"
	PodKindPlain     PodKind = "plain"
	PodKindXMarkdown PodKind = "x_markdown"
	PodKindXPod      PodKind = "x_pod"
)

type Pod struct {
	// Name is the name of the module this Pod belongs to.
	Name string `json:"name"`

	// URLPrefix is the url to be used for the links within the POD
	// documentation.
	URLPrefix string `json:"url_prefix"`

	// XPod is the raw POD for the file.
	XPod string `json:"x_pod"`

	// HTML is the POD in html form.
	HTML string `json:"html"`

	// XMarkdown is the POD in markdown form.
	XMarkdown string `json:"x_markdown"`

	// Plain is the POD formatted as plain text, such as you might find
	Plain string `json:"plain"`
	mc    *Client
}

func (p *Pod) load(kind PodKind) (string, error) {
	qb := strings.Builder{}
	qb.WriteString("content-type=text/")
	qb.WriteString(string(kind))
	if p.URLPrefix != "" {
		qb.WriteString("&url_prefix=")
		qb.WriteString(p.URLPrefix)
	}
	url := buildRequestURL("/pod", p.Name, qb.String())
	rc, err := p.mc.request("", url, nil)
	if err != nil {
		return "", err
	}
	defer pui.CloseBody(rc)
	buf, err := io.ReadAll(rc)
	if err != nil {
		// okay if not under go coverage
		return "", err
	}
	return string(buf), nil
}

func NewPod(name, urlPrefix string, mc *Client) (*Pod, error) {
	if mc == nil {
		return nil, ErrNilClient
	}
	p := &Pod{
		Name:      name,
		URLPrefix: urlPrefix,
		mc:        mc,
	}
	errCh := make(chan error)
	doneCh := make(chan struct{})
	res := make([]string, 4)
	for i, kind := range []PodKind{
		PodKindHTML,
		PodKindPlain,
		PodKindXMarkdown,
		PodKindXPod,
	} {
		go func(i int, kind PodKind) {
			s, err := p.load(kind)
			if err != nil {
				errCh <- err
				return
			}
			res[i] = s
			doneCh <- struct{}{}
		}(i, kind)
	}
	for i := 0; i < 4; i++ {
		select {
		case err := <-errCh:
			return nil, err
		case <-doneCh:
		}
	}
	p.HTML = res[0]
	p.Plain = res[1]
	p.XMarkdown = res[2]
	p.XPod = res[3]
	return p, nil
}

func (p *Pod) get(kind PodKind) (string, error) {
	switch kind {
	case PodKindHTML:
		return p.HTML, nil
	case PodKindPlain:
		return p.Plain, nil
	case PodKindXMarkdown:
		return p.XMarkdown, nil
	case PodKindXPod:
		return p.XPod, nil
	default:
		return "", ErrInvalidPodKind
	}
}

var (
	ErrNilClient = errors.New("nil client")
)
