package biz

import (
	"github.com/PGshen/go-xxl-executor/biz/model"
)

// 全局变量
var (
	DispatchReqQueue map[int]TaskQueue        // 调度请求队列 key为jobId
	ExecutionRetQueue map[int]RetQueue        // 执行结果队列 key为JobId
	JobHandlerRegistry map[string]interface{} // JobHandler

)

// 单个JobHandler的任务队列
type TaskQueue struct {
	Running bool                   // 是否运行中
	TodoTasks []model.TriggerParam // 待执行的任务
}

// 单个JobHandler的待回调的结果队列
type RetQueue struct {
	TobeCallbackRets []model.HandleCallbackParam
}