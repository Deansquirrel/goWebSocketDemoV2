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
)

import log "github.com/Deansquirrel/goToolLog"

type Server struct {
	manager object.IClientManager
}

func (s *Server) Start() {
	s.manager = object.NewClientManager()

	http.HandleFunc("/webSocket", s.wsPage)
	http.HandleFunc("/webSocket/File", s.wsFile)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", global.SysConfig.Iris.Port), nil)
}

func (s *Server) wsFile(res http.ResponseWriter, req *http.Request) {
	//===================================================================================================
	//解析请求
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}
	var fileInfo object.OprMessage
	err = conn.ReadJSON(&fileInfo)
	if err != nil {
		errMsg := fmt.Sprintf("Get request info error,error: %s", err.Error())
		s.writeFileReturnMessage(conn, -1, errMsg, nil)
		return
	}
	log.Debug(fmt.Sprintf("Client [%s] [%s]", fileInfo.Id, fileInfo.Key))
	var f object.File
	err = json.Unmarshal([]byte(fileInfo.Data), &f)
	if err != nil {
		errMsg := fmt.Sprintf("Get FileInfo error: %s", err.Error())
		log.Error(errMsg)
		s.writeFileReturnMessage(conn, -1, errMsg, nil)
		return
	}
	//===================================================================================================
	path, err := goToolCommon.GetCurrPath()
	if err != nil {
		errMsg := fmt.Sprintf("Get CurrPath error: %s", err.Error())
		log.Error(errMsg)
		s.writeFileReturnMessage(conn, -1, errMsg, nil)
		return
	}
	filePath := path + "//" + "File" + "//" + f.Name
	b, err := goToolCommon.PathExists(filePath)
	if err != nil {
		errMsg := fmt.Sprintf("检查文件是否存在时遇到错误，error: %s", err.Error())
		log.Error(errMsg)
		s.writeFileReturnMessage(conn, -1, errMsg, nil)
		return
	}
	if !b {
		errMsg := "请求的文件不存在"
		log.Error(errMsg)
		log.Error(filePath)
		s.writeFileReturnMessage(conn, -1, errMsg, nil)
		return
	}
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		if err != nil {
			errMsg := fmt.Sprintf("读取文件时遇到错误，error: %s", err.Error())
			log.Error(errMsg)
			s.writeFileReturnMessage(conn, -1, errMsg, nil)
			return
		}
	}
	s.writeFileReturnMessage(conn, 0, "", fileData)
	return
}

func (s *Server) writeFileReturnMessage(conn *websocket.Conn, errCode int, errMsg string, data []byte) {
	writeErr := conn.WriteJSON(&object.ReFile{
		ErrCode: errCode,
		ErrMsg:  errMsg,
		Data:    data,
	})
	if writeErr != nil {
		log.Error(fmt.Sprintf("Write Response error: %s", writeErr.Error()))
	}
	return
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
		select {
		case msg := <-c.GetChReceive():
			rm := s.msgHandler(&msg)
			if rm != nil {
				c.GetChSend() <- *rm
			}
		}
	}()
}

func (s *Server) clientTest(c object.IClient) bool {
	return false
}

func (s *Server) msgHandler(msg *object.SocketMessage) *object.SocketMessage {
	comm := common{}

	if msg.MessageType != websocket.TextMessage {
		errMsg := fmt.Sprintf("Unexpected MessageType : %d", msg.MessageType)
		log.Error(errMsg)
		return comm.GetRMessage(msg.ClientId, -1, errMsg)
	}
	var m object.CtrlMessage
	err := json.Unmarshal(msg.Data, &m)
	if err != nil {
		errMsg := fmt.Sprintf("Get Message Object error: %s", err.Error())
		log.Error(errMsg)
		return comm.GetRMessage(msg.ClientId, -1, errMsg)
	}
	log.Debug(fmt.Sprintf("Server rec new message,id: %s,key: %s", m.Id, m.Key))
	//===================================================================================
	switch m.Key {
	case object.CtrlMessageHello:
		return s.handlerHello(msg.ClientId, []byte(m.Data))
	case object.CtrlMessageUpdateId:
		return s.handlerUpdateId(msg.ClientId, []byte(m.Data))
	case object.CtrlMessageReturn:
		s.handlerReturn(msg.ClientId, []byte(m.Data))
		return nil
	default:
		errMsg := fmt.Sprintf("Message Key is not exist,key: %s", m.Key)
		log.Warn(errMsg)
		return comm.GetRMessage(msg.ClientId, -1, errMsg)
	}
}

func (s *Server) handlerHello(clientId string, d []byte) *object.SocketMessage {
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

func (s *Server) handlerUpdateId(clientId string, d []byte) *object.SocketMessage {
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

func (s Server) handlerReturn(clientId string, d []byte) {
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
