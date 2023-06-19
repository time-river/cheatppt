package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"cheatppt/config"
	"cheatppt/utils"
)

func Setup() bool {
	opts := config.LogOpts

	lvl := utils.Must(logrus.ParseLevel(opts.Level))
	log.SetLevel(lvl)

	var output io.Writer

	switch opts.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		output = utils.Must(os.OpenFile(opts.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
	}

	log.SetOutput(output)
	log.SetReportCaller(true)

	if opts.Format == "json" {
		log.SetFormatter(new(logrus.JSONFormatter))
	}

	return lvl > logrus.InfoLevel
}

func GetWriter() io.Writer {
	logger := log.New()
	return logger.WriterLevel(log.GetLevel())
}
