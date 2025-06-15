package Tests

import (
	"CommonDevP1/Server/Core"
	"testing"
)

func TestCore_1(t *testing.T) {

	mgr := Core.NewTskMgr([]int{1, 2, 4, 7})
	tsk := mgr.NewTask([]int{1, 4}, "", "")
	if tsk != nil {
		tsk.Pt()
	}

}
