package service

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goWebSocketDemoV2/object"
)

import log "github.com/Deansquirrel/goToolLog"

type common struct {
}

func (c *common) GetCtrlMessage(clientId string, key string, v interface{}) *object.CtrlMessage {
	var data []byte
	data, err := json.Marshal(v)
	if err != nil {
		errMsg := fmt.Sprintf("Get Object[%s] byte error: %s", key, err.Error())
		log.Error(errMsg)
		data = []byte(errMsg)
	}
	cm := object.CtrlMessage{
		Id:   clientId,
		Key:  key,
		Data: string(data),
	}
	return &cm
}

func (c *common) GetRMessage(clientId string, code int, data string) *object.CtrlMessage {
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
	return &rm
}
