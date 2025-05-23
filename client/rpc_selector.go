package client

import (
	"context"
	"errors"
	"github.com/lolizeppelin/micro"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/selector"
	"github.com/lolizeppelin/micro/tracing"
	"github.com/lolizeppelin/micro/utils"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// next endpoint删选
func (r *rpcClient) next(ctx context.Context, request micro.Request, opts CallOptions) (selector.Next, error) {

	endpoint := request.Endpoint()
	filters := opts.Filters
	version := request.Version()
	protocols := request.Protocols()

	var span oteltrace.Span
	tracer := tracing.GetTracer(CallScope, _version)
	ctx, span = tracer.Start(ctx, "node.selector",
		oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
		oteltrace.WithAttributes(
			attribute.String("name", r.opts.Selector.Name()),
			attribute.String("endpoint", request.Endpoint()),
			attribute.String("name", r.opts.Selector.Name()),
			attribute.String("version", request.Version().Version()),
		),
	)
	if opts.Node != "" {
		sid, _ := utils.FromBase62(opts.Node)
		span.AddEvent("selector", oteltrace.WithAttributes(
			attribute.String("node", opts.Node), attribute.Int("sid", sid),
		))
	}

	defer span.End()

	// 标准过滤器
	filters = utils.InsertSlice(opts.Filters,
		func(services []*micro.Service) ([]*micro.Service, error) {
			var matched []*micro.Service
			span.AddEvent("selector", oteltrace.WithAttributes(attribute.Int("services", len(services))))

			defer span.AddEvent("selector", oteltrace.WithAttributes(attribute.Int("matched", len(matched))))

			for _, s := range services {
				// 主版本不匹配
				if version != nil && version.Major != s.Version {
					continue
				}
				// 服务endpoint匹配
				ep, ok := s.Endpoints[endpoint]
				if !ok {
					return nil, micro.ErrSelectEndpointNotFound
				}

				if !micro.MatchCodec(protocols.Reqeust, ep.Metadata["req"]) ||
					!micro.MatchCodec(protocols.Response, ep.Metadata["res"]) {
					return nil, exc.BadRequest("micro.client.selector", "request or response type mismatch")
				}
				if ep.Internal && !opts.Internal { // 屏蔽内部rpc请求
					return nil, exc.Forbidden("micro.client.selector", "disabled request")
				}
				pk := request.PrimaryKey() == ""
				if (ep.PrimaryKey && !pk) || (pk && !ep.PrimaryKey) {
					return nil, exc.BadRequest("micro.client.selector", "request path param error")
				}
				//for _, ep := range s.Endpoints {
				//	if ep.Name == endpoint {
				//		// 校验请求与返回协议
				//		if !micro.MatchCodec(protocols.Reqeust, ep.Metadata["req"]) ||
				//			!micro.MatchCodec(protocols.Response, ep.Metadata["res"]) {
				//			return nil, exc.BadRequest("micro.client.selector", "request or response type mismatch")
				//		}
				//		if ep.Internal && !opts.Internal { // 屏蔽内部rpc请求
				//			return nil, exc.Forbidden("micro.client.selector", "disabled request")
				//		}
				//		pk := request.PrimaryKey() == ""
				//		if (ep.PrimaryKey && !pk) || (pk && !ep.PrimaryKey) {
				//			return nil, exc.BadRequest("micro.client.selector", "request path param error")
				//		}
				//		found = true
				//		break
				//	}
				//}
				//if !found {
				//	return nil, micro.ErrSelectEndpointNotFound
				//}
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
		span.RecordError(err)
		if errors.Is(err, micro.ErrSelectServiceNotFound) {
			return nil, exc.ServiceUnavailable("micro.client.selector", err.Error())
		}
		if errors.Is(err, micro.ErrNoneServiceAvailable) {
			return nil, exc.ServiceUnavailable("micro.client.selector", err.Error())
		}
		if errors.Is(err, micro.ErrSelectEndpointNotFound) {
			return nil, exc.NotFound("go.micro.client.selector", err.Error())
		}
		return nil, err
	}
	return next, nil
}
