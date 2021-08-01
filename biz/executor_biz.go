package biz

import (
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/common"
	"log"
	"strconv"
)

type ExecutorBiz struct {
}

// Beat 心跳检测
func (e *ExecutorBiz) Beat() ReturnT {
	return NewReturnT(common.SuccessCode, "Success")
}

// IdleBeat 空闲检测
func (e *ExecutorBiz) IdleBeat(param model.IdleBeatParam) ReturnT {
	jobId := param.JobId
	if taskQueue, ok := DispatchReqQueue.JobTaskQueueMap[jobId]; ok {
		// 当前JobHandler有任务正在运行中
		if taskQueue.Running || len(taskQueue.TodoTasks) > 0 {
			return NewReturnT(common.FailCode, "job goroutine is running.")
		} else {
			return NewReturnT(common.SuccessCode, "job goroutine is idle.")
		}
	} else {
		return NewReturnT(common.FailCode, "jobId["+strconv.Itoa(jobId)+"] does not exists.")
	}
}

// Run 运行
func (e *ExecutorBiz) Run(param model.TriggerParam) ReturnT {
	jobId := param.JobId
	// 这里需要注意加锁的位置，操作JobTaskQueueMap和todoTasks分别是不同的锁
	if taskQueue, ok := DispatchReqQueue.JobTaskQueueMap[jobId]; ok {
		// jobId任务已在队列
		taskQueue.Mutex.Lock()
		var todoTasks = taskQueue.TodoTasks
		todoTasks = append(todoTasks, param)
		taskQueue.TodoTasks = todoTasks
		taskQueue.Mutex.Unlock()
	} else {
		// jobId任务不在队列
		DispatchReqQueue.Mutex.Lock()
		todoTasks := []model.TriggerParam{param}
		DispatchReqQueue.JobTaskQueueMap[jobId] = &TaskQueue{Running: false, TodoTasks: todoTasks}
		DispatchReqQueue.Mutex.Unlock()
	}
	log.Println("add a task[jobId=" + strconv.Itoa(jobId) + "] to dispatchReqQueue")
	return NewReturnT(common.SuccessCode, "run success")
}

// Kill 终止
func (e *ExecutorBiz) Kill(param model.KillParam) ReturnT {
	// todo
	return NewReturnT(common.SuccessCode, "kill success")
}

// Log 查看日志
func (e *ExecutorBiz) Log(param model.LogParam) ReturnT {
	// todo
	return NewReturnT(common.SuccessCode, "log success")
}
