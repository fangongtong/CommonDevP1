package bs

import (
	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"CommonDevP1/Projs/S7_1200_AirCom1/structs"
	"MyToolBox/DataDriver"
	"bytes"
	"encoding/binary"
	"fmt"
	"runtime"
	"time"

	syad "MyToolBox/Thread/SyncAppData"

	"github.com/robinson/gos7"
)

//  与各个被监控的子程序进行通信

type CommuProc struct {
	topicDealer *DataDriver.TopicDataDealer
	recvBuf     [1024]byte
	clientS7    gos7.Client
	clientConn  *gos7.TCPClientHandler
	syncMgr     *syad.SyncDataMgr

	maxCmdCnt    int
	targetCmdIdx uint32
}

func (this *CommuProc) Init(topicDealer *DataDriver.TopicDataDealer, syncMgr *syad.SyncDataMgr, plcAddr string, plcRack, plcSlot int, timeout int, maxCmdCnt int) error {
	var err error
	this.clientConn = gos7.NewTCPClientHandler(plcAddr, plcRack, plcSlot)
	this.clientConn.Timeout = time.Millisecond * time.Duration(timeout) //time.Millisecond * 1500

	this.topicDealer = topicDealer
	this.topicDealer.ChangeDealer(this.deal)

	this.syncMgr = syncMgr

	this.maxCmdCnt = maxCmdCnt
	//handler.Close()

	//this.hostIpPort = hostIpPort

	//this.localAddr, _ = net.ResolveUDPAddr("udp", fmt.Sprint("0.0.0.0:", udpListenPort))

	//  从数据库获得所有参数的版本

	//localAddr, err := net.ResolveUDPAddr("udp", hostIpPort)

	return err
	//this.localAddr = localAddr
	//return this.connect()
}

func (this *CommuProc) connect() error {
	err := this.clientConn.Connect()
	if err != nil {
		if err == gos7.Err_IsoInvalidPDU {
			fmt.Println("Error: ", err.Error())
		}
		//fmt.Println("haha ", err.Error())
		time.Sleep(time.Second * 2)
	} else {
		fmt.Println("plc connected")
	}
	return err
}

func (this *CommuProc) deal(dt interface{}, errGetDt error) {
	var err error
	if this.clientS7 == nil {
		if err = this.connect(); err != nil {
			return
		} else {
			this.clientS7 = gos7.NewClient(this.clientConn)
		}
	}
	syncTmp := this.syncMgr.GetBuf(Sync_PlcConnStatus).(*SyncData_PlcConnStatus)
	defer func() {
		syncTmp.Connected = false

		this.clientS7 = nil
		this.clientConn.Close()
	}()

	dbStatus := tmpl.New_DBA_Status(this.maxCmdCnt)
	var statusBuf []byte = make([]byte, dbStatus.Size())
	var cmdIdx uint32 = 0
	hasNewCmd := false

	syncTmp.Connected = true
	startTime := time.Now()
	for {

		//time.Sleep(time.Second * 1)
		//  先获得模式
		//  是读取状态模式, 还是读取调试菜单模式
		syncPosMenuParam := this.syncMgr.GetBuf(Sync_PosMenuParam).(*SyncData_PosMenuParam)
		if syncPosMenuParam.InReq {
			fmt.Println("CommuProc) deal syncPosMenuParam---")
			startPos := syncPosMenuParam.Pos * syncPosMenuParam.RspDt.Size()
			if err = this.clientS7.AGReadDB(103, startPos, syncPosMenuParam.RspDt.Size(), statusBuf[:0]); err != nil {
				fmt.Println("AGReadDB failed:", err.Error())
				syncPosMenuParam.Err = err

				syncPosMenuParam.InReq = false
				return
			} else {
				if err = syncPosMenuParam.RspDt.Unmarshal(bytes.NewBuffer(statusBuf), binary.BigEndian); err != nil {
					syncPosMenuParam.Err = err
				}
				syncPosMenuParam.InReq = false
			}
		} else if false { //

		}

		endTime := time.Now()
		ms := endTime.Sub(startTime).Milliseconds()
		if false {
			fmt.Println("used time:", ms)
		}
		startTime = endTime
		//  读取数据
		dbStatus = tmpl.New_DBA_Status(this.maxCmdCnt)
		if err = this.clientS7.AGReadDB(101, 0, dbStatus.Size(), statusBuf[:0]); err != nil {
			fmt.Println("AGReadDB failed:", err.Error())
			return
		}

		dbStatus.Unmarshal(bytes.NewBuffer(statusBuf), binary.BigEndian)

		if this.targetCmdIdx > 0 {
			if dbStatus.CmdIdx == 0 {
				//  这里假定为plc重启了
			} else {
				if dbStatus.CmdIdx != this.targetCmdIdx {
					continue
				}
			}
		}

		this.topicDealer.GetDriver().PutData(Topic_Task, dbStatus)

		if hasNewCmd && dbStatus.CmdIdx == cmdIdx {
			fmt.Println("cmd pass ok: ", dbStatus.CmdIdx)
			hasNewCmd = false
		}

		//  有指令数据则插入
		var cmdDt interface{}
		for cmdDt == nil {
			cmdDt = this.syncMgr.Lock(Sync_Cmd)
			runtime.Gosched()
		}

		cmdPack := cmdDt.(*structs.PlcCmdContainer)

		if pak := cmdPack.GetCmdPack(); pak != nil {
			//fmt.Println("cmds found !!!")
			this.targetCmdIdx = pak.CmdIdx
			this.syncMgr.Unlock(Sync_Cmd)

			buf := bytes.NewBuffer(nil)
			pak.Marshal(buf, binary.BigEndian)
			byts := buf.Bytes()
			//fmt.Println("sent dt:", len(byts))
			if err = this.clientS7.AGWriteDB(100, 0, len(byts), byts); err != nil {
				fmt.Println("AGWriteDB failed:", err.Error())
				return
			}
		} else {
			this.syncMgr.Unlock(Sync_Cmd)
		}
		//  如果插入成功，这里要检查下是否需要等待5ms

	}

	//  接收来自其他进程的通信数据
	//  目前应该只有优先主机数据(相位等等)
	/*
		this.conn.SetDeadline(time.Now().Add(time.Millisecond * 50))
		n, _, err := this.conn.ReadFromUDP(this.recvBuf[:])

		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				if opErr.Timeout() {
					return
				} else {
					log.Log(true, common.InstMgr, 5, "read from udp failed:"+opErr.Err.Error())
					goto exitProc
				}
			} else {
				log.Log(true, common.InstMgr, 5, "read from udp failed:"+opErr.Err.Error())
				goto exitProc
			}
		}
		log.Log(false, common.InstMgr, 1, fmt.Sprint("read from udp count:", n))
	*/

}
