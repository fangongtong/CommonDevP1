package bs

import (
	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"CommonDevP1/Projs/S7_1200_AirCom1/structs"
	"MyToolBox/DataDriver"
	syad "MyToolBox/Thread/SyncAppData"
	"container/list"
	"fmt"
	"runtime"
	"time"
)

type ITaskNotify interface {
	Notify(uuid string, state int)
}

type _TskNotify struct {
}

func (this *_TskNotify) Notify(uuid string, state int) {}

type TaskProc struct {
	topicDealer *DataDriver.TopicDataDealer
	syncMgr     *syad.SyncDataMgr

	luaTasks    list.List
	collectCmds []tmpl.DBA_Sub_CmdParam
	posResMgr   *structs.PosResMgr

	tskNotifier ITaskNotify

	_TaskHumanSafeMod
}

func (this *TaskProc) Init(topicDealer *DataDriver.TopicDataDealer, syncMgr *syad.SyncDataMgr, tskNotifier ITaskNotify, maxCmdCnt, plcPosCnt, limitSwCnt int, plcRecovSafeTime int) error {
	this.topicDealer = topicDealer
	this.topicDealer.ChangeDealer(this.deal)
	this.tskNotifier = tskNotifier
	if this.tskNotifier == nil {
		this.tskNotifier = &_TskNotify{}
	}

	this.syncMgr = syncMgr

	this.collectCmds = make([]tmpl.DBA_Sub_CmdParam, maxCmdCnt)

	this.posResMgr = structs.New_PosMgr(plcPosCnt, limitSwCnt)

	this._TaskHumanSafeMod.SetTimeout(plcRecovSafeTime)
	this._TaskHumanSafeMod.InitTime()

	return nil
}

func (this *TaskProc) Reg(poses []int, limitSwitch []int, uuid uint32) error {
	//  注册资源
	return this.posResMgr.Reg(poses, limitSwitch, uuid)
}
func (this *TaskProc) UnReg(uuid uint32) {
	//  释放资源
	this.posResMgr.UnReg(uuid)
}

var __Err_Timeout = fmt.Errorf("timeout")

func (this *TaskProc) GetPosMenuParam(pos int) (*tmpl.DBA_Sub_MenuParam, error) {
	fmt.Println("TaskProc) GetPosMenuParam pos:", pos)
	req := &PosReq{Pos: pos}
	req.Init()

	this.topicDealer.GetDriver().PutData(Topic_DbOther, req)
	if req.Req(300) {
		if req.Err != nil {
			return nil, req.Err
		}
		return &req.Dt, nil
	} else {
		if req.Err != nil {
			return nil, req.Err
		} else {
			return nil, __Err_Timeout
		}
	}
	return nil, nil
}

func (this *TaskProc) NotifyMenuChanged(cmds []tmpl.DBA_Sub_CmdParam) {
	for i := 0; i < len(cmds); i++ {
		code := cmds[i].CmdCode
		if (code >= 10 && code <= 34) || code == 42 {
			this.topicDealer.GetDriver().PutData(Topic_DbOther, int(cmds[i].Param1))
		}
	}
}

