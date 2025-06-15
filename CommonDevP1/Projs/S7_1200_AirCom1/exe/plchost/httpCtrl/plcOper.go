package httpCtrl

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"CommonDevP1/Projs/S7_1200_AirCom1/bs"
	def "CommonDevP1/Projs/S7_1200_AirCom1/exe/plchost/define"
	"CommonDevP1/Projs/S7_1200_AirCom1/structs"
	"MyToolBox/DataDriver"
	SModMgr "MyToolBox/ModuleMgr/SimpleModuleMgr"
)

type _StPosBind struct {
	Pos int
}

func GetMenuParam(c *gin.Context) {
	var posbind _StPosBind
	if err := c.BindQuery(&posbind); err != nil {
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_BadParam, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_BadParam], err.Error())})
	}

	taskDealer := SModMgr.GetModInst(def.Mod_TaskDealer).(*bs.TaskProc)
	if menuParam, err := taskDealer.GetPosMenuParam(posbind.Pos - 1); err != nil {
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_ReqMenuFailed, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_ReqMenuFailed], err.Error())})
	} else {
		c.JSON(http.StatusOK, &RtnData{Code: 0, RtnData: menuParam})
	}

}

func SetMenuParam(c *gin.Context) {
	cmd := &tmpl.DBA_Sub_CmdParam{}
	if err := c.BindJSON(cmd); err != nil {
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_BadParam, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_BadParam], err.Error())})
		return
	}

	dataDriver := SModMgr.GetModInst(def.Mod_DataDriver).(DataDriver.ITopicDriver)

	oper := &bs.CmdOper{
		Req: make(chan bool),
		Cmds: []*tmpl.DBA_Sub_CmdParam{
			cmd,
		},
	}

	dataDriver.PutData(bs.Topic_Task, oper)
	select {
	case <-oper.Req:
		if oper.Err != nil {
			c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_PushCmdsFailed, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_PushCmdsFailed], oper.Err.Error())})
		} else {
			c.JSON(http.StatusOK, &RtnData{Code: 0})
		}
	case <-time.After(time.Millisecond * 200):
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_OperOvertime, ErrInfo: __Err_Info[__ReqErr_OperOvertime]})
	}
}
func SetMenuParams(c *gin.Context) {
	var cmds []tmpl.DBA_Sub_CmdParam
	if err := c.BindJSON(&cmds); err != nil {
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_BadParam, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_BadParam], err.Error())})
		return
	}

	dataDriver := SModMgr.GetModInst(def.Mod_DataDriver).(DataDriver.ITopicDriver)

	oper := &bs.CmdOper{
		Req:  make(chan bool),
		Cmds: []*tmpl.DBA_Sub_CmdParam{},
	}

	for i, _ := range cmds {
		oper.Cmds = append(oper.Cmds, &cmds[i])
	}

	dataDriver.PutData(bs.Topic_Task, oper)
	select {
	case <-oper.Req:
		if oper.Err != nil {
			c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_PushCmdsFailed, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_PushCmdsFailed], oper.Err.Error())})
		} else {
			c.JSON(http.StatusOK, &RtnData{Code: 0})
		}
	case <-time.After(time.Millisecond * 200):
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_OperOvertime, ErrInfo: __Err_Info[__ReqErr_OperOvertime]})
	}
}

func SetCmds(c *gin.Context) {
	fmt.Println("SetCmds called ")
	var collectCmds []tmpl.DBA_Sub_CmdParam
	if err := c.BindJSON(&collectCmds); err != nil {
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_BadParam, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_BadParam], err.Error())})
		return
	}
	for _, v := range collectCmds {
		fmt.Printf("%+v \r\n", v)
	}

	dataDriver := SModMgr.GetModInst(def.Mod_DataDriver).(DataDriver.ITopicDriver)

	oper := &bs.CmdOper{
		Req:  make(chan bool),
		Cmds: []*tmpl.DBA_Sub_CmdParam{},
	}

	for i, _ := range collectCmds {
		oper.Cmds = append(oper.Cmds, &collectCmds[i])
	}

	dataDriver.PutData(bs.Topic_Task, oper)
	select {
	case <-oper.Req:
		if oper.Err != nil {
			c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_PushCmdsFailed, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_PushCmdsFailed], oper.Err.Error())})
		} else {
			c.JSON(http.StatusOK, &RtnData{Code: 0})
		}
	case <-time.After(time.Millisecond * 200):
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_OperOvertime, ErrInfo: __Err_Info[__ReqErr_OperOvertime]})
	}
}
func SetCmds2(c *gin.Context) {
	fmt.Println("SetCmds called ")
	var collectCmds []tmpl.DBA_Sub_CmdParam
	c.BindJSON(&collectCmds)

	for _, v := range collectCmds {
		fmt.Printf("%+v \r\n", v)
	}
	fmt.Println("")
	//collectCmds []tmpl.DBA_Sub_CmdParam

	posL := 0
	posR := 0

	for {
		if posL+__MaxCmdCnt > len(collectCmds) {
			posR = len(collectCmds)

			if posR == posL {
				break
			}
		} else {
			posR = posL + __MaxCmdCnt
		}

		var synDt interface{}
		for synDt == nil {
			synDt = __SyncDataMgr.Lock(bs.Sync_Cmd)
			runtime.Gosched()
		}

		sdt := synDt.(*structs.PlcCmdContainer)
		_, err := sdt.PackIn2(collectCmds[posL:posR])
		__SyncDataMgr.Unlock(bs.Sync_Cmd)

		if err != nil {
			c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_PushCmdsFailed, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_PushCmdsFailed], err.Error())})
			break
		} else {
			//  一切顺利结束时通知清除menu缓存
			taskDealer := SModMgr.GetModInst(def.Mod_TaskDealer).(*bs.TaskProc)
			taskDealer.NotifyMenuChanged(collectCmds[posL:posR])
		}

		posL = posR
	}
	c.JSON(http.StatusOK, &RtnData{Code: 0})
}
