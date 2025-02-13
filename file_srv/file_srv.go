package main

import (
	"LearningGuide/file_srv/global"
	"LearningGuide/file_srv/handler"
	proto "LearningGuide/file_srv/proto/.FileProto"
	"github.com/OuterCyrex/Gorra/GorraSrv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 初始化日志
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)

	// 初始化设置
	global.ServerConfig, _ = GorraSrv.InitConfig("file_srv/config.yaml")

	// 初始化数据库
	global.DB, _ = GorraSrv.InitMysql(global.ServerConfig.Mysql)
	// _ = global.DB.AutoMigrate(&model.File{}, &model.Session{}, &model.Message{}, &model.Noun{}, &model.Exercise{}, &model.Summary{})

	// 初始化rpc服务
	server := grpc.NewServer()
	proto.RegisterFileServer(server, &handler.FileServer{})

	//启动服务
	err := GorraSrv.ServerRun(server, global.ServerConfig)
	if err != nil {
		zap.S().Errorf("启动文件服务器失败: %v", err)
	}
}
