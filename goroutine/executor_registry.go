package goroutine

import (
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/common"
	"time"
)

// StartRegistry 执行器注册
func StartRegistry(appname, address string) {
	common.Log.Info("start registry...")
	param := model.RegistryParam{
		RegistryGroup: "EXECUTOR",
		RegistryKey: appname,
		RegistryValue: address,
	}
	for {
		biz.Registry(param)
		common.Log.Info("[" + appname + "][" + address + "] registry beat...")
		time.Sleep(10 * time.Second)
	}
}

// RemoveRegistry 执行器摘除
func RemoveRegistry(appname, address string) {
	common.Log.Info("remove registry...")
	param := model.RegistryParam{
		RegistryGroup: "EXECUTOR",
		RegistryKey: appname,
		RegistryValue: address,
	}
	biz.RegistryRemove(param)
}

