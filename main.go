package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Version             bool   `long:"version" description:"print current version and exit"`
	Address             string `long:"http-listen" short:"l" description:"Address for the HTTP server" default:":3000"`
	RPC                 int    `long:"rpc-listen" short:"r" description:"Port for the RPC Service" default:"3001"`
	StreamServerAddress string `long:"stream-server" description:"address of the stream server" default:"http://localhost/olympus"`
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
	httpServer := &http.Server{
		Addr:    opts.Address,
		Handler: router,
	}
	return NewGracefulServer(httpServer)
}

func setUpRpcServer(o *Olympus, opts Options) GracefulServer {
	rpcRouter := rpc.NewServer()
	rpcRouter.RegisterName("Olympus", (*OlympusRPCWrapper)(o))
	rpcRouter.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	rpcServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.RPC),
		Handler: rpcRouter,
	}
	return NewGracefulServer(rpcServer)
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

	o := NewOlympus(opts.StreamServerAddress)

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
		log.Printf("[rpc]: listening on :%d", opts.RPC)
		err := rpcServer.Run()
		if err != nil {
			log.Printf("[rpc]: unhandled error: %s", err)
		}
		wg.Done()
	}()

	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	if err := rpcServer.Close(); err != nil {
		log.Printf("[rpc]: stop error: %s", err)
	}

	if err := httpServer.Close(); err != nil {
		log.Printf("[http]: stop error: %s", err)
	}

	wg.Wait()

	return o.Close()
}

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("Unhandled error: %s", err)
	}
}
