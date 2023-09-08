package configuration

import (
	"github.com/spf13/viper"
	"github.com/xans-me/authopia/core/environment"
)

// EnvironmentConfigBinderProperties struct
type EnvironmentConfigBinderProperties struct {
	FileName string
	Path     string
}

// EnvironmentConfigBinder struct
type EnvironmentConfigBinder struct {
	fileName string
	path     string
	config   AppConfig
}

// Bind Method of EnvironmentConfigBinder.
func (binder *EnvironmentConfigBinder) Bind() (err error) {
	env, err := environment.FromOsEnv()
	if err != nil {
		return err
	}
	viper.SetConfigName(binder.fileName + "_" + string(env))
	viper.AddConfigPath(binder.path)
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	if err = viper.ReadInConfig(); err != nil {
		return err
	}
	viper.SetDefault("http.address", ":8080")
	err = viper.Unmarshal(&binder.config)
	return err
}

// GetAppConfig Method of EnvironmentConfigBinder.
func (binder *EnvironmentConfigBinder) GetAppConfig() (*AppConfig, error) {
	return &binder.config, nil
}

// NewEnvironmentConfigBinder func
func NewEnvironmentConfigBinder(properties EnvironmentConfigBinderProperties) *EnvironmentConfigBinder {
	return &EnvironmentConfigBinder{fileName: properties.FileName, path: properties.Path}
}
