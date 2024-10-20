package assess

import (
	"context"
	"git.pmx.cn/hci/microservice-app/pkg/utils"
	"io"
	"time"

	"git.pmx.cn/hci/microservice-app/proto/assess"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/tracing/opentracing"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
	stdOpentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

var assessCli assess.AssessServiceClient
var etcdInstancer *etcdv3.Instancer
var prefix = "/services/assess"

func Init(conn *grpc.ClientConn, tracer stdOpentracing.Tracer, logger log.Logger) {
	assessCli = NewAssessClient(conn, tracer, logger)
}

func InitWithSD(sdClient etcdv3.Client, tracer stdOpentracing.Tracer, logger log.Logger) {
	etcdInstancer, _ = etcdv3.NewInstancer(sdClient, prefix, logger)
	assessCli = NewAssessClientWithSD(sdClient, tracer, logger)

}

func GetClient() assess.AssessServiceClient {
	if assessCli == nil {
		panic("assess client is not be initialized!")
	}
	return assessCli
}

type AssessClient struct {
	InitScoreEndpoint endpoint.Endpoint
	ScoreInfoEndpoint endpoint.Endpoint
	BasicSaveEndpoint endpoint.Endpoint
	EducationSaveEndpoint endpoint.Endpoint
	WorkSaveEndpoint endpoint.Endpoint
}

func (p *AssessClient) InitScore(ctx context.Context, in *assess.ScoreRequest, opts ...grpc.CallOption) (*assess.OkResponse, error) {
	resp, err := p.InitScoreEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*assess.OkResponse), nil
}

func (p *AssessClient) ScoreInfo(ctx context.Context, in *assess.ScoreRequest, opts ...grpc.CallOption) (*assess.ScoreResponse, error) {
	resp, err := p.ScoreInfoEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*assess.ScoreResponse), nil
}

func (p *AssessClient) BasicSave(ctx context.Context, in *assess.BasicInfoRequest, opts ...grpc.CallOption) (*assess.OkResponse, error) {
	resp, err := p.BasicSaveEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*assess.OkResponse), nil
}

func (p *AssessClient) EducationSave(ctx context.Context, in *assess.EducationRequest, opts ...grpc.CallOption) (*assess.OkResponse, error) {
	resp, err := p.EducationSaveEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*assess.OkResponse), nil
}

func (p *AssessClient) WorkSave(ctx context.Context, in *assess.WorkinfoRequest, opts ...grpc.CallOption) (*assess.OkResponse, error) {
	resp, err := p.WorkSaveEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*assess.OkResponse), nil
}

func NewAssessClient(conn *grpc.ClientConn, tracer stdOpentracing.Tracer, logger log.Logger) assess.AssessServiceClient {
	limiter := ratelimit.NewDelayingLimiter(rate.NewLimiter(rate.Every(time.Second), 1000))

	var initScoreEndpoint endpoint.Endpoint
	{
		initScoreEndpoint = grpcTransport.NewClient(
			conn,
			"assess.AssessService",
			"InitScore",
			utils.DummyEncode,
			utils.DummyDecode,
			assess.OkResponse{},
			grpcTransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		initScoreEndpoint = opentracing.TraceClient(tracer, "InitScore")(initScoreEndpoint)
		initScoreEndpoint = limiter(initScoreEndpoint)
		initScoreEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "InitScore",
			Timeout: 5 * time.Second,
		}))(initScoreEndpoint)
	}

	var scoreInfoEndpoint endpoint.Endpoint
	{
		scoreInfoEndpoint = grpcTransport.NewClient(
			conn,
			"assess.AssessService",
			"ScoreInfo",
			utils.DummyEncode,
			utils.DummyDecode,
			assess.ScoreResponse{},
			grpcTransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		scoreInfoEndpoint = opentracing.TraceClient(tracer, "ScoreInfo")(scoreInfoEndpoint)
		scoreInfoEndpoint = limiter(scoreInfoEndpoint)
		scoreInfoEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "ScoreInfo",
			Timeout: 5 * time.Second,
		}))(scoreInfoEndpoint)
	}

	var basicInfoEndpoint endpoint.Endpoint
	{
		basicInfoEndpoint = grpcTransport.NewClient(
			conn,
			"assess.AssessService",
			"BasicSave",
			utils.DummyEncode,
			utils.DummyDecode,
			assess.OkResponse{},
			grpcTransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		basicInfoEndpoint = opentracing.TraceClient(tracer, "BasicSave")(basicInfoEndpoint)
		basicInfoEndpoint = limiter(basicInfoEndpoint)
		basicInfoEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "BasicSave",
			Timeout: 5 * time.Second,
		}))(basicInfoEndpoint)
	}

	var educationSaveEndpoint endpoint.Endpoint
	{
		educationSaveEndpoint = grpcTransport.NewClient(
			conn,
			"assess.AssessService",
			"EducationSave",
			utils.DummyEncode,
			utils.DummyDecode,
			assess.OkResponse{},
			grpcTransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		educationSaveEndpoint = opentracing.TraceClient(tracer, "EducationSave")(educationSaveEndpoint)
		educationSaveEndpoint = limiter(educationSaveEndpoint)
		educationSaveEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "EducationSave",
			Timeout: 5 * time.Second,
		}))(educationSaveEndpoint)
	}


	var workSaveEndpoint endpoint.Endpoint
	{
		workSaveEndpoint = grpcTransport.NewClient(
			conn,
			"assess.AssessService",
			"WorkSave",
			utils.DummyEncode,
			utils.DummyDecode,
			assess.OkResponse{},
			grpcTransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		workSaveEndpoint = opentracing.TraceClient(tracer, "WorkSave")(workSaveEndpoint)
		workSaveEndpoint = limiter(workSaveEndpoint)
		workSaveEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "WorkSave",
			Timeout: 5 * time.Second,
		}))(workSaveEndpoint)
	}

	return &AssessClient{
		InitScoreEndpoint: initScoreEndpoint,
		ScoreInfoEndpoint: scoreInfoEndpoint,
		BasicSaveEndpoint: basicInfoEndpoint,
		EducationSaveEndpoint: educationSaveEndpoint,
		WorkSaveEndpoint: workSaveEndpoint,
	}
}

