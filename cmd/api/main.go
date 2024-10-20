package main

import (
	"context"
	"flag"
	"git.pmx.cn/hci/microservice-app/client/assess"
	"git.pmx.cn/hci/microservice-app/client/socialite"
	"git.pmx.cn/hci/microservice-app/client/user"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"git.pmx.cn/hci/microservice-app/api"
	"git.pmx.cn/hci/microservice-app/client/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	stdOpentracing "github.com/opentracing/opentracing-go"
	zipkinOt "github.com/openzipkin-contrib/zipkin-go-opentracing"
	stdZipkin "github.com/openzipkin/zipkin-go"
)

/**
  这个应该没有什么用
*/
var (
	HttpAddr = flag.String("http.addr", ":8080", "HTTP server address")
	etcdAddr = flag.String("etcd.addr", "http://127.0.0.1:2379", "etcd registry address")
	//zipkinAddr = flag.String("zipkin.addr", "http://127.0.0.1:9411", "tracer server address")
	zipkinAddr = flag.String("zipkin.addr", "", "tracer server address")
)

func main() {
	flag.Parse()
	ctx := context.Background()
	// Logging domain.
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// Service discovery domain. In this example we use etcd.
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

	// Transport domain.
	//tracer := stdOpentracing.GlobalTracer() // nop by default
	zipkinTracer, _ := stdZipkin.NewTracer(nil, stdZipkin.WithNoopTracer(true))

	var tracer stdOpentracing.Tracer
	{
		switch {
		case *zipkinAddr != "":
			level.Info(logger).Log("tracer", "Zipkin", "type", "OpenTracing", "URL", *zipkinAddr)
			tracer = zipkinOt.Wrap(zipkinTracer)
			fallthrough

		default:
			tracer = stdOpentracing.GlobalTracer() // no-op
		}
	}

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

		logger.Log("addr", ":6060")
		http.ListenAndServe(":6060", m)
	}()

	auth.InitWithSD(sdClient, tracer, logger)
	assess.InitWithSD(sdClient, tracer, logger)
	user.InitWithSD(sdClient, tracer, logger)
	socialite.InitWithSD(sdClient, tracer, logger)

	//gin.SetMode(gin.DebugMode)
	router := gin.New()
	api.Register(router)

	server := &http.Server{Addr: *HttpAddr, Handler: router}
	if err = server.Shutdown(ctx); err != nil {
		panic(err)
	}
	//if err = gracehttp.Serve(server); err != nil {
	//	panic(err)
	//}
}
