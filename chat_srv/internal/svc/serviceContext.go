package svc

import (
	FileProto "LearningGuide/chat_srv/.FileProto"
	"LearningGuide/chat_srv/internal/config"
	"context"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config     config.Config
	Redis      *redis.Client
	RAG        *RAGEngine
	FileClient FileProto.FileClient
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

	conn, err := consulConn(c.Consul.Host, c.FileSrv, 14)
	if err != nil {
		panic(err)
	}

	fileClient := FileProto.NewFileClient(conn)

	return &ServiceContext{
		Config:     c,
		Redis:      r,
		RAG:        rag,
		FileClient: fileClient,
	}
}
