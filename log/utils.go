package logs

import (
	"github.com/sirupsen/logrus"
)

func GetLevelString(level logrus.Level) string {
	logTAG := "UNKNOWN"
	switch level {
	case logrus.PanicLevel:
		logTAG = "PANIC"
	case logrus.FatalLevel:
		logTAG = "FATAL"
	case logrus.ErrorLevel:
		logTAG = "ERROR"
	case logrus.WarnLevel:
		logTAG = "WARN"
	case logrus.InfoLevel:
		logTAG = "INFO"
	case logrus.DebugLevel:
		logTAG = "DEBUG"
	case logrus.TraceLevel:
		logTAG = "TRACE"
	}
	return logTAG
}
