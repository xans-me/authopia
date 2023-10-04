package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/xans-me/authopia/core/environment"
)

// loggerOptions models struct
type loggerOptions struct {
	basePath         string
	releasePrefix string
	sandboxPrefix        string
	sandboxLevel         logrus.Level
	releaseLevel  logrus.Level
	fileTemplate     string
}

func (options *loggerOptions) getPrefix(env environment.AppEnvironment) string {
	if env.IsSandbox() {
		return options.sandboxPrefix
	} else {
		return options.releasePrefix
	}
}

func (options *loggerOptions) getLevel(env environment.AppEnvironment) logrus.Level {
	if env.IsSandbox() {
		return options.sandboxLevel
	} else {
		return options.releaseLevel
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
		releasePrefix: "release",
		sandboxPrefix:        "sandbox",
		fileTemplate:     "app_%Y-%m-%d-%H%M.log",
		sandboxLevel:         logrus.DebugLevel,
		releaseLevel:  logrus.InfoLevel,
	}
}
