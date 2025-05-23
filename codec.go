package micro

import (
	"strings"
)

var (
	DefaultCodecs = map[string]string{
		"text/html":                "application/grpc+bytes",
		"text/plain":               "application/grpc+bytes",
		"application/grpc+json":    "application/grpc+json",
		"application/grpc+proto":   "application/grpc+proto",
		"application/grpc+bytes":   "application/grpc+bytes",
		"application/grpc+msgpack": "application/grpc+msgpack",

		"application/json":         "application/grpc+json",
		"application/msgpack":      "application/grpc+msgpack",
		"application/grpc":         "application/grpc+proto",
		"application/protobuf":     "application/grpc+proto",
		"application/octet-stream": "application/grpc+bytes",
	}

	DefaultContentType = "application/grpc+bytes"
)

const (
	NodeHeader  = "X-Node-Id"      // 限定 node
	TokenHeader = "X-Auth-Token"   // 认证头
	TokenScope  = "X-Token-Scope"  // token范围
	TokenTenant = "X-Token-Tenant" // token限定租户范围

	ContentType = "Content-Type"
	Accept      = "Accept"
	Host        = "Host"
	Tenant      = "Tenant"
	PrimaryKey  = "PrimaryKey"
)

func MatchCodec(protocol, codec string) bool {
	ss := strings.Split(protocol, "/")
	if len(ss) < 2 {
		return protocol == codec
	}
	s := strings.Split(ss[1], "+")
	if len(s) < 2 {
		return ss[1] == codec
	}
	return s[1] == codec
}

type Protocols struct {
	ContentType string // 原始 ContentType
	Accept      string // 原始 Accept
	Reqeust     string // 请求
	Response    string // 返回
}

func GetProtocol(ContentType string) string {
	if protocol, ok := DefaultCodecs[ContentType]; ok {
		return protocol
	}
	return DefaultContentType
}
