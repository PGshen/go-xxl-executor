package main

import (
	"github.com/PGshen/go-xxl-executor/app"
	"github.com/PGshen/go-xxl-executor/executor"
	"github.com/PGshen/go-xxl-executor/handler"
	"github.com/PGshen/go-xxl-executor/server"
)

func main() {
	server.Start()	// 启动http服务
	// 注册JobHandler
	_ = handler.AddJobHandler("test", &app.TestJobHandler{})
	// 启动执行器服务
	xxlExecutor := executor.NewXxlJobExecutor()
	defer xxlExecutor.Destroy()
	xxlExecutor.Start()
}