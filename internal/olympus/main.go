package olympus

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/formicidae-tracker/olympus/pkg/api"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//go:generate go run generate_version.go

type Options struct {
	Version           bool `long:"version" description:"print current version and exit"`
	GenerateVAPIDKeys bool `long:"generate-vapid-keys" description:"generate and output on stdout a new pair of VAPID Keys"`
	GenerateSecret    bool `long:"generate-secret" description:"generate and output on stdout a secret for HMAC signature"`

	Address   string   `long:"http-listen" short:"l" description:"Address for the HTTP server" default:":3000"`
	RPC       int      `long:"rpc-listen" short:"r" description:"Port for the RPC Service" default:"3001"`
	AllowCORS []string `long:"allow-cors" description:"allow cors from domain"`
}

func setUpHttpServer(o *Olympus, opts Options) GracefulServer {
	router := mux.NewRouter()
	o.setRoutes(router)
	logger := logrus.WithField("group", "http")
	router.Use(RecoverWrap(logger))
	router.Use(HTTPLogWrap(logger))
	if len(opts.AllowCORS) > 0 {
		router.Use(EnableCORS(opts.AllowCORS))
	}
	httpServer := &http.Server{
		Addr:    opts.Address,
		Handler: router,
	}
	return NewGracefulServer(httpServer)
}

func setUpRpcServer(o *Olympus, opts Options) *grpc.Server {
	server := grpc.NewServer(api.DefaultServerOptions...)
	api.RegisterOlympusServer(server, (*OlympusGRPCWrapper)(o))
	return server
}

func outputNewVAPIDKeys() error {
	private, public, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return fmt.Errorf("could not generate VAPID Keys: %w", err)
	}

	_, err = fmt.Printf("OLYMPUS_VAPID_PRIVATE=%s\n", private)
	if err != nil {
		return err
	}
	_, err = fmt.Printf("OLYMPUS_VAPID_PUBLIC=%s\n", public)
	if err != nil {
		return err
	}

	return nil
}

func outputNewSecret() error {
	secret := make([]byte, 64)
	_, err := rand.Read(secret)
	if err != nil {
		return err
	}
	_, err = fmt.Printf("OLYMPUS_SECRET=%s\n", base64.URLEncoding.EncodeToString(secret))
	if err != nil {
		return err
	}
	return nil
}

func Execute() error {
	opts := Options{}

	if _, err := flags.Parse(&opts); err != nil {
		return nil
	}

	if opts.Version == true {
		fmt.Println(OLYMPUS_VERSION)
		return nil
	}

	if opts.GenerateVAPIDKeys == true {
		return outputNewVAPIDKeys()
	}

	if opts.GenerateSecret == true {
		return outputNewSecret()
	}

	o, err := NewOlympus()
	if err != nil {
		return err
	}

	httpServer := setUpHttpServer(o, opts)
	rpcServer := setUpRpcServer(o, opts)

	httpLog := logrus.WithFields(logrus.Fields{
		"group":   "http",
		"address": opts.Address,
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		httpLog.Infof("listening")
		err := httpServer.Run()
		if err != nil {
			httpLog.WithField("error", err).Errorf("unhandled error")
		}
		wg.Done()
	}()

	rpcLog := logrus.WithFields(logrus.Fields{
		"group": "rpc",
		"port":  opts.RPC,
	})

	wg.Add(1)
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.RPC))
		if err != nil {
			rpcLog.WithField("error", err).Errorf("could not listen")
			return
		}

		rpcLog.Infof("listening")
		err = rpcServer.Serve(l)
		if err != nil && err != grpc.ErrServerStopped {
			rpcLog.WithField("error", err).Errorf("unhandled error")
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
		httpLog.WithField("error", err).Errorf("stop error")
	}

	if err := o.Close(); err != nil {
		logrus.WithField("error", err).Errorf("olympus stop error")
	}

	wg.Wait()

	return nil
}
