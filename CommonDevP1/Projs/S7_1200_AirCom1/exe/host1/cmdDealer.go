package main

import (
	"CommonDevP1/PlcSimulator/DeviceMemory"
	"MyToolBox/DataDriver"
	syad "MyToolBox/Thread/SyncAppData"
	"runtime"
)

type CmdProc struct {
	topicDealer *DataDriver.TopicDataDealer
	recvBuf     [1024]byte
	syncMgr     *syad.SyncDataMgr
}

func (this *CmdProc) Init(topicDealer *DataDriver.TopicDataDealer, syncMgr *syad.SyncDataMgr) error {

	this.topicDealer = topicDealer
	this.topicDealer.ChangeDealer(this.deal)

	this.syncMgr = syncMgr

	return nil
}

func (this *CmdProc) deal(dt interface{}, errGetDt error) {

	if dt == nil {
		return
	}

	cmds := dt.([]*DeviceMemory.DBCmdParam)
	var synDt interface{}
	for synDt == nil {
		synDt = this.syncMgr.Lock(Sync_Cmd)
		runtime.Gosched()
	}

	sdt := synDt.(*Sync_CmdPack)
	sdt.cmds = cmds
	this.syncMgr.Unlock(Sync_Cmd)
}

/*

type Sync_CmdPack struct {
	cmds []*DeviceMemory.DBCmdParam
}

*/
