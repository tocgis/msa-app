package user

import (
	"context"
	"fmt"
	"time"

	"git.pmx.cn/hci/microservice-app/proto/user"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"

	grpcTransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	duration metrics.Histogram = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "user",
		Name:      "request_duration_ns",
		Help:      "Request duration in nanoseconds.",
	}, []string{"method", "success"})
)

func EndpointInstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func MakeGetPrifileEndpoint(s user.UserServiceServer, tracer stdopentracing.Tracer, logger log.Logger) endpoint.Endpoint {
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*user.GetProfileRequest)
		return s.GetProfile(ctx, req)
	}
	epduration := duration.With("method", "GetProfile")
	eplog := log.With(logger, "method", "GetProfile")
	ep = opentracing.TraceServer(tracer, "GetProfile")(ep)
	ep = EndpointInstrumentingMiddleware(epduration)(ep)
	ep = EndpointLoggingMiddleware(eplog)(ep)
	return ep
}

// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(ctx context.Context, s user.UserServiceServer, tracer stdopentracing.Tracer, logger log.Logger) user.UserServiceServer {
	options := []grpcTransport.ServerOption{
		grpcTransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		getProfile: grpcTransport.NewServer(
			MakeGetPrifileEndpoint(s, tracer, logger),
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			append(options, grpcTransport.ServerBefore(opentracing.GRPCToContext(tracer, "GetProfile", logger)))...,
		),
	}
}

type grpcServer struct {
	getProfile grpcTransport.Handler
}

func (s *grpcServer) GetProfile(ctx context.Context, req *user.GetProfileRequest) (*user.GetProfileResponse, error) {
	_, rep, err := s.getProfile.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*user.GetProfileResponse), nil
}
