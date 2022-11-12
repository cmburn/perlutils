package metacpanclient

import (
	"strings"
	"time"
)

func buildRequestURL(path, target, query string) string {
	sb := strings.Builder{}
	sb.WriteString(path)
	sb.WriteRune('/')
	sb.WriteString(target)
	if query != "" {
		sb.WriteRune('?')
		sb.WriteString(query)
	}
	return sb.String()
}

func post[T any](mc *Client, path, target, query string,
	params map[string]interface{}) (T, error) {
	req := newRequest[*wrapper[T]](mc)
	url := buildRequestURL(path, target, query)
	res, err := req.Fetch(url, params)
	if err != nil {
		var null T
		return null, err
	}
	return res.Result, nil
}

func get[T any](mc *Client, path, item string) (T, error) {
	return post[T](mc, path, item, "", nil)
}

func getResult[T result](mc *Client, path, item string) (T,
	error) {
	req := newRequest[T](mc)
	url := buildRequestURL(path, item, "")
	v, err := req.Fetch(url, nil)
	if err != nil {
		var null T
		return null, err
	}
	v.setClient(mc)
	return v, nil
}

func query[T any](mc *Client, path, item, query string) (T, error) {
	return post[T](mc, path, item, query, nil)
}

func doSearch[T result](mc *Client, params map[string]interface{}) (*Scroll[T],
	error) {
	req := newRequest[T](mc)
	res, err := req.SSearch(params)
	if err != nil {
		return nil, err
	}
	res.timeout = mc.timeout
	return res, nil
}

func rsSearch[T result](mc *Client, params map[string]interface{}) (
	ResultSet[T], error) {
	s, err := doSearch[T](mc, params)
	if err != nil {
		return nil, err
	}
	return newResultSetScroll[T](s), nil
}

func adjustTimeString(t time.Duration) string {
	str := t.String()
	str = strings.TrimSuffix(str, "0s")
	str = strings.TrimSuffix(str, "0m")
	return str
}
