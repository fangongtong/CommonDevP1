//  本文件用于处理脚本与plc交互的每个实际工作任务
package Core

import (
	"CommonDevP1/Server/Common"
)

type ITask interface {
	//Reg() bool
	//UnReg()
	SetBase(*BaseTask)
	//Wrapper() ITask
	Base() *BaseTask
	Start()
}

type BaseTask struct {
	mgr       *_TskMgr
	script    string
	tskCode   string
	cylinder  []int
	OutputErr error
	wrap      ITask
}

func (this *BaseTask) Base() *BaseTask {
	return nil
}
func (this *BaseTask) wrapper() ITask {
	return this.wrap
}
func (this *BaseTask) SetBase(b ITask) {
}

func (this *BaseTask) reg() bool {
	//return this.mgr.regTask(this)
	return this.mgr.regTask(this.wrap)
}

func (this *BaseTask) unReg() {
	this.mgr.unregTask(this.tskCode)
}

func (this *BaseTask) Destory() {
	this.mgr.unregTask(this.tskCode)
}

func (this *BaseTask) getCylinders() []int {
	return this.cylinder
}
func (this *BaseTask) getUid() string {
	return this.tskCode
}

func (this *BaseTask) GetStartResult() error {
	return this.OutputErr
}

func (this *BaseTask) Start() {
}

func newTask(wrap ITask, cylinder []int, script string, paramsJson string, mgr *_TskMgr) ITask {
	tsk := &BaseTask{mgr: mgr, tskCode: Common.GetUID(), wrap: wrap}
	wrap.SetBase(tsk)
	tsk.cylinder = make([]int, len(cylinder))
	copy(tsk.cylinder, cylinder)

	if !tsk.reg() {
		return nil
	}
	return wrap
}

/*
func newTask(cylinder []int, script string, paramsJson string, mgr *_TskMgr) *BaseTask {
	tsk := &BaseTask{mgr: mgr, tskCode: Common.GetUID()}
	tsk.cylinder = make([]int, len(cylinder))
	copy(tsk.cylinder, cylinder)

	if !tsk.reg() {
		return nil
	}
	return tsk
}
*/
