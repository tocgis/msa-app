package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pbSocialite "git.pmx.cn/hci/microservice-app/proto/socialite"
	"git.pmx.cn/hci/microservice-app/srv/socialite"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"

	stdOpentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	var (
		addr      = flag.String("addr", ":8085", "the microservices grpc address")
		debugAddr = flag.String("debug.addr", ":6065", "the debug and metrics address")
		etcdAddr  = flag.String("etcd.addr", "http://127.0.0.1:2379", "etcd registry address")
	)
	flag.Parse()

	key := "/services/socialite"
	value := *addr
	ctx := context.Background()

	// logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
		logger = log.With(logger, "service", "socialite")
	}

	// Service registrar domain. In this example we use etcd.
	var sdClient etcdv3.Client
	var peers []string
	if len(*etcdAddr) > 0 {
		peers = strings.Split(*etcdAddr, ",")
	}
	sdClient, err := etcdv3.NewClient(ctx, peers, etcdv3.ClientOptions{})
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	// Build the registrar.
	registrar := etcdv3.NewRegistrar(sdClient, etcdv3.Service{
		Key:   key,
		Value: value,
	}, log.NewNopLogger())

	// Register our instance.
	registrar.Register()
	defer registrar.Deregister()

	tracer := stdOpentracing.GlobalTracer() // nop by default

	service := socialite.NewSocialiteService()

	errchan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		errchan <- fmt.Errorf("%f", <-c)
	}()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Log("err", err)
		return
	}
	srv := socialite.MakeGRPCServer(service, tracer, logger)
	s := grpc.NewServer()
	pbSocialite.RegisterSocialiteServiceServer(s, srv)

	go func() {
		//logger := log.NewContext(logger).With("transport", "gRPC")
		logger.Log("addr", *addr)
		errchan <- s.Serve(ln)
	}()

	// Debug listener.
	go func() {
		logger := log.With(logger, "transport", "debug")

		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		m.Handle("/metrics", promhttp.Handler())

		logger.Log("addr", *debugAddr)
		errchan <- http.ListenAndServe(*debugAddr, m)
	}()

	logger.Log("graceful shutdown...", <-errchan)
	s.GracefulStop()

}
