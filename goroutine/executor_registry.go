package goroutine

import (
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"log"
	"time"
)

// StartRegistry 执行器注册
func StartRegistry(appname, address string) {
	log.Println("start registry...")
	param := model.RegistryParam{
		RegistryGroup: "EXECUTOR",
		RegistryKey: appname,
		RegistryValue: address,
	}
	for {
		biz.Registry(param)
		log.Printf("[" + appname + "][" + address + "] registry beat...")
		time.Sleep(10 * time.Second)
	}
}

// RemoveRegistry 执行器摘除
func RemoveRegistry(appname, address string) {
	log.Printf("remove registry...")
	param := model.RegistryParam{
		RegistryGroup: "EXECUTOR",
		RegistryKey: appname,
		RegistryValue: address,
	}
	biz.RegistryRemove(param)
}

