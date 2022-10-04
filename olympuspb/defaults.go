package olympuspb

import (
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var DefaultServerOptions []grpc.ServerOption

var DefaultDialOptions []grpc.DialOption

var DefaultCallOptions []grpc.CallOption

func init() {
	DefaultServerOptions = []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime: 10 * time.Second,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 2 * time.Second,
		}),
	}

	DefaultDialOptions = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 10 * time.Second,
		}),
	}
}
