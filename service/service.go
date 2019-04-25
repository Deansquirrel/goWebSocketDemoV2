package service

import (
	log "github.com/Deansquirrel/goToolLog"
	"github.com/robfig/cron"
	"time"
)

var s *Server
var c *Client

//启动服务内容
func StartServer() error {
	log.Debug("StartServer")
	go func() {
		s = &Server{}
		s.Start()
	}()
	return nil
}

//启动客户端
func StartClient() error {
	log.Debug("StartClient")
	go func() {
		c = &Client{}
		c.Start()
	}()
	for c == nil {
		time.Sleep(time.Second)
	}
	t := cron.New()
	err := t.AddFunc("0 0 * * * ?", c.DownloadList)
	if err != nil {
		return err
	}
	t.Start()
	return nil
}
