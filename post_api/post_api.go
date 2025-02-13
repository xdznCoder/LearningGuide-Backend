package main

import (
	"LearningGuide/post_api/config"
	"LearningGuide/post_api/global"
	"LearningGuide/post_api/initialize"
	PostProto "LearningGuide/post_api/proto/.PostProto"
	UserProto "LearningGuide/post_api/proto/.UserProto"
	"LearningGuide/post_api/router"
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
	postConn, err := GorraAPI.GetSrvConnection(14, cf, global.ServerConfig.SrvList[0])
	if err != nil {
		zap.S().Panicf("get connection error: %s", err.Error())
	}
	userConn, err := GorraAPI.GetSrvConnection(14, cf, global.ServerConfig.SrvList[0])
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
	)

	// 启动路由服务
	err = GorraAPI.RunRouter(r, cf)

	if err != nil {
		zap.S().Panicf("Run Router Failed: %v", err)
	}
}
