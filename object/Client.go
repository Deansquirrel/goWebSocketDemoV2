package object

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

type client struct {
	id        string
	socket    *websocket.Conn
	chReceive chan SocketMessage
	chSend    chan SocketMessage
	lock      sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	//退出控制
	isStopped bool
}

func NewClient(id string, socket *websocket.Conn) *client {
	c := client{
		id:        id,
		socket:    socket,
		chReceive: make(chan SocketMessage),
		chSend:    make(chan SocketMessage),
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.start()
	return &c
}

func (c *client) GetId() string {
	return c.id
}

func (c *client) SetId(id string) {
	c.id = id
}

func (c *client) GetChReceive() <-chan SocketMessage {
	if c.isStopped || c.chReceive == nil {
		return nil
	}
	return c.chReceive
}

func (c *client) GetChSend() chan<- SocketMessage {
	if c.isStopped || c.chSend == nil {
		return nil
	}
	return c.chSend
}

func (c *client) GetChClose() <-chan struct{} {
	if c.isStopped {
		return nil
	}
	return c.ctx.Done()
}

func (c *client) Close() {
	if c.isStopped {
		return
	}
	c.isStopped = true
	c.cancel()
	time.Sleep(time.Second)
	_ = c.socket.Close()
	close(c.chReceive)
	c.chReceive = nil
	close(c.chSend)
	c.chSend = nil
	c.ctx = nil
	c.cancel = nil
}

func (c *client) start() {
	go c.read()
	go c.write()
}

func (c *client) read() {
	//log.Debug(fmt.Sprintf("Client read start,id:%s", c.id))
	//defer log.Debug(fmt.Sprintf("Client read exit,id:%s", c.id))
	for {
		t, d, err := c.socket.ReadMessage()
		if err != nil {
			log.Error(fmt.Sprintf("Read error:%s,ClientID:%s", err.Error(), c.GetId()))
			if !c.isStopped {
				c.Close()
			}
			return
		}
		m := SocketMessage{
			ClientId:    c.id,
			MessageType: t,
			Data:        d,
		}
		c.chReceive <- m
	}
}

func (c *client) write() {
	//log.Debug(fmt.Sprintf("Client write start,id:%s", c.id))
	//defer log.Debug(fmt.Sprintf("Client write exit,id:%s", c.id))
	for {
		select {
		case msg, ok := <-c.chSend:
			if ok {
				c.writeData(&msg)
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *client) writeData(msg *SocketMessage) {
	//log.Debug(fmt.Sprintf("Client write data start,id:%s,dateLength:%d", c.id, len(msg.Data)))
	//defer log.Debug(fmt.Sprintf("Client write data exit,id:%s", c.id))
	err := c.socket.WriteMessage(msg.MessageType, msg.Data)
	if err != nil {
		log.Error(err.Error())
	}
}
