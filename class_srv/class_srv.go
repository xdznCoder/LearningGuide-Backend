package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"net"

	"LearningGuide/class_srv/.ClassProto"
	"LearningGuide/class_srv/internal/config"
	"LearningGuide/class_srv/internal/server"
	"LearningGuide/class_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "class_srv/etc/class-debug.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	c.ListenOn = fmt.Sprintf("%s:%d", c.ListenOn, getFreePort())

	ctx := svc.NewServiceContext(c)

	_ = consul.RegisterService(c.ListenOn, c.Consul)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		__ClassProto.RegisterClassServer(grpcServer, server.NewClassServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

func getFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0
	}
	defer func() {
		_ = l.Close()
	}()

	return l.Addr().(*net.TCPAddr).Port
}
