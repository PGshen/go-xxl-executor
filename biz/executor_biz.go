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
	if AddLogIdToSet(param.LogId) {	// 判断当前logId是否已经在集合里，若是表明已经触发过，不在重新触发
		AddDispatchReqToQueue(param)
		log.Println("add a task[jobId=" + strconv.Itoa(param.JobId) + "] to dispatchReqQueue")
		return NewReturnT(common.SuccessCode, "run success")
	} else {
		log.Println("logId[" + strconv.FormatInt(param.LogId, 10) + "] already in queue")
		return NewFailReturnT("logId[" + strconv.FormatInt(param.LogId, 10) + "] already in queue")
	}
}

// Kill 终止 这里传入的是jobId,目前只终止了当前正在运行的；需要终止所有正在队列里排队的任务吗？或者说传入logId是否更合适？？
func (e *ExecutorBiz) Kill(param model.KillParam) ReturnT {
	jobId := param.JobId
	if runningCtx, ok := PopRunningCtxFromList(jobId); ok {
		runningCtx.Cancel()	// 通过context取消
		log.Println("kill job manually! jobId = " + strconv.Itoa(jobId))
		return NewReturnT(common.SuccessCode, "kill success")
	} else {
		return NewFailReturnT("current Job[jobId = " + strconv.Itoa(jobId) + "] does not running...")
	}
}

// Log 查看日志
func (e *ExecutorBiz) Log(param model.LogParam) ReturnT {
	// todo
	return NewReturnT(common.SuccessCode, "log success")
}
