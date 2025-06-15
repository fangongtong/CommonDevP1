package structs

import (
	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"fmt"

	"github.com/yuin/gopher-lua"
)

// func New_LuaTbKey(tb *lua.LTable, key string) *LuaTbKey{
// 	return &LuaTbKey{
// 		tb:tb,
// 		key:key,
// 	}
// }
type LuaTbKey struct {
	tb  *lua.LTable
	key string
}

func (this *LuaTbKey) Reset(tb *lua.LTable, key string) {
	this.tb, this.key = tb, key
}
func (this *LuaTbKey) SetKeyVal(val lua.LValue) {
	this.tb.RawSetString(this.key, val)
}

type LuaStatePack struct {
	Alarm        LuaTbKey
	RealForce    LuaTbKey
	MaxForce     LuaTbKey
	RealDisplace LuaTbKey
	Threshold    LuaTbKey
}

func (this *LuaStatePack) Reset(tb *lua.LTable) {
	this.Alarm.Reset(tb, "Alarm")
	this.RealForce.Reset(tb, "RealForce")
	this.MaxForce.Reset(tb, "MaxForce")
	this.RealDisplace.Reset(tb, "RealDisplace")
	this.Threshold.Reset(tb, "Threshold")
}

type LuaStatesPath struct {
	DevAlarm LuaTbKey
	Status   []LuaStatePack
}

func (this *LuaStatesPath) Reset(baseTb *lua.LTable, maxPosCnt int) error {
	real := baseTb.RawGet(lua.LString("RealDt")).(*lua.LTable)
	if real == nil {
		return __Err_Init_LuaStatus_Failed
	}
	this.DevAlarm.Reset(real, "DevAlarm")

	this.Status = make([]LuaStatePack, maxPosCnt)
	for i := 0; i < maxPosCnt; i++ {

		// Alarm:=tb.RawGetString("Alarm")
		// RealForce:=tb.RawGetString("RealForce")
		// MaxForce:=tb.RawGetString("MaxForce")
		// RealDisplace:=tb.RawGetString("RealDisplace")
		// Threshold:=tb.RawGetString("Threshold")
		if tmp := real.RawGet(lua.LString("CR")).(*lua.LTable); tmp != nil {

			if tmp := tmp.RawGetInt(i + 1).(*lua.LTable); tmp != nil {
				this.Status[i].Reset(tmp)
				continue
			}

		}
		return __Err_Init_LuaStatus_Failed
	}
	return nil
}
func (this *LuaStatesPath) SetVal(status *tmpl.DBA_Status) {
	this.DevAlarm.SetKeyVal(lua.LNumber(status.DevAlarm))
	for i := 0; i < len(this.Status); i++ {
		this.Status[i].Alarm.SetKeyVal(lua.LNumber(status.CylinderDataAry[i].Alarm))
		this.Status[i].RealForce.SetKeyVal(lua.LNumber(status.CylinderDataAry[i].RealForce))
		this.Status[i].MaxForce.SetKeyVal(lua.LNumber(status.CylinderDataAry[i].PeekForce))
		this.Status[i].RealDisplace.SetKeyVal(lua.LNumber(status.CylinderDataAry[i].RealDisplacement))
		this.Status[i].Threshold.SetKeyVal(lua.LNumber(status.CylinderDataAry[i].DriveDegree))
	}
}

var __Err_Init_LuaStatus_Failed = fmt.Errorf("get lua common status item failed")
