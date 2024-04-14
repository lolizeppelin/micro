package client

import (
	"errors"
	"fmt"
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

	// 标准过滤器
	filters = utils.InsertSlice(opts.Filters,
		func(services []*micro.Service) ([]*micro.Service, error) {
			var matched []*micro.Service
			for _, s := range services {
				// 服务版本匹配
				if version.Main() != s.Version {
					continue
				}
				// 服务endpoint匹配
				for _, ep := range s.Endpoints {
					if ep.Name == endpoint {
						// 返回与内部请求校验
						if opts.Internal || !ep.Internal {
							// 请求与返回协议校验
							if request.ContentType() != ep.Metadata["res"] ||
								request.Accept() != ep.Metadata["req"] {
								return nil, fmt.Errorf("bad request or response content type")
							}
							break
						}
						return nil, fmt.Errorf("disable proxy internal request")
					}
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
			return nil, exc.InternalServerError("go.micro.client", "service %s: %s", service, err.Error())
		}
		return nil, exc.InternalServerError("go.micro.client", "error selecting %s node: %s", service, err.Error())
	}
	return next, nil
}
