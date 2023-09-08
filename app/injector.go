//go:build wireinject
// +build wireinject

package app

import (
	"database/sql"
	"net"
	"search-svc/infrastructure/configuration"
	"search-svc/infrastructure/environment"

	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	AppModule = wire.NewSet(
		configModuleSets,
		ProvideAppEnvironment,
		ProvideAppEnvConfig,
		ProvideLogger,
		ProvidePostgres,
		ProvideGRPC,
		ProvideListener,
	)
)

func InjectAppEnvironment() (environment.AppEnvironment, error) {
	panic(wire.Build(AppModule))
}

func InjectAppConfig() configuration.AppConfig {
	panic(wire.Build(AppModule))
}

func InjectLogger() (*logrus.Logger, error) {
	panic(wire.Build(AppModule))
}

func InjectPostgres() *sql.DB {
	panic(wire.Build(AppModule))
}

func InjectGRPC() *grpc.Server {
	panic(wire.Build(AppModule))
}

func InjectListener() net.Listener {
	panic(wire.Build(AppModule))
}
