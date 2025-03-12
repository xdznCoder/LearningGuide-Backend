package main

import (
	"LearningGuide/post_api/global"
	"LearningGuide/post_api/initialize"
	PostProto "LearningGuide/post_api/proto/.PostProto"
	UserProto "LearningGuide/post_api/proto/userProto"
	"LearningGuide/post_api/router"
	"flag"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

var configFile = flag.String("f", "post_api/config/config-debug.yaml", "the config file")

func main() {
	flag.Parse()

	// 初始化日志
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	// 初始化设置

	global.InitConfig(*configFile)

	// 初始化链路追踪
	tracer, closer := GorraAPI.InitTracer(global.ServerConfig.Name, global.ServerConfig.Jaeger.Host, global.ServerConfig.Jaeger.Port)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// 初始化rpc连接
	postConn, err := GorraAPI.GetSrvConnection(14, global.ServerConfig, global.ServerConfig.SrvList[0])
	userConn, err := GorraAPI.GetSrvConnection(14, global.ServerConfig, global.ServerConfig.SrvList[1])
	if err != nil {
		zap.S().Panicf("get connection error: %s", err.Error())
	}

	// 初始化Redis
	initialize.InitRedis()

	global.PostSrvClient = PostProto.NewPostClient(postConn)
	global.UserSrvClient = UserProto.NewUserClient(userConn)

	// 初始化路由
	r := GorraAPI.KeepAliveRouters("v1",
		router.PostRouter,
		router.LikeRouter,
		router.FavRouter,
		router.CommentRouter,
		router.NoticeRouter,
	)

	// 启动路由服务
	err = GorraAPI.RunRouter(r, global.ServerConfig)

	if err != nil {
		zap.S().Panicf("Run Router Failed: %v", err)
	}
}
