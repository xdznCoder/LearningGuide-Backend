package consul

import (
	"LearningGuide/class_api/global"
	"fmt"
	"github.com/hashicorp/consul/api"
)

func PullServiceByName(name string) (map[string]*api.AgentService, error) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, name))
	if err != nil {
		return nil, err
	}
	return data, nil
}
