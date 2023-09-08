package configuration

// AppConfig struct model, the base configuration struct for this web service
type AppConfig struct {
	App         App
	SQLDatabase SQLDatabase
	KeyCloak    KeyCloak
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

// KeyCloak struct model, contains configuration for keycloak. Keycloak is an open source identity and access management
type KeyCloak struct {
	BaseURLAuth   string
	Realm         string
	AdminUsername string
	AdminPassword string
	ClientID      string
	ClientSecret  string
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
