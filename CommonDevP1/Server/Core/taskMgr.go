//  本文件处理task的注册？、生命周期？、缸实时数据向任务发送
//  配置文件里缸应该是这么写的: cylinders=1,2,3,4,5,6 表示1-6编号的6个缸
package Core

import (
	"container/list"
	"fmt"
	"sync"
)

type _CO struct {
	pos   int
	ocupy bool
}

//-------

type _CylinderHelper struct {
	c []*_CO
	m sync.Mutex
}

func (this *_CylinderHelper) Occupy(cs []int) bool {
	this.m.Lock()
	defer this.m.Unlock()
	y := 0

	fmt.Printf("_CylinderHelper occupy 1: %+v\r\n %+v \r\n", this.c, cs)

	for _, p := range cs {
		for _, v := range this.c {
			//fmt.Println("1.param:", p, " pos:", v.pos)
			if v.pos == p && !v.ocupy {
				y++
				break
			}
		}
	}

	if y == len(cs) {
		for _, p := range cs {
			for _, v := range this.c {
				if v.pos == p {
					v.ocupy = true
					break
				}
			}
		}
		return true
	}
	return false
}
func (this *_CylinderHelper) Release(cs []int) {
	for _, p := range cs {
		for _, v := range this.c {
			if v.pos == p {
				v.ocupy = false
			}
		}
	}
}

func (this *_CylinderHelper) Init(cs []int) {
	for _, v := range cs {
		this.c = append(this.c, &_CO{pos: v})
	}
}

//-------

type _TskMgr struct {
	cylinder _CylinderHelper //  记录缸编号,物理缸的编号应该同plc相应,物理缸出现问题了也能通过配置文件屏蔽缸号来达到合理分配
	//tasks    map[string]*_Task
	tasks map[string]ITask

	dispatchedDt chan []byte
	plcTskNum    uint32
	cmdsLst      list.List
}

func NewTskMgr(cylinders []int) _TskMgr {
	//mgr := _TskMgr{tasks: make(map[string]*_Task)}
	mgr := _TskMgr{tasks: make(map[string]ITask)}
	mgr.cylinder.Init(cylinders)

	return mgr
}

func (this *_TskMgr) NewTask(b ITask, cylinder []int, script string, paramsJson string) ITask {

	if tsk := newTask(b, cylinder, script, paramsJson, this); tsk != nil {
		return tsk
	}
	return nil
}

/*
func (this *_TskMgr) NewTask(cylinder []int, script string, paramsJson string) *_Task {

	if tsk := newTask(cylinder, script, paramsJson, this); tsk != nil {
		return tsk
	}
	return nil
}
*/

//func (this *_TskMgr) regTask(tsk *_Task) bool {
func (this *_TskMgr) regTask(tsk ITask) bool {

	fmt.Println("mgr regtask 1")

	if this.cylinder.Occupy(tsk.Base().getCylinders()) {
		fmt.Println("mgr regtask 2")
		this.tasks[tsk.Base().getUid()] = tsk
		return true
	}
	fmt.Println("mgr regtask 3")
	return false
}
func (this *_TskMgr) unregTask(tskCode string) {
	if tsk, ok := this.tasks[tskCode]; ok {
		this.cylinder.Release(tsk.Base().getCylinders())
		delete(this.tasks, tskCode)
	}
}

func (this *_TskMgr) Pt() {
	fmt.Printf("%+v \r\n", this.cylinder)
}

func (this *_TskMgr) DispatchData(dt []byte) {
	this.dispatchedDt <- dt
}

//  由各个task调用来向mgr添加指令，mgr再向plc发送指令
func (this *_TskMgr) AddCmds(cylinderCode uint32) {
	c := &_CmdsList{
		CylinderCode: cylinderCode,
	}
	this.cmdsLst.PushBack(c)
}

//  高速chan处理,可以统一交给外部来处理
func (this *_TskMgr) GetChanFunc() chan interface{} {

}

func (this *_TskMgr) dispatchTaskDt() {

}

//  低速循环处理,可以统一交给外部来处理
func (this *_TskMgr) Loop() {
	//  获得返回数据
	select {
	case dt, ok := <-this.dispatchedDt:
		if ok {
			//  处理数据

			//  更新plcTskNum
			//  比对plcTskNum

			//  去掉mark了的this.cmdsLst

			// 分发给各个Task
		}
	}

	// 当plcTskNum符合发送要求,这里还要加入计数，一旦超过计数，可能要重新连接，或者重新 设置plcTskNum
	// 判断要发送的数据
	var cPos [16]bool
	for e := this.cmdsLst.Front(); e != nil; e = e.Next() {
		//  第一遍循环判断重复pos

		//  第二遍循环实际占用pos

		//  加入到command里

	}

	// 组装command为byte数据

	// 向plc发送单元 添加待发送数据
}

type _CmdsList struct {
	Mark         bool
	CylinderCode uint32 //  最多支持15缸,同时运行8缸
	// [] 工位指令
}
