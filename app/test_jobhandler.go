package app

import (
	"encoding/json"
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/common"
	"github.com/PGshen/go-xxl-executor/handler"
	"log"
	"time"
)

type TestJobHandler struct {
	handler.MethodJobHandler
}

func (receiver *TestJobHandler) Execute(param handler.Param) biz.ReturnT {
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return biz.ReturnT{}
	}
	paramStr := string(paramBytes)
	log.Println("begin to execute job, receive param: " + paramStr)
	log.Println("Test...")
	log.Println("sleep 10s")
	time.Sleep(5 * time.Second)
	log.Println("Finish work!!!")
	return biz.NewReturnT(common.SuccessCode, "Test JobHandler")
}

func (receiver TestJobHandler) Init() {
	log.Println("init something...")
}

func (receiver TestJobHandler) Destroy() {
	log.Println("destroy...")
}


