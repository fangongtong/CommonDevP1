package bs

import (
	dm "CommonDevP1/PlcSimulator/DeviceMemory"
	"time"
)

type PlcBs struct {
	dba1 dm.DBA_1
	dba2 dm.DBA_2
}

func NewBs() *PlcBs {
	return &PlcBs{}
}

func (this *PlcBs) Run() {

	var cmdParam [dm.CylinderCnt]dm.DBCmdParam
	for {
		//  一个循环周期
		time.Sleep(time.Millisecond * 5)

		if this.dba1.CmdIdx > 0 {
			cmdIdx := this.dba1.CmdIdx
			cmdCnt := this.dba1.CmdCnt
			this.dba1.CmdIdx = 0

			for i := 0; i < int(this.dba1.CmdCnt); i++ {
				cmdParam[i] = this.dba1.CmdParam[i]
			}
			this.dba2.CmdIdx = cmdIdx

			this.dealCmd(cmdParam[:cmdCnt])
		}
	}
}

func (this *PlcBs) dealCmd(cmds []dm.DBCmdParam) {
	//  执行对应指令
	//  这里要根据指令做点假数据
}

//  接收到远端的指令数据后复制到dba1里
func (this *PlcBs) SetCmd(cmds *dm.DBA_1) {
	this.dba1.CmdCnt = cmds.CmdCnt

	for i := 0; i < int(cmds.CmdCnt); i++ {
		this.dba1.CmdParam[i] = cmds.CmdParam[i]
	}

	this.dba1.CmdIdx = cmds.CmdIdx
}

//  接收到远端的实时数据请求后 把dba2数据复制出去
func (this *PlcBs) GetRealData() *dm.DBA_2 {
	return &this.dba2
}
