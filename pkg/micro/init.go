package micro

import (
	"context"
	"git.pmx.cn/hci/microservice-app/client/assess"
	"git.pmx.cn/hci/microservice-app/client/auth"
	"git.pmx.cn/hci/microservice-app/client/socialite"
	"git.pmx.cn/hci/microservice-app/client/user"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd/etcdv3"

	stdOpentracing "github.com/opentracing/opentracing-go"
	zipkinOt "github.com/openzipkin-contrib/zipkin-go-opentracing"
	stdZipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

var (
	etcdAddr string
	zipkinAddr string
)

func init() {
}

func Run() {
	etcdAddr = viper.GetString("micro.etcd")
	zipkinAddr = viper.GetString("micro.zipkin")

	ctx := context.Background()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)


	// Service discovery domain. In this example we use etcd.
	var sdClient etcdv3.Client
	var peers []string
	if len(etcdAddr) > 0 {
		peers = strings.Split(etcdAddr, ",")
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
		case zipkinAddr != "":
			level.Info(logger).Log("tracer", "Zipkin", "type", "OpenTracing", "URL", zipkinAddr)
			tracer = zipkinOt.Wrap(zipkinTracer)
			fallthrough

		default:
			tracer = stdOpentracing.GlobalTracer() // no-op
		}
	}

	// INIT 服务发现
	auth.InitWithSD(sdClient, tracer, logger)
	assess.InitWithSD(sdClient, tracer, logger)
	user.InitWithSD(sdClient, tracer, logger)
	socialite.InitWithSD(sdClient, tracer, logger)

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

}