package handler

import (
	"encoding/json"
	"errors"
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/common"
	"log"
)

var (
	JobHandlerRegistry map[string]IJobHandler = make(map[string]IJobHandler) // JobHandler列表
)

type Param struct {
	JobParam string `json:"jobParam"`
	ShardIndex int `json:"shardIndex"`
	ShardTotal int `json:"shardTotal"`
}

type IJobHandler interface {
	Init()
	Execute(param Param) biz.ReturnT
	Destroy()
}

// AddJobHandler 注册JobHandler
func AddJobHandler(name string, jobHandler IJobHandler) (err error) {
	if _, ok := JobHandlerRegistry[name]; ok {
		return errors.New("JobHandler[" + name + "] already exists!")
	}
	JobHandlerRegistry[name] = jobHandler
	return nil
}

// GetJobHandler 获取JobHandler
func GetJobHandler(name string) IJobHandler {
	if jobHandler, ok := JobHandlerRegistry[name]; ok {
		return jobHandler
	}
	return nil
}

type MethodJobHandler struct {
}

func (receiver *MethodJobHandler) Execute(param Param) biz.ReturnT {
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return biz.ReturnT{}
	}
	paramStr := string(paramBytes)
	log.Println("begin to execute job, receive param: " + paramStr)
	return biz.NewReturnT(common.SuccessCode, "success")
}

func (receiver *MethodJobHandler) Init() {
	log.Println("init...")
}

func (receiver *MethodJobHandler) Destroy() {
	log.Println("destroy...")
}
