package biz

import (
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/common"
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
	if taskQueue, ok := DispatchReqQueue[jobId]; ok {
		// 当前JobHandler有任务正在运行中
		if taskQueue.Running || len(taskQueue.TodoTasks) > 0 {
			return NewReturnT(common.FailCode, "job goroutine is running.")
		} else {
			return NewReturnT(common.SuccessCode, "job goroutine is idle.")
		}
	} else {
		return NewReturnT(common.FailCode, "jobId[" + strconv.Itoa(jobId) + "] does not exists.")
	}
}

// Run 运行
func (e *ExecutorBiz) Run(param model.TriggerParam) ReturnT {
	// todo
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
