package main

import (
	"github.com/PGshen/go-xxl-executor/app"
	"github.com/PGshen/go-xxl-executor/executor"
	"github.com/PGshen/go-xxl-executor/handler"
	"github.com/PGshen/go-xxl-executor/server"
)

func main() {
	// 注册JobHandler
	_ = handler.AddJobHandler("test", &app.TestJobHandler{})
	xxlExecutor := executor.NewXxlJobExecutor()
	server.Start(xxlExecutor.GetIp(), xxlExecutor.GetPort()) // 启动http服务
	xxlExecutor.Start()                    // 启动执行器服务
	defer xxlExecutor.Destroy()
}
