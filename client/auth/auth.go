package auth

import (
	"context"
	"git.pmx.cn/hci/microservice-app/pkg/utils"
	"io"
	"time"

	// this app
	"git.pmx.cn/hci/microservice-app/proto/auth"
	// go-kit
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

var authClient auth.AuthServiceClient
var authInstancer *etcdv3.Instancer
var prefix = "/services/auth"

func Init(conn *grpc.ClientConn, tracer stdOpentracing.Tracer, logger log.Logger) {
	authClient = NewAuthClient(conn, tracer, logger)
}

func InitWithSD(sdClient etcdv3.Client, tracer stdOpentracing.Tracer, logger log.Logger) {
	authInstancer, _ = etcdv3.NewInstancer(sdClient, prefix, logger)
	authClient = NewAuthClientWithSD(sdClient, tracer, logger)

}

func GetClient() auth.AuthServiceClient {
	if authClient == nil {
		panic("auth client is not be initialized!")
	}
	return authClient
}

type AuthClient struct {
	LoginEndpoint   endpoint.Endpoint
	AuthEndpoint 	endpoint.Endpoint
}

func (f *AuthClient) Login(ctx context.Context, in *auth.LoginRequest, opts ...grpc.CallOption) (*auth.LoginResponse, error) {
	resp, err := f.LoginEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*auth.LoginResponse), nil
}

func (f *AuthClient) Auth(ctx context.Context, in *auth.AuthRequest, opts ...grpc.CallOption) (*auth.AuthResponse, error) {
	resp, err := f.AuthEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*auth.AuthResponse), nil
}

func NewAuthClient(conn *grpc.ClientConn, tracer stdOpentracing.Tracer, logger log.Logger) auth.AuthServiceClient {

	limiter := ratelimit.NewDelayingLimiter(rate.NewLimiter(rate.Every(time.Second), 1000))

	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = grpcTransport.NewClient(
			conn,
			"auth.AuthService",
			"Login",
			utils.DummyEncode,
			utils.DummyDecode,
			auth.LoginResponse{},
			grpcTransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		loginEndpoint = opentracing.TraceClient(tracer, "Login")(loginEndpoint)
		loginEndpoint = limiter(loginEndpoint)
		loginEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Login",
			Timeout: 5 * time.Second,
		}))(loginEndpoint)
	}

	var authEndpoint endpoint.Endpoint
	{
		authEndpoint = grpcTransport.NewClient(
			conn,
			"auth.AuthService",
			"Auth",
			utils.DummyEncode,
			utils.DummyDecode,
			auth.AuthResponse{},
			grpcTransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		authEndpoint = opentracing.TraceClient(tracer, "Auth")(authEndpoint)
		authEndpoint = limiter(authEndpoint)
		authEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Auth",
			Timeout: 5 * time.Second,
		}))(authEndpoint)
	}

	return &AuthClient{
		LoginEndpoint:   loginEndpoint,
		AuthEndpoint:  authEndpoint,
	}
}


func MakeLoginEndpoint(f auth.AuthServiceClient) endpoint.Endpoint {
	return f.(*AuthClient).LoginEndpoint
}

func MakeAuthEndpoint(f auth.AuthServiceClient) endpoint.Endpoint {
	return f.(*AuthClient).AuthEndpoint
}

func NewAuthClientWithSD(sdClient etcdv3.Client, tracer stdOpentracing.Tracer, logger log.Logger) auth.AuthServiceClient {
	res := &AuthClient{}

	factory := AuthFactory(MakeLoginEndpoint, tracer, logger)
	endpointer := sd.NewEndpointer(authInstancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, time.Second, balancer)
	res.LoginEndpoint = retry

	factory = AuthFactory(MakeAuthEndpoint, tracer, logger)
	endpointer = sd.NewEndpointer(authInstancer, factory, logger)
	balancer = lb.NewRoundRobin(endpointer)
	retry = lb.Retry(3, time.Second, balancer)
	res.AuthEndpoint = retry

	return res
}

// Todo: use connect pool, and reference counting to one connection.
func AuthFactory(makeEndpoint func(f auth.AuthServiceClient) endpoint.Endpoint, tracer stdOpentracing.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := NewAuthClient(conn, tracer, logger)
		endpoint := makeEndpoint(service)

		return endpoint, conn, nil
	}
}
