package service

import "time"

import log "github.com/Deansquirrel/goToolLog"

var s Server
var c Client

//启动服务内容
func StartServer() error {
	log.Debug("StartServer")
	go func() {
		s = Server{}
		s.Start()
	}()
	return nil
}

//启动客户端
func StartClient() error {
	log.Debug("StartClient")
	go func() {
		c = Client{}
		c.Start()
	}()
	go func() {
		time.AfterFunc(time.Second*5, func() {
			c.SayHello()
		})
	}()
	go func() {
		time.AfterFunc(time.Second*10, func() {
			c.GetFile()
		})
	}()
	return nil
}
