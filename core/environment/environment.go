package environment

import (
	"errors"
	"flag"
	"os"
	"strings"
)

type AppEnvironment string

func (env AppEnvironment) IsSandbox() bool {
	return env == SANDBOX
}

func (env AppEnvironment) IsRelease() bool {
	return env == RELEASE
}

const (
	SANDBOX   AppEnvironment = "sandbox"
	RELEASE                    = "release"
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
			strEnv = getEnvFlag("env", "", "env Mode project sandbox / release")
		}
	}

	switch strEnv {
	case "release":
		return RELEASE, nil
}
	return SANDBOX, nil
}
