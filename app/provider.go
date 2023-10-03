package app

import (
	"database/sql"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/xans-me/authopia/core/configuration"
	"github.com/xans-me/authopia/core/db"
	"github.com/xans-me/authopia/core/environment"
	"github.com/xans-me/authopia/core/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// ProvidePostgres is a function to init postgres connection
func ProvidePostgres(config *configuration.AppConfig) *sql.DB {
	return db.NewPostgres(config)
}

func ProvideListener(config *configuration.AppConfig) net.Listener {
	return configuration.ListenGRPC(config)
}

func ProvideGRPC() *grpc.Server {
	return configuration.NewGRPCServer()
}

// ProvideAppEnvironment is a function to provide app environment data
func ProvideAppEnvironment() (environment.AppEnvironment, error) {
	return environment.FromOsEnv()
}

// ProvideAppEnvConfig is a function to get AppConfig struct data
func ProvideAppEnvConfig(conf *configuration.AppConfig) configuration.AppConfig {
	return *conf
}

// ProvideLogger is a function to log http request and deployment to a file
func ProvideLogger(env environment.AppEnvironment) *logrus.Logger {
	l := logger.New(env, logger.FileTemplate("authopia-app-%Y_%m_%d"))
	return l
}

// ProvideKeycloakConfig is a function to provide keycloak environment
func ProvideKeycloakConfig(conf *configuration.AppConfig) configuration.KeyCloak {
	return conf.KeyCloak
}

func ProvideTracer() trace.Tracer {
	return otel.Tracer("authopia-tracer")
}
