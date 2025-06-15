package structs

import (
	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"encoding/json"
	"sync/atomic"

	"fmt"
	"time"

	"github.com/yuin/gopher-lua"
)

var __Std_Time = time.Now()

var __Uid uint32 = 1

// func NewTask(luaFilePth string) *Task {
// 	return &Task{
// 		luaFilePath: luaFilePth,
// 	}
// }

type IResReg interface {
	Reg(poses []int, limitSwitch []int, resUid uint32) error
	UnReg(resUid uint32)
	GetPosMenuParam(pos int) (*tmpl.DBA_Sub_MenuParam, error)
}

type Task struct {
	TskUid     string
	taskStatus int // 0: uninit 1: inited 2: pause 3: running 4: over
	//luaFilePath string
	lState        *lua.LState
	coState       *lua.LState
	fStartWork    *lua.LFunction
	sendingCmdIdx uint32
	resReger      IResReg
	tskResUid     uint32
	luaStatusTB   LuaStatesPath
	errInLuaCommu error

	usingPoses []int
}

var __Err_Init_GetStartWork_Failed = fmt.Errorf("get startwork function failed")

func (this *Task) Init(uuid string, reger IResReg, luaFilePth string, jsCfg string) error {
	isOk := false
	this.TskUid = uuid
	this.resReger = reger
	this.lState = lua.NewState()
	this.coState, _ = this.lState.NewThread()
	defer func() {
		if !isOk {
			this.lState.Close()
			this.lState = nil
		}
	}()

	if err := this.lState.DoFile(luaFilePth); err != nil {
		return err
	}

	if fReg := this.lState.GetGlobal("RegMe").(*lua.LFunction); fReg != nil {

		if err := this.lState.CallByParam(lua.P{
			Fn:      fReg,
			NRet:    0,
			Protect: true,
		}); err != nil {
			return err
		}
	}

	if comTb := this.lState.GetGlobal("com").(*lua.LTable); comTb != nil {
		if err := this.luaStatusTB.Reset(comTb, 8); err != nil {
			return err
		}

		this.fStartWork = comTb.RawGetString("StartWork").(*lua.LFunction)

		fGetTaskUsingRes := comTb.RawGetString("GetTaskUsingRes").(*lua.LFunction)
		fSetTaskConfig := comTb.RawGetString("SetTaskConfig").(*lua.LFunction)
		if fGetTaskUsingRes == nil || fSetTaskConfig == nil {
			return fmt.Errorf("get function failed")
		}

		if err := this.lState.CallByParam(lua.P{
			Fn:      fSetTaskConfig,
			NRet:    0,
			Protect: true,
		}, lua.LString(jsCfg)); err != nil {
			return err
		}

		if err := this.lState.CallByParam(lua.P{
			Fn:      fGetTaskUsingRes,
			NRet:    2,
			Protect: true,
		}); err != nil {
			return err
		}

		//  从lua获得所占资源
		var poses, limitSws []int
		if err := json.Unmarshal([]byte(this.lState.Get(1).(lua.LString).String()), &poses); err != nil { //posesStr
			return err
		}
		if err := json.Unmarshal([]byte(this.lState.Get(2).(lua.LString).String()), &limitSws); err != nil { //limitSwsStr
			return err
		}
		//  注册占用资源
		if err := this.regRes(poses, limitSws); err != nil {
			return err
		}

		//  如果init失败则马上释放资源
		defer func() {
			if !isOk {
				this.unregRes()
			}
		}()
		//fmt.P

	} else {
		return __Err_Init_LuaStatus_Failed
	}

	this.lState.SetGlobal("GetMilliseconds", this.lState.NewFunction(GetMilliseconds))
	this.lState.SetGlobal("LoadMenuParam", this.lState.NewFunction(this.loadMenuParam))

	//this.fStartWork = this.lState.GetGlobal("StartWork").(*lua.LFunction)
	if this.fStartWork != nil {
		isOk = true
		this.taskStatus = 1
	} else {
		return __Err_Init_GetStartWork_Failed
	}
	return nil
}

