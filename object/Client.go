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
	chReceive chan CtrlMessage
	chSend    chan CtrlMessage
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
		chReceive: make(chan CtrlMessage),
		chSend:    make(chan CtrlMessage),
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

func (c *client) GetChReceive() <-chan CtrlMessage {
	if c.isStopped || c.chReceive == nil {
		return nil
	}
	return c.chReceive
}

func (c *client) GetChSend() chan<- CtrlMessage {
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
	defer func() {
		if !c.isStopped {
			c.Close()
		}
	}()
	for {
		var cm CtrlMessage
		err := c.socket.ReadJSON(&cm)
		if err != nil {
			log.Error(fmt.Sprintf("Read error:%s,ClientID:%s", err.Error(), c.GetId()))
			return
		}
		cm.Id = c.id
		c.chReceive <- cm
	}
}

func (c *client) write() {
	//log.Debug(fmt.Sprintf("Client write start,id:%s", c.id))
	//defer log.Debug(fmt.Sprintf("Client write exit,id:%s", c.id))
	for {
		select {
		case msg, ok := <-c.chSend:
			log.Debug(fmt.Sprintf("Get Send message %v %v", msg, ok))
			if ok {
				//c.writeData(&msg)
				err := c.socket.WriteJSON(msg)
				if err != nil {
					log.Error(err.Error())
				}
			}
		case <-c.ctx.Done():
			return
		}
	}
}
