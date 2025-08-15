//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

//go:generate wire

import (
	"github.com/lyb88999/resume_helper/backend/services/ai-service/internal/biz"
	"github.com/lyb88999/resume_helper/backend/services/ai-service/internal/conf"
	"github.com/lyb88999/resume_helper/backend/services/ai-service/internal/data"
	"github.com/lyb88999/resume_helper/backend/services/ai-service/internal/server"
	"github.com/lyb88999/resume_helper/backend/services/ai-service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.AI, *conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
