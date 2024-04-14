package client

import (
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/codec"
)

type rpcRequest struct {
	codec     codec.Codec
	body      interface{}
	service   string
	method    string
	endpoint  string
	protocols *micro.Protocols
	version   *micro.Version
}

func (r *rpcRequest) Protocols() *micro.Protocols {
	return r.protocols
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
