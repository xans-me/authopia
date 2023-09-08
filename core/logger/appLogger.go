package logger

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/xans-me/authopia/core/environment"
)

// New constructor of appLogger
func New(env environment.AppEnvironment, options ...loggerOption) *logrus.Logger {
	loggerOptions := createDefaultOptions()
	logger := logrus.New()
	if options != nil {
		loggerOptions.apply(options)
	}

	logger.SetLevel(loggerOptions.getLevel(env))
	path := loggerOptions.basePath + "/[" + loggerOptions.getPrefix(env) + "]"
	writer, err := rotatelogs.New(path+loggerOptions.fileTemplate,
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)
	if err != nil {
		panic(err)
	}

	writerMap := lfshook.WriterMap{}

	for _, level := range logrus.AllLevels {
		writerMap[level] = writer
	}

	logger.Hooks.Add(lfshook.NewHook(
		writerMap,
		&logrus.TextFormatter{},
	))

	return logger
}
