package server

import (
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"

	"github.com/liyubo06/resumeOptim_claude/backend/services/file-service/internal/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewRegistrar)

// NewRegistrar 创建服务注册器
func NewRegistrar(c *conf.Bootstrap) registry.Registrar {
	// 创建consul配置
	config := api.DefaultConfig()

	// 如果有配置信息，则使用配置信息
	if c.Registry != nil && c.Registry.Consul != nil {
		config.Address = c.Registry.Consul.Address
		config.Scheme = c.Registry.Consul.Scheme
	}

	// 按照官方教程创建consul client
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	// 使用consul client创建registry，返回kratos标准接口
	reg := consul.New(client)
	return reg
}
