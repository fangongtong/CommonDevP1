package httpCtrl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"CommonDevP1/Projs/S7_1200_AirCom1/bs"
	"MyToolBox/DataDriver"

	"CommonDevP1/Projs/S7_1200_AirCom1/structs"

	SModMgr "MyToolBox/ModuleMgr/SimpleModuleMgr"

	def "CommonDevP1/Projs/S7_1200_AirCom1/exe/plchost/define"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RegClient(w http.ResponseWriter, r *http.Request, td DataDriver.ITopicDriver) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	client := bs.New_WebSockClient(conn, 64)
	td.PutData(bs.Topic_WebSock, client)
	// client := (&client{conn: conn)}).Init(uuid.New().String(), this, 35)
	// this.regChan <- client
	return nil
}

type StandTaskJson struct {
	ScriptFileName string
	Pos            []int
	PosLimitSw     []int
}

func StartNewTask(c *gin.Context) {
	jsCfgRaw, _ := c.GetRawData()
	fmt.Println("json is:", string(jsCfgRaw))
	jsTskObj := &StandTaskJson{}

	if err := json.Unmarshal(jsCfgRaw, jsTskObj); err != nil {
		fmt.Println("here 1")
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_BadParam, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_BadParam], err.Error())})
		return
	}

	if len(strings.Trim(jsTskObj.ScriptFileName, " /\\")) == 0 {
		fmt.Println("here 2")
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_ScriptFileNotFound, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_ScriptFileNotFound], jsTskObj.ScriptFileName)})
		return
	}

	fp := path.Join("./tasks", jsTskObj.ScriptFileName)
	_, err := os.Stat(fp)
	if err != nil {
		fmt.Println("here 3")
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_ScriptFileNotFound, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_ScriptFileNotFound], jsTskObj.ScriptFileName)})
		return
	}

	task := &structs.Task{}
	//jsCfg := `{"TaskType":3,"Force":0,"TotalTimes":100,"Pos":[2,0],"PosLimitSw":[0,0,0,0],"abc":{"a1":{"b":[{"b1":1,"b2":"abc"},{"b1":1,"b2":"def"}]}}}`

	dataDriver := SModMgr.GetModInst(def.Mod_DataDriver).(DataDriver.ITopicDriver)
	taskDealer := SModMgr.GetModInst(def.Mod_TaskDealer).(*bs.TaskProc)
	if err := task.Init(time.Now().Format("20060102150405"), taskDealer, fp, string(jsCfgRaw)); err == nil {
		dataDriver.PutData(bs.Topic_Task, task)
		dataDriver.PutData(bs.Topic_Task, &bs.TaskOper{
			Uid:  task.TskUid,
			Oper: 3,
		})
		c.JSON(http.StatusOK, &RtnData{Code: 0, RtnData: task.TskUid})
		// if task.Run() {
		// 	fmt.Println("here 4")
		// 	c.JSON(http.StatusOK, &RtnData{Code: 0})
		// } else {
		// 	fmt.Println("here 5")
		// 	c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_StartTaskFailed, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_StartTaskFailed], "unknown. not running")})
		// }
	} else {
		fmt.Println("here 6")
		c.JSON(http.StatusBadRequest, &RtnData{Code: __ReqErr_StartTaskFailed, ErrInfo: fmt.Sprintf(__Err_Info[__ReqErr_StartTaskFailed], err.Error())})
	}
}
func QueryRunningTasks(c *gin.Context) {
	fmt.Println("query running tasks")
	rtnVal := __TaskList.Tasks()

	for _, v := range rtnVal {
		fmt.Printf("%+v  | ", v)
	}
	fmt.Println("")

	c.JSON(http.StatusOK, &RtnData{Code: 0, RtnData: rtnVal})
}
