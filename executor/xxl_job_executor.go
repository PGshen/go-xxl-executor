package executor

import (
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
	wg.Add(1)	// todo 待完善终止机制
	go startWorker()	// 单独一个线程轮询
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
	defer func() {taskQueue.Running = false}()
	// 串行，暂未实现阻塞处理策略
	for {
		if len(taskQueue.TodoTasks) == 0 {
			break
		}
		taskQueue.Mutex.Lock()
		task := taskQueue.TodoTasks[0]
		taskQueue.TodoTasks = taskQueue.TodoTasks[1:]
		taskQueue.Mutex.Unlock()
		// todo使用waitGroup阻塞在这儿
		// todo任务参数替换为Context
		_ = triggerTask(task)
	}
	// 当前jobId对于的任务都跑完了
	biz.DispatchReqQueue.Mutex.Lock()
	delete(biz.DispatchReqQueue.JobTaskQueueMap, jobId)
	log.Println("JobTaskQueueMap remove jobId: " + strconv.Itoa(jobId))
	biz.DispatchReqQueue.Mutex.Unlock()
	log.Println("TodoTasks is empty, exists current goroutine...")
}

// 跑任务
func triggerTask(triggerParam model.TriggerParam) error {
	executorHandler := triggerParam.ExecutorHandler
	jobHandler := handler.GetJobHandler(executorHandler)
	if jobHandler == nil {
		log.Println("can not found the related executorHandler")
		return errors.New("can not found the related executorHandler")
	}
	// 传递给任务的参数，其他参数不应给到任务
	param := handler.Param{
		JobParam: triggerParam.ExecutorParams,
		ShardIndex: triggerParam.BroadcastIndex,
		ShardTotal: triggerParam.BroadcastTotal,
	}
	jobHandler.Init()
	jobHandler.Execute(param)
	jobHandler.Destroy()
	return nil
}
