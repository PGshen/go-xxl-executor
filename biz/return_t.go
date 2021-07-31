package biz

import (
	"encoding/json"
	"github.com/PGshen/go-xxl-executor/common"
)

type ReturnT struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Content interface{} `json:"content"`
}

func NewReturnT(code int, msg string) ReturnT {
	return ReturnT{Code: code, Msg: msg}
}

func NewFailReturnT(msg string) ReturnT {
	return ReturnT{Code: common.FailCode, Msg: msg}
}

func (receiver *ReturnT) SetCode(code int) {
	receiver.Code = code
}

func (receiver *ReturnT) SetMsg(msg string) {
	receiver.Msg = msg
}

func (receiver *ReturnT) SetContent(content interface{}) {
	receiver.Content = content
}

func (receiver ReturnT) String() string {
	retByte, _ := json.Marshal(receiver)
	return string(retByte)
}