func (this *Task) UnInit() {
	if this.lState != nil {
		this.lState.Close()
		this.lState = nil
		this.taskStatus = 0
		this.unregRes()
	}
}

func (this *Task) unregRes() {
	//if this.tskResUid > 0 {
	this.resReger.UnReg(this.tskResUid)
	//	this.tskResUid = 0
	//}
}

// func (this *Task) resetConfigDt(jsCfg string) error {
// 	// 1. 向lua虚拟机设置config参数

// 	// 2. 通知注册占用的工位和行程开关
// 	this.tskResUid = atomic.AddUint32(&__Uid, 1)

// 	return this.resReger.Reg(nil, nil, this.tskResUid)
// }
func (this *Task) regRes(poses, limitSws []int) error {
	// 通知注册占用的工位和行程开关
	this.tskResUid = atomic.AddUint32(&__Uid, 1)

	if err := this.resReger.Reg(poses, limitSws, this.tskResUid); err != nil {
		return err
	} else {
		this.usingPoses = poses
	}
	return nil
}
func (this *Task) Run() bool {
	switch this.taskStatus {
	case 1, 2, 3:
		this.taskStatus = 3
		return true
	default:
		return false
	}
}

func (this *Task) Pause() bool {
	switch this.taskStatus {
	case 2, 3:
		this.taskStatus = 2
		return true
	default:
		return false
	}
}

func (this *Task) Status() int {
	return this.taskStatus
}

func (this *Task) SetCmdIdx_SendingCmd(idx uint32) {
	this.sendingCmdIdx = idx
}

type RecordInfo struct {
	NeedRecord bool
	TaskUid    string
	//TaskStatus    int
	TaskInfo      string
	TaskPosFields string
	Poses         []int
}

