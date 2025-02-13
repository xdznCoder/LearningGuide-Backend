package main

import (
	"LearningGuide/file_api/config"
	"LearningGuide/file_api/global"
	"LearningGuide/file_api/initialize"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"LearningGuide/file_api/router"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	// 初始化设置
	c := config.MainConfig{}

	cf, err := GorraAPI.InitConfig("post_api/config/config.yaml", &c)
	if err != nil {
		zap.S().Panicf("init config error: %s", err.Error())
	}

	global.ServerConfig = cf.(config.MainConfig)

	// 初始化链路追踪
	tracer, closer := GorraAPI.InitTracer(global.ServerConfig.Name, global.ServerConfig.Jaeger.Host, global.ServerConfig.Jaeger.Port)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// 初始化rpc连接
	conn, err := GorraAPI.GetSrvConnection(14, cf, global.ServerConfig.SrvList[0])
	if err != nil {
		zap.S().Panicf("get connection error: %s", err.Error())
	}

	// 初始化Redis
	initialize.InitRedis()

	global.FileSrvClient = FileProto.NewFileClient(conn)

	// 初始化路由
	r := GorraAPI.KeepAliveRouters("v1",
		router.FileRouter,
		router.MessageRouter,
		router.NounRouter,
		router.ExerciseRouter,
		router.SummaryRouter,
	)

	// 启动路由服务
	err = GorraAPI.RunRouter(r, cf)

	if err != nil {
		zap.S().Panicf("Run Router Failed: %v", err)
	}
}
