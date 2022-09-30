package proto

import (
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var DefaultServerOptions []grpc.ServerOption

var DefaultDialOptions []grpc.DialOption

var DefaultCallOptions []grpc.CallOption

func init() {
	DefaultDialOptions = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}
