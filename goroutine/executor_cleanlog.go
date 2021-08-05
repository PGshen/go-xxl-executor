package goroutine

import (
	"github.com/PGshen/go-xxl-executor/common"
	"io/ioutil"
	"os"
	"time"
)

func StartCleanLog() {
	for {
		// 扫面目录
		logPath := common.Config.XxlJob.Executor.LogPath
		logRetentionDays := common.Config.XxlJob.Executor.LogRetentionDays
		dirs, _ := ioutil.ReadDir(logPath)
		for _, dir := range dirs {
			// 超过指定时间了，删掉
			if int((time.Now().Unix() - dir.ModTime().Unix())/86400) > logRetentionDays {
				_ = os.RemoveAll(dir.Name())
			}
		}
		time.Sleep(24 * time.Hour)
	}
}
