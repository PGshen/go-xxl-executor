package biz

import (
	"bytes"
	"encoding/json"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/common"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	timeout            time.Duration
	adminBizClientList []AdminBizClient
	callbackUri        = "/api/callback"
	registryUri        = "/api/registry"
	registryRemoveUri  = "/api/registryRemove"
)

const XXL_JOB_ACCESS_TOKEN = "XXL-JOB-ACCESS-TOKEN"

type AdminBizClient struct {
	AdminAddress string
	AccessToken string
}

// 初始化加载admin client
func InitAdminBizClient(adminAddresses, accessToken string, httpTimeout int) {
	if adminAddresses != "" && len(strings.TrimSpace(adminAddresses)) > 0 {
		for _, adminAddress := range strings.Split(adminAddresses, ",") {
			adminBiz := AdminBizClient{
				AdminAddress: adminAddress,
				AccessToken: accessToken,
			}
			adminBizClientList = append(adminBizClientList, adminBiz)
		}
	}
	timeout = time.Duration(httpTimeout)
}

// GetAdminBizClientList 获取adminBiz列表
func GetAdminBizClientList() []AdminBizClient {
	return adminBizClientList
}

// GetAdminBizClient 随机获取一个adminBiz
func GetAdminBizClient() AdminBizClient {
	rand.Seed(time.Now().Unix())
	return adminBizClientList[rand.Intn(len(adminBizClientList))]
}

// Callback 结果回调
func Callback(paramList []model.HandleCallbackParam) ReturnT {
	return post(callbackUri, "POST", paramList)
}

// Registry 注册
func Registry(param model.RegistryParam) ReturnT {
	return post(registryUri, "POST", param)
}

// RegistryRemove 注册摘除
func RegistryRemove(param model.RegistryParam) ReturnT {
	return post(registryRemoveUri, "POST", param)
}

// post请求
func post(uri string, method string, param interface{}) ReturnT {
	adminBizList := GetAdminBizClientList()
	for _, adminBiz := range adminBizList {

		url := adminBiz.AdminAddress + uri
		// 超时时间：5秒
		client := &http.Client{Timeout: timeout * time.Second}
		jsonStr, _ := json.Marshal(param)
		request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
		if request == nil || err != nil {
			common.Log.Error(err)
			continue
		}
		request.Header.Add("ContentType", "application/json")
		request.Header.Add(XXL_JOB_ACCESS_TOKEN, adminBiz.AccessToken)
		resp, err := client.Do(request)
		if err != nil {
			common.Log.Error(err)
			continue
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				common.Log.Info(err)
			}
		}(resp.Body)

		result, _ := ioutil.ReadAll(resp.Body)
		return NewReturnT(common.SuccessCode, string(result))
	}
	return NewFailReturnT("post error")
}