package client

import (
	"errors"
	"github.com/lolizeppelin/micro"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/selector"
	"github.com/lolizeppelin/micro/utils"
)

// next endpoint删选
func (r *rpcClient) next(request micro.Request, opts CallOptions) (selector.Next, error) {

	endpoint := request.Endpoint()
	filters := opts.Filters
	version := request.Version()
	protocols := request.Protocols()

	// 标准过滤器
	filters = utils.InsertSlice(opts.Filters,
		func(services []*micro.Service) ([]*micro.Service, error) {
			var matched []*micro.Service
			for _, s := range services {
				// 主版本不匹配
				if version != nil && version.Major != s.Version {
					continue
				}
				found := false
				// 服务endpoint匹配
				for _, ep := range s.Endpoints {
					if ep.Name == endpoint {
						// 校验请求与返回协议
						if !micro.MatchCodec(protocols.Reqeust, ep.Metadata["req"]) ||
							!micro.MatchCodec(protocols.Response, ep.Metadata["res"]) {
							return nil, exc.BadRequest("go.micro.client.selector", "request or response type mismatch")
						}
						if ep.Internal && !opts.Internal { // 屏蔽内部rpc请求
							return nil, exc.Forbidden("go.micro.client.selector", "disabled request")
						}
						pk := request.PrimaryKey() == ""
						if (ep.PrimaryKey && !pk) || (pk && !ep.PrimaryKey) {
							return nil, exc.BadRequest("go.micro.client.selector", "request path param error")
						}
						found = true
						break
					}
				}
				if !found {
					return nil, micro.ErrSelectEndpointNotFound
				}
				// 节点过滤
				if opts.Node != "" {
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
						} else { // 节点版本过滤
							if version != nil {
								if node.Max != nil && version.Compare(*node.Min) > 0 { // 超过最大兼容
									continue
								}
								if node.Min != nil && version.Compare(*node.Min) < 0 { // 低于最小兼容
									continue
								}
								if node.Max == nil && node.Min == nil && version.Compare(node.Version) != 0 { // 无兼容,不匹配
									continue
								}
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
