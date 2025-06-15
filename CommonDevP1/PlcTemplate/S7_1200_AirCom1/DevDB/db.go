package DevDB

import (
	"bytes"
	"encoding/binary"
)

//---------------------------------

type DBA_Sub_CmdParam struct {
	CmdCode uint16
	Param1  uint16
	Param2  int32
	Param3  int32
	Param5  float32
	Param6  float32
	Param7  float32
}

func (this *DBA_Sub_CmdParam) Marshal(buf *bytes.Buffer, order binary.ByteOrder) {
	binary.Write(buf, order, this.CmdCode)
	binary.Write(buf, order, this.Param1)
	binary.Write(buf, order, this.Param2)
	binary.Write(buf, order, this.Param3)
	binary.Write(buf, order, this.Param5)
	binary.Write(buf, order, this.Param6)
	binary.Write(buf, order, this.Param7)
}

//  command area
type DBA_Command struct {
	CmdCnt      uint16
	CmdParamAry []DBA_Sub_CmdParam
	CmdIdx      uint32
}

func (this *DBA_Command) Marshal(buf *bytes.Buffer, order binary.ByteOrder) {
	binary.Write(buf, order, this.CmdCnt)
	for i := 0; i < len(this.CmdParamAry); i++ {
		this.CmdParamAry[i].Marshal(buf, order)
	}

	binary.Write(buf, order, this.CmdIdx)
}

//----------------------------------

type DBA_Sub_CylinderDatas struct {
	Alarm            uint16
	RealForce        float32
	PeekForce        float32
	RealDisplacement float32
	DriveDegree      float32 //  电缸的情况是转速,气缸的情况是气阀阀值
}

func (this *DBA_Sub_CylinderDatas) Marshal(buf *bytes.Buffer, order binary.ByteOrder) {
	binary.Write(buf, order, this.Alarm)
	binary.Write(buf, order, this.RealForce)
	binary.Write(buf, order, this.PeekForce)
	binary.Write(buf, order, this.RealDisplacement)
	binary.Write(buf, order, this.DriveDegree)
}
func (this *DBA_Sub_CylinderDatas) Unmarshal(buf *bytes.Buffer, order binary.ByteOrder) error {
	var err error
	if err = binary.Read(buf, order, &this.Alarm); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.RealForce); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.PeekForce); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.RealDisplacement); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.DriveDegree); err != nil {
		return err
	}
	return nil
}

//  real data area
type DBA_Status struct {
	CmdIdx          uint32
	DevAlarm        uint16
	CylinderDataAry []DBA_Sub_CylinderDatas
}

func (this *DBA_Status) Marshal(buf *bytes.Buffer, order binary.ByteOrder) {
	binary.Write(buf, order, this.CmdIdx)
	binary.Write(buf, order, this.DevAlarm)
	for i := 0; i < len(this.CylinderDataAry); i++ {
		this.CylinderDataAry[i].Marshal(buf, order)
	}
}
func (this *DBA_Status) Unmarshal(buf *bytes.Buffer, order binary.ByteOrder) error {
	var err error
	if err = binary.Read(buf, order, &this.CmdIdx); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.DevAlarm); err != nil {
		return err
	}
	for i := 0; i < len(this.CylinderDataAry); i++ {
		if err = this.CylinderDataAry[i].Unmarshal(buf, order); err != nil {
			return err
		}
	}
	return nil
}

// -------------------------------------------

type DBA_Sub_OtherMenu struct {
	RealForce  float32
	AdOriginal float32
}

func (this *DBA_Sub_OtherMenu) Marshal(buf *bytes.Buffer, order binary.ByteOrder) {
	binary.Write(buf, order, this.RealForce)
	binary.Write(buf, order, this.AdOriginal)
}

func (this *DBA_Sub_OtherMenu) Unmarshal(buf *bytes.Buffer, order binary.ByteOrder) error {
	var err error
	if err = binary.Read(buf, order, &this.RealForce); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.AdOriginal); err != nil {
		return err
	}
	return nil
}

type DBA_OtherMenu struct {
	LimitSwitch      byte
	Reserved1        byte
	OtherMenuDataAry []DBA_Sub_OtherMenu
}

func (this *DBA_OtherMenu) Marshal(buf *bytes.Buffer, order binary.ByteOrder) {
	buf.WriteByte(this.LimitSwitch)
	buf.WriteByte(this.Reserved1)
	for i := 0; i < len(this.OtherMenuDataAry); i++ {
		this.OtherMenuDataAry[i].Marshal(buf, order)
	}
}

func (this *DBA_OtherMenu) Unmarshal(buf *bytes.Buffer, order binary.ByteOrder) error {
	var err error
	if this.LimitSwitch, err = buf.ReadByte(); err != nil {
		return err
	}
	if this.Reserved1, err = buf.ReadByte(); err != nil {
		return err
	}
	for i := 0; i < len(this.OtherMenuDataAry); i++ {
		if err = this.OtherMenuDataAry[i].Unmarshal(buf, order); err != nil {
			return err
		}
	}
	return nil
}
