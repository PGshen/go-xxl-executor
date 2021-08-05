package biz

import (
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/common"
	"log"
	"strconv"
	"strings"
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
	ok, msg := CheckJobHandlerIsIdle(jobId)
	var code int
	if ok {
		code = common.SuccessCode
	} else {
		code = common.FailCode
	}
	return NewReturnT(code, msg)
}

// Run 运行
func (e *ExecutorBiz) Run(param model.TriggerParam) ReturnT {
	if !CheckLogIdIsInSet(param.LogId) {	// 判断当前logId是否已经在集合里，若是表明已经触发过，不在重新触发
		// 阻塞策略处理
		blockStrategy := param.ExecutorBlockStrategy
		logId := strconv.FormatInt(param.LogId, 10)
		if strings.EqualFold(blockStrategy, BLOCK_STRATEGY_SERIAL_EXECUTION) {
			// 串行
			AddDispatchReqToQueue(param)
			AddLogIdToSet(param.LogId)
			log.Println("add a task[jobId=" + strconv.Itoa(param.JobId) + ", logId= " + logId + "] to dispatchReqQueue")
			return NewReturnT(common.SuccessCode, "run success")
		} else if strings.EqualFold(blockStrategy, BLOCK_STRATEGY_DISCARD_LATER) {
			// 丢弃后续,如果当前jobId有任务正在运行，那么丢弃当前调度
			idle, _ := CheckJobHandlerIsIdle(param.JobId)
			if idle {
				AddDispatchReqToQueue(param)
				AddLogIdToSet(param.LogId)
				log.Println("add a task[jobId=" + strconv.Itoa(param.JobId) + ", logId= " + logId + "] to dispatchReqQueue")
				return NewReturnT(common.SuccessCode, "run success")
			} else {
				log.Println("jobHandler is busy, discard current dispatch[logId="+ logId + "]")
				return NewFailReturnT("jobHandler is busy, discard current dispatch[logId="+ logId + "]")
			}
		} else if strings.EqualFold(blockStrategy, BLOCK_STRATEGY_COVER_EARLY) {
			// 覆盖之前，如果当前jobId有任务正在运行，那么终止它，把当前调度加入队列
			idle, _ := CheckJobHandlerIsIdle(param.JobId)
			if !idle {
				// 停止当前运行的
				jobId := param.JobId
				if runningCtx, ok := TakeRunningCtxFromList(jobId); ok {
					runningCtx.Cancel()	// 通过context取消
					// 清空之前在排队的logId
					taskQueue, yes := GetTaskQueue(jobId)
					if yes {
						for _, task := range taskQueue.TodoTasks {
							RemoveLogIdFromSet(task.LogId)
						}
					}
					RemoveDispatchReqFromQueue(jobId)
					log.Println("kill job due to block strategy! jobId = " + strconv.Itoa(jobId))
				}
			}
			AddDispatchReqToQueue(param)
			AddLogIdToSet(param.LogId)
			log.Println("add a task[jobId=" + strconv.Itoa(param.JobId) + ", logId= " + logId + "] to dispatchReqQueue")
			return NewReturnT(common.SuccessCode, "run success")
		} else {
			log.Println("unknown block strategy!")
			return NewFailReturnT("unknown block strategy!")
		}
	} else {
		log.Println("logId[" + strconv.FormatInt(param.LogId, 10) + "] already in queue")
		return NewFailReturnT("logId[" + strconv.FormatInt(param.LogId, 10) + "] already in queue")
	}
}

// Kill 终止 这里传入的是jobId,目前只终止了当前正在运行的；需要终止所有正在队列里排队的任务吗？或者说传入logId是否更合适？？
func (e *ExecutorBiz) Kill(param model.KillParam) ReturnT {
	jobId := param.JobId
	if runningCtx, ok := TakeRunningCtxFromList(jobId); ok {
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
