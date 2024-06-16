package __

import "net/url"

func (x *Message) QueryParams() url.Values {
	v, _ := url.ParseQuery(x.Query)
	return v
}
