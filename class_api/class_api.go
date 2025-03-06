package main

import (
	"LearningGuide/class_api/gateway/consul"
	"LearningGuide/class_api/gateway/policy"
	"LearningGuide/class_api/global"
	"LearningGuide/class_api/initialize"
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var configFile = flag.String("f", "class_api/config/config-debug.yaml", "the config file")

func main() {
	flag.Parse()

	initialize.InitLogger()
	initialize.InitConfig(*configFile)
	tracerCloser := initialize.InitTracer()
	initialize.InitRedis()
	initialize.InitSrvConnection(14, policy.RoundRobin)
	err := initialize.InitTrans("zh")
	if err != nil {
		zap.S().Panicf("init trans failed: %v", zap.Error(err))
	}
	defer tracerCloser.Close()

	registryId := uuid.NewV4().String()

	R := initialize.Routers()

	registryClient := consul.NewRegistryClient(global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	err = registryClient.Register(
		global.ServerConfig.Address,
		global.ServerConfig.Port,
		global.ServerConfig.Name,
		global.ServerConfig.Tags,
		registryId,
	)
	if err != nil {
		zap.S().Panicf("Connect to Register Center Failed: %v", err)
	}

	zap.S().Debugf("server start... port: %d", global.ServerConfig.Port)

	go func() {
		if err := R.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panicf("server start failed : %v", err)
		}
	}()

	//终止时注销服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = registryClient.DeRegister(registryId)
	if err == nil {
		zap.S().Infof("API Gateway Deregistry Success")
	}
}
