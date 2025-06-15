//  中转获取设置菜单数据
//  缓存设置菜单plc数据
package bs

import (
	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"MyToolBox/DataDriver"
	syad "MyToolBox/Thread/SyncAppData"
	"fmt"
	"runtime"
	"time"
)

type DbOtherProc struct {
	topicDealer *DataDriver.TopicDataDealer
	syncMgr     *syad.SyncDataMgr

	menuItems   []tmpl.DBA_Sub_MenuParam
	menuRefresh []bool
}

func (this *DbOtherProc) Init(topicDealer *DataDriver.TopicDataDealer, syncMgr *syad.SyncDataMgr, maxPackCnt int) error {
	this.topicDealer = topicDealer
	this.topicDealer.ChangeDealer(this.deal)

	this.syncMgr = syncMgr

	this.menuRefresh = make([]bool, maxPackCnt)
	for i := 0; i < maxPackCnt; i++ {
		this.menuRefresh[i] = true
	}

	this.menuItems = make([]tmpl.DBA_Sub_MenuParam, maxPackCnt)

	return nil
}

var __Err_PlcDisconnected = fmt.Errorf("plc not connect")

func (this *DbOtherProc) deal(dt interface{}, errGetDt error) {
	switch dt.(type) {
	case nil:
		return
	case int:
		pos := dt.(int)
		//fmt.Println(pos)
		this.menuRefresh[pos-1] = true
	// case *MenuChange:
	// 	mc := dt.(*MenuChange)
	// 	this.menuRefresh[mc.Pos-1] = true
	case *PosReq:

		startTime := time.Now()
		syncTmp := this.syncMgr.GetBuf(Sync_PlcConnStatus).(*SyncData_PlcConnStatus)
		req := dt.(*PosReq)
		posIdx := req.Pos

		if this.menuRefresh[posIdx] == true {
			//  请求数据

			for {
				if !syncTmp.Connected { //  如果plc没连，返回
					req.Err = __Err_PlcDisconnected
					req.Rsp()
					return
				}

				if val := this.syncMgr.Lock(Sync_PosMenuParam); val == nil {
					runtime.Gosched()
					continue
				} else {
					reqS := val.(*SyncData_PosMenuParam)
					reqS.Err = nil
					reqS.Pos = posIdx
					reqS.InReq = true

					this.syncMgr.Unlock(Sync_PosMenuParam)

					fmt.Println("DbOtherProc) deal 1")

					for reqS.InReq {
						time.Sleep(time.Millisecond * 10)

						if !syncTmp.Connected { //  如果plc没连，返回
							req.Err = __Err_PlcDisconnected
							req.Rsp()
							return
						}
					}

					fmt.Println("DbOtherProc) deal 2")
					if reqS.Err != nil {
						req.Err = reqS.Err
						req.Rsp()
						return
					}

					this.menuItems[posIdx].CopyIn(&reqS.RspDt)

					fmt.Println("DbOtherProc) deal 3")
					break
				}
			}

			//  重置刷新标志
			this.menuRefresh[posIdx] = false
		}

		if this.menuRefresh[posIdx] == false {
			req.Dt.CopyIn(&this.menuItems[posIdx])
			req.Rsp()
		}
		fmt.Println("used time: ", time.Now().Sub(startTime).Milliseconds())
	}

}

// type MenuChange struct {
// 	Pos int
// }
// type LoadMenu struct {
// 	Req chan bool
// 	Pos int
// 	//Menu []*tmpl.DBA_Sub_CmdParam
// 	Err error
// }

type PosReq struct {
	syad.ChanBaseSyncData
	Pos int
	Dt  tmpl.DBA_Sub_MenuParam
	Err error
}
