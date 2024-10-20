package assess_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	client "git.pmx.cn/hci/microservice-app/client/assess"
	pbTopic "git.pmx.cn/hci/microservice-app/proto/assess"
	"git.pmx.cn/hci/microservice-app/srv/assess"

	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

)

func runAssessServer(addr string) *grpc.Server {
	service := assess.NewAssessService()
	ctx := context.Background()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	srv := assess.MakeGRPCServer(ctx, service, opentracing.NoopTracer{}, log.NewNopLogger())
	s := grpc.NewServer()
	pbTopic.RegisterAssessServiceServer(s, srv)

	go func() {
		s.Serve(ln)
	}()
	time.Sleep(time.Second)
	return s
}

func TestNewAssessClient(t *testing.T) {
	s := runAssessServer(":8003")
	defer s.GracefulStop()
	conn, err := grpc.Dial(":8003", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	service := client.NewAssessClient(conn, opentracing.NoopTracer{}, log.NewNopLogger())
	req := &pbTopic.ScoreRequest{
		Phone: "13800138000",
	}
	resp, err := service.ScoreInfo(context.Background(), req)
	if err != nil {
		fmt.Println(resp, err)
	}
}
