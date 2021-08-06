package goroutine

import (
	"github.com/PGshen/go-xxl-executor/common"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func StartCleanLog(logPath string, logRetentionDays int) {
	for {
		if !strings.HasSuffix(logPath, "/") {
			logPath += "/"
		}
		// 扫描目录
		dirs, _ := ioutil.ReadDir(logPath)
		for _, dir := range dirs {
			// 超过指定时间了，删掉
			if int((time.Now().Unix() - dir.ModTime().Unix())/86400) >= logRetentionDays {
				common.Log.Info("remove: " + dir.Name())
				if dir.IsDir() {
					_ = os.RemoveAll(logPath + dir.Name())
				} else {
					_ = os.Remove(logPath + dir.Name())
				}
			}
		}
		time.Sleep(24 * time.Hour)
	}
}
