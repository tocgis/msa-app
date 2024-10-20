package user

import (
	"context"
	"git.pmx.cn/hci/microservice-app/pkg/utils"
	"io"
	"time"

	"git.pmx.cn/hci/microservice-app/proto/user"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/tracing/opentracing"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

var userCli user.UserServiceClient
var etcdInstancer *etcdv3.Instancer
var prefix = "/services/user"

func Init(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) {
	userCli = NewUserClient(conn, tracer, logger)
}

func InitWithSD(sdClient etcdv3.Client, tracer stdopentracing.Tracer, logger log.Logger) {
	etcdInstancer, _ = etcdv3.NewInstancer(sdClient, prefix, logger)
	userCli = NewUserClientWithSD(sdClient, tracer, logger)
}

func GetClient() user.UserServiceClient {
	if userCli == nil {
		panic("profile client is not be initialized!")
	}
	return userCli
}

type UserClient struct {
	GetProfileEndpoint endpoint.Endpoint
}

func (p *UserClient) GetProfile(ctx context.Context, in *user.GetProfileRequest, opts ...grpc.CallOption) (*user.GetProfileResponse, error) {
	resp, err := p.GetProfileEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*user.GetProfileResponse), nil
}

func NewUserClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) user.UserServiceClient {
	limiter := ratelimit.NewDelayingLimiter(rate.NewLimiter(rate.Every(time.Second), 1000))

	var getProfileEndpoint endpoint.Endpoint
	{
		getProfileEndpoint = grpctransport.NewClient(
			conn,
			"user.UserService",
			"GetProfile",
			utils.DummyEncode,
			utils.DummyDecode,
			user.GetProfileResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		getProfileEndpoint = opentracing.TraceClient(tracer, "GetProfile")(getProfileEndpoint)
		getProfileEndpoint = limiter(getProfileEndpoint)
		getProfileEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetProfile",
			Timeout: 5 * time.Second,
		}))(getProfileEndpoint)
	}

	return &UserClient{
		GetProfileEndpoint: getProfileEndpoint,
	}
}

func MakeGetProfileEndpoint(f user.UserServiceClient) endpoint.Endpoint {
	return f.(*UserClient).GetProfileEndpoint
}

func NewUserClientWithSD(sdClient etcdv3.Client, tracer stdopentracing.Tracer, logger log.Logger) user.UserServiceClient {
	res := &UserClient{}

	factory := ProfileFactory(MakeGetProfileEndpoint, tracer, logger)
	endpointer := sd.NewEndpointer(etcdInstancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, time.Second, balancer)
	res.GetProfileEndpoint = retry

	return res
}

// Todo: use connect pool, and reference counting to one connection.
func ProfileFactory(makeEndpoint func(f user.UserServiceClient) endpoint.Endpoint, tracer stdopentracing.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := NewUserClient(conn, tracer, logger)
		endpoint := makeEndpoint(service)

		return endpoint, conn, nil
	}
}
