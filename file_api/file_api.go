package main

import (
	"LearningGuide/file_api/global"
	"LearningGuide/file_api/initialize"
	ChatProto "LearningGuide/file_api/proto/.ChatProto"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"LearningGuide/file_api/router"
	"flag"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

var configFile = flag.String("f", "file_api/config/config-debug.yaml", "the config file")

func main() {
	flag.Parse()

	// 初始化日志
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	// 初始化设置

	global.InitConfig(*configFile)

	initialize.InitClient(&global.ServerConfig.AliyunOss)

	// 初始化链路追踪
	tracer, closer := GorraAPI.InitTracer(global.ServerConfig.Name, global.ServerConfig.Jaeger.Host, global.ServerConfig.Jaeger.Port)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// 初始化rpc连接
	fileConn, err := GorraAPI.GetSrvConnection(14, global.ServerConfig, global.ServerConfig.SrvList[0])
	if err != nil {
		zap.S().Panicf("get connection error: %s", err.Error())
	}
	chatConn, err := GorraAPI.GetSrvConnection(14, global.ServerConfig, global.ServerConfig.SrvList[1])
	if err != nil {
		zap.S().Panicf("get connection error: %s", err.Error())
	}

	// 初始化Redis
	initialize.InitRedis()

	global.FileSrvClient = FileProto.NewFileClient(fileConn)
	global.ChatSrvClient = ChatProto.NewChatClient(chatConn)

	// 初始化路由
	r := GorraAPI.KeepAliveRouters("v1",
		router.FileRouter,
		router.MessageRouter,
		router.NounRouter,
		router.ExerciseRouter,
		router.SummaryRouter,
		router.MiscRouter,
	)

	// 启动路由服务
	err = GorraAPI.RunRouter(r, global.ServerConfig)

	if err != nil {
		zap.S().Panicf("Run Router Failed: %v", err)
	}
}
