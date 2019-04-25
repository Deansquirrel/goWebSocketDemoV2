package service

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goWebSocketDemoV2/global"
	"github.com/Deansquirrel/goWebSocketDemoV2/object"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strings"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	WebPathTextMessage  = "/webSocket"
	WebPathDownloadFile = "/webSocket/File"
)

type Server struct {
	manager object.IClientManager
}

func (s *Server) Start() {
	s.manager = object.NewClientManager()

	http.HandleFunc(WebPathTextMessage, s.wsPage)
	http.HandleFunc(WebPathDownloadFile, s.wsFile)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", global.SysConfig.Iris.Port), nil)
}

func (s *Server) wsFile(res http.ResponseWriter, req *http.Request) {
	//解析请求
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		errMsg := fmt.Sprintf("解析请求时发生错误，error：%s", err.Error())
		log.Error(errMsg)
		http.NotFound(res, req)
		return
	}
	var fileInfo object.DownloadFile
	err = conn.ReadJSON(&fileInfo)
	if err != nil {
		errMsg := fmt.Sprintf("Get request info error,error: %s", err.Error())
		s.wsFileWriteResponse(conn, -1, errMsg, &fileInfo, nil)
		return
	}

	currPath, err := goToolCommon.GetCurrPath()
	if err != nil {
		errMsg := fmt.Sprintf("server get currPath error,error: %s", err.Error())
		s.wsFileWriteResponse(conn, -1, errMsg, &fileInfo, nil)
		return
	}

	filePath := currPath + "\\" + fileInfo.SubPath + "\\" + fileInfo.Name
	b, err := goToolCommon.PathExists(filePath)
	if err != nil {
		errMsg := fmt.Sprintf("server check file exist error: %s", err.Error())
		s.wsFileWriteResponse(conn, -1, errMsg, &fileInfo, nil)
		return
	}
	if !b {
		errMsg := fmt.Sprintf("file not exist error: %s", fileInfo.SubPath+"\\"+fileInfo.Name)
		s.wsFileWriteResponse(conn, -1, errMsg, &fileInfo, nil)
		return
	}

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		errMsg := fmt.Sprintf("读取文件时遇到错误,error: %s", err.Error())
		s.wsFileWriteResponse(conn, -1, errMsg, &fileInfo, nil)
		return
	}
	s.wsFileWriteResponse(conn, 0, "", &fileInfo, fileData)
	return
}

func (s *Server) wsFileWriteResponse(conn *websocket.Conn, errCode int, errMsg string, fileInfo *object.DownloadFile, data []byte) {
	if errCode != 0 || errMsg != "" {
		log.Error(errMsg)
	}
	rData := object.DownloadFileData{
		ErrCode: errCode,
		ErrMsg:  errMsg,
		Info:    *fileInfo,
		Data:    data,
	}
	err := conn.WriteJSON(rData)
	if err != nil {
		log.Error(fmt.Sprintf("websocket write json error: %s", err.Error()))
	}
}

func (s *Server) wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	c := object.NewClient(goToolCommon.Guid(), conn)

	s.manager.GetChRegister() <- c
	go func() {
		select {
		case <-c.GetChClose():
			s.manager.GetChUnregister() <- c.GetId()
		}
	}()
	go func() {
		for {
			select {
			case msg, ok := <-c.GetChReceive():
				if ok {
					rm := s.msgHandler(&msg)
					if rm != nil {
						c.GetChSend() <- *rm
					}
				}
			case <-c.GetChClose():
				return
			}
		}
	}()
}

func (s *Server) msgHandler(msg *object.CtrlMessage) *object.CtrlMessage {
	comm := common{}
	log.Debug(fmt.Sprintf("Server rec new message,id: %s,key: %s", msg.Id, msg.Key))
	switch msg.Key {
	case object.CtrlMessageHello:
		return s.handlerHello(msg.Id, []byte(msg.Data))
	case object.CtrlMessageUpdateId:
		return s.handlerUpdateId(msg.Id, []byte(msg.Data))
	case object.CtrlMessageDownloadFileList:
		return s.handlerDownloadFileList(msg.Id, []byte(msg.Data))
	case object.CtrlMessageReturn:
		s.handlerReturn(msg.Id, []byte(msg.Data))
		return nil
	default:
		errMsg := fmt.Sprintf("Message Key is not exist,key: %s", msg.Key)
		log.Warn(errMsg)
		return comm.GetRMessage(msg.Id, -1, errMsg)
	}
}

