package main

import (
	"time"
	//"CommonDevP1/ScriptMod"
	"fmt"

	"github.com/yuin/gopher-lua"
)

var __Std_Time = time.Now()

func GetMilliseconds(L *lua.LState) int {
	//lv := L.ToInt(time.Now().Sub(__Std_Time).Milliseconds())             /* get argument */
	L.Push(lua.LNumber(time.Now().Sub(__Std_Time).Milliseconds())) /* push result */
	return 1
}

func main() {
	fmt.Println("test start")

	L := lua.NewState()
	defer func() {
		fmt.Println("defer 1")
		L.Close()
	}()

	if err := L.DoFile("./tasks/1AC_pullOrPush.lua"); err != nil {
		panic(err)
	}

	L.SetGlobal("GetMilliseconds", L.NewFunction(GetMilliseconds)) /* Original lua_setglobal uses stack... */

	top := L.GetTop()

	// f_startWork := lua.P{
	// 	Fn:   L.GetGlobal("StartWork"),
	// 	NRet: 1,
	// }
	//cylinderConfig := L.GetGlobal("CylinderConfig")

	co, _ := L.NewThread()
	defer func() {
		fmt.Println("defer 2")
		//cocancel()
	}()
	f_startWork := L.GetGlobal("StartWork").(*lua.LFunction)

	L.SetTop(top)

	needSleep := false
	inForCnt := 0
	for {
		if needSleep {
			time.Sleep(time.Millisecond * 10)
		}
		inForCnt++
		fmt.Println("In For Cnt:", inForCnt)
		_, err, values := L.Resume(co, f_startWork) // err is nil
		if err != nil {
			fmt.Println("something wrong:", err.Error())
			break
		}
		fmt.Println("value cnt:", len(values))
		switch lua.LVAsNumber(values[1]) {
		case 0:
			if lua.LVAsNumber(values[2]) == 1 {
				needSleep = true
			} else {
				fmt.Println("over~")
				return
			}
		case 1:
			fmt.Println("Call Plc Release")
		case 3:
			fmt.Println("Call Plc Pull")
		default:
			fmt.Printf("unknown value: %v \r\n", values[1])
			needSleep = true
		}
		if !lua.LVAsBool(values[0]) {
			needSleep = true
		}
	}

}
func main2() {
	fmt.Println("test")

	L := lua.NewState()
	defer L.Close()
	/*
		if err := L.DoString(`print("hello")`); err != nil {
			panic(err)
		}
	*/
	if err := L.DoFile("f_startWork.lua"); err != nil {
		panic(err)
	}

	fmt.Println(L.GetTop())

	f_startWork := lua.P{
		Fn:   L.GetGlobal("StartWork"),
		NRet: 1,
	}

	fmt.Println(L.GetTop())

	f_stopWork := lua.P{
		Fn:   L.GetGlobal("StopWork"),
		NRet: 1,
	}

	cylinderConfig := L.GetGlobal("CylinderConfig")

	cylinders := L.GetField(cylinderConfig, "Cylinder").(*lua.LTable)

	cylinderCnt := L.ObjLen(cylinders)

	for i := 0; i < cylinderCnt; i++ {
		cylinder := L.RawGetInt(cylinders, i+1)
		L.SetField(cylinder, "Pos", lua.LNumber(i))

		//cylinderObj := L.GetField(cylinder, "Name")
		//fmt.Println(cylinderObj.String())
	}
	fmt.Println(L.GetTop())

	//f_startWork := L.GetGlobal("StartWork")

	if err := L.CallByParam(f_startWork); err != nil {
		fmt.Println("StartWork call failed")
	}
	L.Pop(1)
	fmt.Println(L.GetTop())
	if err := L.CallByParam(f_stopWork); err != nil {
		fmt.Println("StopWork call failed")
	}
	L.Pop(1)
	fmt.Println(L.GetTop())

	for i := 0; i < cylinderCnt; i++ {
		cylinder := L.RawGetInt(cylinders, i+1)
		L.SetField(cylinder, "Pos", lua.LNumber(i+1))
	}
	fmt.Println(L.GetTop())
	if err := L.CallByParam(f_startWork); err != nil {
		fmt.Println("StartWork call failed")
	}
	L.Pop(1)
	fmt.Println(L.GetTop())
	if err := L.CallByParam(f_stopWork); err != nil {
		fmt.Println("StopWork call failed")
	}
	L.Pop(1)

	fmt.Println(L.GetTop())

}
