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

func InitLogger(path string) {
	logPath = path
}

func init() {
	path := "./executorlog"
	os.RemoveAll(path)

	os.Mkdir(path, 0777)

	//fileName := path + "/executor.log"
	//h, err := log.NewTimeRotatingFileHandler(fileName, log.WhenMinute, 1)

	h, err := log.NewStreamHandler(os.Stdout)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	logger := log.NewDefault(h)
	//logger.SetLevel(log.LevelTrace)
	return logger
}