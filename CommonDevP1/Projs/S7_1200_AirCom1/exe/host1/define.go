package main

import (
	"CommonDevP1/PlcSimulator/DeviceMemory"
	"bytes"
	"encoding/binary"

	//"MyToolBox/Thread/SyncAppData"
	"container/list"
	"fmt"
)

const (
	Topic_Main = iota
	Topic_Conn
	Topic_Cmd
	Topic_Lua
)

const Sync_Data_Len = 8
const (
	Sync_ReadMode = iota
	Sync_Cmd
)

const ReadMode_Status int = 0
const ReadMode_Menu int = 1

type ICmdPack interface {
	Poses() []int
	ContainsPoses([]int) bool
	PackIn(ICmdPack) bool
	CmdCnt() int
	CmdPlcIdx() int
	Pack() []byte
}

const __MaxCmdPackCnt = 8

type CmdList struct {
	cmdPacks list.List
}

var __Err_OverMaxCnt = fmt.Errorf("Too many cmd packs in container")

func (this *CmdList) PackIn(cmdpack ICmdPack) error {
	//var err error

	if bak := this.cmdPacks.Back(); bak != nil {
		p := bak.Value.(ICmdPack)
		if !p.PackIn(cmdpack) {
			if this.cmdPacks.Len() >= __MaxCmdPackCnt {
				return __Err_OverMaxCnt
			}
		} else {
			return nil
		}
	}
	this.cmdPacks.PushBack(cmdpack)
	return nil
}

func (this *CmdList) GetCmdPack() []byte {
	if first := this.cmdPacks.Front(); first != nil {
		itm := this.cmdPacks.Remove(first)
		p := itm.(ICmdPack)
		return p.Pack()
	}
	return nil
}

/*
type Cmd struct {
	CmdCode uint16
	Param1  uint16
	Param2  uint32
	Param3  uint32
	Param4  float32
	Param5  float32
	Param6  float32
}
func (this *Cmd) GetPos() int {
	switch this.CmdCode {
	case 41:
		return -1
	default:
		return int(this.CmdCode)
	}
}

//func (this *Cmd) CheckPos()

*/
/*
type CmdPack struct {
	cmdCnt int
	cmds   []Cmd
}

func (this *CmdPack) Poses() []int {
	var poses []int
	for i := 0; i < this.cmdCnt; i++ {
		if pos := cmds[i].GetPos(); pos != -1 {
			poses = append(poses, i)
		}
	}
	return poses
}
func (this *CmdPack) ContainsPoses(poses []int) bool {

	// for i := 0; i < this.cmdCnt; i++ {
	// 	if pos := cmds[i].GetPos(); pos != -1 {

	// 	}
	// }
	return false
}

func (this *CmdPack) PackIn(cmdPack ICmdPack) bool {
	if cmdPack.CmdCnt()+this.cmdCnt > __MaxCmdPackCnt {
		return false
	}
	if !this.ContainsPoses(cmdPack.Poses()) {
		this.cmdPack
	}
	return false
}
func (this *CmdPack) CmdCnt() int {
	return this.cmdCnt
}
func (this *CmdPack) CmdPlcIdx() int {
	return 0
}
func (this *CmdPack) Pack() []byte {
	return nil
}
*/
type CmdPack_Simple struct {
	size   int
	cnt    uint16
	cmds   [8]DeviceMemory.DBCmdParam
	cmdIdx uint32
}

func (this *CmdPack_Simple) TurnBytes(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, this.cnt)
	for i := 0; i < len(this.cmds); i++ {
		this.cmds[i].TurnBytes(buf)
	}
	binary.Write(buf, binary.BigEndian, this.cmdIdx)
}

func (this *CmdPack_Simple) Clear() {
	this.cnt = 0
}

func (this *CmdPack_Simple) Size() int {
	if this.size == 0 {
		this.size = 2 + 8*this.cmds[0].Size() + 4
	}
	return this.size
}

func (this *CmdPack_Simple) Set(cmds []*DeviceMemory.DBCmdParam) uint32 {
	i := 0
	for ; i < len(cmds); i++ {
		//copy(this.cmds[i], cmds[i])
		this.cmds[i].CopyIn(cmds[i])
	}
	for ; i < len(this.cmds); i++ {
		this.cmds[i].Clear()
	}
	this.cnt = uint16(len(cmds))
	this.cmdIdx++
	return this.cmdIdx
}

func (this *CmdPack_Simple) CmdCount() int {
	return int(this.cnt)
}

type Sync_CmdPack struct {
	cmds []*DeviceMemory.DBCmdParam
}
