package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type Registry struct {
	Host   string
	Port   int
	Client *api.Client
}

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

func NewRegistryClient(host string, port int) RegistryClient {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", host, port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Panicf("Register Failed: %v", err)
		return nil
	}
	return &Registry{
		Host:   host,
		Port:   port,
		Client: client,
	}
}

func (r *Registry) Register(address string, port int, name string, tags []string, id string) error {

	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "15s",
	}

	registration := &api.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Port:    port,
		Tags:    tags,
		Address: address,
		Check:   check,
	}

	err := r.Client.Agent().ServiceRegister(registration)
	return err
}

func (r *Registry) DeRegister(serviceId string) error {
	err := r.Client.Agent().ServiceDeregister(serviceId)
	return err
}
