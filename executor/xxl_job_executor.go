package executor

import (
	"github.com/PGshen/go-xxl-executor/common"
	"github.com/PGshen/go-xxl-executor/goroutine"
	"log"
	"strconv"
	"strings"
	"sync"
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

// Start 启动
func (executor XxlJobExecutor) Start() {
	log.Println("executor start...")
	var wg sync.WaitGroup
	wg.Add(1)                  // todo 待完善终止机制
	go goroutine.StartRegistry(executor.appname, executor.address)	// 注册协程
	go goroutine.StartWorker() // 单独一个线程轮询
	go goroutine.StartCallback()	// 回调协程
	wg.Wait()
}

// Destroy 销毁
func (executor XxlJobExecutor) Destroy() {
	log.Println("executor destroy...")
	goroutine.RemoveRegistry(executor.appname, executor.address)
}

func NewXxlJobExecutor() XxlJobExecutor {
	conf := common.Config.XxlJob
	adminAddress := conf.Admin.Address
	accessToken := conf.AccessToken
	appname := conf.Executor.Appname
	address := conf.Executor.Address
	ip := conf.Executor.Ip
	port := conf.Executor.Port
	logPath := conf.Executor.LogPath
	logRetentionDays := conf.Executor.LogRetentionDays
	// 如果address没填写则自动获取本机IP
	if strings.TrimSpace(address) == "" {
		if strings.TrimSpace(ip) == "" {
			ip = common.GetInternalIp()
		}
		ipPort := ip + ":" + strconv.Itoa(port)
		address = strings.ReplaceAll("http://ip:port", "ip:port", ipPort)
	}
	return XxlJobExecutor{
		adminAddress: adminAddress,
		accessToken: accessToken,
		appname: appname,
		address: address,
		ip: ip,
		port: port,
		logPath: logPath,
		logRetentionDays: logRetentionDays,
	}
}
