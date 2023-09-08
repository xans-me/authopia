//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/xans-me/authopia/core/configuration"
)

var (
	configModuleSets = wire.NewSet(wire.Value(configuration.EnvironmentConfigBinderProperties{
		FileName: "app-config",
		Path:     "./env",
	}),
		ProvideAppConfig)
)

func ProvideAppConfig(properties configuration.EnvironmentConfigBinderProperties) *configuration.AppConfig {
	environmentConfigBinder := configuration.NewEnvironmentConfigBinder(properties)
	environmentConfigBinder.Bind()
	config, err := environmentConfigBinder.GetAppConfig()
	if err != nil {
		panic(err)
	}
	return config
}

func InitializeAppConfig() (*configuration.AppConfig, error) {
	panic(wire.Build(configModuleSets))
}
