package bs

import (
	"fmt"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	api "github.com/influxdata/influxdb-client-go/v2/api"

	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"MyToolBox/DataDriver"
	syad "MyToolBox/Thread/SyncAppData"
)

type InfluxDbProc struct {
	httpAddr string
	writeAPI api.WriteAPI

	topicDealer   *DataDriver.TopicDataDealer
	syncMgr       *syad.SyncDataMgr
	connLastCheck time.Time

	dbTag      map[string]string
	dbTaskInfo map[string]interface{}
	//dbFieldPoses []map[string]interface{}
}

func (this *InfluxDbProc) Init(topicDealer *DataDriver.TopicDataDealer, syncMgr *syad.SyncDataMgr, influxHttpAddr string, posCnt int) error {

	this.httpAddr = influxHttpAddr

	this.topicDealer = topicDealer
	this.topicDealer.ChangeDealer(this.deal)
	this.syncMgr = syncMgr

	this.dbTag = map[string]string{
		"Uid": "",
	}
	this.dbTaskInfo = map[string]interface{}{
		"TaskInfo": "{}",
	}
	// this.dbFieldPoses = make([]map[string]interface{}, posCnt)
	// for i := 0; i < posCnt; i++ {
	// 	this.dbFieldPoses[i] = map[string]interface{}
	// }

	this.connect()
	return nil
}

func (this *InfluxDbProc) connect() {
	client := influxdb2.NewClient(this.httpAddr, "")
	client.Options().SetPrecision(time.Millisecond)
	this.writeAPI = client.WriteAPI("", "db")
}

func setMapPosData(mapItm map[string]interface{}, dt *tmpl.DBA_Sub_CylinderDatas, targetFields string) {
	fields := strings.Split(targetFields, "|")

	for _, v := range fields {
		switch v {
		case "RealForce":
			mapItm["RealForce"] = dt.RealForce
		case "PeekForce":
			mapItm["PeekForce"] = dt.PeekForce
		case "RealDisplacement":
			mapItm["RealDisplacement"] = dt.RealDisplacement
		case "DriveDegree":
			mapItm["DriveDegree"] = dt.DriveDegree
		}
	}

}

func (this *InfluxDbProc) deal(dt interface{}, errGetDt error) {
	if dt == nil {
		return
	}

	if this.writeAPI == nil {
		tNow := time.Now()
		if tNow.Sub(this.connLastCheck).Seconds() > 3 {
			this.connLastCheck = tNow
			this.connect()
			if this.writeAPI == nil {
				return
			}
		}
	}

	dtInfo := dt.(*LogTasksInfo)

	for _, v := range dtInfo.TaskInfo {

		this.dbTag["Uid"] = v.TaskUid
		this.dbTaskInfo["TaskInfo"] = v.TaskInfo

		p := influxdb2.NewPoint("TaskPointInfo", this.dbTag, this.dbTaskInfo, dtInfo.TimeFlg)
		this.writeAPI.WritePoint(p)

		for _, pos := range v.Poses {
			if pos == 0 {
				continue
			}
			mp := make(map[string]interface{})
			setMapPosData(mp, &dtInfo.Status[pos-1], v.PosFields)

			p = influxdb2.NewPoint(fmt.Sprint("TaskPointPos", pos), this.dbTag, mp, dtInfo.TimeFlg)
			this.writeAPI.WritePoint(p)
		}
	}
}
