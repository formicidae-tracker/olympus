package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	olympuspb "github.com/formicidae-tracker/olympus/api"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
)

//go:generate go run generate_version.go
//go:generate go run generate_client_utils.go
//go:generate go fmt api/client_utils.go
//go:generate protoc --experimental_allow_proto3_optional  --go_out=api --go-grpc_out=api ./api/olympus_service.proto
//go:generate go run ./api/examples/generate.go

type Options struct {
	Version   bool   `long:"version" description:"print current version and exit"`
	Address   string `long:"http-listen" short:"l" description:"Address for the HTTP server" default:":3000"`
	RPC       int    `long:"rpc-listen" short:"r" description:"Port for the RPC Service" default:"3001"`
	AllowCORS string `long:"allow-cors" description:"allow cors from domain"`
}

type spaHandler struct {
	root       string
	index      string
	fileServer http.Handler
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.root, path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.root, h.index))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.fileServer.ServeHTTP(w, r)
}

func NewSpaHandler(root string, cacheTTL time.Duration) http.Handler {
	return spaHandler{
		root:       root,
		index:      "index.html",
		fileServer: CacheControl(cacheTTL)(http.FileServer(http.Dir(root))),
	}

}

func setAngularRoute(router *mux.Router) {
	router.PathPrefix("/").Handler(
		NewSpaHandler("./webapp/dist/olympus/browser", 7*24*time.Hour))
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

	o, err := NewOlympus()
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
