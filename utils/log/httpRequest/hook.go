package loggerHttpRequest

import (
	"dingtalk-push/conf"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

type Hook interface {
	Levels() []logrus.Level
	Fire(*logrus.Entry) error
}

type DefaultFieldHook struct {
}

func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	entry.Data["appName"] = conf.ConfigYamlInstance.AppConfig.AppName
	return nil
}

func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

var writers map[string]*rotatelogs.RotateLogs

func initWriters(rotationTime int, maxRemainCnt uint) {
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		logrus.Errorf("创建目录失败！ [logDir]=%s [error]=%v", logDir, err)
	}

	rotation := time.Duration(rotationTime)
	writers = make(map[string]*rotatelogs.RotateLogs)
	// 记录系统日志
	fullPath := path.Join(logDir, logName)
	logWriter, err := rotatelogs.New(
		fullPath+".%Y%m%d%H"+".log",
		// 为最新的日志建立软连接
		//rotatelogs.WithLinkName(fullPath),
		// 设置日志分割时间
		rotatelogs.WithRotationTime(time.Hour*rotation),
		// 最多保存个数
		rotatelogs.WithRotationCount(maxRemainCnt),
	)
	if err != nil {
		logrus.Errorf("日志初始化失败！[error]=%v", err)
	}
	writers["log"] = logWriter
}

func newLfsHook(logLevel *string, rotationTime int, maxRemainCnt uint) logrus.Hook {
	initWriters(rotationTime, maxRemainCnt)
	level, ok := logLevels[*logLevel]
	if ok {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writers["log"],
		logrus.InfoLevel:  writers["log"],
		logrus.WarnLevel:  writers["log"],
		logrus.ErrorLevel: writers["log"],
		logrus.FatalLevel: writers["log"],
		logrus.PanicLevel: writers["log"],
	}, &logrus.JSONFormatter{DisableTimestamp: false})

	return lfsHook
}
