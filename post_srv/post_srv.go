package main

import (
	"LearningGuide/post_srv/global"
	"LearningGuide/post_srv/handler"
	proto "LearningGuide/post_srv/proto/.PostProto"
	"github.com/OuterCyrex/Gorra/GorraSrv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 初始化日志
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)

	// 初始化设置
	global.ServerConfig, _ = GorraSrv.InitConfig("post_srv/config.yaml")

	// 初始化数据库
	global.DB, _ = GorraSrv.InitMysql(global.ServerConfig.Mysql)
	// _ = global.DB.AutoMigrate(&model.Fav{}, &model.Like{}, &model.Comment{}, &model.Post{})

	// 初始化rpc服务
	server := grpc.NewServer()
	proto.RegisterPostServer(server, &handler.PostServer{})

	//启动服务
	err := GorraSrv.ServerRun(server, global.ServerConfig)
	if err != nil {
		zap.S().Errorf("启动文件服务器失败: %v", err)
	}
}
