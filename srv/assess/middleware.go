package assess

import (
	"context"
	"fmt"
	"time"

	"git.pmx.cn/hci/microservice-app/proto/assess"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"

	grpcTransport "github.com/go-kit/kit/transport/grpc"
	stdOpentracing "github.com/opentracing/opentracing-go"
	stdPrometheus "github.com/prometheus/client_golang/prometheus"

)

var (
	duration metrics.Histogram = prometheus.NewSummaryFrom(stdPrometheus.SummaryOpts{
		Namespace: "assess",
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

func MakeInitScoreEndpoint(s assess.AssessServiceServer, tracer stdOpentracing.Tracer, logger log.Logger) endpoint.Endpoint {
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*assess.ScoreRequest)
		return s.InitScore(ctx, req)
	}
	epduration := duration.With("method", "InitScore")
	eplog := log.With(logger, "method", "InitScore")
	ep = opentracing.TraceServer(tracer, "InitScore")(ep)
	ep = EndpointInstrumentingMiddleware(epduration)(ep)
	ep = EndpointLoggingMiddleware(eplog)(ep)
	return ep
}

func MakeScoreInfoEndpoint(s assess.AssessServiceServer, tracer stdOpentracing.Tracer, logger log.Logger) endpoint.Endpoint {
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*assess.ScoreRequest)
		return s.ScoreInfo(ctx, req)
	}
	epduration := duration.With("method", "ScoreInfo")
	eplog := log.With(logger, "method", "ScoreInfo")
	ep = opentracing.TraceServer(tracer, "ScoreInfo")(ep)
	ep = EndpointInstrumentingMiddleware(epduration)(ep)
	ep = EndpointLoggingMiddleware(eplog)(ep)
	return ep
}

func MakeBasicSaveEndpoint(s assess.AssessServiceServer, tracer stdOpentracing.Tracer, logger log.Logger) endpoint.Endpoint {
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*assess.BasicInfoRequest)
		return s.BasicSave(ctx, req)
	}
	epduration := duration.With("method", "BasicSave")
	eplog := log.With(logger, "method", "BasicSave")
	ep = opentracing.TraceServer(tracer, "BasicSave")(ep)
	ep = EndpointInstrumentingMiddleware(epduration)(ep)
	ep = EndpointLoggingMiddleware(eplog)(ep)
	return ep
}

func MakeEducationSaveEndpoint(s assess.AssessServiceServer, tracer stdOpentracing.Tracer, logger log.Logger) endpoint.Endpoint {
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*assess.EducationRequest)
		return s.EducationSave(ctx, req)
	}
	epduration := duration.With("method", "EducationSave")
	eplog := log.With(logger, "method", "EducationSave")
	ep = opentracing.TraceServer(tracer, "EducationSave")(ep)
	ep = EndpointInstrumentingMiddleware(epduration)(ep)
	ep = EndpointLoggingMiddleware(eplog)(ep)
	return ep
}

func MakeWorkSaveEndpoint(s assess.AssessServiceServer, tracer stdOpentracing.Tracer, logger log.Logger) endpoint.Endpoint {
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*assess.WorkinfoRequest)
		return s.WorkSave(ctx, req)
	}
	epduration := duration.With("method", "WorkSave")
	eplog := log.With(logger, "method", "WorkSave")
	ep = opentracing.TraceServer(tracer, "WorkSave")(ep)
	ep = EndpointInstrumentingMiddleware(epduration)(ep)
	ep = EndpointLoggingMiddleware(eplog)(ep)
	return ep
}


// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(ctx context.Context, s assess.AssessServiceServer, tracer stdOpentracing.Tracer, logger log.Logger) assess.AssessServiceServer {
	options := []grpcTransport.ServerOption{
		grpcTransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		initScore: grpcTransport.NewServer(
			MakeInitScoreEndpoint(s, tracer, logger),
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			append(options, grpcTransport.ServerBefore(opentracing.GRPCToContext(tracer, "InitScore", logger)))...,
		),
		scoreInfo: grpcTransport.NewServer(
			MakeScoreInfoEndpoint(s, tracer, logger),
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			append(options, grpcTransport.ServerBefore(opentracing.GRPCToContext(tracer, "ScoreInfo", logger)))...,
			),
		basicSave: grpcTransport.NewServer(
			MakeBasicSaveEndpoint(s, tracer, logger),
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			append(options, grpcTransport.ServerBefore(opentracing.GRPCToContext(tracer, "BasicSave", logger)))...,
		),
		educationSave: grpcTransport.NewServer(
			MakeEducationSaveEndpoint(s, tracer, logger),
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			append(options, grpcTransport.ServerBefore(opentracing.GRPCToContext(tracer, "EducationSave", logger)))...,
		),
		workSave: grpcTransport.NewServer(
			MakeWorkSaveEndpoint(s, tracer, logger),
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			func(_ context.Context, request interface{}) (interface{}, error) { return request, nil },
			append(options, grpcTransport.ServerBefore(opentracing.GRPCToContext(tracer, "WorkSave", logger)))...,
		),

	}
}

type grpcServer struct {
	initScore 		grpcTransport.Handler
	scoreInfo 		grpcTransport.Handler
	basicSave 		grpcTransport.Handler
	educationSave 	grpcTransport.Handler
	workSave 		grpcTransport.Handler
}

func (s *grpcServer) ScoreInfo(ctx context.Context, request *assess.ScoreRequest) (*assess.ScoreResponse, error) {
	_, rep, err := s.scoreInfo.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*assess.ScoreResponse), nil
}

func (s *grpcServer) BasicSave(ctx context.Context, request *assess.BasicInfoRequest) (*assess.OkResponse, error) {
	_, rep, err := s.basicSave.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*assess.OkResponse), nil
}

func (s *grpcServer) EducationSave(ctx context.Context, request *assess.EducationRequest) (*assess.OkResponse, error) {
	_, rep, err := s.educationSave.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*assess.OkResponse), nil
}

func (s *grpcServer) WorkSave(ctx context.Context, request *assess.WorkinfoRequest) (*assess.OkResponse, error) {
	_, rep, err := s.workSave.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*assess.OkResponse), nil
}

func (s *grpcServer) InitScore(ctx context.Context, req *assess.ScoreRequest) (*assess.OkResponse, error) {
	_, rep, err := s.initScore.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*assess.OkResponse), nil
}
