package goroutine

import (
	"encoding/json"
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/common"
	"time"
)

// StartCallback 启动回调
func StartCallback() {
	common.Log.Info("start callback...")
	for {
		if params, ok := biz.PopExecutionRetFromQueue(); ok {
			paramsStr, _ := json.Marshal(params)
			common.Log.Info("callback: " + string(paramsStr))
			biz.Callback(params)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}