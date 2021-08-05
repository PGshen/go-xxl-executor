package goroutine

import (
	"context"
	"errors"
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/common"
	"github.com/PGshen/go-xxl-executor/handler"
	"log"
	"strconv"
	"strings"
	"time"
)

func StartWorker() {
	// 轮询DispatchReqQueue
	for {
		// 这里的遍历是否安全？？待确认
		for jobId, taskQueue := range biz.DispatchReqQueue.JobTaskQueueMap {
			log.Println("check jobId: " + strconv.Itoa(jobId))
			// 当前jobId没有goroutine跑，那就启动一个
			if !taskQueue.Running && len(taskQueue.TodoTasks) > 0 {
				log.Println("start a goroutine for jobId[" + strconv.Itoa(jobId) + "]")
				go doTask(jobId, taskQueue)
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

// 一个协程跑一个jobId对应的任务
func doTask(jobId int, taskQueue *biz.TaskQueue) {
	taskQueue.Running = true
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
	// 当前jobId对于的任务都跑完了,那么移除,这里可能有bug,会不会影响上面的遍历？？
	if biz.RemoveDispatchReqFromQueue(jobId) {
		log.Println("JobTaskQueueMap remove jobId: " + strconv.Itoa(jobId))
		log.Println("TodoTasks is empty, exists current goroutine...")
	} else {
		// 再次进入，是否真的需要？？？
		doTask(jobId, taskQueue)
	}
}

func trigger(task model.TriggerParam) error {
	jobId := task.JobId
	// todo使用waitGroup阻塞在这儿
	// todo任务参数替换为Context
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
	var runningCtx = &biz.RunningContext{Ctx: ctx, Cancel: cancel}
	biz.AddRunningToList(jobId, runningCtx)
	go func() {
		ret, err := execTask(cancel, task)
		if err != nil {
			log.Println(err)
		}
		// 正常执行完成吗
		if biz.RemoveLogIdFromSet(task.LogId) {
			// 是的
			log.Println("Task[" + strconv.FormatInt(task.LogId,10) + "]: complete normally")
			biz.AddExecutionRetToQueue(model.HandleCallbackParam{
				LogId:      task.LogId,
				LogDateTim: time.Now().Unix(),
				HandleCode: ret.Code,
				HandleMsg:  ret.Msg,
			})
		} else {
			log.Println("Task[" + strconv.FormatInt(task.LogId,10) + "] has been terminated due to timeout or killed")
		}
	}()
	// 这里会阻塞等待
	log.Println("000999")
	select {
	case <-ctx.Done():
		log.Println("ctx.Done()")
		err := ctx.Err()
		if err == nil || strings.Contains(err.Error(), "context canceled") {
			// 正常退出or手动取消
			if biz.RemoveLogIdFromSet(task.LogId) {
				log.Println("Task[" + strconv.FormatInt(task.LogId,10) + "]: kill manually")
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
			log.Println("Task[" + strconv.FormatInt(task.LogId,10) + "]: timeout exit")
			biz.RemoveLogIdFromSet(task.LogId)
			biz.AddExecutionRetToQueue(model.HandleCallbackParam{
				LogId:      task.LogId,
				LogDateTim: time.Now().Unix(),
				HandleCode: common.FailCode,
				HandleMsg:  "timeout",
			})
		}
	}
	log.Println("..--..")
	// 任务完成，从队列里删除
	biz.PopRunningCtxFromList(jobId)
	return nil
}

// 跑任务
func execTask(cancel context.CancelFunc, triggerParam model.TriggerParam) (biz.ReturnT, error) {
	// 找到相应的JobHandler
	executorHandler := triggerParam.ExecutorHandler
	jobHandler := handler.GetJobHandler(executorHandler)
	if jobHandler == nil {
		log.Println("can not found the related executorHandler")
		return biz.NewFailReturnT("jobHandler not exists"), errors.New("can not found the related executorHandler")
	}
	// 传递给任务的参数，其他参数不应给到任务
	param := handler.Param{
		JobParam:   triggerParam.ExecutorParams,
		ShardIndex: triggerParam.BroadcastIndex,
		ShardTotal: triggerParam.BroadcastTotal,
	}
	jobHandler.Init()
	ret := jobHandler.Execute(param)
	jobHandler.Destroy()
	cancel()	// 任务完成
	return ret, nil
}

