package socialite

import (
	"context"
	"git.pmx.cn/hci/microservice-app/proto/socialite"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
	stdOpentracing "github.com/opentracing/opentracing-go"
)

func MakeWxJsLoginEndpoint(s socialite.SocialiteServiceServer, tracer stdOpentracing.Tracer, logger log.Logger) endpoint.Endpoint {
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*socialite.WxJsLoginRequest)
		return s.WxJsLogin(ctx, req)
	}
	epduration := duration.With("method", "WxJsLogin")
	eplog := log.With(logger, "method", "WxJsLogin")
	ep = opentracing.TraceServer(tracer, "WxJsLogin")(ep)
	ep = EndpointInstrumentingMiddleware(epduration)(ep)
	ep = EndpointLoggingMiddleware(eplog)(ep)
	return ep
}

func MakeWxJsConfigEndpoint(s socialite.SocialiteServiceServer, tracer stdOpentracing.Tracer, logger log.Logger) endpoint.Endpoint {
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*socialite.NoParam)
		return s.WxJsConfig(ctx, req)
	}
	epduration := duration.With("method", "WxJsConfig")
	eplog := log.With(logger, "method", "WxJsConfig")
	ep = opentracing.TraceServer(tracer, "WxJsConfig")(ep)
	ep = EndpointInstrumentingMiddleware(epduration)(ep)
	ep = EndpointLoggingMiddleware(eplog)(ep)
	return ep
}




type grpcServer struct {
	wxjslogin grpcTransport.Handler
	wxjsconfig grpcTransport.Handler
}

func (g grpcServer) WxJsLogin(ctx context.Context, request *socialite.WxJsLoginRequest) (*socialite.WxJsLoginResponse, error) {
	_, rep, err := g.wxjslogin.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*socialite.WxJsLoginResponse), nil
}

func (g grpcServer) WxJsConfig(ctx context.Context, request *socialite.NoParam) (*socialite.WxJsConfigResponse, error) {
	_, rep, err := g.wxjsconfig.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*socialite.WxJsConfigResponse), nil
}

func MakeGRPCServer(s socialite.SocialiteServiceServer, tracer stdOpentracing.Tracer, logger log.Logger) socialite.SocialiteServiceServer {
	options := []grpcTransport.ServerOption{
		grpcTransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		wxjslogin: grpcTransport.NewServer(
			MakeWxJsLoginEndpoint(s, tracer, logger),
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			append(options, grpcTransport.ServerBefore(opentracing.GRPCToContext(tracer, "Login", logger)))...,
		),
		wxjsconfig: grpcTransport.NewServer(
			MakeWxJsConfigEndpoint(s, tracer, logger),
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			append(options, grpcTransport.ServerBefore(opentracing.GRPCToContext(tracer, "Login", logger)))...,
		),
	}
}