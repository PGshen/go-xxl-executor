package biz

import (
	"context"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"log"
	"strconv"
	"sync"
	"time"
)

// 全局变量
var (
	DispatchReqQueue  = DispatchReq{JobTaskQueueMap: make(map[int]*TaskQueue)}            // 调度请求队列 key为jobId
	ExecutionRetQueue = RetQueue{TodoCallbackRets: []model.HandleCallbackParam{}} // 执行结果队列 key为JobId
	RunningList       = Running{RunningContextMap: make(map[int]*RunningContext)}         // 运行队列,是为了可以手动终止
	TriggerLogIdSet   = LogIdSet{LogIdSet: map[int64]int64{}}                      // 触发ID，即logId集合，避免重复触发和重复回调
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

// AddDispatchReqToQueue 添加调度请求到队列
func AddDispatchReqToQueue(param model.TriggerParam) {
	jobId := param.JobId
	// 这里需要注意加锁的位置，操作JobTaskQueueMap和todoTasks分别是不同的锁
	if taskQueue, ok := DispatchReqQueue.JobTaskQueueMap[jobId]; ok {
		// jobId任务已在队列
		taskQueue.Lock()
		var todoTasks = taskQueue.TodoTasks
		todoTasks = append(todoTasks, param)
		taskQueue.TodoTasks = todoTasks
		taskQueue.Unlock()
	} else {
		// jobId任务不在队列
		DispatchReqQueue.Lock()
		todoTasks := []model.TriggerParam{param}
		DispatchReqQueue.JobTaskQueueMap[jobId] = &TaskQueue{Running: false, TodoTasks: todoTasks}
		DispatchReqQueue.Unlock()
	}
}

// RemoveDispatchReqFromQueue 只有当任务队列为空时，才用map中移除
func RemoveDispatchReqFromQueue(jobId int) bool {
	DispatchReqQueue.Lock()
	if taskQueue, ok := DispatchReqQueue.JobTaskQueueMap[jobId]; ok {
		if len(taskQueue.TodoTasks) == 0 {
			delete(DispatchReqQueue.JobTaskQueueMap, jobId)
			DispatchReqQueue.Unlock()
			return true
		} else {
			// 任务队列不为空
			DispatchReqQueue.Unlock()
			return false
		}
	} else {
		DispatchReqQueue.Unlock()
		return true
	}
}

// 取一个未被协程领取的调度任务
func GetDispatchReqFromQueue() (jobId int, queue *TaskQueue){
	// 通过running队列判断
}

// AddExecutionRetToQueue 添加执行结果到队列
func AddExecutionRetToQueue(item model.HandleCallbackParam) {
	log.Println("AddExecutionRetToQueue")
	ExecutionRetQueue.Lock()
	var todoCallbackRets = ExecutionRetQueue.TodoCallbackRets
	todoCallbackRets = append(todoCallbackRets, item)
	ExecutionRetQueue.TodoCallbackRets = todoCallbackRets
	ExecutionRetQueue.Unlock()
}

// PopExecutionRetFromQueue 从执行结果队列获取
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

// AddRunningToList 添加运行中的任务
func AddRunningToList(jobId int, runningContext *RunningContext) bool {
	RunningList.Lock()
	if _, ok := RunningList.RunningContextMap[jobId]; ok {
		log.Println("jobId[" + strconv.Itoa(jobId) + "] already in list")
		RunningList.Unlock()
		return false
	} else {
		RunningList.RunningContextMap[jobId] = runningContext
		RunningList.Unlock()
		return true
	}
}

// PopRunningCtxFromList 弹出运行中的任务
func PopRunningCtxFromList(jobId int) (*RunningContext,bool) {
	RunningList.Lock()
	if runningCtx, ok := RunningList.RunningContextMap[jobId]; ok {
		delete(RunningList.RunningContextMap, jobId)	// 从运行中队列里移除
		RunningList.Unlock()
		return runningCtx, true
	} else {
		RunningList.Unlock()
		return nil, false
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
