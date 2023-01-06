package logs

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
	"tool/iotools"

	"github.com/sirupsen/logrus"
)

var (
	logFileCleaner *LogFileCleaner
)

const (
	defaultLogDirName = "logs"   //日志目录
	maxLogFileSize    = 10 << 20 //最大日志文件大小10MB
	maxLogDirSize     = 1 << 30  //最大日志文件目录大小1GB
)

///日志清理器
type LogFileCleaner struct {
	logDir         string
	maxLogFileSize int64
	maxLogDirSize  int64
	wg             *sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	log            *logrus.Entry
}

///工厂方法
func NewLogFileCleaner(logDir string, maxLogFileSize, maxLogDirSize int64) *LogFileCleaner {
	object := &LogFileCleaner{
		logDir:         logDir,
		maxLogFileSize: maxLogFileSize,
		maxLogDirSize:  maxLogDirSize,
		wg:             &sync.WaitGroup{},
		log:            mlog.WithField("module", "LogFileCleaner"),
	}
	object.ctx, object.cancel = context.WithCancel(context.Background())
	return object
}

///名字
//func (object *LogFileCleaner) Name() string {
//	return "LogFileCleaner"
//}
//
/////关闭优先级
//func (object *LogFileCleaner) ShutdownPriority() int {
//	return shutdown.LogPriority
//}
//
/////关闭钩子
//func (object *LogFileCleaner) BeforeShutdown() {
//	object.close()
//}

/// 关闭
func (object *LogFileCleaner) close() {
	object.cancel()
	object.wg.Wait()
	mlog.Infof("log file cleaner closed")
}

/// 删除历史日志
func (object *LogFileCleaner) CheckLogFileSizeAndRemove() {
	object.wg.Add(1)
loop:
	for {
		select {
		case <-object.ctx.Done():
			break loop
		case <-time.After(5 * time.Minute):
			size, err := iotools.DirSize(object.logDir)
			if err != nil {
				object.log.Errorf("dir size: %s, error: %s\n", object.logDir, err)
				continue
			}
			if size < object.maxLogDirSize {
				continue
			}

			arr, err := ioutil.ReadDir(object.logDir)
			if nil != err {
				object.log.Errorf("read dir: %s, error: %s\n", object.logDir, err)
				continue
			}
			fiArr := FileInfoArray(arr)
			sort.Sort(fiArr)

			for _, fi := range fiArr {
				size, _ := iotools.DirSize(object.logDir)
				if size < object.maxLogDirSize {
					break
				}
				filePath, err := filepath.Abs(filepath.Join(object.logDir, fi.Name()))
				if nil != err {
					object.log.Errorf("filepath abs error: %s\n", err)
					continue
				}
				fi, err = os.Stat(filePath)
				if err != nil {
					continue
				}
				if fi.IsDir() {
					err = os.RemoveAll(filePath)
				} else {
					err = os.Remove(filePath)
				}
				object.log.Warn("remove:%s error:%s", filePath, err)
			}
		}
	}
	object.wg.Done()
}

/// 文件信息排序
type FileInfoArray []os.FileInfo

func (object FileInfoArray) Len() int { return len(object) }
func (object FileInfoArray) Less(i, j int) bool {
	return object[i].ModTime().Unix() < object[j].ModTime().Unix()
}
func (object FileInfoArray) Swap(i, j int) { object[i], object[j] = object[j], object[i] }
