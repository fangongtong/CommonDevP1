package structs

import (
	"fmt"
)

func New_PosMgr(posCnt, limitSwCnt int) *PosResMgr {
	return &PosResMgr{
		posRes:     make([]uint32, posCnt),
		limitSwRes: make([]uint32, limitSwCnt),
	}
}

type PosResMgr struct {
	posRes     []uint32
	limitSwRes []uint32
}

var __Err_PosOccupied = fmt.Errorf("pos is not exists or occupied")
var __Err_LmSwOccupied = fmt.Errorf("LimitSwitch is not exists or occupied")

func (this *PosResMgr) Reg(poses, limitSws []int, uuid uint32) error {
	for _, v := range poses {
		if v == 0 {
			continue
		}
		v--

		if v >= len(this.posRes) || this.posRes[v] != 0 {
			return __Err_PosOccupied
		}
	}
	for _, v := range limitSws {
		if v == 0 {
			continue
		}
		v--

		if v >= len(this.limitSwRes) || this.limitSwRes[v] != 0 {
			return __Err_LmSwOccupied
		}
	}

	for _, v := range poses {
		if v == 0 {
			continue
		}
		v--

		this.posRes[v] = uuid
	}
	for _, v := range limitSws {
		if v == 0 {
			continue
		}
		v--

		this.limitSwRes[v] = uuid
	}

	return nil
}

func (this *PosResMgr) UnReg(uuid uint32) {

	for i, v := range this.posRes {
		if v == uuid {
			this.posRes[i] = 0
		}
	}
	for i, v := range this.limitSwRes {
		if v == uuid {
			this.limitSwRes[i] = 0
		}
	}
}
