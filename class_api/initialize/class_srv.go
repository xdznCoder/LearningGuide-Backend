package initialize

import (
	"LearningGuide/class_api/global"
	"LearningGuide/class_api/middlewares"
	proto "LearningGuide/class_api/proto/.ClassProto"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConnection(wait uint, policy string) {
	Conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=%ds",
			global.ServerConfig.Consul.Host,
			global.ServerConfig.Consul.Port,
			global.ServerConfig.ClassSrv.Name,
			wait,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy": "%s"}`, policy)),
		grpc.WithUnaryInterceptor(middlewares.GrpcTracerInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		zap.S().Panicf("Load Balance Init Failed: %v", err)
		return
	}

	global.ClassSrvClient = proto.NewClassClient(Conn)
}
