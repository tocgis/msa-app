package auth_test

import (
	client "git.pmx.cn/hci/microservice-app/client/auth"
	pbAuth "git.pmx.cn/hci/microservice-app/proto/auth"
	"git.pmx.cn/hci/microservice-app/srv/auth"
	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"context"
	"net"
	"testing"
	"time"
)

func runAuthServer(addr string) *grpc.Server {
	service := auth.NewAuthService()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	srv := auth.MakeGRPCServer(service, opentracing.NoopTracer{}, log.NewNopLogger())
	s := grpc.NewServer()
	pbAuth.RegisterAuthServiceServer(s, srv)

	go func() {
		s.Serve(ln)
	}()
	time.Sleep(time.Second)
	return s
}

func TestNewAuthClient(t *testing.T) {
	s := runAuthServer(":8001")
	defer s.GracefulStop()
	conn, err := grpc.Dial(":8001", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	service := client.NewAuthClient(conn, opentracing.NoopTracer{}, log.NewNopLogger())
	req := &pbAuth.LoginRequest{
		Name:   "xiaocheng",
		Code:  	"123456",
		Type: 	1,
	}
	_, err = service.Login(context.Background(), req)
	if err != nil {
		panic(err)
	}

	req2 := &pbAuth.AuthRequest{
		Token:   "xiaocheng",
	}
	_, err = service.Auth(context.Background(), req2)
	if err != nil {
		panic(err)
	}
	//if len(resp.Message()) <= 0 {
	//	panic(resp)
	//}
}