func (s *Server) handlerHello(clientId string, d []byte) *object.CtrlMessage {
	comm := common{}
	var data object.Hello
	err := json.Unmarshal(d, &data)
	if err != nil {
		errMsg := fmt.Sprintf("Get Message Data Hello error: %s", err.Error())
		log.Error(errMsg)
		return comm.GetRMessage(clientId, -1, errMsg)
	}
	log.Debug(fmt.Sprintf("Client %s say hello,msg: %s", clientId, data.Msg))
	return comm.GetRMessage(clientId, 0, fmt.Sprintf("Hello %s", clientId))
}

func (s *Server) handlerUpdateId(clientId string, d []byte) *object.CtrlMessage {
	comm := common{}
	var data object.UpdateId
	err := json.Unmarshal(d, &data)
	if err != nil {
		errMsg := fmt.Sprintf("Get Message Data UpdateId error: %s", err.Error())
		log.Error(errMsg)
		return comm.GetRMessage(clientId, -1, errMsg)
	}
	c := s.manager.GetClient(clientId)
	if c != nil {
		s.manager.GetChUnregister() <- clientId
		c.SetId(data.Id)
		s.manager.GetChRegister() <- c
		return comm.GetRMessage(data.Id, 0, "success")
	}
	return comm.GetRMessage(data.Id, -1, fmt.Sprintf("ClientId is not exists,id: %s", clientId))
}

func (s *Server) handlerDownloadFileList(clientId string, d []byte) *object.CtrlMessage {
	comm := common{}
	var data object.DownloadFileList
	err := json.Unmarshal(d, &data)
	if err != nil {
		errMsg := fmt.Sprintf("Get Message Data DownloadFileList error: %s", err.Error())
		log.Error(errMsg)
		return comm.GetRMessage(clientId, -1, errMsg)
	}
	currPath, err := goToolCommon.GetCurrPath()
	if err != nil {
		errMsg := fmt.Sprintf("Server get currPath error: %s", err.Error())
		log.Error(errMsg)
		return comm.GetRMessage(clientId, -1, errMsg)
	}
	subPath := currPath + data.SubPath
	isPathExist, err := goToolCommon.PathExists(subPath)
	if err != nil {
		errMsg := fmt.Sprintf("Server check subPath[%s] error: %s", data.SubPath, err.Error())
		log.Error(errMsg)
		return comm.GetRMessage(clientId, -1, errMsg)
	}
	if !isPathExist {
		errMsg := fmt.Sprintf("Server subPath[%s] is not exist", data.SubPath)
		log.Error(errMsg)
		return comm.GetRMessage(clientId, -1, errMsg)
	}
	list := s.getDownloadFileList(subPath)
	for _, f := range list {
		f.SubPath = strings.Replace(f.SubPath, currPath, "", -1)
	}
	rList := make([]object.DownloadFile, 0)
	for _, f := range list {
		rList = append(rList, *f)
	}
	return comm.GetCtrlMessage(clientId, object.CtrlMessageDownloadFileList, rList)
}

func (s *Server) getDownloadFileList(path string) []*object.DownloadFile {
	list := make([]*object.DownloadFile, 0)
	folderList, fileList, err := goToolCommon.GetFolderAndFileList(path)
	if err != nil {
		log.Error(fmt.Sprintf("getDownloadFileList %s", err.Error()))
		return list
	}
	for _, f := range fileList {
		list = append(list, &object.DownloadFile{
			Name:    f,
			SubPath: path,
			MD5:     "",
		})
	}
	for _, folder := range folderList {
		tList := s.getDownloadFileList(path + "\\" + folder)
		for _, f := range tList {
			list = append(list, f)
		}
	}
	return list
}

func (s *Server) handlerReturn(clientId string, d []byte) {
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

func (s *Server) sendCtrlMessage(clientId, key string, v interface{}) {
	c := s.manager.GetClient(clientId)
	if c != nil {
		comm := common{}
		c.GetChSend() <- *comm.GetCtrlMessage(c.GetId(), key, v)
	}
}
