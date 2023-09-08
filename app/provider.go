package app

import (
	"database/sql"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/xans-me/authopia/core/configuration"
	"github.com/xans-me/authopia/core/db"
	"github.com/xans-me/authopia/core/environment"
	"github.com/xans-me/authopia/core/logger"
	"google.golang.org/grpc"
)

// ProvidePostgres is function to init postgres connection
func ProvidePostgres(config *configuration.AppConfig) *sql.DB {
	return db.NewPostgres(config)
}

func ProvideListener(config *configuration.AppConfig) net.Listener {
	return configuration.ListenGRPC(config)
}

func ProvideGRPC() *grpc.Server {
	return configuration.NewGRPCServer()
}

// ProvideAppEnvironment is a function to provide app enviroment data
func ProvideAppEnvironment() (environment.AppEnvironment, error) {
	return environment.FromOsEnv()
}

// ProvideAppEnvConfig is a function to get AppConfig struct data
func ProvideAppEnvConfig(conf *configuration.AppConfig) configuration.AppConfig {
	return *conf
}

// ProvideLogger is a function to log http request and deployment to a file
func ProvideLogger(env environment.AppEnvironment) *logrus.Logger {
	logger := logger.New(env, logger.FileTemplate("authopia-app-%Y_%m_%d"))
	return logger
}
