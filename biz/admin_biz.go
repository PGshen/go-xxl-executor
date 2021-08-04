package biz

import (
	"bytes"
	"encoding/json"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/common"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	timeout time.Duration
	callbackUrl string
	registryUrl string
	registryRemoveUrl string
	contentType string
)

func init() {
	adminAddress := common.Config.XxlJob.Admin.Address
	timeout = time.Duration(common.Config.Http.Timeout)
	callbackUrl = adminAddress + "/callback"
	registryUrl = adminAddress + "/registry"
	registryRemoveUrl = adminAddress + "/registryRemove"
	contentType = "application/json"
}

type AdminBiz struct {
}

// Callback 结果回调
func (a AdminBiz) Callback(paramList []model.HandleCallbackParam) ReturnT {
	return ReturnT{}
}

// Registry 注册
func (a AdminBiz) Registry(param model.RegistryParam) {
	// 超时时间：5秒
	client := &http.Client{Timeout: timeout * time.Second}
	jsonStr, _ := json.Marshal(param)
	resp, err := client.Post(registryUrl, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	result, _ := ioutil.ReadAll(resp.Body)
	log.Printf(string(result))
}

// RegistryRemove 注册摘除
func (a AdminBiz) RegistryRemove(param model.RegistryParam) {

}