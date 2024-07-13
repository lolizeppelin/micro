package utils

import (
	"golang.org/x/exp/maps"
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
	buf := strings.Builder{}
	keys := maps.Keys(v)
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

func PopQuery(query string, skip ...string) (string, error) {
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

/*
ParseQuery query string转化为map
*/
func ParseQuery(query string) (values map[string]string, err error) {
	m, err := url.ParseQuery(query)
	if err != nil {
		return
	}
	values = make(map[string]string)
	for key, value := range m {
		values[key] = value[0]
	}
	return
}

/*
SortedQuery 按照key升序的顺序生成的query string
*/
func SortedQuery(values map[string]string) string {
	q := queryValues{}
	for k, v := range values {
		q.Add(k, v)
	}
	return q.Query()
}

/*
OrderedQuery 按照字典写入顺的生成query string
*/
func OrderedQuery(values *OrderedMap[string, string]) string {
	buf := strings.Builder{}
	for el := values.Front(); el != nil; el = el.Next() {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(el.Key)
		buf.WriteByte('=')
		buf.WriteString(el.Value)

	}
	return buf.String()

}
