package bootstrap

import (
	"fmt"
	"github.com/gin-generator/sugar/foundation"
	"google.golang.org/grpc"
	"net"
)

// RegisterGrpcService gRPC service registration function type
type RegisterGrpcService func(*grpc.Server)

// Grpc server struct
type Grpc struct {
	Server *grpc.Server
}

// newGrpc creates a new Grpc server instance
func newGrpc() *Grpc {
	return &Grpc{
		Server: grpc.NewServer(),
	}
}

// Run starts the gRPC server
func (g *Grpc) Run(app *foundation.Application) {
	cfg := app.Config

	name := cfg.App.Name
	host := cfg.App.Host
	port := cfg.App.Port

	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("Failed to listen: " + err.Error())
	}

	fmt.Printf("%s gRPC server start: %s...\n", name, address)
	if err = g.Server.Serve(listener); err != nil {
		panic("Failed to serve: " + err.Error())
	}
}
