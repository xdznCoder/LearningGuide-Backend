package initialize

import (
	"LearningGuide/user_api/global"
	"LearningGuide/user_api/middlewares"
	proto "LearningGuide/user_api/proto/userProto"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitSrvConnection 用于实现负载均衡
func InitSrvConnection(wait uint, policy string) {
	userConn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=%ds",
			global.ServerConfig.Consul.Host,
			global.ServerConfig.Consul.Port,
			global.ServerConfig.UserSrv.Name,
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

	global.UserSrvClient = proto.NewUserClient(userConn)
}
