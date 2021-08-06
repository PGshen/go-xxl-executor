package goroutine

import (
	"io/ioutil"
	"log"
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
				log.Println("remove: " + dir.Name())
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
