package main

import (
	"CommonDevP1/Projs/S7_1200_AirCom1/bs"
	httpCtrl "CommonDevP1/Projs/S7_1200_AirCom1/exe/plchost/httpCtrl"
	"CommonDevP1/Projs/S7_1200_AirCom1/exe/plchost/other"
	"CommonDevP1/Projs/S7_1200_AirCom1/structs"
	"MyToolBox/DataDriver"
	log "MyToolBox/Log/udp-log"
	syad "MyToolBox/Thread/SyncAppData"
	"fmt"
	"net/http"
	"time"

	def "CommonDevP1/Projs/S7_1200_AirCom1/exe/plchost/define"
	SModMgr "MyToolBox/ModuleMgr/SimpleModuleMgr"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	__MaxPackCnt = 5
	__MaxCmdCnt  = 8
	__PlcPosCnt  = 8
	__LimitSwCnt = 2
)

var httpRouter = gin.Default()

func main() {
	log.SetUdpCfg("127.0.0.1:7039")
	sdmgr := initSyncDataMgr()
	taskNotifyList := &other.TaskList{}
	taskNotifyList.Init()
	httpCtrl.SetTaskList(taskNotifyList)
	httpCtrl.SetSyncMgr(sdmgr)

	threadDriver := DataDriver.Create_TopicDataDriver2(10)
	SModMgr.Init(8)
	SModMgr.Reg(def.Mod_DataDriver, threadDriver)

	pak_commu := DataDriver.Create_TopicDataDealer(bs.Topic_PlcCmmu, 64, nil, int64(time.Second*3))
	connDealer := &bs.CommuProc{}
	connDealer.Init(pak_commu, sdmgr, "192.168.1.180", 0, 0, 3*1000, __MaxCmdCnt)
	threadDriver.AddDealer(pak_commu)

	// pak_lua := DataDriver.Create_TopicDataDealer(bs.Topic_Task, 64, nil, int(time.Millisecond*10))
	// luaDealer := &LuaProc{}
	// luaDealer.Init(pak_lua, sdmgr)
	// threadDriver.AddDealer(pak_lua)

	pak_cmd := DataDriver.Create_TopicDataDealer(bs.Topic_Task, 64, nil, int64(time.Second*10))
	tskDealer := &bs.TaskProc{}
	tskDealer.Init(pak_cmd, sdmgr, taskNotifyList, __MaxCmdCnt, __PlcPosCnt, __LimitSwCnt, 5) //maxCmdCnt, plcPosCnt, limitSwCnt
	threadDriver.AddDealer(pak_cmd)
	SModMgr.Reg(def.Mod_TaskDealer, tskDealer)

	pak_dbOther := DataDriver.Create_TopicDataDealer(bs.Topic_DbOther, 32, nil, int64(time.Second*10))
	dbOtherDealer := &bs.DbOtherProc{}
	dbOtherDealer.Init(pak_dbOther, sdmgr, __MaxPackCnt)
	threadDriver.AddDealer(pak_dbOther)

	pak_influxDb := DataDriver.Create_TopicDataDealer(bs.Topic_Influxdb, 128, nil, int64(time.Second*10))
	dbInfluxDealer := &bs.InfluxDbProc{}
	dbInfluxDealer.Init(pak_influxDb, sdmgr, "http://192.168.1.133:8086", __PlcPosCnt)
	threadDriver.AddDealer(pak_influxDb)

	pak_webSock := DataDriver.Create_TopicDataDealer(bs.Topic_WebSock, 64, nil, int64(time.Second*3))
	webSockDealer := &bs.WebSockProc{}
	webSockDealer.Init(pak_webSock, sdmgr)
	threadDriver.AddDealer(pak_webSock)

	initRouter(threadDriver)
	// httpCtrl.SetDtDriver(threadDriver)
	// httpCtrl.SetTaskProcDealer(tskDealer)
	go httpRouter.Run(":9088")

	go func() {
		return
		time.Sleep(time.Second * 1)
		task := &structs.Task{}
		jsCfg := `{"TaskType":3,"Force":0,"TotalTimes":100,"Pos":[2,0],"PosLimitSw":[0,0,0,0],"abc":{"a1":{"b":[{"b1":1,"b2":"abc"},{"b1":1,"b2":"def"}]}}}`
		if err := task.Init(time.Now().Format("20060102150405"), tskDealer, "./tasks/1AC_pullOrPush.lua", jsCfg); err == nil {
			threadDriver.PutData(bs.Topic_Task, task)
			if task.Run() {
				fmt.Println("task.run")
			}
		} else {
			fmt.Println("task.init failed:", err.Error())
		}

		time.Sleep(time.Second * 1)
		task = &structs.Task{}
		jsCfg = `{"TaskType":3,"Force":0,"TotalTimes":99,"Pos":[1,0],"PosLimitSw":[0,0,0,0],"abc":{"a1":{"b":[{"b1":1,"b2":"abc"},{"b1":1,"b2":"def"}]}}}`
		if err := task.Init(time.Now().Format("20060102150405"), tskDealer, "./tasks/1AC_pullOrPush.lua", jsCfg); err == nil {
			threadDriver.PutData(bs.Topic_Task, task)
			if task.Run() {
				fmt.Println("task.run")
			}
		} else {
			fmt.Println("task.init failed:", err.Error())
		}
	}()

	threadDriver.Serve()
}

func initSyncDataMgr() *syad.SyncDataMgr {

	syncDataMgr := syad.New_SyncDataMgr(8)
	syncDataMgr.RegBuf(bs.Sync_Cmd, structs.NewPlcCmdContainer(__MaxPackCnt, __MaxCmdCnt, 1))
	syncDataMgr.RegBuf(bs.Sync_PosMenuParam, &bs.SyncData_PosMenuParam{})
	syncDataMgr.RegBuf(bs.Sync_PlcConnStatus, &bs.SyncData_PlcConnStatus{})
	syncDataMgr.RegBuf(bs.Sync_SysAlarmInfo, &bs.SyncData_SysAlarmInfo{})

	return syncDataMgr
}

func initRouter(td DataDriver.ITopicDriver) {

	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AllowCredentials = true
	httpRouter.Use(cors.New(cfg))

	httpRouter.POST("/setMenuParam", httpCtrl.SetMenuParam)
	httpRouter.POST("/newTask", httpCtrl.StartNewTask)
	httpRouter.POST("/SetCmds", httpCtrl.SetCmds)
	httpRouter.GET("/queryRunningTasks", httpCtrl.QueryRunningTasks)
	httpRouter.GET("/getMenuParam", httpCtrl.GetMenuParam)
	httpRouter.POST("/setMenuParams", httpCtrl.SetMenuParams)

	httpRouter.GET("/ws", func(c *gin.Context) {
		err := httpCtrl.RegClient(c.Writer, c.Request, td)
		if err != nil {
			c.JSON(http.StatusConflict, err.Error())
			return
		}
	})
}
