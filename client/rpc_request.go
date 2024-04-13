package client

import (
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/codec"
)

type rpcRequest struct {
	codec    codec.Codec
	body     interface{}
	service  string
	method   string
	endpoint string
	protocol string
	accept   string
	version  *micro.Version
}

func (r *rpcRequest) ContentType() string {
	return r.protocol
}

func (r *rpcRequest) Accept() string {
	return r.protocol
}

func (r *rpcRequest) Service() string {
	return r.service
}

func (r *rpcRequest) Method() string {
	return r.method
}

func (r *rpcRequest) Endpoint() string {
	return r.endpoint
}

func (r *rpcRequest) Body() interface{} {
	return r.body
}

func (r *rpcRequest) Version() *micro.Version {
	return r.version
}
