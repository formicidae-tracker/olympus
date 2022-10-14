package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/formicidae-tracker/olympus/olympuspb"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
)

//go:generate go run generate_version.go
//go:generate go run generate_client_utils.go
//go:generate go fmt olympuspb/client_utils.go
//go:generate protoc --experimental_allow_proto3_optional  --go_out=olympuspb --go-grpc_out=olympuspb ./olympuspb/olympus_service.proto

type Options struct {
	Version   bool   `long:"version" description:"print current version and exit"`
	Address   string `long:"http-listen" short:"l" description:"Address for the HTTP server" default:":3000"`
	RPC       int    `long:"rpc-listen" short:"r" description:"Port for the RPC Service" default:"3001"`
	AllowCORS string `long:"allow-cors" description:"allow cors from domain"`
	SlackURL  string `long:"slack-url" description:"slack webhook URL, overiden by OLYMPUS_SLACK_URL"`
}

func setAngularRoute(router *mux.Router) {
	angularPaths := []string{
		"/host/{h}/zone/{z}",
		"/logs",
	}

	angularAssetsPath := "./webapp/dist/webapp"
	for _, p := range angularPaths {
		router.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			indexBytes, err := ioutil.ReadFile(filepath.Join(angularAssetsPath, "index.html"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Write(indexBytes)
		}).Methods("GET")
	}

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./webapp/dist/webapp")))
}

func setUpHttpServer(o *Olympus, opts Options) GracefulServer {
	router := mux.NewRouter()
	o.route(router)
	setAngularRoute(router)
	router.Use(HTTPLogWrap)
	router.Use(RecoverWrap)
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
	server := grpc.NewServer(olympuspb.DefaultServerOptions...)
	olympuspb.RegisterOlympusServer(server, (*OlympusGRPCWrapper)(o))
	return server
}

func Execute() error {
	opts := Options{}

	if _, err := flags.Parse(&opts); err != nil {
		return err
	}

	if opts.Version == true {
		fmt.Println(OLYMPUS_VERSION)
		return nil
	}

	slackURL := os.Getenv("OLYMPUS_SLACK_URL")
	if len(slackURL) > 0 {
		opts.SlackURL = slackURL
	}

	o, err := NewOlympus(opts.SlackURL)
	if err != nil {
		return err
	}

	httpServer := setUpHttpServer(o, opts)
	rpcServer := setUpRpcServer(o, opts)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Printf("[http]: listening on %s", opts.Address)
		err := httpServer.Run()
		if err != nil {
			log.Printf("[http]: unhandled error: %s", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.RPC))
		if err != nil {
			log.Printf("[rpc]: could not listen on :%d : %s", opts.RPC, err)
			return
		}

		log.Printf("[rpc]: listening on :%d", opts.RPC)
		err = rpcServer.Serve(l)
		if err != nil && err != grpc.ErrServerStopped {
			log.Printf("[rpc]: unhandled error: %s", err)
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
		log.Printf("[http]: stop error: %s", err)
	}

	if err := o.Close(); err != nil {
		log.Printf("[olympus]: close failure: %s", err)
	}

	wg.Wait()

	return nil
}

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("Unhandled error: %s", err)
	}
}
