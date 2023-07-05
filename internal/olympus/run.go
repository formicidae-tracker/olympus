package olympus

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/formicidae-tracker/olympus/pkg/api"
	"github.com/formicidae-tracker/olympus/pkg/tm"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

//go:generate go run generate_version.go

func Execute() error {
	_, err := parser.Parse()
	return err
}

type RunCommand struct {
	Verbose []bool `long:"verbose" short:"v" description:"enables verbose logging, set multiple time to increase the level"`

	Address      string   `long:"http-listen" short:"l" description:"Address for the HTTP server" default:":3000"`
	RPC          int      `long:"rpc-listen" short:"r" description:"Port for the RPC Service" default:"3001"`
	AllowCORS    []string `long:"allow-cors" description:"allow cors from domain"`
	OtelEndpoint string   `long:"otel-exporter" description:"Open Telemetry exporter endpoint" env:"OLYMPUS_OTEL_ENDPOINT"`
}

func (c *RunCommand) Execute([]string) error {
	c.setLogger()
	defer tm.Shutdown(context.Background())

	o, err := NewOlympus()
	if err != nil {
		return err
	}

	httpServer := c.setUpHttpServer(o)
	rpcServer := c.setUpRpcServer(o)

	httpLog := tm.NewLogger("http").WithField("address", c.Address)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		httpLog.Infof("listening")
		err := httpServer.Run()
		if err != nil {
			httpLog.WithError(err).Errorf("unhandled error")
		}
		wg.Done()
	}()

	rpcLog := tm.NewLogger("gRPC").WithField("port", c.RPC)

	wg.Add(1)
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", c.RPC))
		if err != nil {
			rpcLog.WithError(err).Errorf("could not listen")
			return
		}

		rpcLog.Infof("listening")
		err = rpcServer.Serve(l)
		if err != nil && err != grpc.ErrServerStopped {
			rpcLog.WithError(err).Errorf("unhandled error")
		}
		wg.Done()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	wg.Add(1)
	go func() {
		rpcServer.GracefulStop()
		wg.Done()
	}()

	if err := httpServer.Close(); err != nil {
		httpLog.WithError(err).Errorf("stop error")
	}

	if err := o.Close(); err != nil {
		o.log.WithError(err).Errorf("stop error")
	}

	wg.Wait()

	return nil

}

func (c *RunCommand) setUpHttpServer(o *Olympus) GracefulServer {
	router := mux.NewRouter()
	o.setRoutes(router)
	logger := tm.NewLogger("http")
	router.Use(RecoverWrap(logger))
	if tm.Enabled() == true {
		router.Use(otelmux.Middleware("olympus-backend"))
	} else {
		router.Use(HTTPLogWrap(logger))
	}
	if len(c.AllowCORS) > 0 {
		router.Use(EnableCORS(c.AllowCORS))
	}
	httpServer := &http.Server{
		Addr:    c.Address,
		Handler: router,
	}
	return NewGracefulServer(httpServer)
}

func (c *RunCommand) setUpRpcServer(o *Olympus) *grpc.Server {
	options := append([]grpc.ServerOption{}, api.DefaultServerOptions...)
	if tm.Enabled() {
		options = append(options,
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		)
	}

	server := grpc.NewServer(options...)

	api.RegisterOlympusServer(server, (*OlympusGRPCWrapper)(o))
	return server
}

func (c *RunCommand) setLogger() {
	if len(c.OtelEndpoint) > 0 {
		tm.SetUpTelemetry(tm.OtelProviderArgs{
			CollectorURL:   c.OtelEndpoint,
			ServiceName:    "olympus",
			ServiceVersion: OLYMPUS_VERSION,
			Level:          tm.VerboseLevel(len(c.Verbose)),
		})
	} else {
		tm.SetUpLocal(tm.VerboseLevel(len(c.Verbose)))
	}
}

func init() {
	parser.AddCommand("run", "run olympus service", "runs olympus service", &RunCommand{})
}
