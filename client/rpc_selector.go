package client

import (
	"errors"
	"github.com/lolizeppelin/micro"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/selector"
	"github.com/lolizeppelin/micro/utils"
)

func (r *rpcClient) next(request micro.Request, opts CallOptions) (selector.Next, error) {

	endpoint := request.Endpoint()
	filters := opts.Filters
	version := request.Version()
	minor := version.Minor()
	protocols := request.Protocols()

	// 标准过滤器
	filters = utils.InsertSlice(opts.Filters,
		func(services []*micro.Service) ([]*micro.Service, error) {
			var matched []*micro.Service
			for _, s := range services {
				// 服务版本匹配
				if version.Main() != s.Version {
					continue
				}
				found := false
				// 服务endpoint匹配
				for _, ep := range s.Endpoints {
					if ep.Name == endpoint {
						// 返回与内部请求校验
						if opts.Internal || !ep.Internal {
							// 请求与返回协议校验
							if !micro.MatchCodec(protocols.Reqeust, ep.Metadata["res"]) ||
								!micro.MatchCodec(protocols.Response, ep.Metadata["req"]) {
								return nil, exc.BadRequest("go.micro.client.selector", "request or response type mismatch")
							}
							found = true
							break
						}
						return nil, exc.NotFound("go.micro.client.selector", "endpoint not found")
					}
				}
				if !found {
					return nil, micro.ErrSelectEndpointNotFound
				}
				// 节点过滤
				if opts.Node != "" || minor > 0 {
					service := &micro.Service{
						Name:      s.Name,
						Version:   s.Version,
						Metadata:  s.Metadata,
						Endpoints: s.Endpoints,
					}
					for _, node := range s.Nodes {
						if opts.Node != "" && node.Id == opts.Node { // 节点ID过滤
							service.Nodes = append(service.Nodes, node)
							return []*micro.Service{service}, nil
						} else if minor > 0 { // 节点版本过滤
							v, err := micro.NewVersion(node.Version)
							if err != nil {
								log.Errorf("node: %s version value error, skip minor filter", node.Id)
								continue
							}
							// 次要版本不匹配
							if minor > v.Minor() {
								continue
							}
							service.Nodes = append(service.Nodes, node)
						}
					}
					if len(service.Nodes) > 0 {
						matched = append(matched, service)
					}
				} else {
					matched = append(matched, s)
				}
			}
			return matched, nil
		})

	service := request.Service()
	// get next nodes from the selector
	next, err := r.opts.Selector.Select(service, filters...)
	if err != nil {
		if errors.Is(err, micro.ErrSelectServiceNotFound) {
			return nil, exc.ServiceUnavailable("go.micro.client", err.Error())
		}
		if errors.Is(err, micro.ErrNoneServiceAvailable) {
			return nil, exc.ServiceUnavailable("go.micro.client", err.Error())
		}
		if errors.Is(err, micro.ErrSelectEndpointNotFound) {
			return nil, exc.NotFound("go.micro.client", err.Error())
		}
		return nil, err
	}
	return next, nil
}
