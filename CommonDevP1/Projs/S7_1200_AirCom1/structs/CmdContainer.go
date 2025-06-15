package structs

import (
	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"container/list"
	"fmt"
)

func NewPlcCmdContainer(maxCmdPackCnt, maxCmdCntInPack int, startCmdIdx uint32) *PlcCmdContainer {
	return &PlcCmdContainer{
		startCmdIdx: startCmdIdx,
		maxCmds:     maxCmdCntInPack,
		maxPacks:    maxCmdPackCnt,
	}
}

type PlcCmdContainer struct {
	startCmdIdx uint32
	maxPacks    int
	maxCmds     int
	cmdPacks    list.List
}

var __Err_OverMaxCnt = fmt.Errorf("Too many cmd packs in container")

func (this *PlcCmdContainer) PrintPacks() {
	fmt.Println("------", this.maxPacks)
	for it := this.cmdPacks.Front(); it != nil; it = it.Next() {
		fmt.Println("------")
		cmdPak := it.Value.(*tmpl.DBA_Command)
		for i := 0; i < cmdPak.CmdCount(); i++ {
			fmt.Printf("%+v \r\n", cmdPak.CmdParamAry[i])
		}
	}
}

func (this *PlcCmdContainer) PackIn(outCmds *tmpl.DBA_Command) (uint32, error) { //cmdpack ICmdPack

	if itm := this.cmdPacks.Back(); itm != nil {
		inCmds := itm.Value.(*tmpl.DBA_Command)
		if inCmds.CmdCount()+outCmds.CmdCount() <= this.maxCmds {
			if checkPosInCmds(inCmds, outCmds) {
				inCmds.AddCmds(outCmds.GetCmds())
				return inCmds.CmdIdx, nil
			}
		}
	}

	if this.cmdPacks.Len() == this.maxPacks {
		this.PrintPacks()
		return 0, __Err_OverMaxCnt
	}

	newCmdPack := tmpl.New_DBA_Command(this.maxCmds)
	newCmdPack.AddCmds(outCmds.GetCmds())
	this.startCmdIdx = this.startCmdIdx + 1
	newCmdPack.CmdIdx = this.startCmdIdx
	this.cmdPacks.PushBack(newCmdPack)

	return this.startCmdIdx, nil
}
func (this *PlcCmdContainer) PackIn2(cmds []tmpl.DBA_Sub_CmdParam) (uint32, error) { //cmdpack ICmdPack

	if itm := this.cmdPacks.Back(); itm != nil {
		inCmds := itm.Value.(*tmpl.DBA_Command)
		if inCmds.CmdCount()+len(cmds) <= this.maxCmds {
			if checkPosInCmds2(inCmds, cmds) {
				inCmds.AddCmds(cmds)
				return inCmds.CmdIdx, nil
			}
		}
	}

	if this.cmdPacks.Len() == this.maxPacks {
		this.PrintPacks()
		return 0, __Err_OverMaxCnt
	}

	newCmdPack := tmpl.New_DBA_Command(this.maxCmds)
	newCmdPack.AddCmds(cmds)
	this.startCmdIdx = this.startCmdIdx + 1
	newCmdPack.CmdIdx = this.startCmdIdx
	this.cmdPacks.PushBack(newCmdPack)

	return this.startCmdIdx, nil
}

func (this *PlcCmdContainer) GetCmdPack() *tmpl.DBA_Command {
	if first := this.cmdPacks.Front(); first != nil {
		itm := this.cmdPacks.Remove(first)
		return itm.(*tmpl.DBA_Command)
	}
	return nil
}

func checkPosInCmds(inCmdsPack, outCmdsPack *tmpl.DBA_Command) bool {
	inCmds := inCmdsPack.GetCmds()
	outCmds := outCmdsPack.GetCmds()
	for i := 0; i < len(inCmds); i++ {
		for j := 0; j < len(outCmds); j++ {
			if inCmds[i].Param1 > 0 && inCmds[i].Param1 == outCmds[j].Param1 {
				return false
			}
		}
	}
	return true
}

func checkPosInCmds2(inCmdsPack *tmpl.DBA_Command, outCmds []tmpl.DBA_Sub_CmdParam) bool {
	inCmds := inCmdsPack.GetCmds()
	for i := 0; i < len(inCmds); i++ {
		for j := 0; j < len(outCmds); j++ {
			if inCmds[i].Param1 > 0 && inCmds[i].Param1 == outCmds[j].Param1 {
				return false
			}
		}
	}
	return true
}
