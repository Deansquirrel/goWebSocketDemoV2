package service

import (
	"github.com/Deansquirrel/goWebSocketDemoV2/global"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

//启动服务内容
func StartServer() error {
	go func() {
		for i := 0; i < 5; i++ {
			log.Debug("Server")
			time.Sleep(time.Second)
		}
		global.Cancel()
	}()
	return nil
}

func StartClient() error {
	go func() {
		for i := 0; i < 5; i++ {
			log.Debug("Client")
			time.Sleep(time.Second)
		}
		global.Cancel()
	}()
	return nil
}
