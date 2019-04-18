package service

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goWebSocketDemoV2/global"
	"github.com/Deansquirrel/goWebSocketDemoV2/object"
	"github.com/gorilla/websocket"
	"net/url"
	"strings"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

type Client struct {
	client object.IClient
}

//启动服务
func (c *Client) Start() {
	c.client = nil
	//u := url.URL{Scheme: "ws", Host: "127.0.0.1:1234", Path: "/websocket"}
	u := url.URL{Scheme: "ws", Host: global.SysConfig.Server.Address, Path: "/webSocket"}
	var dialer = &websocket.Dialer{
		HandshakeTimeout: global.HttpConnectTimeout * time.Second,
	}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Error(fmt.Sprintf("WebSocket Dial error: %s", err.Error()))
		time.AfterFunc(global.ReConnectDuration*time.Second, c.Start)
		return
	}
	c.client = object.NewClient(strings.ToUpper(goToolCommon.Guid()), conn)
	go c.msgHandler()
}

//消息处理
func (c *Client) msgHandler() {
	for {
		select {
		case msg, ok := <-c.client.GetChReceive():
			if ok {
				rm := c.msgHandlerWorker(&msg)
				if rm != nil {
					c.client.GetChSend() <- *rm
				}
			}
		case <-c.client.GetChClose():
			time.Sleep(global.ReConnectDuration * time.Second)
			c.Start()
		}
	}
}

//消息处理实体
func (c *Client) msgHandlerWorker(msg *object.CtrlMessage) *object.CtrlMessage {
	comm := common{}
	//if msg.MessageType != websocket.TextMessage {
	//	errMsg := fmt.Sprintf("Unexpected MessageType : %d", msg.MessageType)
	//	log.Error(errMsg)
	//	return comm.GetRMessage(msg.ClientId, -1, errMsg)
	//}
	//var m object.CtrlMessage
	//err := json.Unmarshal(msg.Data, &m)
	//if err != nil {
	//	errMsg := fmt.Sprintf("Get Message Object error: %s", err.Error())
	//	log.Error(errMsg)
	//	return comm.GetRMessage(msg.ClientId, -1, errMsg)
	//}
	log.Debug(fmt.Sprintf("Client rec new message,id: %s,key: %s", msg.Id, msg.Key))
	//===================================================================================
	switch msg.Key {
	case object.CtrlMessageReturn:
		c.handlerReturn(msg.Id, []byte(msg.Data))
		return nil
	case object.CtrlMessageTest:
		return comm.GetRMessage(msg.Id, 0, "ok")
	default:
		errMsg := fmt.Sprintf("Message Key is not exist,key: %s", msg.Key)
		log.Warn(errMsg)
		return comm.GetRMessage(msg.Id, -1, errMsg)
	}
}

//返回消息处理
func (c *Client) handlerReturn(clientId string, d []byte) {
	var data object.ReturnMessage
	err := json.Unmarshal(d, &data)
	if err != nil {
		log.Error(fmt.Sprintf("Get Message Data ReturnMessage error: %s", err.Error()))
		return
	}
	if data.ErrCode != 0 {
		log.Warn(fmt.Sprintf("错误[%d]: %s", data.ErrCode, data.ErrMsg))
	} else {
		log.Debug(fmt.Sprintf("Return Message: [%d]%s", data.ErrCode, data.ErrMsg))
	}
	return
}

//sayHello
func (c *Client) SayHello() {
	log.Debug("Client say hello")
	h := object.Hello{
		Msg: "Hello Server",
	}
	c.sendCtrlMessage(object.CtrlMessageHello, h)
}

//获取文件
func (c *Client) GetFile() {
	log.Debug("Client get file")
	df := object.DownloadFileList{
		SubPath: "\\clientFile",
	}
	c.sendCtrlMessage(object.CtrlMessageDownloadFileList, df)
}

func (c *Client) sendCtrlMessage(key string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Error(fmt.Sprintf("Get Object[%s] byte error: %s", key, err.Error()))
		return
	}
	cm := object.CtrlMessage{
		Id:   c.client.GetId(),
		Key:  key,
		Data: string(data),
	}
	c.client.GetChSend() <- cm
}

//
////获取文件
//func (c *Client) GetFile() {
//	u := url.URL{Scheme: "ws", Host: global.SysConfig.Server.Address, Path: "/webSocket/File"}
//	var dialer = &websocket.Dialer{
//		HandshakeTimeout: global.HttpConnectTimeout * time.Second,
//	}
//	conn, _, err := dialer.Dial(u.String(), nil)
//	if err != nil {
//		log.Error(fmt.Sprintf("WebSocket Dial error: %s", err.Error()))
//		time.AfterFunc(global.ReConnectDuration*time.Second, c.Start)
//		return
//	}
//	f := object.File{
//		Name: "test.exe",
//	}
//	d, err := json.Marshal(f)
//	if err != nil {
//		log.Error(fmt.Sprintf("Get filename byte error: %s", err.Error()))
//		return
//	}
//	om := object.OprMessage{
//		Id:   c.client.GetId(),
//		Key:  "",
//		Data: string(d),
//	}
//	err = conn.WriteJSON(om)
//	if err != nil {
//		log.Error(fmt.Sprintf("Socket write json error: %s", err.Error()))
//		return
//	}
//	var rf object.ReFile
//	err = conn.ReadJSON(&rf)
//	if err != nil {
//		log.Error(fmt.Sprintf("Socket read json error: %s", err.Error()))
//		return
//	}
//	if rf.ErrCode != 0 {
//		log.Error(fmt.Sprintf("Server return error:[%d]%s", rf.ErrCode, rf.ErrMsg))
//		return
//	}
//	path, err := goToolCommon.GetCurrPath()
//	if err != nil {
//		errMsg := fmt.Sprintf("Get CurrPath error: %s", err.Error())
//		log.Error(errMsg)
//		return
//	}
//	filePath := path + "\\" + f.Name
//	err = ioutil.WriteFile(filePath, rf.Data, 0644)
//	if err != nil {
//		log.Error(fmt.Sprintf("write file error: %s", err.Error()))
//	}
//	return
//}
