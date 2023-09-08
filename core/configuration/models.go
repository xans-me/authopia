package configuration

// AppConfig struct model, the base configuration struct for this web service
type AppConfig struct {
	App         App
	SQLDatabase SQLDatabase
}

// App struct model, the configuration file related to basic webservice info such as name environtment
type App struct {
	Name        string
	Environment string
	Debug       bool
	Host        string
	Port        string
	Protocol    string
}

// SQLDatabase struct model, this is configuration for postgresql, storing relational data
type SQLDatabase struct {
	Name                  string
	User                  string
	Password              string
	Port                  string
	Connection            string
	Host                  string
	MaximumOpenConnection int
	MaximumIdleConnection int
}
