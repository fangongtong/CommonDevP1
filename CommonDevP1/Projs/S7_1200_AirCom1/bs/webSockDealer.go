package bs

import (
	"encoding/json"
	"fmt"

	//"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	tmpl "CommonDevP1/Projs/S7_1200_AirCom1/PlcDbTemplate"
	"CommonDevP1/Projs/S7_1200_AirCom1/structs"
	"MyToolBox/DataDriver"
	syad "MyToolBox/Thread/SyncAppData"
)

type IUnRegWebSockClient interface {
	UnRegWsClient(client *WebSockClient)
}

type WebSockProc struct {
	topicDealer *DataDriver.TopicDataDealer
	syncMgr     *syad.SyncDataMgr

	collectCmds []tmpl.DBA_Sub_CmdParam
	posResMgr   *structs.PosResMgr
	clients     map[string]*WebSockClient
}

func (this *WebSockProc) Init(topicDealer *DataDriver.TopicDataDealer, syncMgr *syad.SyncDataMgr) error {
	this.topicDealer = topicDealer
	this.topicDealer.ChangeDealer(this.deal)
	this.syncMgr = syncMgr

	this.clients = make(map[string]*WebSockClient)

	return nil
}

func (this *WebSockProc) UnRegWsClient(client *WebSockClient) {
	this.topicDealer.GetDriver().PutData(Topic_WebSock, client)
}

//  status + task{pos/lmSw/runStatus}
func (this *WebSockProc) deal(dt interface{}, errGetDt error) {
	switch dt.(type) {
	case nil:
		return
	case *WebSockClient:
		client := dt.(*WebSockClient)

		if len(client.uid) == 0 { // 没有uid的为新建
			client.Init(uuid.New().String(), this)

			this.clients[client.uid] = client
			go client.write()
		} else { //  有uid的则为删除
			if _, ok := this.clients[client.uid]; ok {
				delete(this.clients, client.uid)
			}
		}

	case *RealTasksInfo: //  收到数据
		//fmt.Println("get ws data....")
		msg, _ := json.Marshal(dt)
		for _, client := range this.clients {
			client.AddData(msg)
		}
	}
}

// var upgrader = &websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func (this *ClientMgr) RegClient(w http.ResponseWriter, r *http.Request) error {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		return errpostma
// 	}
// 	client := (&client{conn: conn)}).Init(uuid.New().String(), this, 35)
// 	this.regChan <- client
// 	return nil
// }

func New_WebSockClient(conn *websocket.Conn, chanLen int) *WebSockClient {
	return &WebSockClient{
		conn:        conn,
		dataChanLen: chanLen,
		dataChannel: make(chan []byte, chanLen),
	}
}

// ------------------
type WebSockClient struct {
	uid         string
	mc          IUnRegWebSockClient
	dataChannel chan []byte
	conn        *websocket.Conn
	//chanFlg byte
	dataChanLen int
	chanLenIdx  int32
}

func (this *WebSockClient) Init(uid string, mc IUnRegWebSockClient) {
	this.uid, this.mc = uid, mc
}

func (this *WebSockClient) write() {
	fmt.Println("start write ws")
	defer func() {
		this.conn.Close()
		this.mc.UnRegWsClient(this)
		fmt.Println("write defer 3")
	}()
	var err error
	for {
		select {
		case data := <-this.dataChannel:
			atomic.AddInt32(&this.chanLenIdx, -1)
			//panic("nothing")
			this.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 300))
			if err = this.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				fmt.Println("ws data error:", err.Error())
				return
			}
		case <-time.After(time.Second):
			this.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 300))
			if err = this.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (this *WebSockClient) AddData(msg []byte) {
	if this.chanLenIdx < int32(this.dataChanLen-3) {
		//fmt.Println("AddData...")
		atomic.AddInt32(&this.chanLenIdx, 1)
		this.dataChannel <- msg
	}
}
