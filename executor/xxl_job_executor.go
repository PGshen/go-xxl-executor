package executor

import (
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/common"
	"github.com/PGshen/go-xxl-executor/goroutine"
	"github.com/PGshen/go-xxl-executor/server"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type XxlJobExecutor struct {
	env              string
	adminAddress     string
	accessToken      string
	appname          string
	address          string
	ip               string
	port             int
	logPath          string
	logRetentionDays int
	httpTimeout      int
}

func NewXxlJobExecutor(conf XxlJobConfig) XxlJobExecutor {
	adminAddress := conf.AdminAddress
	accessToken := conf.AccessToken
	appname := conf.Appname
	address := conf.Address
	ip := conf.Ip
	port := conf.Port
	logPath := conf.LogPath
	logRetentionDays := conf.LogRetentionDays
	httpTime := conf.HttpTimeout
	env := conf.Env
	// 如果address没填写则自动获取本机IP
	if strings.TrimSpace(address) == "" {
		if strings.TrimSpace(ip) == "" {
			ip = common.GetInternalIp()
		}
		ipPort := ip + ":" + strconv.Itoa(port)
		address = strings.ReplaceAll("http://ip:port", "ip:port", ipPort)
	}
	common.InitLogger(logPath, env)	// 提前初始化Logger
	return XxlJobExecutor{
		adminAddress:     adminAddress,
		accessToken:      accessToken,
		appname:          appname,
		address:          address,
		ip:               ip,
		port:             port,
		logPath:          logPath,
		logRetentionDays: logRetentionDays,
		httpTimeout:      httpTime,
	}
}

// Start 启动
func (executor XxlJobExecutor) Start() {
	defer executor.Destroy()
	common.Log.Info("executor start...")
	biz.InitAdminBizClient(executor.adminAddress, executor.accessToken, executor.httpTimeout)
	biz.InitExecutorBiz(executor.logPath)
	go server.StartServer(executor.ip, executor.port)                       // 启动http服务
	go goroutine.StartRegistry(executor.appname, executor.address)          // 注册协程
	go goroutine.StartWorker()                                              // 单独一个线程轮询
	go goroutine.StartCallback()                                            // 回调协程
	go goroutine.StartCleanLog(executor.logPath, executor.logRetentionDays) // 日志定期清理
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// Destroy 销毁
func (executor XxlJobExecutor) Destroy() {
	common.Log.Info("executor destroy...")
	goroutine.RemoveRegistry(executor.appname, executor.address)
}
