//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/xans-me/authopia/core/configuration"
	"github.com/xans-me/authopia/core/environment"
	"github.com/xans-me/authopia/core/obs"

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
		ProvideKeycloakConfig,
		ProvideTracer,
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

func InjectProvider() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	shutdown, err := obs.InitProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()
}
