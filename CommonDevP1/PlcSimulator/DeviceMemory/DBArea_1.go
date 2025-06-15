package DeviceMemory

import (
	"bytes"
	"encoding/binary"
)

/*
type DBCmdParam struct {
	CmdCode uint32
	Param1  int32
	Param2  int32
	Param3  int32
	Param4  int32
	Param5  float32
	Param6  float32
	Param7  float32
	Param8  float32
}
*/
type DBCmdParam struct {
	CmdCode uint16
	Param1  uint16
	Param2  int32
	Param3  int32
	Param4  float32
	Param5  float32
	Param6  float32
	//Param7  float32
}

func (this *DBCmdParam) TurnBytes(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, this.CmdCode)
	binary.Write(buf, binary.BigEndian, this.Param1)
	binary.Write(buf, binary.BigEndian, this.Param2)
	binary.Write(buf, binary.BigEndian, this.Param3)
	binary.Write(buf, binary.BigEndian, this.Param4)
	binary.Write(buf, binary.BigEndian, this.Param5)
	binary.Write(buf, binary.BigEndian, this.Param6)
}
func (this *DBCmdParam) Clear() {
	this.CmdCode = 0
	this.Param1 = 0
	this.Param2 = 0
	this.Param3 = 0
	this.Param4 = 0
	this.Param5 = 0
	this.Param6 = 0
}

func (this *DBCmdParam) Size() int {
	return 4 + 4*2 + 4*3
}

func (this *DBCmdParam) CopyIn(param *DBCmdParam) {
	this.CmdCode = param.CmdCode
	this.Param1 = param.Param1
	this.Param2 = param.Param2
	this.Param3 = param.Param3
	this.Param4 = param.Param4
	this.Param5 = param.Param5
	this.Param6 = param.Param6
}

//  command area
type DBA_1 struct {
	CmdCnt   uint32
	CmdParam [CylinderCnt]DBCmdParam
	CmdIdx   uint32
}

func (this *DBA_1) Area() int {
	return 1
}
