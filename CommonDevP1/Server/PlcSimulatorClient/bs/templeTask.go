package bs

import (
	"CommonDevP1/Server/Core"
	"fmt"

	"math/rand"
	"time"

	"github.com/yuin/gopher-lua"
)

type BType Core.BaseTask

type MyTask struct {
	*BType
}

func (this *MyTask) Base() *Core.BaseTask {
	return (*Core.BaseTask)(this.BType)
}

func (this *MyTask) SetBase(b *Core.BaseTask) {
	this.BType = (*BType)(b)
}

func (this *MyTask) Start() {
	//luaPlcStateType := "plcState"
	fmt.Println("Task.Start 0")
	this.OutputErr = nil
	L := lua.NewState()
	defer L.Close()

	fmt.Println("Task.Start 1")

	L.SetGlobal("GoSleep", L.NewFunction(Core.Api2L_Sleep))

	if err := L.DoFile("./lua/1AC_test.lua"); err != nil {
		this.OutputErr = err
		return
	}

	fmt.Println("Task.Start ")

	f_initCylinder := lua.P{
		Fn:   L.GetGlobal("InitCylinder"),
		NRet: 0,
	}
	//fmt.Println(L.GetTop())

	f_startWork := lua.P{
		Fn:   L.GetGlobal("StartWork"),
		NRet: 1,
	}

	// cylinderConfig
	cylinderConfig := L.GetGlobal("CylinderConfig")
	cylinders := L.GetField(cylinderConfig, "Cylinder").(*lua.LTable)
	cylinderCnt := L.ObjLen(cylinders)

	//  加入lua协程
	co, _ := L.NewThread() /* create a new thread */
	//fn := co.GetGlobal("DealReal").(*lua.LFunction)
	fn := L.GetGlobal("DealReal").(*lua.LFunction) /* get function from lua */

	L.Mark = "mainState"
	co.Mark = "coState"
	/*
		ud := L.NewUserData()
		ud.Value = &dm.DBA_2{}
		L.SetMetatable(ud, L.GetTypeMetatable(luaPlcStateType))
		L.Push(ud)
	*/

	for i := 0; i < cylinderCnt; i++ {
		cylinder := L.RawGetInt(cylinders, i+1)
		L.SetField(cylinder, "Pos", lua.LNumber(i))

		//cylinderObj := L.GetField(cylinder, "Name")
		//fmt.Println(cylinderObj.String())
	}

	//---- init -----
	if err := L.CallByParam(f_initCylinder); err != nil {
		//fmt.Println("StartWork call failed")
		this.OutputErr = err
		return
	}

	//---- go coroutine ----
	exCtrl := &struct {
		Ex     bool
		ChQuit chan bool
	}{}
	exCtrl.ChQuit = make(chan bool)

	defer func() {
		exCtrl.Ex = true
		fmt.Println("waiting...")
		<-exCtrl.ChQuit
	}()

	go func(l *lua.LState) {

		defer func() {
			exCtrl.ChQuit <- true
		}()

		p2 := 50
		var coroutineArgs []lua.LValue = make([]lua.LValue, 8)
		coroutineArgs[0], coroutineArgs[1], coroutineArgs[2], coroutineArgs[3] = lua.LBool(exCtrl.Ex), lua.LNumber(1000), lua.LNumber(p2), lua.LNumber(55)

		tmpRealForce := 0

		defer func(t time.Time) {
			fmt.Println("total used:", time.Since(t).Milliseconds())
		}(time.Now())
		max := 9000
		for i := 0; i < 9000; i++ {

			if i == max-1 {
				exCtrl.Ex = true
			}

			//time.Sleep(time.Millisecond * 20)
			p2++
			coroutineArgs[0], coroutineArgs[2] = lua.LBool(exCtrl.Ex), lua.LNumber(p2)

			if tmpRealForce > 200 {
				tmpRealForce = 0
			}

			tmpRealForce++

			coroutineArgs[4], coroutineArgs[5] = lua.LNumber(tmpRealForce+rand.Intn(50)), lua.LNumber(0)
			coroutineArgs[6], coroutineArgs[7] = lua.LNumber(tmpRealForce), lua.LNumber(0)
			fmt.Println("do resume")
			st, err, values := l.Resume(co, fn, coroutineArgs...)
			if st == lua.ResumeError {
				fmt.Println("yield break(error):", err.Error())
				//fmt.Println(err)
				break
			}

			if false {
				for i, lv := range values {
					fmt.Printf("%v : %+v\n", i, lv)
				}
			}

			if st == lua.ResumeOK {
				fmt.Println("yield break(ok)")
				break
			}
		}
		fmt.Println("Func exit")
	}(L)

	//---- start -----
	if err := L.CallByParam(f_startWork); err != nil {
		fmt.Println("StartWork call failed")
		this.OutputErr = err
		return
	}
	fmt.Println("StartWork call exit")

	//_Host.Send()
}

//---------------------------

type MyTask2 struct {
	*BType
}

func (this *MyTask2) Base() *Core.BaseTask {
	return (*Core.BaseTask)(this.BType)
}

func (this *MyTask2) SetBase(b *Core.BaseTask) {
	this.BType = (*BType)(b)
}

