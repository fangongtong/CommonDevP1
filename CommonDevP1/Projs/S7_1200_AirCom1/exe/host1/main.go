package main

import (
	"MyToolBox/DataDriver"
	log "MyToolBox/Log/udp-log"
	syad "MyToolBox/Thread/SyncAppData"
	"fmt"
	"time"
)

func main() {
	log.SetUdpCfg("127.0.0.1:7039")
	sdmgr := initSyncDataMgr()

	threadDriver := DataDriver.Create_TopicDataDriver2(10)

	pak_commu := DataDriver.Create_TopicDataDealer(Topic_Conn, 64, nil, int(time.Second))
	connDealer := &CommuProc{}
	connDealer.Init(pak_commu, sdmgr, "192.168.1.180", 0, 0, 3*1000)
	threadDriver.AddDealer(pak_commu)

	pak_lua := DataDriver.Create_TopicDataDealer(Topic_Lua, 64, nil, int(time.Millisecond*10))
	luaDealer := &LuaProc{}
	if err := luaDealer.Init(pak_lua, sdmgr); err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	threadDriver.AddDealer(pak_lua)

	pak_cmd := DataDriver.Create_TopicDataDealer(Topic_Cmd, 64, nil, int(time.Millisecond*10))
	cmdDealer := &CmdProc{}
	cmdDealer.Init(pak_cmd, sdmgr)
	threadDriver.AddDealer(pak_cmd)

	threadDriver.Serve()
}

func initSyncDataMgr() *syad.SyncDataMgr {

	syncDataMgr := syad.New_SyncDataMgr(Sync_Data_Len)
	syncDataMgr.RegBuf(Sync_ReadMode, ReadMode_Status)
	syncDataMgr.RegBuf(Sync_Cmd, &Sync_CmdPack{})

	return syncDataMgr
}
