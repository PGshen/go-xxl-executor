package executor

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
func (executor XxlJobExecutor) start() {
}

// 销毁
func (executor XxlJobExecutor) destroy() {

}



func NewXxlJobExecutor(adminAddress, appname, address, ip string, port int) XxlJobExecutor {
	return XxlJobExecutor{adminAddress: adminAddress, appname: appname, address: address, ip: ip, port: port}
}

func (executor *XxlJobExecutor) SetAdminAddress(adminAddress string) {
	executor.adminAddress = adminAddress
}

func (executor *XxlJobExecutor) SetAccessToken(accessToken string) {
	executor.accessToken = accessToken
}

func (executor *XxlJobExecutor) SetAppname(appname string) {
	executor.appname = appname
}

func (executor *XxlJobExecutor) SetAddress(address string) {
	executor.address = address
}

func (executor *XxlJobExecutor) SetIp(ip string) {
	executor.ip = ip
}

func (executor *XxlJobExecutor) SetPort(port int) {
	executor.port = port
}

func (executor *XxlJobExecutor) SetLogPath(logPath string) {
	executor.logPath = logPath
}

func (executor *XxlJobExecutor) SetLogRetentionDays(logRetentionDays int) {
	executor.logRetentionDays = logRetentionDays
}
