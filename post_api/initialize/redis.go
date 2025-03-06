package initialize

import (
	"LearningGuide/post_api/global"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.Redis.Host, global.ServerConfig.Redis.Port),
		DB:       global.ServerConfig.Redis.DB,
		Password: global.ServerConfig.Redis.Password,
	})
	if rdb.Ping(context.Background()).Err() != nil {
		zap.S().Panicw("redis init failed", "err", "redis init failed")
	}
	global.RDB = *rdb
}
