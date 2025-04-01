package svc

import (
	"LearningGuide/chat_srv/internal/config"
	"context"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Client
	RAG    *RAGEngine
}

func NewServiceContext(c config.Config) *ServiceContext {
	r := redis.NewClient(&redis.Options{
		Addr:          c.Redis.Host,
		Protocol:      2,
		UnstableResp3: true,
	})

	rag, err := initRAG(&c, r)
	if err != nil {
		panic(err)
	}
	if _, err := r.Do(context.Background(), "Ping").Result(); err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config: c,
		Redis:  r,
		RAG:    rag,
	}
}
