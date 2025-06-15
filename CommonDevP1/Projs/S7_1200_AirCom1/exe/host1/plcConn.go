package main

import (
	"MyToolBox/DataDriver"
	"bytes"
	"fmt"
	"runtime"
	"time"

	syad "MyToolBox/Thread/SyncAppData"

	"CommonDevP1/PlcSimulator/DeviceMemory"

	"github.com/robinson/gos7"
)

//  与各个被监控的子程序进行通信

type CommuProc struct {
	topicDealer *DataDriver.TopicDataDealer
	recvBuf     [1024]byte
	clientS7    gos7.Client
	clientConn  *gos7.TCPClientHandler
	syncMgr     *syad.SyncDataMgr

	dbCmd    CmdPack_Simple
	dbStatus DeviceMemory.DBA_2
}

func (this *CommuProc) Init(topicDealer *DataDriver.TopicDataDealer, syncMgr *syad.SyncDataMgr, plcAddr string, plcRack, plcSlot int, timeout int) error {
	var err error
	this.clientConn = gos7.NewTCPClientHandler(plcAddr, plcRack, plcSlot)
	this.clientConn.Timeout = time.Millisecond * time.Duration(timeout) //time.Millisecond * 1500

	this.topicDealer = topicDealer
	this.topicDealer.ChangeDealer(this.deal)

	this.syncMgr = syncMgr

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
	defer func() {
		this.clientS7 = nil
		this.clientConn.Close()
	}()

	var statusBuf []byte = make([]byte, this.dbStatus.Size())
	var cmdIdx uint32 = 0
	hasNewCmd := false

	startTime := time.Now()
	for {
		//this.dbStatus = &DeviceMemory.DBA_2{}

		//  先获得模式
		//  是读取状态模式, 还是读取调试菜单模式

		endTime := time.Now()
		ms := endTime.Sub(startTime).Milliseconds()
		fmt.Println("used time:", ms)
		startTime = endTime
		//  读取数据
		if err = this.clientS7.AGReadDB(101, 0, this.dbStatus.Size(), statusBuf[:0]); err != nil {
			fmt.Println("AGReadDB failed:", err.Error())
			return
		}

		this.dbStatus.Retrive(bytes.NewBuffer(statusBuf))

		this.topicDealer.GetDriver().PutData(Topic_Lua, this.dbStatus)

		if hasNewCmd && this.dbStatus.CmdIdx == cmdIdx {
			fmt.Println("cmd pass ok: ", this.dbStatus.CmdIdx)
			hasNewCmd = false
		}

		//  有指令数据则插入
		var cmdDt interface{}
		for cmdDt == nil {
			cmdDt = this.syncMgr.Lock(Sync_Cmd)
			runtime.Gosched()
		}

		cmdPack := cmdDt.(*Sync_CmdPack)
		if len(cmdPack.cmds) > 0 {
			fmt.Println("cmds found !!!")
			cmdIdx = this.dbCmd.Set(cmdPack.cmds)
		}
		cmdPack.cmds = []*DeviceMemory.DBCmdParam{}
		this.syncMgr.Unlock(Sync_Cmd)

		if this.dbCmd.CmdCount() > 0 {

			fmt.Printf("%+v \r\n", this.dbCmd)
			buf := bytes.NewBuffer(nil)
			this.dbCmd.TurnBytes(buf)
			byts := buf.Bytes()
			fmt.Println("sent dt:", len(byts))
			if err = this.clientS7.AGWriteDB(100, 0, len(byts), byts); err != nil {
				fmt.Println("AGWriteDB failed:", err.Error())
				return
			}
			hasNewCmd = true
			this.dbCmd.Clear()
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
