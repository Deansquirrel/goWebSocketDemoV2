package service

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goWebSocketDemoV2/global"
	"github.com/Deansquirrel/goWebSocketDemoV2/object"
	"github.com/gorilla/websocket"
	"io/ioutil"
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
	//chList := make(chan struct{}, global.MaxThread)
	//defer close(chList)
	//for i := 0; i < global.MaxThread; i++ {
	//	chList <- struct{}{}
	//}
	for _, f := range list {
		c.downFile(currPath, &f)
	}
}

func (c *Client) downFile(currPath string, f *object.DownloadFile) {
	log.Debug(fmt.Sprintf("Download %s Start", f.Name))
	defer log.Debug(fmt.Sprintf("Download %s Complete", f.Name))
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

	err = conn.WriteJSON(&f)
	if err != nil {
		log.Error(fmt.Sprintf("WebSocket write error: %s", err.Error()))
		return
	}

	var rData object.DownloadFileData
	err = conn.ReadJSON(&rData)
	if err != nil {
		log.Error(fmt.Sprintf("WebSocket read error: %s", err.Error()))
		return
	}

	fullPath := currPath + rData.Info.SubPath
	err = goToolCommon.CheckAndCreateFolder(fullPath)
	if err != nil {
		log.Error(fmt.Sprintf("检查并创建路径时遇到错误：%s", err.Error()))
		return
	}

	err = ioutil.WriteFile(fullPath+"\\"+rData.Info.Name, rData.Data, 0644)
	if err != nil {
		log.Error(fmt.Sprintf("保存文件时遇到错误：%s", err.Error()))
		return
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
