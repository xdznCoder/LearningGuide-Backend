package main

import (
	"LearningGuide/class_srv/gateway"
	"LearningGuide/class_srv/global"
	"LearningGuide/class_srv/handler"
	"LearningGuide/class_srv/initialize"
	proto "LearningGuide/class_srv/proto/.ClassProto"
	"LearningGuide/class_srv/utils"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitMysql()

	port, err := utils.GetFreePort()
	if err != nil {
		zap.S().Panicf("get free port error:%v", err)
	}

	server := grpc.NewServer()
	proto.RegisterClassServer(server, &handler.ClassServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Addr, port))
	if err != nil {
		zap.S().Panicf("failed to listen: %v", err)
	}

	zap.S().Infof("Server Runs On Port %d", port)

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//健康检查
	gateway.HealthCheck(fmt.Sprintf("%s:%d", global.ServerConfig.Addr, port), 15)

	err = server.Serve(lis)
	if err != nil {
		zap.S().Panicf("failed to serve: %v", err)
	}
}
