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
	time.AfterFunc(time.Second*1, c.SayHello)
	time.AfterFunc(time.Second*2, c.UpdateId)
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
			time.AfterFunc(global.ReConnectDuration*time.Second, c.Start)
			log.Warn(fmt.Sprintf("客户端连接断开，%d秒后重连", global.ReConnectDuration))
		}
	}
}

//消息处理实体
func (c *Client) msgHandlerWorker(msg *object.CtrlMessage) *object.CtrlMessage {
	comm := common{}
	log.Debug(fmt.Sprintf("Client rec new message,id: %s,key: %s", msg.Id, msg.Key))
	switch msg.Key {
	case object.CtrlMessageDownloadFileList:
		c.handlerDownloadList(msg.Id, []byte(msg.Data))
		return nil
	case object.CtrlMessageReturn:
		c.handlerReturn(msg.Id, []byte(msg.Data))
		return nil
	default:
		errMsg := fmt.Sprintf("Message Key is not exist,key: %s", msg.Key)
		log.Warn(errMsg)
		return comm.GetRMessage(msg.Id, -1, errMsg)
	}
}

func (c *Client) handlerDownloadList(clientId string, d []byte) {
	var data []object.DownloadFile
	err := json.Unmarshal(d, &data)
	if err != nil {
		log.Error(fmt.Sprintf("Get Message Data DownloadFile error: %s", err.Error()))
		return
	}
	c.downFileList(data)
}

func (c *Client) downFileList(list []object.DownloadFile) {
	currPath, err := goToolCommon.GetCurrPath()
	if err != nil {
		log.Error(fmt.Sprintf("Get CurrPath error: %s", err.Error()))
		log.Warn("Download Stopped")
		return
	}
	for _, f := range list {
		go c.downFile(currPath, &f)
	}
}

func (c *Client) downFile(currPath string, f *object.DownloadFile) {
	fullPath := currPath + "\\" + f.SubPath
	fullPath = "C:\\Users\\yuansong\\go\\pkg\\goWebSocketDemoV2\\Client\\AA\\BB\\CC"
	err := goToolCommon.CheckAndCreateFolder(fullPath)
	if err != nil {
		log.Error(fmt.Sprintf("检查并创建路径时遇到错误：%s", err.Error()))
		return
	}
	//下载文件
	u := url.URL{Scheme: "ws", Host: global.SysConfig.Server.Address, Path: WebPathDownloadFile}
	var dialer = &websocket.Dialer{
		HandshakeTimeout: global.HttpConnectTimeout * time.Second,
	}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Error(fmt.Sprintf("WebSocket Dial error: %s", err.Error()))
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	//TODO
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

//updateId
func (c *Client) UpdateId() {
	log.Debug("Update id")
	u := object.UpdateId{
		Id: c.client.GetId(),
	}
	c.sendCtrlMessage(object.CtrlMessageUpdateId, u)
}

//获取文件
func (c *Client) DownloadList() {
	log.Debug("Client get file")
	df := object.DownloadFileList{
		SubPath: "\\" + global.ClientFileFolder,
	}
	c.sendCtrlMessage(object.CtrlMessageDownloadFileList, df)
}

func (c *Client) sendCtrlMessage(key string, v interface{}) {
	comm := common{}
	c.client.GetChSend() <- *comm.GetCtrlMessage(c.client.GetId(), key, v)
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