func (this *TaskProc) deal(dt interface{}, errGetDt error) {
	switch dt.(type) {
	case nil:
		return
	case *CmdOper:
		//  应该先检查一下设定的参数所在工位是否正在运行中
		//  运行中的话应该告知错误
		cmdOper := dt.(*CmdOper)
		for i, _ := range this.collectCmds {
			this.collectCmds[i].Clear()
		}

		for i, v := range cmdOper.Cmds {
			this.collectCmds[i].CopyIn(v)
		}

		var synDt interface{}
		for synDt == nil {
			synDt = this.syncMgr.Lock(Sync_Cmd)
			runtime.Gosched()
		}

		sdt := synDt.(*structs.PlcCmdContainer)
		_, err := sdt.PackIn2(this.collectCmds[:len(cmdOper.Cmds)])
		this.syncMgr.Unlock(Sync_Cmd)
		if err != nil {
			cmdOper.Err = err
		}
		cmdOper.Req <- true

		//  一切顺利结束时通知清除menu缓存
		this.NotifyMenuChanged(this.collectCmds[:len(cmdOper.Cmds)])

	case *structs.Task: //  添加任务
		fmt.Println("new task")
		this.luaTasks.PushBack(dt)

	case *TaskOper: //  操作任务暂停继续和终止
		oper := dt.(*TaskOper)
		for itm := this.luaTasks.Front(); itm != nil; itm = itm.Next() {
			task := itm.Value.(*structs.Task)
			if task.TskUid == oper.Uid {
				switch oper.Oper {
				case 2:
					if task.Pause() {
						this.tskNotifier.Notify(oper.Uid, 2)
					}
				case 3:
					if task.Run() {
						this.tskNotifier.Notify(oper.Uid, 3)
					}
				case 4:
					task.UnInit()
					this.luaTasks.Remove(itm)
					this.tskNotifier.Notify(oper.Uid, 4)
				}
				break
			}
		}

	case *tmpl.DBA_Status: //  check任务

		startTime := time.Now()

		//  安全检查
		if !this._TaskHumanSafeMod.Check(startTime) {
			this.syncMgr.GetBuf(Sync_SysAlarmInfo).(*SyncData_SysAlarmInfo).AlarmInfo |= 1 << 1
			return
		} else {
			this.syncMgr.GetBuf(Sync_SysAlarmInfo).(*SyncData_SysAlarmInfo).AlarmInfo &= ^(uint32(1) << 1)
		}

		//fmt.Println(" ------- ")
		//var stpTime1, stpTime2, stpTime3 int64

		status := dt.(*tmpl.DBA_Status)

		//var rmTsk []*list.Element

		var logTsksInfo *LogTasksInfo

		var usedTime0 int64
		{
			if usedTime := time.Now().Sub(startTime).Milliseconds(); usedTime > 0 {
				usedTime0 = usedTime
				// fmt.Println("deal tasks time1:", stpTime1)
				// fmt.Println("deal tasks time2:", stpTime2)
				// fmt.Println("deal tasks time3:", stpTime3)
			}
		}

		for itm := this.luaTasks.Front(); itm != nil; itm = itm.Next() {
			task := itm.Value.(*structs.Task)

			if task.Status() == 4 {
				/*
					task.UnInit()
					rmTsk = append(rmTsk, itm)
				*/
				this.topicDealer.GetDriver().PutData(Topic_Task, &TaskOper{
					Uid:  task.TskUid,
					Oper: 4,
				})
				continue
			}
			if task.Status() != 3 {
				continue
			}

			for i, _ := range this.collectCmds {
				this.collectCmds[i].Clear()
			}

			rdInfo := &structs.RecordInfo{}

			if cmdCnt, err := task.Check(status, this.collectCmds, rdInfo); err != nil {

			} else if cmdCnt > 0 {

				//stpTime1 = time.Now().Sub(startTime).Milliseconds()
				//fmt.Println("cmd cnt:", cmdCnt)

				var synDt interface{}
				for synDt == nil {
					synDt = this.syncMgr.Lock(Sync_Cmd)
					runtime.Gosched()
				}

				sdt := synDt.(*structs.PlcCmdContainer)
				newIdx, err := sdt.PackIn2(this.collectCmds[:cmdCnt])
				this.syncMgr.Unlock(Sync_Cmd)

				if err != nil {
					fmt.Println("PackIn2 failed:", err.Error())
				} else {
					task.SetCmdIdx_SendingCmd(newIdx)

					//  一切顺利结束时通知清除menu缓存
					this.NotifyMenuChanged(this.collectCmds[:cmdCnt])
				}
			}

			//stpTime2 = time.Now().Sub(startTime).Milliseconds()
			//  当需要记录数据
			if rdInfo.NeedRecord {
				if logTsksInfo == nil {
					logTsksInfo = &LogTasksInfo{}
					logTsksInfo.Status = status.CylinderDataAry
					logTsksInfo.TimeFlg = startTime
					logTsksInfo.TaskInfo = append(logTsksInfo.TaskInfo, &TaskLowInfo{
						TaskUid:   rdInfo.TaskUid,
						Poses:     rdInfo.Poses,
						PosFields: rdInfo.TaskPosFields,
						TaskInfo:  rdInfo.TaskInfo,
					})
				} else {
					logTsksInfo.TaskInfo = append(logTsksInfo.TaskInfo, &TaskLowInfo{
						TaskUid:   rdInfo.TaskUid,
						Poses:     rdInfo.Poses,
						PosFields: rdInfo.TaskPosFields,
						TaskInfo:  rdInfo.TaskInfo,
					})
				}
			}
			//stpTime3 = time.Now().Sub(startTime).Milliseconds()

		}

		var usedTime1 int64
		{
			if usedTime := time.Now().Sub(startTime).Milliseconds(); usedTime > 0 {
				usedTime1 = usedTime
				// fmt.Println("deal tasks time1:", stpTime1)
				// fmt.Println("deal tasks time2:", stpTime2)
				// fmt.Println("deal tasks time3:", stpTime3)
			}
		}

		if logTsksInfo != nil {
			this.topicDealer.GetDriver().PutData(Topic_Influxdb, logTsksInfo)

			if usedTime := time.Now().Sub(startTime).Milliseconds(); usedTime > 0 {
				fmt.Println("deal tasks 2 time:", usedTime)
				// fmt.Println("deal tasks time1:", stpTime1)
				// fmt.Println("deal tasks time2:", stpTime2)
				// fmt.Println("deal tasks time3:", stpTime3)
			}
		}

		realDt := &RealTasksInfo{
			DevAlarm: status.DevAlarm,
			Status:   status.CylinderDataAry,
		}
		this.topicDealer.GetDriver().PutData(Topic_WebSock, realDt)

		// for _, v := range rmTsk {

		// 	this.luaTasks.Remove(v)
		// }

		if usedTime := time.Now().Sub(startTime).Milliseconds(); usedTime > 3 {
			fmt.Println("deal tasks time:", usedTime, usedTime1, usedTime0)
			// fmt.Println("deal tasks time1:", stpTime1)
			// fmt.Println("deal tasks time2:", stpTime2)
			// fmt.Println("deal tasks time3:", stpTime3)
		}
	}
}

type TaskOper struct {
	Uid  string
	Oper int // 2: pause 3: running 4: over
}
type CmdOper struct {
	Req  chan bool
	Cmds []*tmpl.DBA_Sub_CmdParam
	Err  error
}

// ---------------------------------------------------------------
//  人员安全模块
//  当plc长时间失联后恢复连接，不能让任务突然执行，否则可能伤害到人
type _TaskHumanSafeMod struct {
	startTime      time.Time
	timeoutSeconds int
}

func (this *_TaskHumanSafeMod) SetTimeout(seconds int) {
	this.timeoutSeconds = seconds
}
func (this *_TaskHumanSafeMod) InitTime() {
	this.startTime = time.Now()
}
func (this *_TaskHumanSafeMod) Check(tNow time.Time) bool {
	if tNow.Sub(this.startTime).Seconds() > float64(this.timeoutSeconds) {
		return false
	} else {
		this.startTime = tNow
		return true
	}
}

// ---------------------------------------------------------------
