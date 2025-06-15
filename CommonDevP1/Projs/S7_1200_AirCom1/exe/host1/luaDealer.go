package main

import (
	"CommonDevP1/PlcSimulator/DeviceMemory"
	"MyToolBox/DataDriver"
	syad "MyToolBox/Thread/SyncAppData"
	"encoding/json"

	"fmt"
	"time"

	"github.com/yuin/gopher-lua"
)

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

	if f_reg := this.L.GetGlobal("RegMe").(*lua.LFunction); f_reg != nil {

		if err := this.L.CallByParam(lua.P{
			Fn:      f_reg,
			NRet:    0,
			Protect: true,
		}); err != nil {
			return err
		}
	}

	if comTb := this.L.GetGlobal("com").(*lua.LTable); comTb != nil {
		fmt.Println(comTb.RawGetString("Description").(lua.LString))
		this.f_startWork = comTb.RawGetString("StartWork").(*lua.LFunction)

		f_GetTaskUsingRes := comTb.RawGetString("GetTaskUsingRes").(*lua.LFunction)
		f_SetTaskConfig := comTb.RawGetString("SetTaskConfig").(*lua.LFunction)
		if f_GetTaskUsingRes == nil || f_SetTaskConfig == nil {
			return fmt.Errorf("get function failed")
		}

		if err := this.L.CallByParam(lua.P{
			Fn:      f_SetTaskConfig,
			NRet:    0,
			Protect: true,
		}, lua.LString(`{"TaskType":3,"Force":0,"TotalTimes":5,"Pos":[2,0],"PosLimitSw":[0,0,0,0],"abc":{"a1":{"b":[{"b1":1,"b2":"abc"},{"b1":1,"b2":"def"}]}}}`)); err != nil {
			return err
		}

		if err := this.L.CallByParam(lua.P{
			Fn:      f_GetTaskUsingRes,
			NRet:    2,
			Protect: true,
		}); err != nil {
			return err
		}

		var poses, limitSws []int
		fmt.Println(this.L.Get(1).(lua.LString))                                //posesStr
		fmt.Println(this.L.Get(2).(lua.LString))                                //limitSwsStr
		json.Unmarshal([]byte(this.L.Get(1).(lua.LString).String()), &poses)    //posesStr
		json.Unmarshal([]byte(this.L.Get(2).(lua.LString).String()), &limitSws) //limitSwsStr

		fmt.Printf("%+v \r\n", poses)    //posesStr
		fmt.Printf("%+v \r\n", limitSws) //limitSwsStr

		//return nil

	} else {
		return fmt.Errorf("failed...")
	}

	this.L.SetGlobal("GetMilliseconds", this.L.NewFunction(GetMilliseconds)) /* Original lua_setglobal uses stack... */

	//top := L.GetTop()
	//this.f_startWork = this.L.GetGlobal("StartWork").(*lua.LFunction)

	this.co, _ = this.L.NewThread()
	this.tskRunFlg = true

	return nil
}

func (this *LuaProc) deal(dt interface{}, errGetDt error) {
	if this.tskRunFlg && dt != nil {
		//status := dt.(*DeviceMemory.DBA_2)

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

//---------------------
var __Std_Time = time.Now()

func GetMilliseconds(L *lua.LState) int {
	//lv := L.ToInt(time.Now().Sub(__Std_Time).Milliseconds())             /* get argument */
	L.Push(lua.LNumber(time.Now().Sub(__Std_Time).Milliseconds())) /* push result */
	return 1
}
