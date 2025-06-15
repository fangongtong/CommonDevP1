package bs

import (
	dtf "CommonDevP1/PlcSimulator/DataTransformer"
	dm "CommonDevP1/PlcSimulator/DeviceMemory"
	"CommonDevP1/PlcSimulator/Trans"
	"CommonDevP1/PlcSimulator/Trans/Udp"
	"container/list"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var CommuHost Commu = Commu{rtnCh: make(chan *dm.DBA_2, 10)}

type Cmd2SendPak struct {
	CmdHIdx uint32
	Cmds    []*dm.DBCmdParam
}

type Commu struct {
	host   Trans.ITrans
	cmdLst list.List
	mutex  sync.Mutex
	rtnCh  chan *dm.DBA_2
}

func (this *Commu) GetCh() <-chan *dm.DBA_2 {
	return this.rtnCh
}

func (this *Commu) Add(dt2Snd *Cmd2SendPak) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.cmdLst.PushBack(dt2Snd)

	//  如果当前连接失败，返回false
	//  如果当前list里指令太多，返回false
	return true
}
func (this *Commu) getData() *dtf.MsgCmd {
	if this.cmdLst.Len() == 0 {
		return nil
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()

	msgCmd := &dtf.MsgCmd{}
	var cmdHIdx []uint32
	//paramIdx := 0
	for {
		if e := this.cmdLst.Front(); e != nil {
			tmp := e.Value.(*Cmd2SendPak)
			for _, v := range cmdHIdx {
				if v == tmp.CmdHIdx {
					msgCmd.CmdCnt = 1
					return msgCmd
				}
			}
			this.cmdLst.Remove(e)
			for _, v := range tmp.Cmds {
				msgCmd.Cmds = append(msgCmd.Cmds, &dtf.StCmd{
					CmdCode: v.CmdCode,
					Param1:  v.Param1,
					Param2:  v.Param2,
					Param3:  v.Param3,
					Param4:  v.Param4,
					Param5:  v.Param5,
					Param6:  v.Param6,
					Param7:  v.Param7,
					Param8:  v.Param8,
				})
				//paramIdx++
			}
		}
	}

	msgCmd.CmdCnt = uint32(len(msgCmd.Cmds))
	return msgCmd
}

func (this *Commu) StartCommu(localIpPort, remoteIpPort string) {
	defer func() {
		fmt.Println("communication routine crash!")
		os.Exit(-1)
		return
	}()
	host := Udp.NewUdpHost(localIpPort)
	if host == nil {
		return
	}
	var rAddr *net.UDPAddr
	var err error
	if rAddr, err = net.ResolveUDPAddr("udp", remoteIpPort); err != nil {
		return
	}
	var recvBuf [1048]byte
	for {

		if dba1 := this.getData(); dba1 != nil {
			dt := dtf.PackCmd(dba1)
			if err = host.SendTo(dt, rAddr); err != nil {
				fmt.Println("send failed:", err.Error())
				return
			}
			continue
		}

		dt := dtf.PackData(0, nil)
		fmt.Print("to send.. ")

		fmt.Printf("datasize: %d", len(dt))
		if err = host.SendTo(dt, rAddr); err != nil {
			fmt.Println("send failed:", err.Error())
			return
		}
		fmt.Println("ok")
		n, e := host.Recv(recvBuf[:])
		if e != nil {
			fmt.Printf("Error on recv: %s \r\n", e.Error())
			return
		}

		//code, rsp := dtf.UnPackData1(recvBuf[:n])
		code, _ := dtf.UnPackData1(recvBuf[:n])

		switch code {
		case 2: //  2: real data return
			//fmt.Printf("response: %+v \r\n", rsp.(*dm.DBA_2))
			//this.rtnCh <- rsp.(*dm.DBA_2)
		//host.Send(dt)
		default:
			fmt.Printf("recv len is %d, code is:%d \r\n", n, code)
		}
		time.Sleep(time.Second * 5)
	}
}
