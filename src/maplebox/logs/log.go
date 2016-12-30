package logs

import (
	"gopkg.in/Sirupsen/logrus.v0"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"maplebox/conf"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type AppLogger struct {
	*logrus.Entry
}

type MapleboxLogger struct {
	LogWriter []io.Writer
}

func (log *MapleboxLogger) Write(p []byte) (n int, err error) {
	for _, lg := range log.LogWriter {
		n, err = lg.Write(p)
	}
	return
}

var (
	AppLog   *AppLogger = newAppLog()
	RouteLog *AppLogger = newRouteLog()
)

func setDefaultLog(log *logrus.Logger, filePath string) {
	logConf := conf.AppConfig.Log

	switch strings.ToLower(logConf.Level) {
	case "debug":
		log.Level = logrus.DebugLevel
	case "info":
		log.Level = logrus.InfoLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	default:
		log.Level = logrus.DebugLevel
	}

	switch logConf.Formatter {
	case "json":
		log.Formatter = &logrus.JSONFormatter{}
	case "text":
		log.Formatter = &logrus.TextFormatter{}
	default:
		log.Formatter = &logrus.JSONFormatter{}
	}

	logList := []io.Writer{}
	logList = append(logList, os.Stdout)

	if logConf.WriteFile {
		logList = append(logList, newFileLog(filePath))
	}

	mapleboxLog := &MapleboxLogger{
		LogWriter: logList,
	}

	log.Out = mapleboxLog

	log.Hooks.Add(&ContextHook{})
}

func initAppLog() {
	logConf := conf.AppConfig.Log
	setDefaultLog(AppLog.Logger, logConf.AppLogFileFormat)
}

func initRouteLog() {
	logConf := conf.AppConfig.Log
	setDefaultLog(RouteLog.Logger, logConf.RouteLogFileFormat)
}

func newAppLog() *AppLogger {
	logger := logrus.New()
	return &AppLogger{
		logger.WithFields(
			logrus.Fields{
				"tag": "app",
			},
		),
	}
}

func newRouteLog() *AppLogger {
	logger := logrus.New()
	return &AppLogger{
		logger.WithFields(
			logrus.Fields{
				"tag": "route",
			},
		),
	}
}

func newFileLog(fileFormat string) *lumberjack.Logger {
	logConf := conf.AppConfig.Log
	flog := &lumberjack.Logger{
		Dir:        logConf.LogDir,
		NameFormat: fileFormat,
		LocalTime:  logConf.LocalTime,
		MaxAge:     logConf.MaxAge,
		MaxBackups: logConf.MaxBackups,
		MaxSize:    logConf.MaxSize,
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			flog.Rotate()
		}
	}()

	return flog
}

func InitLog() {
	initAppLog()
	initRouteLog()
}
