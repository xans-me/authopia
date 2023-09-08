package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/xans-me/authopia/core/environment"
)

// loggerOptions models struct
type loggerOptions struct {
	basePath         string
	stagingPrefix    string
	productionPrefix string
	devPrefix        string
	devLevel         logrus.Level
	stagingLevel     logrus.Level
	productionLevel  logrus.Level
	fileTemplate     string
}

func (options *loggerOptions) getPrefix(env environment.AppEnvironment) string {
	if env.IsDev() {
		return options.devPrefix
	} else if env.IsStaging() {
		return options.stagingPrefix
	} else {
		return options.productionPrefix
	}
}

func (options *loggerOptions) getLevel(env environment.AppEnvironment) logrus.Level {
	if env.IsDev() {
		return options.devLevel
	} else if env.IsStaging() {
		return options.stagingLevel
	} else {
		return options.productionLevel
	}
}

// loggerOption models
type loggerOption func(*loggerOptions)

func FileTemplate(template string) loggerOption {
	return func(options *loggerOptions) {
		options.fileTemplate = template
	}
}
func (options *loggerOptions) apply(setters []loggerOption) {
	for _, setter := range setters {
		setter(options)
	}
}

func createDefaultOptions() *loggerOptions {
	return &loggerOptions{
		basePath:         "./logs",
		stagingPrefix:    "staging",
		productionPrefix: "prod",
		devPrefix:        "dev",
		fileTemplate:     "app_%Y-%m-%d-%H%M.log",
		devLevel:         logrus.DebugLevel,
		stagingLevel:     logrus.InfoLevel,
		productionLevel:  logrus.InfoLevel,
	}
}
