package structs

import (
	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	cc := NewPlcCmdContainer(3, 8, 0)

	cmds := tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  1,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  2,
	})
	idx, err := cc.PackIn(cmds)
	fmt.Printf("1 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  3,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  4,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("2 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 4,
		Param1:  1,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 4,
		Param1:  2,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("3 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  3,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  4,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("4 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 4,
		Param1:  5,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 4,
		Param1:  6,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("5 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 9,
		Param1:  0,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 9,
		Param1:  0,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("6 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 10,
		Param1:  0,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 10,
		Param1:  0,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("7 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  1,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 4,
		Param1:  1,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("8 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 10,
		Param1:  0,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 10,
		Param1:  0,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("9 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  1,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  2,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("10 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  1,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  2,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("11 idx: %v, err:%v \r\n", idx, err)

	cmds = tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  1,
	})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{
		CmdCode: 3,
		Param1:  2,
	})
	idx, err = cc.PackIn(cmds)
	fmt.Printf("12 idx: %v, err:%v \r\n", idx, err)
}

func Test2(t *testing.T) {
	cmds := tmpl.New_DBA_Command(8)
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{})
	cmds.AddCmd(&tmpl.DBA_Sub_CmdParam{})
	vcmds := cmds.GetCmds()
	tt(vcmds)
	fmt.Printf("%+v \r\n", vcmds[1])
}

func tt(cmds []tmpl.DBA_Sub_CmdParam) {
	cmds[1].CmdCode = 2
}

func Test3(t *testing.T) {
	var ch2 chan bool
	if ch2 == nil {
		fmt.Println("ch2 is nil")
	}
	ch := make(chan bool)
	go func() {
		time.Sleep(time.Second)
		ch <- true
	}()

	select {
	case <-ch:
		fmt.Println("???")
	case <-time.After(time.Millisecond * 500):
		fmt.Println(<-ch)
	}

}
