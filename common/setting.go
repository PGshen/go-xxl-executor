package common

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var (
	Config *Conf
)

type Conf struct {
	XxlJob XxlJob `yaml:"xxl_job"`
	Http Http `yaml:"http"`
}

type XxlJob struct {
	Admin XxlJobAdmin `yaml:"admin"`
	AccessToken string `yaml:"access_token"`
	Executor XxlJobExecutor `yaml:"executor"`
}

type XxlJobAdmin struct {
	Address string `yaml:"addresses"`
}

type XxlJobExecutor struct {
	Appname string `yaml:"appname"`
	Address string `yaml:"address"`
	Ip string `yaml:"ip"`
	Port int `yaml:"port"`
	LogPath string `yaml:"log_path"`
	LogRetentionDays int `yaml:"log_retention_days"`
}

type Http struct {
	Timeout int `yaml:"timeout"`
}

func init() {
	Config = getConf()
	log.Println("[Setting] config init success")
}

// 读取配置
func getConf() *Conf {
	var c *Conf
	file, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalln("[Setting] config error: ", err)
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Fatalln("[Setting] yaml unmarshal error: ", err)
	}
	return c
}