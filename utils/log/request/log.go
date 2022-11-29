package loggerRequest

import (
	"dingtalk-push/conf"
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
)

var logger = logrus.New()
var logDir string
var logName = "request"
var logLevel = ""
var logLevels = map[string]logrus.Level{
	"PanicLevel": logrus.PanicLevel,
	"FatalLevel": logrus.FatalLevel,
	"ErrorLevel": logrus.ErrorLevel,
	"WarnLevel":  logrus.WarnLevel,
	"InfoLevel":  logrus.InfoLevel,
	"DebugLevel": logrus.DebugLevel,
}

func init() {
	logDir = path.Join(conf.ConfigYamlInstance.LogConfig.LogsDir, logName)
	logRotationTime := conf.ConfigYamlInstance.LogConfig.LogsRotationTime
	logRotationCount := conf.ConfigYamlInstance.LogConfig.LogsRotationCount
	SetLogFormatter(&logrus.JSONFormatter{})
	hook := newLfsHook(&logLevel, logRotationTime, uint(logRotationCount))
	logger.SetOutput(&NilWriter{})
	logger.AddHook(hook)
}

// 不将日志记录到控制台
type NilWriter struct {
}

func (w *NilWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

// 封装logrus.Fields
type Fields logrus.Fields

func SetLogLevel(level string) {
	logLevel = level
}
func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

// Debug
func Debug(args ...interface{}) {
	SetLogLevel("DebugLevel")
	entry := logger.WithFields(logrus.Fields{})
	//entry.Data["file"] = fileInfo(2)
	entry.Debug(args)
}

// 带有field的Debug
func DebugWithFields(l interface{}, f Fields) {
	entry := logger.WithFields(logrus.Fields(f))
	entry.Debug(l)
}

// Info
func Info(args ...interface{}) {
	SetLogLevel("InfoLevel")
	entry := logger.WithFields(logrus.Fields{})
	entry.Info(args...)
}

// 带有field的Info
func InfoWithFields(l interface{}, f Fields) {
	entry := logger.WithFields(logrus.Fields(f))
	entry.Info(l)
}

// Warn
func Warn(args ...interface{}) {
	SetLogLevel("WarnLevel")
	entry := logger.WithFields(logrus.Fields{})
	entry.Warn(args...)
}

// 带有Field的Warn
func WarnWithFields(l interface{}, f Fields) {
	entry := logger.WithFields(logrus.Fields(f))
	entry.Warn(l)
}

// ErrMsg
func Error(args ...interface{}) {
	SetLogLevel("ErrorLevel")
	entry := logger.WithFields(logrus.Fields{})
	entry.Error(args...)

}

// 带有Fields的Error
func ErrorWithFields(l interface{}, f Fields) {
	entry := logger.WithFields(logrus.Fields(f))
	entry.Error(l)
}

// Fatal
func Fatal(args ...interface{}) {
	SetLogLevel("FatalLevel")
	entry := logger.WithFields(logrus.Fields{})
	entry.Fatal(args...)
}

// 带有Field的Fatal
func FatalWithFields(l interface{}, f Fields) {
	entry := logger.WithFields(logrus.Fields(f))
	entry.Fatal(l)
}

// Panic
func Panic(args ...interface{}) {
	SetLogLevel("PanicLevel")
	entry := logger.WithFields(logrus.Fields{})
	entry.Panic(args...)
}

// 带有Field的Panic
func PanicWithFields(l interface{}, f Fields) {
	entry := logger.WithFields(logrus.Fields(f))
	entry.Panic(l)
}
func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
