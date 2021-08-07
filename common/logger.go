package common

import (
	"github.com/PGshen/go-xxl-executor/common/log"
	"os"
	"strconv"
	"strings"
	"time"
)

var Log *log.Logger
var logPath string

func InitLogger(path string, env string) {
	logPath = path
	executorLogPath := "./executorlog"
	//os.RemoveAll(executorLogPath)
	_ = os.Mkdir(executorLogPath, 0777)

	var h log.Handler
	var err error
	if strings.EqualFold(env, "dev") {
		// 开发环境日志输出到控制台
		h, err = log.NewStreamHandler(os.Stdout)
	} else {
		fileName := executorLogPath + "/executor.log"
		h, err = log.NewTimeRotatingFileHandler(fileName, log.WhenDay, 1)
	}

	if err != nil {
		panic(err.Error())
	}
	Log = log.NewDefault(h)
}

func GetXxlLogger(logId int64) *log.Logger {
	curLogPath := logPath
	if !strings.HasSuffix(curLogPath, "/") {
		curLogPath += "/"
	}
	today := time.Now().Format("20060102")
	curLogPath = curLogPath + today
	if !IsDirExists(curLogPath) {
		// 目录不存在则创建
		_ = os.Mkdir(curLogPath, 0777)
	}
	fileName := curLogPath + "/" + strconv.FormatInt(logId, 10) + ".log"
	h, err := log.NewRotatingFileHandler(fileName, 20971520, 2)
	if err != nil {
		panic(err.Error())
	}
	logger := log.NewDefault(h)
	//logger.SetLevel(log.LevelTrace)
	return logger
}
