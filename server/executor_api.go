package server

import (
	"encoding/json"
	"fmt"
	"github.com/PGshen/go-xxl-executor/biz"
	"github.com/PGshen/go-xxl-executor/biz/model"
	"github.com/PGshen/go-xxl-executor/common"
	"github.com/PGshen/go-xxl-executor/common/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

var executorBiz = biz.ExecutorBiz{}

// 启动http服务
func StartServer(ip string, port int) {
	addr := ip + ":" + strconv.Itoa(port)
	common.Log.Info("http server start...")
	http.HandleFunc("/beat", beat)
	http.HandleFunc("/idleBeat", idleBeat)
	http.HandleFunc("/run", run)
	http.HandleFunc("/kill", kill)
	http.HandleFunc("/log", loglog)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// 心跳检测
func beat(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, executorBiz.Beat().String())
}

// 空闲检测
func idleBeat(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Log.Info("read body err, %v\n", err)
		return
	}

	var param model.IdleBeatParam
	if err = json.Unmarshal(body, &param); err != nil {
		common.Log.Info("Unmarshal err, %v\n", err)
		return
	}
	_, _ = fmt.Fprintln(w, executorBiz.IdleBeat(param))
}

// 运行任务
func run(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Log.Info("read body err, %v\n", err)
		return
	}

	var param model.TriggerParam
	if err = json.Unmarshal(body, &param); err != nil {
		common.Log.Info("Unmarshal err, %v\n", err)
		return
	}
	_, _ = fmt.Fprintln(w, executorBiz.Run(param))
}

// 终止任务
func kill(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Log.Info("read body err, %v\n", err)
		return
	}

	var param model.KillParam
	if err = json.Unmarshal(body, &param); err != nil {
		common.Log.Info("Unmarshal err, %v\n", err)
		return
	}
	_, _ = fmt.Fprintln(w, executorBiz.Kill(param))
}

// 查日志
func loglog(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Log.Info("read body err, %v\n", err)
		return
	}

	var param model.LogParam
	if err = json.Unmarshal(body, &param); err != nil {
		common.Log.Info("Unmarshal err, %v\n", err)
		return
	}
	_, _ = fmt.Fprintln(w, executorBiz.Log(param))
}
