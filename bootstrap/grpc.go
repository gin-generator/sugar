package bootstrap

import (
	"fmt"
	"github.com/gin-generator/sugar/foundation"
	"google.golang.org/grpc"
	"net"
)

// RegisterGrpcService gRPC 服务注册函数类型
type RegisterGrpcService func(*grpc.Server)

// Grpc gRPC 服务器结构
type Grpc struct {
	Server *grpc.Server
}

// newGrpc 创建新的 gRPC 服务器实例
func newGrpc() *Grpc {
	return &Grpc{
		Server: grpc.NewServer(),
	}
}

// Run 实现 Server 接口
func (g *Grpc) Run(app *foundation.Application) {
	cfg, _ := app.GetConfig("app")
	appCfg := cfg.(map[string]interface{})

	name := appCfg["name"].(string)
	host := appCfg["host"].(string)
	port := appCfg["port"].(int)

	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("Failed to listen: " + err.Error())
	}

	fmt.Printf("%s gRPC server start: %s...\n", name, address)
	if err := g.Server.Serve(listener); err != nil {
		panic("Failed to serve: " + err.Error())
	}
}
