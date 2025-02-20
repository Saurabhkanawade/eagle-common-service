package logwriter

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLoggers() *logrus.Logger {
	log := logrus.New()

	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter)
	log.Formatter.(*logrus.TextFormatter).ForceColors = true
	log.Formatter.(*logrus.TextFormatter).DisableTimestamp = false
	log.Level = logrus.TraceLevel
	log.Out = os.Stdout

	return log
}
