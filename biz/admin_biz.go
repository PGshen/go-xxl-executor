package biz

import (
	"github.com/PGshen/go-xxl-executor/biz/model"
)

type AdminBiz interface {
	// Callback 结果回调
	Callback(paramList []model.HandleCallbackParam) ReturnT
	// Registry 注册
	Registry(param model.RegistryParam)
	// RegistryRemove 注册摘除
	RegistryRemove(param model.RegistryParam)
}