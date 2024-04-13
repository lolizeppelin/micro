package micro

var (
	DefaultCodecs = map[string]string{

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