func (this *MyTask2) Start() {
	//luaPlcStateType := "plcState"
	fmt.Println("Task.Start 0")
	this.OutputErr = nil
	L := lua.NewState()
	defer L.Close()

	fmt.Println("Task.Start 1")

	L.SetGlobal("GoSleep", L.NewFunction(Core.Api2L_Sleep))

	if err := L.DoFile("./lua/1AC_test2.lua"); err != nil {
		this.OutputErr = err
		return
	}

	fmt.Println("Task.Start ")

	f_initCylinder := lua.P{
		Fn:   L.GetGlobal("InitCylinder"),
		NRet: 0,
	}
	//fmt.Println(L.GetTop())

	f_startWork := lua.P{
		Fn:   L.GetGlobal("StartWork"),
		NRet: 1,
	}

	// cylinderConfig
	cylinderConfig := L.GetGlobal("CylinderConfig")
	cylinders := L.GetField(cylinderConfig, "Cylinder").(*lua.LTable)
	cylinderCnt := L.ObjLen(cylinders)

	// common
	/*
		realDtTranser := struct {
			Break
			Com1
			Com2
			Com3
		}
	*/
	com := L.GetGlobal("Common")
	//fmt.Printf("t: %+v\r\n", t)
	realDtTranser := L.GetField(com, "RealDt").(*lua.LTable)
	realDtTranserCR := L.GetField(realDtTranser, "CR").(*lua.LTable)
	//realDtTranser := L.GetGlobal("Common.RealDt").(*lua.LTable)
	//return

	//  加入lua协程
	co, _ := L.NewThread() /* create a new thread */
	//fn := co.GetGlobal("DealReal").(*lua.LFunction)
	fn := L.GetGlobal("DealReal").(*lua.LFunction) /* get function from lua */

	L.Mark = "mainState"
	co.Mark = "coState"
	/*
		ud := L.NewUserData()
		ud.Value = &dm.DBA_2{}
		L.SetMetatable(ud, L.GetTypeMetatable(luaPlcStateType))
		L.Push(ud)
	*/

	for i := 0; i < cylinderCnt; i++ {
		cylinder := L.RawGetInt(cylinders, i+1)
		L.SetField(cylinder, "Pos", lua.LNumber(i))

		//cylinderObj := L.GetField(cylinder, "Name")
		//fmt.Println(cylinderObj.String())
	}

	//---- init -----
	if err := L.CallByParam(f_initCylinder); err != nil {
		//fmt.Println("StartWork call failed")
		this.OutputErr = err
		return
	}

	//---- go coroutine ----
	exCtrl := &struct {
		Ex     bool
		ChQuit chan bool
	}{}
	exCtrl.ChQuit = make(chan bool)

	defer func() {
		exCtrl.Ex = true
		fmt.Println("waiting...")
		<-exCtrl.ChQuit
	}()

	go func(l *lua.LState) {

		fmt.Println("go 1")
		defer func() {
			exCtrl.ChQuit <- true
		}()

		p2 := 50
		//coroutineArgs[0], coroutineArgs[1], coroutineArgs[2], coroutineArgs[3] = lua.LBool(exCtrl.Ex), lua.LNumber(1000),

		tmpRealForce := 0

		//realDtTranserCR := realDtTranser.RawGet(lua.LString("CR")).(*lua.LTable)

		fmt.Println("go 2")
		cr1 := realDtTranserCR.RawGetInt(1).(*lua.LTable)
		cr2 := realDtTranserCR.RawGetInt(2).(*lua.LTable)
		//cr3 := realDtTranser.RawGetInt(3).(*lua.LTable)
		//cr4 := realDtTranser.RawGetInt(4).(*lua.LTable)
		//cr5 := realDtTranser.RawGetInt(5).(*lua.LTable)
		//cr6 := realDtTranser.RawGetInt(6).(*lua.LTable)
		//cr7 := realDtTranser.RawGetInt(7).(*lua.LTable)
		//cr8 := realDtTranser.RawGetInt(8).(*lua.LTable)

		fmt.Println("go 3")

		defer func(t time.Time) {
			fmt.Println("total used:", time.Since(t).Milliseconds())
		}(time.Now())
		max := 9000
		for i := 0; i < 9000; i++ {

			if i == max-1 {
				exCtrl.Ex = true
			}
			time.Sleep(time.Millisecond * 20)
			p2++
			if tmpRealForce > 200 {
				tmpRealForce = 0
			}

			tmpRealForce++

			realDtTranser.RawSetString("Break", lua.LBool(exCtrl.Ex))
			realDtTranser.RawSetString("Com1", lua.LNumber(1000))
			realDtTranser.RawSetString("Com2", lua.LNumber(p2))
			realDtTranser.RawSetString("Com3", lua.LNumber(55))
			cr1.RawSetString("RealForce", lua.LNumber(tmpRealForce+rand.Intn(50)))
			cr1.RawSetString("DriveDegree", lua.LNumber(0))

			cr2.RawSetString("RealForce", lua.LNumber(tmpRealForce))
			cr2.RawSetString("DriveDegree", lua.LNumber(0))

			//realDtTranserCR.RawSetInt(1, lua.LNumber(tmpRealForce+rand.Intn(50)))
			//realDtTranserCR.RawSetInt(1, lua.LNumber(tmpRealForce+rand.Intn(50)))
			//realDtTranserCR.Metatable

			L.SetField(com, "RealDt", realDtTranser)
			fmt.Println("do resume")
			st, err, values := l.Resume(co, fn)
			if st == lua.ResumeError {
				fmt.Println("yield break(error):", err.Error())
				//fmt.Println(err)
				break
			}

			if false {
				for i, lv := range values {
					fmt.Printf("%v : %+v\n", i, lv)
				}
			}

			if st == lua.ResumeOK {
				fmt.Println("yield break(ok)")
				break
			}
		}
		fmt.Println("Func exit")
	}(L)

	//---- start -----
	if err := L.CallByParam(f_startWork); err != nil {
		fmt.Println("StartWork call failed")
		this.OutputErr = err
		return
	}
	fmt.Println("StartWork call exit")

	//_Host.Send()
}
