package micro

import (
	"strings"
)

var (
	DefaultCodecs = map[string]string{

		"text/html":              "application/grpc+bytes",
		"text/plain":             "application/grpc+bytes",
		"application/grpc+json":  "application/grpc+json",
		"application/grpc+proto": "application/grpc+proto",
		"application/grpc+bytes": "application/grpc+bytes",

		"application/json":         "application/grpc+json",
		"application/grpc":         "application/grpc+proto",
		"application/protobuf":     "application/grpc+proto",
		"application/octet-stream": "application/grpc+bytes",
	}
)

const (
	HeaderNode  = "X-Node-Id" // 限定 node
	ContentType = "Content-Type"
	Accept      = "Accept"
)

func MatchCodec(protocol, codec string) bool {
	s := strings.Split(protocol, "+")
	if len(s) < 2 {
		return protocol == codec
	}
	return s[1] == codec
}

type Protocols struct {
	ContentType string // 原始
	Accept      string
	Reqeust     string
	Response    string
}
