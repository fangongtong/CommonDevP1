package bs

import (
	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"time"
)

var _InstName string = "Host"

func SetInstName(name string) {
	_InstName = name
}

const (
	Topic_Main = iota
	Topic_PlcCmmu
	Topic_Task
	Topic_DbOther
	Topic_WebSock
	Topic_Influxdb
)

const (
	Sync_Cmd = iota
	Sync_PosMenuParam
	Sync_PlcConnStatus
	Sync_SysAlarmInfo
)

//var SyncMgr = common.New_SyncDataMgr(8)

type SyncData_Base struct {
	InReq  bool
	ReqSNo uint32
}

// type SyncData_HttpReq_PirorEnable struct {
// 	SyncData_Base
// 	Enable bool
// 	RspDt  interface{}
// }

// type SyncData_HttpReq_Interpose struct {
// 	SyncData_Base
// 	RspDt interface{}
// }

type SyncData_PosMenuParam struct {
	SyncData_Base
	Pos   int
	Err   error
	RspDt tmpl.DBA_Sub_MenuParam
}

type SyncData_PlcConnStatus struct {
	Connected bool
}

type SyncData_SysAlarmInfo struct {
	// 0:plc conn failed
	// 1:plc recv timeout for task
	// 2:
	AlarmInfo uint32
}

//---------

type TaskLowInfo struct {
	TaskUid    string
	TaskStatus int
	Poses      []int
	LimitSws   []int
	PosFields  string
	TaskInfo   string
}
type RealTasksInfo struct {
	DevAlarm uint16
	Status   []tmpl.DBA_Sub_CylinderDatas
	TaskInfo []TaskLowInfo
}

type LogTasksInfo struct {
	TimeFlg  time.Time
	TaskInfo []*TaskLowInfo
	Status   []tmpl.DBA_Sub_CylinderDatas
}
