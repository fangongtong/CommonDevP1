package httpCtrl

import (
	"CommonDevP1/Projs/S7_1200_AirCom1/exe/plchost/other"
	syad "MyToolBox/Thread/SyncAppData"
)

type RtnData struct {
	Code    int
	ErrInfo string
	RtnData interface{}
}

const (
	__ReqErr_BadParam = 1 + iota
	__ReqErr_PushCmdsFailed
	__ReqErr_OperOvertime
	__ReqErr_ScriptFileNotFound
	__ReqErr_StartTaskFailed
	__ReqErr_ReqMenuFailed
	__ReqErr_MAX
)

var __Err_Info = make([]string, __ReqErr_MAX)

func init() {
	__Err_Info[__ReqErr_BadParam] = "request param error: %v"
	__Err_Info[__ReqErr_PushCmdsFailed] = "push cmds failed: %v"
	__Err_Info[__ReqErr_OperOvertime] = "oper overtime"
	__Err_Info[__ReqErr_ScriptFileNotFound] = `task script file "%v" not found`
	__Err_Info[__ReqErr_StartTaskFailed] = "start task failed: %v"
	__Err_Info[__ReqErr_ReqMenuFailed] = "request menu failed: %v"
}

// var __DtDriver DataDriver.ITopicDriver

// func SetDtDriver(driver DataDriver.ITopicDriver) {
// 	__DtDriver = driver
// }

// var __TaskProcDealer *bs.TaskProc

// func SetTaskProcDealer(dealer *bs.TaskProc) {
// 	__TaskProcDealer = dealer
// }

var __TaskList *other.TaskList

func SetTaskList(list *other.TaskList) {
	__TaskList = list
}

var __SyncDataMgr *syad.SyncDataMgr

func SetSyncMgr(syncMgr *syad.SyncDataMgr) {
	__SyncDataMgr = syncMgr
}

var __MaxCmdCnt int

func SetMaxCmdCnt(val int) {
	__MaxCmdCnt = val
}
