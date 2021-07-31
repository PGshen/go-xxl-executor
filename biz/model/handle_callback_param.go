package model

type HandleCallbackParam struct {
	LogId int64 `json:"logId"`
	LogDateTim int64 `json:"logDateTim"`
	HandleCode int `json:"handleCode"`
	HandleMsg string `json:"handleMsg"`
}