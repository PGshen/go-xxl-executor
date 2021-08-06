package goroutine

import (
	"context"
	"errors"
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/common"
	"github.com/PGshen/go-xxl-executor/common/log"
	"github.com/PGshen/go-xxl-executor/handler"
	"strconv"
	"strings"
	"time"
)

func StartWorker() {
	// 轮询DispatchReqQueue
	for {
		for {
			// 这里是以一个jobId为单位，一个jobId被一个协程领取
			jobId, taskQueue, ok := biz.GetDispatchReqFromQueue()
			if !ok {
				// 没有未被领取的DispatchReq了，退出内循环，外面循环有个休眠，避免过度消耗
				break
			}
			// 当前jobId没有goroutine跑，那就启动一个协程。再次判断队列里是否有任务
			taskQueue.Lock()
			if !taskQueue.Running && len(taskQueue.TodoTasks) > 0 {
				taskQueue.Running = true
				common.Log.Info("start a goroutine for jobId[" + strconv.Itoa(jobId) + "]")
				go doTask(jobId, taskQueue)
			}
			taskQueue.Unlock()
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

// 一个协程跑一个jobId对应的任务
func doTask(jobId int, taskQueue *biz.TaskQueue) {
	defer func() { taskQueue.Running = false }()
	// 串行，暂未实现阻塞处理策略
	for {
		if len(taskQueue.TodoTasks) == 0 {
			break	// 所有任务跑完，退出任务循环
		}
		taskQueue.Lock()
		task := taskQueue.TodoTasks[0]
		taskQueue.TodoTasks = taskQueue.TodoTasks[1:]
		taskQueue.Unlock()
		_ = trigger(task)
	}
	// 当前jobId对于的任务都跑完了,那么移除，但有从上面循环退出到现在之间又有新的任务加入，所以再此判断
	if biz.RemoveDispatchReqFromQueue(jobId) {
		common.Log.Info("JobTaskQueueMap remove jobId: " + strconv.Itoa(jobId))
		common.Log.Info("TodoTasks is empty, exists current goroutine...")
	} else {
		// 再次进入，是否真的需要？？？
		doTask(jobId, taskQueue)
	}
}

func trigger(task model.TriggerParam) error {
	logger := common.GetXxlLogger(task.LogId)
	defer logger.Close()	// todo 待确认是否合适
	jobId := task.JobId
	var ctx context.Context
	var cancel context.CancelFunc
	if task.ExecutorTimeout > 0 {
		// 有超时
		timeout := time.Second * time.Duration(task.ExecutorTimeout)
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		// 无超时
		ctx, cancel = context.WithCancel(context.Background())
	}
	// 将ctx, cancel加入到运行队列，方便手动kill
	var runningCtx = &biz.RunningContext{Ctx: ctx, Cancel: cancel}
	biz.AddRunningToList(jobId, runningCtx)
	go func() {
		ret, err := execTask(logger, cancel, task)
		if err != nil {
			common.Log.Info(err)
		}
		// 正常执行完成吗
		if biz.RemoveLogIdFromSet(task.LogId) {
			// 是的
			common.Log.Info("Task[" + strconv.FormatInt(task.LogId,10) + "]: complete normally")
			logger.Info("Task[" + strconv.FormatInt(task.LogId,10) + "]: complete normally")
			biz.AddExecutionRetToQueue(model.HandleCallbackParam{
				LogId:      task.LogId,
				LogDateTim: time.Now().Unix(),
				HandleCode: ret.Code,
				HandleMsg:  ret.Msg,
			})
		} else {
			common.Log.Info("Task[" + strconv.FormatInt(task.LogId,10) + "] has been terminated due to timeout or killed")
			logger.Warn("Task[" + strconv.FormatInt(task.LogId,10) + "] has been terminated due to timeout or killed")
		}
	}()
	// 这里会阻塞等待，直到ctx.Done()
	common.Log.Info("000999")
	select {
	case <-ctx.Done():
		common.Log.Info("ctx.Done()")
		err := ctx.Err()
		if err == nil || strings.Contains(err.Error(), "context canceled") {
			// 正常退出or手动取消
			if biz.RemoveLogIdFromSet(task.LogId) {
				common.Log.Info("Task[" + strconv.FormatInt(task.LogId,10) + "]: kill manually")
				logger.Warn("Task[" + strconv.FormatInt(task.LogId,10) + "]: kill manually")
				// 手动取消，因为如果是正常退出的话，logId已被移除
				biz.AddExecutionRetToQueue(model.HandleCallbackParam{
					LogId:      task.LogId,
					LogDateTim: time.Now().Unix(),
					HandleCode: common.FailCode,
					HandleMsg:  "kill manually",
				})
			}
			// 正常退出，之前已处理过，不必再其他操作
		} else if strings.Contains(err.Error(), "context deadline exceeded") {
			// 超时退出
			common.Log.Info("Task[" + strconv.FormatInt(task.LogId,10) + "]: timeout exit")
			logger.Warn("Task[" + strconv.FormatInt(task.LogId,10) + "]: timeout exit")
			biz.RemoveLogIdFromSet(task.LogId)
			biz.AddExecutionRetToQueue(model.HandleCallbackParam{
				LogId:      task.LogId,
				LogDateTim: time.Now().Unix(),
				HandleCode: common.FailCode,
				HandleMsg:  "timeout",
			})
		}
	}
	common.Log.Info("..--..")
	// 任务完成or超时or被杀，统一这里从队列里删除
	biz.PopRunningCtxFromList(jobId)
	return nil
}

// 跑任务
func execTask(logger *log.Logger, cancel context.CancelFunc, triggerParam model.TriggerParam) (biz.ReturnT, error) {
	// 找到相应的JobHandler
	executorHandler := triggerParam.ExecutorHandler
	jobHandler := handler.GetJobHandler(executorHandler)
	if jobHandler == nil {
		common.Log.Info("can not found the related executorHandler")
		logger.Error("can not found the related executorHandler")
		return biz.NewFailReturnT("jobHandler not exists"), errors.New("can not found the related executorHandler")
	}
	// 传递给任务的参数，其他参数不应给到任务
	param := handler.Param{
		JobParam:   triggerParam.ExecutorParams,
		ShardIndex: triggerParam.BroadcastIndex,
		ShardTotal: triggerParam.BroadcastTotal,
	}
	jobHandler.SetLogger(logger)	// jobHandler完成后需要移除不？？
	jobHandler.Init()
	ret := jobHandler.Execute(param)
	jobHandler.Destroy()
	cancel()	// 任务完成
	return ret, nil
}

