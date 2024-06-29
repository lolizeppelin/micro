package jsonschema

type Request struct {
	Method string         `json:"method" description:"请求方法"`
	Query  map[string]any `json:"query,omitempty" description:"请求参数结构"`
	Body   ContentBody    `json:"body,omitempty" description:"请求载荷结构"`
}

type ContentBody struct {
	Type   string         `json:"type,omitempty" description:"返回类型"`
	Schema map[string]any `json:"schema,omitempty" description:"返回结构"`
}

type Response struct {
	Code        int          `json:"code" description:"返回码"`
	Description string       `json:"description,omitempty" description:"返回说明"`
	Body        *ContentBody `json:"body,omitempty" description:"返回结构"`
}

type APIPath struct {
	Path        string      `json:"path"  description:"接口"`
	Summary     string      `json:"summary,omitempty"  description:"接口概述"`
	Description string      `json:"description,omitempty"  description:"接口详细说明"`
	Request     *Request    `json:"request,omitempty" description:"请求"`
	Responses   []*Response `json:"responses,omitempty" description:"返回"`
}
