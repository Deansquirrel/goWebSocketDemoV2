package service

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goWebSocketDemoV2/object"
	"github.com/gorilla/websocket"
)

import log "github.com/Deansquirrel/goToolLog"

type common struct {
}

func (c *common) GetRMessage(clientId string, code int, data string) *object.SocketMessage {
	m := object.ReturnMessage{
		ErrCode: code,
		ErrMsg:  data,
	}
	d, err := json.Marshal(m)
	if err != nil {
		log.Error(fmt.Sprintf("getRMessage error: %s", err.Error()))
		return nil
	}
	rm := object.CtrlMessage{
		Id:   clientId,
		Key:  object.CtrlMessageReturn,
		Data: string(d),
	}
	rd, err := json.Marshal(rm)
	if err != nil {
		log.Error(fmt.Sprintf("getRMessage getRd error: %s", err.Error()))
		return nil
	}
	return &object.SocketMessage{
		ClientId:    clientId,
		MessageType: websocket.TextMessage,
		Data:        rd,
	}
}
