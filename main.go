package main

import (
	"github.com/PGshen/go-xxl-executor/app"
	"github.com/PGshen/go-xxl-executor/executor"
	"github.com/PGshen/go-xxl-executor/handler"
)

func main() {
	// 初始化配置，这里根据自己的应用，可以用配置文件加载
	xxlJobConfig := executor.XxlJobConfig{
		Env:              "dev",
		AdminAddress:     "http://127.0.0.1:8080/xxl-job-admin",
		AccessToken:      "",
		Appname:          "go-xxl-executor-sample",
		Address:          "",
		Ip:               "",
		Port:             9998,
		LogPath:          "/Users/shen/Me/Study/Operation/Go/go-xxl-executor/log",
		LogRetentionDays: 7,
		HttpTimeout:      5,
	}
	// 注册JobHandler
	_ = handler.AddJobHandler("test", &app.TestJobHandler{})
	xxlExecutor := executor.NewXxlJobExecutor(xxlJobConfig)
	xxlExecutor.Start() // 启动执行器服务
}
