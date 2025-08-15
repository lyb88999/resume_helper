//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/lyb88999/resume_helper/backend/services/user-service/internal/biz"
	"github.com/lyb88999/resume_helper/backend/services/user-service/internal/conf"
	"github.com/lyb88999/resume_helper/backend/services/user-service/internal/data"
	"github.com/lyb88999/resume_helper/backend/services/user-service/internal/server"
	"github.com/lyb88999/resume_helper/backend/services/user-service/internal/service"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Auth, *conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
