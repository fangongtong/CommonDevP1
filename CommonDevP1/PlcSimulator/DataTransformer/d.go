package DataTransformer

import (
	dm "CommonDevP1/PlcSimulator/DeviceMemory"

	"github.com/golang/protobuf/proto"
)

//  0:real data request
//  1:cmd
//  2:real data return

var _Buf_RealDt_Req []byte = []byte{0}

func PackData(cmd int, dt dm.IDevMemory) []byte {

	switch cmd {
	case 0:
		return _Buf_RealDt_Req[:]
	case 2:
		{
			org := dt.(*dm.DBA_2)
			trg := &MsgRealDt{}
			trg.CmdIdx = org.CmdIdx
			for _, v := range org.CylinderDatas {
				trg.RealDts = append(trg.RealDts, &StRealDt{
					RealForce: v.RealForce,
				})
			}

			protoDt, _ := proto.Marshal(trg)
			rtn := make([]byte, len(protoDt)+1)
			rtn[0] = 2
			copy(rtn[1:], protoDt)
			return rtn
		}
	default:
		return nil
	}

	/*
		switch dt.(type) {
		case *dm.DBA_2:
		default:
			return nil
		}
	*/
}
func PackCmd(cmd *MsgCmd) []byte {
	protoDt, _ := proto.Marshal(cmd)
	rtn := make([]byte, len(protoDt)+1)
	rtn[0] = 1
	copy(rtn[1:], protoDt)
	return rtn
}

func UnPackData1(dt []byte) (int, dm.IDevMemory) {
	//fmt.Printf("Unpackdata: %X \r\n", dt)
	switch dt[0] {
	case 0:
		{
			return 0, nil
		}
	case 1:
		{
			org := &MsgCmd{}
			if proto.Unmarshal(dt[1:], org) != nil {
				return -1, nil
			}
			trg := &dm.DBA_1{}
			trg.CmdCnt = org.CmdCnt
			trg.CmdIdx = org.CmdIdx
			for i, v := range org.Cmds {
				trg.CmdParam[i].CmdCode = v.CmdCode
				trg.CmdParam[i].Param1 = v.Param1
				trg.CmdParam[i].Param2 = v.Param2
				trg.CmdParam[i].Param3 = v.Param3
				trg.CmdParam[i].Param4 = v.Param4
				trg.CmdParam[i].Param5 = v.Param5
				trg.CmdParam[i].Param6 = v.Param6
				trg.CmdParam[i].Param7 = v.Param7
				trg.CmdParam[i].Param8 = v.Param8
			}
			return 1, trg
		}
	case 2:
		{
			org := &MsgRealDt{}
			if proto.Unmarshal(dt[1:], org) != nil {
				return -1, nil
			}
			trg := &dm.DBA_2{}
			trg.CmdIdx = org.CmdIdx
			for i, v := range org.RealDts {
				trg.CylinderDatas[i].RealForce = v.RealForce
			}
			return 2, trg
		}
	default:
		{
			return -1, nil
		}
	}
}
func UnPackData2(dt []byte) int {
	switch dt[0] {
	case 0:
		{
			org := &MsgCmd{}
			if proto.Unmarshal(dt[1:], org) != nil {
				return 0
			}
		}
	}
	return -1
}
