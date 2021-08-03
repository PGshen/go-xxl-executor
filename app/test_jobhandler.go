package app

import (
	"encoding/json"
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/common"
	"github.com/PGshen/go-xxl-executor/handler"
	"log"
	"strconv"
	"time"
)

type TestJobHandler struct {
	handler.MethodJobHandler
}

func (receiver *TestJobHandler) Execute(param handler.Param) biz.ReturnT {
	receiver.MethodJobHandler.Execute(param)
	log.Println("Test...")
	log.Println("sleep 30s")
	jobParams := make(map[string]interface{})
	_ = json.Unmarshal([]byte(param.JobParam), &jobParams)
	times := int(jobParams["times"].(float64))
	for i := 0; i < times; i++ {
		log.Println("Test running: " + strconv.Itoa(i))
		time.Sleep(time.Second)
	}
	log.Println("Finish work!!!")
	return biz.NewReturnT(common.SuccessCode, "Test JobHandler")
}

func (receiver TestJobHandler) Init() {
	log.Println("init something...")
}

func (receiver TestJobHandler) Destroy() {
	log.Println("destroy...")
}


