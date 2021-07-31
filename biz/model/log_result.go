package model

type LogResult struct {
	FromLineNum int `json:"fromLineNum"`
	ToLineNum int `json:"toLineNum"`
	LogContent int64 `json:"logContent"`
	IsEnd bool `json:"isEnd"`
}