func MakeInitScoreEndpoint(f assess.AssessServiceClient) endpoint.Endpoint {
	return f.(*AssessClient).InitScoreEndpoint
}

func MakeScoreInfoEndpoint(f assess.AssessServiceClient) endpoint.Endpoint {
	return f.(*AssessClient).ScoreInfoEndpoint
}

func MakeBasicSaveEndpoint(f assess.AssessServiceClient) endpoint.Endpoint {
	return f.(*AssessClient).BasicSaveEndpoint
}

func MakeEducationSaveEndpoint(f assess.AssessServiceClient) endpoint.Endpoint {
	return f.(*AssessClient).EducationSaveEndpoint
}

func MakeWorkSaveEndpoint(f assess.AssessServiceClient) endpoint.Endpoint {
	return  f.(*AssessClient).WorkSaveEndpoint
}

func NewAssessClientWithSD(sdClient etcdv3.Client, tracer stdOpentracing.Tracer, logger log.Logger) assess.AssessServiceClient {
	res := &AssessClient{}

	factory := AssessFactory(MakeInitScoreEndpoint, tracer, logger)
	endpointer := sd.NewEndpointer(etcdInstancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, time.Second, balancer)
	res.InitScoreEndpoint = retry

	factory = AssessFactory(MakeScoreInfoEndpoint, tracer, logger)
	endpointer = sd.NewEndpointer(etcdInstancer, factory, logger)
	balancer = lb.NewRoundRobin(endpointer)
	retry = lb.Retry(3, time.Second, balancer)
	res.ScoreInfoEndpoint = retry

	factory = AssessFactory(MakeBasicSaveEndpoint, tracer, logger)
	endpointer = sd.NewEndpointer(etcdInstancer, factory, logger)
	balancer = lb.NewRoundRobin(endpointer)
	retry = lb.Retry(3, time.Second, balancer)
	res.BasicSaveEndpoint = retry

	factory = AssessFactory(MakeEducationSaveEndpoint, tracer, logger)
	endpointer = sd.NewEndpointer(etcdInstancer, factory, logger)
	balancer = lb.NewRoundRobin(endpointer)
	retry = lb.Retry(3, time.Second, balancer)
	res.EducationSaveEndpoint = retry

	factory = AssessFactory(MakeWorkSaveEndpoint, tracer, logger)
	endpointer = sd.NewEndpointer(etcdInstancer, factory, logger)
	balancer = lb.NewRoundRobin(endpointer)
	retry = lb.Retry(3, time.Second, balancer)
	res.WorkSaveEndpoint = retry


	return res
}

// Todo: use connect pool, and reference counting to one connection.
func AssessFactory(makeEndpoint func(f assess.AssessServiceClient) endpoint.Endpoint, tracer stdOpentracing.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := NewAssessClient(conn, tracer, logger)
		endpoint := makeEndpoint(service)

		return endpoint, conn, nil
	}
}