func (this *Task) Check(status *tmpl.DBA_Status, rtnCmds []tmpl.DBA_Sub_CmdParam, rdInf *RecordInfo) (int, error) {

	//  首先应该检查工位报警,有报警则不工作
	for i := 0; i < len(this.usingPoses); i++ {
		if this.usingPoses[i] == 0 {
			continue
		}
		if status.CylinderDataAry[this.usingPoses[i]-1].Alarm > 0 {
			return 0, nil
		}
	}

	//return 0, nil

	//startTime := time.Now()
	var lastTime time.Time
	var t2 int64

	apdCmdCnt := 0
	//status := dt.(*tmpl.DBA_Status)

	//this.lState.SetTable(this.luaStatusTB, , )
	this.luaStatusTB.SetVal(status)

	for {
		lastTime = time.Now()
		_, err, values := this.lState.Resume(this.coState, this.fStartWork) // err is nil

		t2 = time.Now().Sub(lastTime).Milliseconds()
		if t2 > 1 {
			fmt.Println("t2 overtime:", t2, " and code:", lua.LVAsNumber(values[1]))
		}
		//lastTime = time.Now()

		if err != nil {
			this.taskStatus = 4
			fmt.Println("something wrong:", err.Error())
			return 0, err
		}
		//fmt.Println("value cnt:", len(values))
		switch lua.LVAsNumber(values[1]) {
		case 0:
			if lua.LVAsNumber(values[2]) == 1 {
				//this.tskRunFlg = false
				return apdCmdCnt, nil
			} else {
				this.taskStatus = 4
				fmt.Println("over~")
				return apdCmdCnt, nil
			}
		case 1:
			if apdCmdCnt == len(rtnCmds) {
				return 0, __Err_OverMaxCnt
			}

			//var cmds []*DeviceMemory.DBCmdParam
			// cmds = append(cmds, &DeviceMemory.DBCmdParam{
			// 	CmdCode: 1,
			// 	Param1:  uint16(int(lua.LVAsNumber(values[2]))),
			// })
			rtnCmds[apdCmdCnt].CmdCode = 1
			rtnCmds[apdCmdCnt].Param1 = uint16(int(lua.LVAsNumber(values[2])))

			apdCmdCnt++

			//this.topicDealer.GetDriver().PutData(Topic_Cmd, cmds)
		case 3, 4:
			if apdCmdCnt == len(rtnCmds) {
				return 0, __Err_OverMaxCnt
			}

			// var cmds []*DeviceMemory.DBCmdParam
			// cmds = append(cmds, &DeviceMemory.DBCmdParam{
			// 	CmdCode: uint16(int(lua.LVAsNumber(values[1]))),
			// 	Param1:  uint16(int(lua.LVAsNumber(values[2]))),
			// 	Param2:  int32(lua.LVAsNumber(values[3])),
			// 	Param3:  int32(lua.LVAsNumber(values[4])),
			// 	Param5:  float32(lua.LVAsNumber(values[5])),
			// })

			rtnCmds[apdCmdCnt].CmdCode = uint16(int(lua.LVAsNumber(values[1])))
			rtnCmds[apdCmdCnt].Param1 = uint16(int(lua.LVAsNumber(values[2])))
			rtnCmds[apdCmdCnt].Param2 = int32(lua.LVAsNumber(values[3]))
			rtnCmds[apdCmdCnt].Param3 = int32(lua.LVAsNumber(values[4]))
			rtnCmds[apdCmdCnt].Param5 = float32(lua.LVAsNumber(values[5]))

			apdCmdCnt++
		case 5, 6:
			if apdCmdCnt == len(rtnCmds) {
				return 0, __Err_OverMaxCnt
			}

			// var cmds []*DeviceMemory.DBCmdParam
			// cmds = append(cmds, &DeviceMemory.DBCmdParam{
			// 	CmdCode: uint16(int(lua.LVAsNumber(values[1]))),
			// 	Param1:  uint16(int(lua.LVAsNumber(values[2]))),
			// 	Param2:  int32(lua.LVAsNumber(values[3])),
			// 	Param3:  int32(lua.LVAsNumber(values[4])),
			// 	Param5:  float32(lua.LVAsNumber(values[5])),
			// })

			rtnCmds[apdCmdCnt].CmdCode = uint16(int(lua.LVAsNumber(values[1])))
			rtnCmds[apdCmdCnt].Param1 = uint16(int(lua.LVAsNumber(values[2])))
			rtnCmds[apdCmdCnt].Param2 = int32(lua.LVAsNumber(values[3]))
			rtnCmds[apdCmdCnt].Param3 = int32(lua.LVAsNumber(values[4]))
			rtnCmds[apdCmdCnt].Param4 = float32(lua.LVAsNumber(values[5]))
			rtnCmds[apdCmdCnt].Param5 = float32(lua.LVAsNumber(values[6]))

			apdCmdCnt++
		case 20, 15, 16:
			if apdCmdCnt == len(rtnCmds) {
				return 0, __Err_OverMaxCnt
			}

			rtnCmds[apdCmdCnt].CmdCode = uint16(int(lua.LVAsNumber(values[1])))
			rtnCmds[apdCmdCnt].Param1 = uint16(int(lua.LVAsNumber(values[2])))
			rtnCmds[apdCmdCnt].Param4 = float32(lua.LVAsNumber(values[3]))

			apdCmdCnt++
			//this.topicDealer.GetDriver().PutData(Topic_Cmd, cmds)

		case 10, 17:
			if apdCmdCnt == len(rtnCmds) {
				return 0, __Err_OverMaxCnt
			}

			rtnCmds[apdCmdCnt].CmdCode = uint16(int(lua.LVAsNumber(values[1])))
			rtnCmds[apdCmdCnt].Param1 = uint16(int(lua.LVAsNumber(values[2])))
			rtnCmds[apdCmdCnt].Param2 = int32(lua.LVAsNumber(values[3]))

			apdCmdCnt++
			//this.topicDealer.GetDriver().PutData(Topic_Cmd, cmds)

		case 255:
			rdInf.NeedRecord = true
			rdInf.TaskInfo = lua.LVAsString(values[2])
			rdInf.TaskPosFields = lua.LVAsString(values[3])
			rdInf.Poses = this.usingPoses
			rdInf.TaskUid = this.TskUid
		default:
			fmt.Printf("unknown value: %v \r\n", values[1])
			//needSleep = true
			//return apdCmdCnt
		}
		if !lua.LVAsBool(values[0]) {
			//needSleep = true
			return apdCmdCnt, nil
		}
	}
}

