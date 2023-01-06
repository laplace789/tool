package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"tool/iotools"

	"github.com/sirupsen/logrus"
)

///文件日志Hook
type LogrusFileLoggerHook struct {
	maxLogSize int64
	logDir     string
	logPath    string
	level      logrus.Level
	logFile    *os.File
}

///工厂方法
func NewLogrusFileLoggerHook(logDir string, maxLogSize int64, level logrus.Level) (hook *LogrusFileLoggerHook, err error) {
	object := &LogrusFileLoggerHook{
		maxLogSize: maxLogSize,
		logDir:     logDir,
		level:      level,
	}
	return object, object.makeLogFile()
}

///创建日志文件
func (object *LogrusFileLoggerHook) makeLogFile() error {
	var err error
	if filepath.IsAbs(object.logDir) {
		object.logPath = object.logDir
	} else {
		dir, err := os.Getwd()
		if nil != err {
			panic(err)
		}
		object.logPath = filepath.Join(dir, object.logDir)
	}
	err = os.MkdirAll(object.logDir, 0777)
	if nil != err {
		panic(err)
	}

	now := time.Now()
	object.logPath += fmt.Sprintf("%s%s-%04d-%02d-%02dT%02d-%02d-%02d.%d.log",
		string(filepath.Separator), GetLevelString(object.level),
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), os.Getpid())
	if nil != object.logFile {
		err = object.logFile.Close()
		if nil != err {
			return err
		}
	}
	object.logFile, err = os.OpenFile(object.logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if nil != err {
		return err
	}
	return nil
}

///日志等级回调
func (object *LogrusFileLoggerHook) Levels() []logrus.Level {
	return []logrus.Level{object.level}
}

///激发
func (object *LogrusFileLoggerHook) Fire(entry *logrus.Entry) error {
	if size, err := iotools.FileSize(object.logPath); nil != err {
		return err
	} else if size > object.maxLogSize {
		if err = object.makeLogFile(); nil != err {
			return err
		}
	}
	content, err := entry.String()
	if nil != err {
		return err
	}
	_, err = object.logFile.Write([]byte(content))
	if nil != err {
		return err
	}
	return nil
}
