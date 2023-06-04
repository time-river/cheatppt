package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"cheatppt/config"
	"cheatppt/utils"
)

var logger *logrus.Logger

func Setup() bool {
	opts := config.LogOpts

	logger = logrus.New()

	lvl := utils.Must(logrus.ParseLevel(opts.Level))
	logger.SetLevel(lvl)

	var output io.Writer

	switch opts.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		output = utils.Must(os.OpenFile(opts.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	}

	logger.SetOutput(output)

	if opts.Format == "json" {
		logger.SetFormatter(new(logrus.JSONFormatter))
	}

	return lvl > logrus.InfoLevel
}

func GetWriter() io.Writer {
	return logger.WriterLevel(logger.GetLevel())
}

func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}
