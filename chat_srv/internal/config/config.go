package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul    consul.Conf
	ChatModel ChatModelConfig `json:"chatModelConfig"`
	FileSrv   string          `json:"FileSrv"`
}

type ChatModelConfig struct {
	ChatModel string `json:"chatModel"`
	Embedder  string `json:"embedder"`
	APIKey    string `json:"apiKey"`
	Endpoint  string `json:"endpoint"`
	Index     string `json:"index"`
	Prefix    string `json:"prefix"`
	Dimension int64  `json:"dimension"`
	TopK      int    `json:"topK"`
}
