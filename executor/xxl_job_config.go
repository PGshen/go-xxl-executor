package executor

type XxlJobConfig struct {
	Env              string // 环境，dev开发，test测试，prod生产
	AdminAddress     string // 调度中心地址
	AccessToken      string // 访问token
	Appname          string // appname
	Address          string // 执行器地址，不填则默认获取本机地址
	Ip               string // 本机IP，不填会自动获取
	Port             int    // 端口，必填
	LogPath          string // xxl日志，注意这个是xxl日志，不是程序自身日志
	LogRetentionDays int    // 日志保存时间
	HttpTimeout      int    // server超时时间
}
