package biz

import (
	"context"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"sync"
	"time"
)

// 全局变量
var (
	DispatchReqQueue  = DispatchReq{JobTaskQueueMap: make(map[int]*TaskQueue)}            // 调度请求队列 key为jobId
	ExecutionRetQueue = RetQueue{TodoCallbackRets: make([]model.HandleCallbackParam, 50)} // 执行结果队列 key为JobId
	RunningList       = Running{RunningContextMap: make(map[int]*RunningContext)}         // 运行队列
	TriggerLogIdSet   = LogIdSet{LogIdSet: make(map[int64]int64, 200)}                      // 触发ID，即logId集合，避免重复触发和重复回调
)

type DispatchReq struct {
	sync.Mutex
	JobTaskQueueMap map[int]*TaskQueue
}

type LogIdSet struct {
	sync.Mutex
	LogIdSet map[int64]int64
}

type Running struct {
	sync.Mutex
	RunningContextMap map[int]*RunningContext
}

type RunningContext struct {
	Ctx    context.Context    // 上下文
	Cancel context.CancelFunc // 取消函数
}

// 单个JobHandler的任务队列
type TaskQueue struct {
	sync.Mutex
	Running   bool                 // 是否运行中
	TodoTasks []model.TriggerParam // 待执行的任务
}

// 所有JobHandler的待回调的结果队列，不需要区分jobHandler，因为大家都是面向同一个调度中心
type RetQueue struct {
	sync.Mutex
	TodoCallbackRets []model.HandleCallbackParam
}

// 添加执行结果到队列
func AddExecutionRetToQueue(item model.HandleCallbackParam) {
	ExecutionRetQueue.Lock()
	var todoCallbackRets = ExecutionRetQueue.TodoCallbackRets
	todoCallbackRets = append(todoCallbackRets, item)
	ExecutionRetQueue.TodoCallbackRets = todoCallbackRets
	ExecutionRetQueue.Unlock()
}

// 从执行结果队列获取
func PopExecutionRetFromQueue() ([]model.HandleCallbackParam, bool) {
	ExecutionRetQueue.Lock()
	if len(ExecutionRetQueue.TodoCallbackRets) == 0 {
		ExecutionRetQueue.Unlock()
		return nil, false
	} else {
		params := ExecutionRetQueue.TodoCallbackRets[0:]
		ExecutionRetQueue.TodoCallbackRets = ExecutionRetQueue.TodoCallbackRets[:0]
		ExecutionRetQueue.Unlock()
		return params, true
	}
}

// AddLogIdToSet 加入
func AddLogIdToSet(logId int64) bool {
	TriggerLogIdSet.Lock()
	if _, ok := TriggerLogIdSet.LogIdSet[logId]; ok {
		// 已存在
		TriggerLogIdSet.Unlock()
		return false
	}
	TriggerLogIdSet.LogIdSet[logId] = time.Now().Unix()
	TriggerLogIdSet.Unlock()
	return true
}

// RemoveLogIdFromSet 移除
func RemoveLogIdFromSet(logId int64) bool {
	TriggerLogIdSet.Lock()
	if _, ok := TriggerLogIdSet.LogIdSet[logId]; ok {
		// 存在,移除
		delete(TriggerLogIdSet.LogIdSet, logId)
		TriggerLogIdSet.Unlock()
		return true
	}
	TriggerLogIdSet.Unlock()
	// 不存在
	return false
}
