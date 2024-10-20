package auth

import (
	"context"
	"fmt"
	"time"

	"git.pmx.cn/hci/microservice-app/proto/auth"

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
		Namespace: "auth",
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

func MakeLoginEndpoint(s auth.AuthServiceServer, tracer stdopentracing.Tracer, logger log.Logger) endpoint.Endpoint {
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*auth.LoginRequest)
		return s.Login(ctx, req)
	}
	epduration := duration.With("method", "Login")
	eplog := log.With(logger, "method", "Login")
	ep = opentracing.TraceServer(tracer, "Login")(ep)
	ep = EndpointInstrumentingMiddleware(epduration)(ep)
	ep = EndpointLoggingMiddleware(eplog)(ep)
	return ep
}

func MakeAuthEndpoint(s auth.AuthServiceServer, tracer stdopentracing.Tracer, logger log.Logger) endpoint.Endpoint {
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*auth.AuthRequest)
		return s.Auth(ctx, req)
	}
	epduration := duration.With("method", "Auth")
	eplog := log.With(logger, "method", "Auth")
	ep = opentracing.TraceServer(tracer, "Auth")(ep)
	ep = EndpointInstrumentingMiddleware(epduration)(ep)
	ep = EndpointLoggingMiddleware(eplog)(ep)
	return ep
}

// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(s auth.AuthServiceServer, tracer stdopentracing.Tracer, logger log.Logger) auth.AuthServiceServer {
	options := []grpcTransport.ServerOption{
		grpcTransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		login: grpcTransport.NewServer(
			MakeLoginEndpoint(s, tracer, logger),
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			append(options, grpcTransport.ServerBefore(opentracing.GRPCToContext(tracer, "Login", logger)))...,
		),
		auth: grpcTransport.NewServer(
			MakeAuthEndpoint(s, tracer, logger),
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			append(options, grpcTransport.ServerBefore(opentracing.GRPCToContext(tracer, "Auth", logger)))...,
		),
	}
}

type grpcServer struct {
	login grpcTransport.Handler
	auth  grpcTransport.Handler
}

func (s *grpcServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	_, rep, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*auth.LoginResponse), nil
}

func (s *grpcServer) Auth(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	_, rep, err := s.auth.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*auth.AuthResponse), nil
}
