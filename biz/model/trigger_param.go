package model

type TriggerParam struct {
	JobId int `json:"jobId"`
	ExecutorHandler string `json:"executorHandler"`
	ExecutorParams string `json:"executorParams"`
	ExecutorBlockStrategy string `json:"executorBlockStrategy"`
	ExecutorTimeout int `json:"executorTimeout"`

	LogId int64 `json:"logId"`
	LogDateTime int64 `json:"logDateTime"`

	GlueType string `json:"glueType"`
	GlueSource string `json:"glueSource"`
	GlueUpdatetime int64 `json:"glueUpdatetime"`

	BroadcastIndex int `json:"broadcastIndex"`
	BroadcastTotal int `json:"broadcastTotal"`
}
