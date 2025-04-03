package svc

import (
	"LearningGuide/chat_srv/internal/config"
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	arkEmbed "github.com/cloudwego/eino-ext/components/embedding/ark"
	redisIndexer "github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino-ext/components/model/ark"
	redisRetriever "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
	"net/http"
	net_url "net/url"
	"path/filepath"
)

type RAGEngine struct {
	Indexer   *redisIndexer.Indexer
	Retriever *redisRetriever.Retriever
	Loader    *url.Loader
	Splitter  document.Transformer
	Chatter   *ark.ChatModel
}

func initRAG(c *config.Config, client *redis.Client) (*RAGEngine, error) {
	ctx := context.Background()

	e, err := arkEmbed.NewEmbedder(ctx, &arkEmbed.EmbeddingConfig{
		APIKey: c.ChatModel.APIKey,
		Model:  c.ChatModel.Embedder,
	})
	if err != nil {
		return nil, err
	}

	i, err := redisIndexer.NewIndexer(ctx, &redisIndexer.IndexerConfig{
		Client:           client,
		KeyPrefix:        c.ChatModel.Prefix,
		DocumentToHashes: nil,
		BatchSize:        10,
		Embedding:        e,
	})
	if err != nil {
		return nil, err
	}

	r, err := redisRetriever.NewRetriever(ctx, &redisRetriever.RetrieverConfig{
		Client:            client,
		Index:             c.ChatModel.Index,
		VectorField:       "vector_content",
		DistanceThreshold: nil,
		Dialect:           2,
		ReturnFields:      []string{"vector_content", "content"},
		DocumentConverter: nil,
		TopK:              1,
		Embedding:         e,
	})
	if err != nil {
		return nil, err
	}

	l, err := url.NewLoader(ctx, &url.LoaderConfig{
		Parser:         &Parser{},
		Client:         &http.Client{},
		RequestBuilder: nil,
	})
	if err != nil {
		return nil, err
	}

	s, err := recursive.NewSplitter(ctx, &recursive.Config{
		ChunkSize:   4000,
		OverlapSize: 500,
		Separators:  []string{",", ".", "。"},
		KeepType:    recursive.KeepTypeEnd,
	})

	ch, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: c.ChatModel.APIKey,
		Model:  c.ChatModel.ChatModel,
	})

	return &RAGEngine{
		Indexer:   i,
		Retriever: r,
		Loader:    l,
		Splitter:  s,
		Chatter:   ch,
	}, nil
}

func (r *RAGEngine) InitVectorIndex(ctx context.Context, client *redis.Client, c *config.Config, courseID int32) error {
	index := fmt.Sprintf("%s-%d", c.ChatModel.Index, courseID)
	if _, err := client.Do(ctx, "FT.INFO", index).Result(); err == nil {
		return nil
	}

	createIndexArgs := []interface{}{
		"FT.CREATE", index,
		"ON", "HASH",
		"PREFIX", "1", c.ChatModel.Prefix,
		"SCHEMA",
		"content", "TEXT",
		"vector_content", "VECTOR", "FLAT",
		"6",
		"TYPE", "FLOAT32",
		"DIM", c.ChatModel.Dimension,
		"DISTANCE_METRIC", "COSINE",
	}

	if err := client.Do(ctx, createIndexArgs...).Err(); err != nil {
		return err
	}

	if _, err := client.Do(ctx, "FT.INFO", index).Result(); err != nil {
		return err
	}
	return nil
}

func (r *RAGEngine) URLToFileInfo(ctx context.Context, URL string) (string, error) {
	parsedURL, err := net_url.Parse(URL)
	if err != nil {
		return "", err
	}

	path := parsedURL.Path
	filename := filepath.Base(path)

	return net_url.QueryUnescape(filename)
}

func (r *RAGEngine) LoadURLFile(ctx context.Context, url string) ([]*schema.Document, error) {
	n, err := r.URLToFileInfo(ctx, url)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, ContextFileNameKey, n)
	return r.Loader.Load(ctx, document.Source{
		URI: url,
	})
}
