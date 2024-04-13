package utils

import (
	"golang.org/x/exp/slices"
	"net/url"
	"sort"
	"strings"
)

type queryValues map[string][]string

func (v queryValues) Add(key, value string) {
	v[key] = append(v[key], value)
}

func (v queryValues) Query() string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		for _, _v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(_v)
		}
	}
	return buf.String()
}

func StripQuery(query string, skip ...string) (string, error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return "", err
	}
	q := queryValues{}
	for k, v := range values {
		if slices.Contains(skip, k) {
			continue
		}
		for _, _v := range v {
			q.Add(k, _v)
		}
	}

	return q.Query(), nil
}

func BuildQuery(values map[string]string) string {
	q := queryValues{}
	for k, v := range values {
		q.Add(k, v)
	}
	return q.Query()
}
