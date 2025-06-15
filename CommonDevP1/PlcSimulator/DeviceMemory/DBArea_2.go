package DeviceMemory

import (
	"bytes"
	"encoding/binary"
)

type CylinderDatas struct {
	Alarm            uint16
	RealForce        float32
	MaxForce         float32
	RealDisplacement float32
	DriveDegree      float32 //  电缸的情况是转速,气缸的情况是气阀阀值
	// 还差个报警信息
}

func (this *CylinderDatas) Retrive(buf *bytes.Buffer) {
	binary.Read(buf, binary.BigEndian, &this.Alarm)
	binary.Read(buf, binary.BigEndian, &this.RealForce)
	binary.Read(buf, binary.BigEndian, &this.MaxForce)
	binary.Read(buf, binary.BigEndian, &this.RealDisplacement)
	binary.Read(buf, binary.BigEndian, &this.DriveDegree)
}

func (this *CylinderDatas) Size() int {
	return 2 + 4*4
}

//  real data area
type DBA_2 struct {
	size          int
	CmdIdx        uint32
	Alarm         uint16
	CylinderDatas [CylinderCnt]CylinderDatas
}

func (this *DBA_2) Size() int {
	if this.size == 0 {
		this.size = 4 + 2 + CylinderCnt*this.CylinderDatas[0].Size()
	}
	return this.size
}

func (this *DBA_2) Area() int {
	return 2
}

func (this *DBA_2) Retrive(buf *bytes.Buffer) {
	binary.Read(buf, binary.BigEndian, &this.CmdIdx)
	binary.Read(buf, binary.BigEndian, &this.Alarm)
	for i := 0; i < len(this.CylinderDatas); i++ {
		this.CylinderDatas[i].Retrive(buf)
	}
}
