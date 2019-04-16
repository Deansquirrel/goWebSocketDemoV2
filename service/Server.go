package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goWebSocketDemoV2/object"
	"github.com/gorilla/websocket"
	"net/http"
)

import log "github.com/Deansquirrel/goToolLog"

type Server struct {
	manager object.IClientManager
}

func (s *Server) Start() {
	s.manager = object.NewClientManager()
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
			rm := s.msgHandler(msg)
			if rm != nil {
				c.GetChSend() <- rm
			}
		}
	}()
}

func (s *Server) msgHandler(msg *object.SocketMessage) *object.SocketMessage {
	if msg.MessageType != websocket.TextMessage {
		errMsg := fmt.Sprintf("Unexpected MessageType : %d", msg.MessageType)
		log.Error(errMsg)
		return s.getRMessage(msg.ClientId, -1, errMsg)
	}
}

func (s *Server) getRMessage(clientId string, code int, data string) *object.SocketMessage {
	m := object.ReturnMessage{
		ErrCode: code,
		ErrMsg:  data,
	}
	d, err := json.Marshal(m)
	if err != nil {
		log.Error(fmt.Sprintf("getRMessage error: %s", err.Error()))
		return nil
	}
	return &object.SocketMessage{
		ClientId:    clientId,
		MessageType: websocket.TextMessage,
		Data:        d,
	}
}

func (s *Server) updateClientId(old, new string) {
	c := s.manager.GetClient(old)
	if c != nil {
		s.manager.GetChUnregister() <- old
		c.SetId(new)
		s.manager.GetChRegister() <- c
	}
}
