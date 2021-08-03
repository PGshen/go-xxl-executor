package executor

import (
	"context"
	"errors"
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/common"
	"github.com/PGshen/go-xxl-executor/handler"
	"log"
	"strconv"
	"sync"
	"time"
)

type XxlJobExecutor struct {
	adminAddress     string
	accessToken      string
	appname          string
	address          string
	ip               string
	port             int
	logPath          string
	logRetentionDays int
}

// 启动
func (executor XxlJobExecutor) Start() {
	log.Println("executor start...")
	var wg sync.WaitGroup
	wg.Add(1)        // todo 待完善终止机制
	go startWorker() // 单独一个线程轮询
	wg.Wait()
}

// 销毁
func (executor XxlJobExecutor) Destroy() {
	log.Println("executor destroy...")
}

func NewXxlJobExecutor() XxlJobExecutor {
	conf := common.Config.XxlJob
	adminAddress := conf.Admin.Address
	appname := conf.Executor.Appname
	address := conf.Executor.Address
	ip := conf.Executor.Ip
	port := conf.Executor.Port
	return XxlJobExecutor{adminAddress: adminAddress, appname: appname, address: address, ip: ip, port: port}
}

func startWorker() {
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
			break
		}
		taskQueue.Mutex.Lock()
		task := taskQueue.TodoTasks[0]
		taskQueue.TodoTasks = taskQueue.TodoTasks[1:]
		taskQueue.Mutex.Unlock()
		_ = trigger(task)
	}
	// 当前jobId对于的任务都跑完了,这里可能有bug
	biz.DispatchReqQueue.Mutex.Lock()
	delete(biz.DispatchReqQueue.JobTaskQueueMap, jobId)
	log.Println("JobTaskQueueMap remove jobId: " + strconv.Itoa(jobId))
	biz.DispatchReqQueue.Mutex.Unlock()
	log.Println("TodoTasks is empty, exists current goroutine...")
}

func trigger(task model.TriggerParam) error {
	jobId := task.JobId
	// todo使用waitGroup阻塞在这儿
	// todo任务参数替换为Context
	timeout := time.Second * time.Duration(task.ExecutorTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	biz.RunningList.Mutex.Lock()
	biz.RunningList.RunningContextMap[jobId] = &biz.RunningContext{Ctx: ctx, Cancel: cancel}
	biz.RunningList.Mutex.Unlock()
	go func() {
		err := execTask(cancel, task)
		if err != nil {
			log.Println(err)
		}
	}()
	// 这里会阻塞等待
	select {
	case <-ctx.Done():
		log.Println("ctx.Done()")
	}
	biz.RunningList.Mutex.Lock()
	delete(biz.RunningList.RunningContextMap, jobId)
	biz.RunningList.Mutex.Unlock()
	return nil
}

// 跑任务
func execTask(cancel context.CancelFunc, triggerParam model.TriggerParam) error {
	executorHandler := triggerParam.ExecutorHandler
	jobHandler := handler.GetJobHandler(executorHandler)
	if jobHandler == nil {
		log.Println("can not found the related executorHandler")
		return errors.New("can not found the related executorHandler")
	}
	// 传递给任务的参数，其他参数不应给到任务
	param := handler.Param{
		JobParam:   triggerParam.ExecutorParams,
		ShardIndex: triggerParam.BroadcastIndex,
		ShardTotal: triggerParam.BroadcastTotal,
	}
	jobHandler.Init()
	jobHandler.Execute(param)
	jobHandler.Destroy()
	cancel()	// 任务完成
	return nil
}
