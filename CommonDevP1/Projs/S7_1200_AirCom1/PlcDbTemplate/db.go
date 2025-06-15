package PlcDbTemplate

import (
	"bytes"
	"encoding/binary"
)

//---------------------------------

type IPlcDB interface {
	Marshal(buf *bytes.Buffer, order binary.ByteOrder)
	Unmarshal(buf *bytes.Buffer, order binary.ByteOrder) error
	Size() int
}

type DBA_Sub_CmdParam struct {
	CmdCode uint16  `json:"CmdCode"`
	Param1  uint16  `json:"Param1"`
	Param2  int32   `json:"Param2"`
	Param3  int32   `json:"Param3"`
	Param4  float32 `json:"Param4"`
	Param5  float32 `json:"Param5"`
	Param6  float32 `json:"Param6"`
}

func (this *DBA_Sub_CmdParam) Marshal(buf *bytes.Buffer, order binary.ByteOrder) {
	binary.Write(buf, order, this.CmdCode)
	binary.Write(buf, order, this.Param1)
	binary.Write(buf, order, this.Param2)
	binary.Write(buf, order, this.Param3)
	binary.Write(buf, order, this.Param4)
	binary.Write(buf, order, this.Param5)
	binary.Write(buf, order, this.Param6)
}

func (this *DBA_Sub_CmdParam) Clear() {
	this.CmdCode = 0
	this.Param1 = 0
	this.Param2 = 0
	this.Param3 = 0
	this.Param4 = 0
	this.Param5 = 0
	this.Param6 = 0
}

func (this *DBA_Sub_CmdParam) Size() int {
	return 4 + 4*2 + 4*3
}

func (this *DBA_Sub_CmdParam) CopyIn(param *DBA_Sub_CmdParam) {
	this.CmdCode = param.CmdCode
	this.Param1 = param.Param1
	this.Param2 = param.Param2
	this.Param3 = param.Param3
	this.Param4 = param.Param4
	this.Param5 = param.Param5
	this.Param6 = param.Param6
}

func New_DBA_Command(cmdCnt int) *DBA_Command {
	return &DBA_Command{
		CmdParamAry: make([]DBA_Sub_CmdParam, cmdCnt),
	}
}

//  command area
type DBA_Command struct {
	size        int
	cmdCnt      uint16
	CmdParamAry []DBA_Sub_CmdParam
	CmdIdx      uint32
}

func (this *DBA_Command) Marshal(buf *bytes.Buffer, order binary.ByteOrder) {
	binary.Write(buf, order, this.cmdCnt)
	for i := 0; i < len(this.CmdParamAry); i++ {
		this.CmdParamAry[i].Marshal(buf, order)
	}

	binary.Write(buf, order, this.CmdIdx)
}
func (this *DBA_Command) Size() int {
	if this.size == 0 {
		this.size = 2 + len(this.CmdParamAry)*this.CmdParamAry[0].Size() + 4
	}
	return this.size
}

func (this *DBA_Command) AddCmd(cmd *DBA_Sub_CmdParam) bool {
	if int(this.cmdCnt)+1 > len(this.CmdParamAry) {
		return false
	} else {
		this.CmdParamAry[this.cmdCnt].CopyIn(cmd)
		this.cmdCnt++
		return true
	}
}
func (this *DBA_Command) AddCmd2(cmd DBA_Sub_CmdParam) bool {
	return this.AddCmd(&cmd)
}
func (this *DBA_Command) AddCmds(cmds []DBA_Sub_CmdParam) bool {
	if int(this.cmdCnt)+len(cmds) > len(this.CmdParamAry) {
		return false
	} else {
		for i := 0; i < len(cmds); i++ {
			this.CmdParamAry[i+int(this.cmdCnt)].CopyIn(&cmds[i])
		}
		this.cmdCnt = this.cmdCnt + uint16(len(cmds))
		return true
	}
}
func (this *DBA_Command) GetCmds() []DBA_Sub_CmdParam {
	return this.CmdParamAry[:this.cmdCnt]
}
func (this *DBA_Command) Clear() {
	for i := 0; i < len(this.CmdParamAry); i++ {
		this.CmdParamAry[i].Clear()
	}
	this.cmdCnt = 0
	this.CmdIdx = 0
}

