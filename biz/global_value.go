package biz

import (
	"context"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"sync"
)

// 全局变量
var (
	DispatchReqQueue  = DispatchReq{JobTaskQueueMap: make(map[int]*TaskQueue)}    // 调度请求队列 key为jobId
	ExecutionRetQueue = ExecutionRet{JobRetQueueMap: make(map[int]*RetQueue)}     // 执行结果队列 key为JobId
	RunningList       = Running{RunningContextMap: make(map[int]*RunningContext)} // 运行队列
)

type DispatchReq struct {
	Mutex           sync.Mutex
	JobTaskQueueMap map[int]*TaskQueue
}

type ExecutionRet struct {
	Mutex          sync.Mutex
	JobRetQueueMap map[int]*RetQueue
}

type Running struct {
	Mutex             sync.Mutex
	RunningContextMap map[int]*RunningContext
}

type RunningContext struct {
	Ctx    context.Context    // 上下文
	Cancel context.CancelFunc // 取消函数
}

// 单个JobHandler的任务队列
type TaskQueue struct {
	Mutex     sync.Mutex
	Running   bool                 // 是否运行中
	TodoTasks []model.TriggerParam // 待执行的任务
}

// 单个JobHandler的待回调的结果队列
type RetQueue struct {
	mutex            sync.Mutex
	TobeCallbackRets []model.HandleCallbackParam
}
