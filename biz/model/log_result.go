package model

type LogResult struct {
	FromLineNum int `json:"fromLineNum"`
	ToLineNum int `json:"toLineNum"`
	LogContent string `json:"logContent"`
	IsEnd bool `json:"isEnd"`
}