// func (this *DBA_Command) Add(cmds []*DBA_Sub_CmdParam) (uint32, bool) {
// 	if int(this.cmdCnt)+len(cmds) > len(this.CmdParamAry) {
// 		return 0, false
// 	}

// 	i := 0 //int(this.cmdCnt)
// 	for ; i < len(cmds); i++ {
// 		//copy(this.cmds[i], cmds[i])
// 		this.CmdParamAry[i].CopyIn(cmds[i])
// 	}
// 	for ; i < len(this.CmdParamAry); i++ {
// 		this.CmdParamAry[i].Clear()
// 	}
// 	this.cmdCnt = uint16(len(cmds))
// 	this.cmdIdx++
// 	return this.cmdIdx, true
// }

func (this *DBA_Command) CmdCount() int {
	return int(this.cmdCnt)
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
func (this *DBA_Sub_CylinderDatas) Size() int {
	return 4*4 + 2
}

//  real data area
func New_DBA_Status(posCnt int) *DBA_Status {
	return &DBA_Status{
		CylinderDataAry: make([]DBA_Sub_CylinderDatas, posCnt),
	}
}

type DBA_Status struct {
	size            int
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

func (this *DBA_Status) Size() int {
	if this.size == 0 {
		this.size = 4 + 2 + len(this.CylinderDataAry)*this.CylinderDataAry[0].Size()
	}
	return this.size
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

type DBA_Sub_MenuParam struct {
	OrderMaxRunSeconds            uint16
	ForceSensorPushOverloadRate   float32
	ForceSensorPullOverloadRate   float32
	SampleForceOverloadRate       float32
	ForceFallProtectRate          float32
	MaxSecondsBeforeReachForce    uint32
	ProportionalValve             float32
	TargetForceUpperDeviationRate float32
	TargetForceLowerDeviationRate float32
	ForceFactorB                  float32
	ForceSensorCapacity           float32
}

func (this *DBA_Sub_MenuParam) CopyIn(param *DBA_Sub_MenuParam) {
	this.OrderMaxRunSeconds = param.OrderMaxRunSeconds
	this.ForceSensorPushOverloadRate = param.ForceSensorPushOverloadRate
	this.ForceSensorPullOverloadRate = param.ForceSensorPullOverloadRate
	this.SampleForceOverloadRate = param.SampleForceOverloadRate
	this.ForceFallProtectRate = param.ForceFallProtectRate
	this.MaxSecondsBeforeReachForce = param.MaxSecondsBeforeReachForce
	this.ProportionalValve = param.ProportionalValve
	this.TargetForceUpperDeviationRate = param.TargetForceUpperDeviationRate
	this.TargetForceLowerDeviationRate = param.TargetForceLowerDeviationRate
	this.ForceFactorB = param.ForceFactorB
	this.ForceSensorCapacity = param.ForceSensorCapacity
}

func (this *DBA_Sub_MenuParam) Unmarshal(buf *bytes.Buffer, order binary.ByteOrder) error {
	var err error
	if err = binary.Read(buf, order, &this.OrderMaxRunSeconds); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.ForceSensorPushOverloadRate); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.ForceSensorPullOverloadRate); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.SampleForceOverloadRate); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.ForceFallProtectRate); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.MaxSecondsBeforeReachForce); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.ProportionalValve); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.TargetForceUpperDeviationRate); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.TargetForceLowerDeviationRate); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.ForceFactorB); err != nil {
		return err
	}
	if err = binary.Read(buf, order, &this.ForceSensorCapacity); err != nil {
		return err
	}
	return nil
}

func (this *DBA_Sub_MenuParam) Size() int {
	return 2 + 4*10
}

type DBA_MenuParam struct {
	size         int
	MenuParamAry []DBA_Sub_MenuParam
}

func (this *DBA_MenuParam) Unmarshal(buf *bytes.Buffer, order binary.ByteOrder) error {
	for i := 0; i < len(this.MenuParamAry); i++ {
		if err := this.MenuParamAry[i].Unmarshal(buf, order); err != nil {
			return err
		}
	}
	return nil
}

func (this *DBA_MenuParam) Size() int {
	if this.size == 0 {
		this.size = len(this.MenuParamAry) * this.MenuParamAry[0].Size()
	}
	return this.size
}
