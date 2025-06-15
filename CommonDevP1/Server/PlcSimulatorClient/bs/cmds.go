package bs

import (
	//dm "CommonDevP1/PlcSimulator/DeviceMemory"
	"fmt"

	"CommonDevP1/Server/Core"

	"github.com/yuin/gopher-lua"
)

type ICmd interface {
	Text() string
	Desc() string
	Resolve(input string) error
	ArgsCheck(args []string) bool
}

var cmd_1Cylinder = &Cmd_1Cylinder{}
var cmd_2Cylinder = &Cmd_2Cylinder{}
var cmd_1Test = &Cmd_1Test{}

type Cmd_1Cylinder struct {
}

func (this *Cmd_1Cylinder) ArgsCheck(args []string) bool {
	return true
}
func (this *Cmd_1Cylinder) Text() string {
	return "1Cylinder"
}
func (this *Cmd_1Cylinder) Desc() string {
	return "[pos] [pull/push/pullpush]"
}
func (this *Cmd_1Cylinder) Resolve(input string) error {
	fmt.Println("Cmd_1Cylinder Resolve called")

	//  start lua
	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("GoSleep", L.NewFunction(Core.Api2L_Sleep))

	if err := L.DoFile("./lua/1AC_pullOrPush.lua"); err != nil {
		return err
	}

	f_initCylinder := lua.P{
		Fn:   L.GetGlobal("InitCylinder"),
		NRet: 0,
	}
	//fmt.Println(L.GetTop())

	f_startWork := lua.P{
		Fn:   L.GetGlobal("StartWork"),
		NRet: 1,
	}

	//fmt.Println(L.GetTop())

	// f_stopWork := lua.P{
	// 	Fn:   L.GetGlobal("StopWork"),
	// 	NRet: 1,
	// }

	// cylinderConfig
	cylinderConfig := L.GetGlobal("CylinderConfig")
	cylinders := L.GetField(cylinderConfig, "Cylinder").(*lua.LTable)
	cylinderCnt := L.ObjLen(cylinders)
	for i := 0; i < cylinderCnt; i++ {
		cylinder := L.RawGetInt(cylinders, i+1)
		L.SetField(cylinder, "Pos", lua.LNumber(i))

		//cylinderObj := L.GetField(cylinder, "Name")
		//fmt.Println(cylinderObj.String())
	}
	/*
		TaskConfig = {
			TaskType = 0,	--1:pull 2:push
			Force = 0,
			Frequency = 0,
			TotalTimes = 0,
		}
	*/
	// TaskConfig
	taskConfig := L.GetGlobal("TaskConfig")
	//tskTyp := L.RawGet(taskConfig, lua.LString("TaskType"))
	L.SetField(taskConfig, "TaskType", lua.LNumber(1))
	L.SetField(taskConfig, "Force", lua.LNumber(1000))
	L.SetField(taskConfig, "Frequency", lua.LNumber(2))
	L.SetField(taskConfig, "TotalTimes", lua.LNumber(100))

	//f_startWork := L.GetGlobal("StartWork")

	if err := L.CallByParam(f_initCylinder); err != nil {
		//fmt.Println("StartWork call failed")
		return err
	}

	if err := L.CallByParam(f_startWork); err != nil {
		//fmt.Println("StartWork call failed")
		return err
	}

	//_Host.Send()
	return nil
}

//-------------------
type Cmd_2Cylinder struct {
}

func (this *Cmd_2Cylinder) Text() string {
	return "2Cylinder"
}
func (this *Cmd_2Cylinder) Desc() string {
	return "[pos1] [pos2] [sync/async] [pull/push/pullpush]"
}
func (this *Cmd_2Cylinder) Resolve(input string) error {
	fmt.Println("Cmd_2Cylinder Resolve called")
	return nil
}

func (this *Cmd_2Cylinder) ArgsCheck(args []string) bool {
	return true
}

//-------------------

type Cmd_1Test struct {
}

func (this *Cmd_1Test) Text() string {
	return "Cmd_1Test"
}
func (this *Cmd_1Test) Desc() string {
	return ""
}
func (this *Cmd_1Test) Resolve(input string) error {
	fmt.Println("Cmd_1Test Resolve called")

	var tsk Core.ITask = &MyTask2{}
	tsk = _TskMgr.NewTask(tsk, []int{1}, "", "")
	defer tsk.Base().Destory()

	tsk.Start()

	if err := tsk.Base().GetStartResult(); err != nil {
		fmt.Println("task start failed:", err.Error())
	}

	return nil
}

func (this *Cmd_1Test) ArgsCheck(args []string) bool {
	return true
}
