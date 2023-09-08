package configuration

// ConfigBinder interface
type ConfigBinder interface {
	Bind() (err error)
	GetAppConfig() (*AppConfig, error)
}
