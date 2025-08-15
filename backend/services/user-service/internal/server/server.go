package server

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"

	"github.com/lyb88999/resume_helper/backend/shared/proto/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewRegistrar)

// NewRegistrar 创建服务注册器
func NewRegistrar(c *conf.Bootstrap) registry.Registrar {
	consulConfig := api.DefaultConfig()
	if c.Registry != nil && c.Registry.Consul != nil {
		consulConfig.Address = c.Registry.Consul.Address
		consulConfig.Scheme = c.Registry.Consul.Scheme
	} else {
		consulConfig.Address = "consul:8500" // 默认值
		consulConfig.Scheme = "http"
	}

	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		panic(err)
	}

	r := consul.New(consulClient)
	return r
}
