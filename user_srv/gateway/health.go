package gateway

import (
	"LearningGuide/user_srv/global"
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func HealthCheck(grpcAddr string, checkInterval uint) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Fatal("consul connect fail", zap.Error(err))
		return
	}

	check := &api.AgentServiceCheck{
		GRPC:                           grpcAddr,
		Timeout:                        "5s",
		Interval:                       fmt.Sprintf("%ds", checkInterval),
		DeregisterCriticalServiceAfter: "15s",
	}

	addr := strings.Split(grpcAddr, ":")
	port, _ := strconv.Atoi(addr[1])

	serviceUUID := uuid.NewV4().String()

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      serviceUUID,
		Port:    port,
		Tags:    global.ServerConfig.Tags,
		Address: global.ServerConfig.Addr,
		Check:   check,
	})

	if err != nil {
		zap.S().Panicf("Service Register Failed: %v", err)
	}

	//接收终止信号后退出
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		if err = client.Agent().ServiceDeregister(serviceUUID); err != nil {
			zap.S().Infof("Deregister Service %s Failed: %v", serviceUUID, err)
		}
		zap.S().Infof("Deregister Service %s Success", serviceUUID)

		os.Exit(200)
	}()
}
