package logs

import (
	"sync"
	"tool/iotools"

	"github.com/sirupsen/logrus"
)

var (
	mlog    *logrus.Logger
	logConf *LogConf
	logOnce sync.Once
)

func GetLogger() *logrus.Logger {
	if mlog == nil {
		InitializeLogger(LogConf{
			ENV:       "production",
			Path:      defaultLogDirName,
			LogLevel:  "debug",
			LogPretty: false,
		})
	}
	return mlog
}

type LogConf struct {
	ENV       string
	Path      string
	LogLevel  string
	LogPretty bool
}

func InitializeLogger(conf LogConf) {
	logOnce.Do(func() {
		logConf = &conf
		if logConf.Path == "" {
			logConf.Path = defaultLogDirName
		}
		iotools.CleanEmptyLogFile(logConf.Path)

		logLevel, err := logrus.ParseLevel(conf.LogLevel)
		if err != nil {
			logLevel = logrus.DebugLevel
		}

		//初始化日志对象
		mlog = logrus.New()
		mlog.SetLevel(logLevel)
		mlog.SetReportCaller(true)

		mlog.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: logConf.LogPretty,
		})

		//创建Hook
		for level := logrus.PanicLevel; level <= logLevel; level++ {
			hook, err := NewLogrusFileLoggerHook(logConf.Path, maxLogFileSize, level)
			if nil != err {
				panic(err)
			}
			mlog.AddHook(hook)
		}

		//enable log cleaner
		logFileCleaner = NewLogFileCleaner(logConf.Path, maxLogFileSize, maxLogDirSize)
		go logFileCleaner.CheckLogFileSizeAndRemove()
	})
}

func ReloadLevel(logLevel string) {
	if mlog == nil {
		return
	}

	if logConf.LogLevel == logLevel {
		return
	}

	mlog.WithFields(logrus.Fields{
		"conf":            logConf,
		"input_log_level": logLevel,
	}).Info("reload Level Before")

	ll, err := logrus.ParseLevel(logLevel)
	if err != nil {
		mlog.WithFields(logrus.Fields{
			"conf":            logConf,
			"input_log_level": logLevel,
		}).Error(err)
		ll = logrus.DebugLevel
	} else {
		logConf.LogLevel = logLevel
	}

	mlog.WithFields(logrus.Fields{
		"conf":            logConf,
		"input_log_level": logLevel,
	}).Info("reload Level After")

	mlog.SetLevel(ll)
}
