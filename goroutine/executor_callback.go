package goroutine

import (
	"encoding/json"
	"github.com/PGshen/go-xxl-executor/biz"
	"log"
	"time"
)

// StartCallback 启动回调
func StartCallback() {
	for {
		if params, ok := biz.PopExecutionRetFromQueue(); ok {
			paramsStr, _ := json.Marshal(params)
			log.Println("callback: " + string(paramsStr))
			biz.Callback(params)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}