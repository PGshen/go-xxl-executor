package common

import (
	log2 "github.com/PGshen/go-xxl-executor/common/log"
	"log"
	"os"
	"strconv"
)

func GetLogger(logId int64) *log2.Logger {
	logPath := Config.XxlJob.Executor.LogPath
	if !IsDirExists(logPath) {
		// 目录不存在则创建
		_ = os.Mkdir(logPath, 0777)
	}
	fileName := logPath + "/" + strconv.FormatInt(logId, 10) + ".log"
	h, err := log2.NewRotatingFileHandler(fileName, 20971520, 2)
	if err != nil {
		log.Fatal(err)
	}
	logger := log2.NewDefault(h)
	//logger.SetLevel(log2.LevelTrace)
	return logger
}