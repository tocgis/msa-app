package user_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	client "git.pmx.cn/hci/microservice-app/client/user"
	pbUser "git.pmx.cn/hci/microservice-app/proto/user"
	"git.pmx.cn/hci/microservice-app/srv/user"

	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func runProfileServer(addr string) *grpc.Server {
	service := user.NewUserService()
	ctx := context.Background()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	srv := user.MakeGRPCServer(ctx, service, opentracing.NoopTracer{}, log.NewNopLogger())
	s := grpc.NewServer()
	pbUser.RegisterUserServiceServer(s, srv)

	go func() {
		s.Serve(ln)
	}()
	time.Sleep(time.Second)
	return s
}

func TestNewProfileClient(t *testing.T) {
	s := runProfileServer(":8002")
	defer s.GracefulStop()
	conn, err := grpc.Dial(":8002", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	service := client.NewUserClient(conn, opentracing.NoopTracer{}, log.NewNopLogger())
	req := &pbUser.GetProfileRequest{
		UserId: 123,
	}
	resp, err := service.GetProfile(context.Background(), req)
	if err != nil {
		fmt.Println(resp, err)
	}
}
