package environment

import (
	"errors"
	"flag"
	"os"
	"strings"
)

type AppEnvironment string

func (env AppEnvironment) IsLocal() bool {
	return env == LOCAL
}

func (env AppEnvironment) IsDev() bool {
	return env == DEV
}

func (env AppEnvironment) IsStaging() bool {
	return env == STAGING
}

func (env AppEnvironment) IsProd() bool {
	return env == PROD
}

const (
	LOCAL   AppEnvironment = "local"
	DEV                    = "development"
	PROD                   = "production"
	STAGING                = "staging"
)

var (
	ErrEnvironmentNotFound = errors.New("environment not found")
	strEnv                 = ""
)

// GetEnvFlag func to get command line flag "env"
func getEnvFlag(name, value, usage string) string {
	var flagEnv string
	flag.StringVar(&flagEnv, name, value, usage)
	flag.Parse()
	return flagEnv
}

func FromOsEnv() (AppEnvironment, error) {
	if strEnv == "" {
		strEnv = strings.Trim(strings.ToLower(os.Getenv("APP_ENV")), " ")
		if flag.Lookup("env") == nil {
			strEnv = getEnvFlag("env", "", "env Mode project local, development, staging, production")
		}
	}

	switch strEnv {
	case "development":
		return DEV, nil
	case "production":
		return PROD, nil
	case "staging":
		return STAGING, nil
	}

	return LOCAL, nil
}
