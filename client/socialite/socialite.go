package socialite

import (
	"context"
	"git.pmx.cn/hci/microservice-app/pkg/utils"
	"io"
	"time"

	"git.pmx.cn/hci/microservice-app/proto/socialite"
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

var (
	socialiteClient socialite.SocialiteServiceClient
	etcdInstancer   *etcdv3.Instancer
	prefix          = "/services/socialite"
)

func GetClient() socialite.SocialiteServiceClient {
	if socialiteClient == nil {
		panic("socialite client is not be initialized!")
	}
	return socialiteClient
}

func Init(conn *grpc.ClientConn, tracer stdOpentracing.Tracer, logger log.Logger) {
	socialiteClient = NewSocialiteClient(conn, tracer, logger)
}

func InitWithSD(sdClient etcdv3.Client, tracer stdOpentracing.Tracer, logger log.Logger) {
	etcdInstancer, _ = etcdv3.NewInstancer(sdClient, prefix, logger)
	socialiteClient = NewSocialiteClientWithSD(sdClient, tracer, logger)

}

type SocialiteClient struct {
	WxJsLoginEndpoint  endpoint.Endpoint
	WxJsConfigEndpoint endpoint.Endpoint
}

func (f *SocialiteClient) WxJsLogin(ctx context.Context, in *socialite.WxJsLoginRequest, opts ...grpc.CallOption) (*socialite.WxJsLoginResponse, error) {
	resp, err := f.WxJsLoginEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*socialite.WxJsLoginResponse), nil
}

func (f *SocialiteClient) WxJsConfig(ctx context.Context, in *socialite.NoParam, opts ...grpc.CallOption) (*socialite.WxJsConfigResponse, error) {
	resp, err := f.WxJsConfigEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*socialite.WxJsConfigResponse), nil
}

func NewSocialiteClient(conn *grpc.ClientConn, tracer stdOpentracing.Tracer, logger log.Logger) socialite.SocialiteServiceClient {
	limiter := ratelimit.NewDelayingLimiter(rate.NewLimiter(rate.Every(time.Second), 1000))

	var wxjsloginEndpoint endpoint.Endpoint
	{
		wxjsloginEndpoint = grpcTransport.NewClient(
			conn,
			"socialite.SocialiteService",
			"WxJsLogin",
			utils.DummyEncode,
			utils.DummyDecode,
			socialite.WxJsLoginResponse{},
			grpcTransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		wxjsloginEndpoint = opentracing.TraceClient(tracer, "WxJsLogin")(wxjsloginEndpoint)
		wxjsloginEndpoint = limiter(wxjsloginEndpoint)
		wxjsloginEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "WxJsLogin",
			Timeout: 5 * time.Second,
		}))(wxjsloginEndpoint)
	}

	var wxjsconfigEndpoint endpoint.Endpoint
	{
		wxjsconfigEndpoint = grpcTransport.NewClient(
			conn,
			"socialite.SocialiteService",
			"WxJsConfig",
			utils.DummyEncode,
			utils.DummyDecode,
			socialite.WxJsConfigResponse{},
			grpcTransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		wxjsconfigEndpoint = opentracing.TraceClient(tracer, "WxJsConfig")(wxjsconfigEndpoint)
		wxjsconfigEndpoint = limiter(wxjsconfigEndpoint)
		wxjsconfigEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "WxJsConfig",
			Timeout: 5 * time.Second,
		}))(wxjsconfigEndpoint)
	}

	return &SocialiteClient{
		WxJsLoginEndpoint:  wxjsloginEndpoint,
		WxJsConfigEndpoint: wxjsconfigEndpoint,
	}
}

func NewSocialiteClientWithSD(sdClient etcdv3.Client, tracer stdOpentracing.Tracer, logger log.Logger) socialite.SocialiteServiceClient {
	res := &SocialiteClient{}

	factory := SocialiteFactory(MakeWxJsLoginEndpoint, tracer, logger)
	endpointer := sd.NewEndpointer(etcdInstancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, time.Second, balancer)
	res.WxJsLoginEndpoint = retry

	factory = SocialiteFactory(MakeWxJsConfigEndpoint, tracer, logger)
	endpointer = sd.NewEndpointer(etcdInstancer, factory, logger)
	balancer = lb.NewRoundRobin(endpointer)
	retry = lb.Retry(3, time.Second, balancer)
	res.WxJsConfigEndpoint = retry

	return res
}

func MakeWxJsLoginEndpoint(f socialite.SocialiteServiceClient) endpoint.Endpoint {
	return f.(*SocialiteClient).WxJsLoginEndpoint
}

func MakeWxJsConfigEndpoint(f socialite.SocialiteServiceClient) endpoint.Endpoint {
	return f.(*SocialiteClient).WxJsConfigEndpoint
}

// Todo: use connect pool, and reference counting to one connection.
func SocialiteFactory(makeEndpoint func(f socialite.SocialiteServiceClient) endpoint.Endpoint, tracer stdOpentracing.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := NewSocialiteClient(conn, tracer, logger)
		endpoint := makeEndpoint(service)

		return endpoint, conn, nil
	}
}