func (this *Task) loadMenuParam(L *lua.LState) int {
	pos := int(L.Get(1).(lua.LNumber))
	menu, err := this.resReger.GetPosMenuParam(pos - 1)
	if err != nil {
		this.errInLuaCommu = err
		fmt.Println("loadMenuParam failed: ", err.Error())
		L.Push(lua.LString(""))
		L.Push(lua.LBool(false))
	} else {
		str, _ := json.Marshal(menu)
		L.Push(lua.LString(str))
		L.Push(lua.LBool(true))
	}

	return 2
}
func GetMilliseconds(L *lua.LState) int {
	//lv := L.ToInt(time.Now().Sub(__Std_Time).Milliseconds())             /* get argument */
	L.Push(lua.LNumber(time.Now().Sub(__Std_Time).Milliseconds())) /* push result */
	return 1
}

//--------------
/*
type LuaProc struct {
	topicDealer *DataDriver.TopicDataDealer
	recvBuf     [1024]byte
	syncMgr     *syad.SyncDataMgr
	L           *lua.LState
	co          *lua.LState
	f_startWork *lua.LFunction
	tskRunFlg   bool
}

func (this *LuaProc) Init(topicDealer *DataDriver.TopicDataDealer, syncMgr *syad.SyncDataMgr) error {

	this.topicDealer = topicDealer
	this.topicDealer.ChangeDealer(this.deal)

	this.syncMgr = syncMgr

	this.L = lua.NewState()
	// defer func() {
	// 	fmt.Println("defer 1")
	// 	L.Close()
	// }()

	if err := this.L.DoFile("./tasks/1AC_pullOrPush.lua"); err != nil {
		return err
	}

	this.L.SetGlobal("GetMilliseconds", this.L.NewFunction(GetMilliseconds))

	//top := L.GetTop()
	this.f_startWork = this.L.GetGlobal("StartWork").(*lua.LFunction)

	this.co, _ = this.L.NewThread()
	this.tskRunFlg = true

	return nil
}

func (this *LuaProc) deal(dt interface{}, errGetDt error) {
	if this.tskRunFlg && dt != nil {
		status := dt.(*DeviceMemory.DBA_2)

		_, err, values := this.L.Resume(this.co, this.f_startWork) // err is nil
		if err != nil {
			fmt.Println("something wrong:", err.Error())
			return
		}
		//fmt.Println("value cnt:", len(values))
		switch lua.LVAsNumber(values[1]) {
		case 0:
			if lua.LVAsNumber(values[2]) == 1 {
				//this.tskRunFlg = false
				return
			} else {
				this.tskRunFlg = false
				fmt.Println("over~")
				return
			}
		case 1:

			var cmds []*DeviceMemory.DBCmdParam

			cmds = append(cmds, &DeviceMemory.DBCmdParam{
				CmdCode: 1,
				Param1:  uint16(int(lua.LVAsNumber(values[2]))),
			})

			//fmt.Printf("%+v \r\n", cmds[0])
			this.topicDealer.GetDriver().PutData(Topic_Cmd, cmds)

			//fmt.Println("Call Plc Release")
		case 3, 4:

			var cmds []*DeviceMemory.DBCmdParam

			cmds = append(cmds, &DeviceMemory.DBCmdParam{
				CmdCode: uint16(int(lua.LVAsNumber(values[1]))),
				Param1:  uint16(int(lua.LVAsNumber(values[2]))),
				Param2:  int32(lua.LVAsNumber(values[3])),
				Param3:  int32(lua.LVAsNumber(values[4])),
				Param5:  float32(lua.LVAsNumber(values[5])),
			})
			//fmt.Printf("%+v \r\n", cmds[0])
			this.topicDealer.GetDriver().PutData(Topic_Cmd, cmds)
			//fmt.Println("Call Plc Pull")
		default:
			fmt.Printf("unknown value: %v \r\n", values[1])
			//needSleep = true
			return
		}
		if !lua.LVAsBool(values[0]) {
			//needSleep = true
			return
		}
	}
}
*/
