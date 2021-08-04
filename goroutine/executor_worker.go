package goroutine

import (
	"context"
	"errors"
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/handler"
	"log"
	"strconv"
	"strings"
	"time"
)

func StartWorker() {
	// 轮询DispatchReqQueue
	for {
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
	// 当前jobId对于的任务都跑完了,这里可能有bug
	biz.DispatchReqQueue.Lock()
	delete(biz.DispatchReqQueue.JobTaskQueueMap, jobId)
	log.Println("JobTaskQueueMap remove jobId: " + strconv.Itoa(jobId))
	biz.DispatchReqQueue.Unlock()
	log.Println("TodoTasks is empty, exists current goroutine...")
}

func trigger(task model.TriggerParam) error {
	jobId := task.JobId
	// todo使用waitGroup阻塞在这儿
	// todo任务参数替换为Context
	timeout := time.Second * time.Duration(task.ExecutorTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	biz.RunningList.Lock()
	biz.RunningList.RunningContextMap[jobId] = &biz.RunningContext{Ctx: ctx, Cancel: cancel}
	biz.RunningList.Unlock()
	go func() {
		_, err := execTask(cancel, task)
		if err != nil {
			log.Println(err)
		}
		// 正常执行完成吗
		if biz.RemoveLogIdFromSet(task.LogId) {
			// 是的
			biz.AddExecutionRetToQueue(model.HandleCallbackParam{
				LogId:      0,
				LogDateTim: 0,
				HandleCode: 0,
				HandleMsg:  "",
			})
		}
	}()
	// 这里会阻塞等待
	log.Println("000999")
	select {
	case <-ctx.Done():
		log.Println("ctx.Done()")
		err := ctx.Err()
		if err == nil {
			// 正常退出or手动取消
			if biz.RemoveLogIdFromSet(task.LogId) {
				// 手动取消，因为如果是正常退出的话，logId已被移除
			}

		}
		if strings.Contains(err.Error(), "context deadline exceeded") {
			// 超时退出
		}
	}
	log.Println("..--..")
	// 任务完成，从队列里删除
	biz.RunningList.Lock()
	delete(biz.RunningList.RunningContextMap, jobId)
	biz.RunningList.Unlock()
	return nil
}

// 跑任务
func execTask(cancel context.CancelFunc, triggerParam model.TriggerParam) (biz.ReturnT, error) {
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

