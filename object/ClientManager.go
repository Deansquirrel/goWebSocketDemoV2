package object

import (
	"fmt"
	"sync"
)

import log "github.com/Deansquirrel/goToolLog"

type clientManager struct {
	clients      map[string]IClient
	chRegister   chan IClient
	chUnregister chan string
	chBroadcast  chan SocketMessage
	lock         sync.Mutex
}

func NewClientManager() *clientManager {
	cm := clientManager{
		clients:      make(map[string]IClient),
		chRegister:   make(chan IClient),
		chUnregister: make(chan string),
		chBroadcast:  make(chan SocketMessage),
	}
	cm.start()
	return &cm
}

func (manager *clientManager) GetIdList() []string {
	list := make([]string, 0)
	for key := range manager.clients {
		list = append(list, key)
	}
	return list
}

func (manager *clientManager) GetClient(id string) IClient {
	client, ok := manager.clients[id]
	if ok {
		return client
	}
	return nil
}

func (manager *clientManager) GetChRegister() chan<- IClient {
	return manager.chRegister
}

func (manager *clientManager) GetChUnregister() chan<- string {
	return manager.chUnregister
}

func (manager *clientManager) GetChBroadcast() chan<- SocketMessage {
	return manager.chBroadcast
}

func (manager *clientManager) start() {
	go func() {
		for {
			select {
			case c := <-manager.chRegister:
				manager.register(c)
			case c := <-manager.chUnregister:
				manager.unregister(c)
			case m := <-manager.chBroadcast:
				manager.broad(m)
			}
		}
	}()
}

func (manager *clientManager) register(c IClient) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	_, ok := manager.clients[c.GetId()]
	if ok {
		log.Error(fmt.Sprintf("ClientId %s is already exist", c.GetId()))
		return
	}
	manager.clients[c.GetId()] = c
	log.Info(fmt.Sprintf("Client Register: %s，CurrClientNum: %d", c.GetId(), len(manager.clients)))
	return
}

func (manager *clientManager) unregister(id string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	log.Info(fmt.Sprintf("Client Unregister: %s", id))
	_, ok := manager.clients[id]
	if ok {
		delete(manager.clients, id)
	}
}

func (manager *clientManager) broad(msg SocketMessage) {
	go func() {
		for _, c := range manager.clients {
			log.Debug(fmt.Sprintf("Board msg: Type:%d,Msg:%s", msg.MessageType, msg.Data))
			c.GetChSend() <- msg
		}
	}()
}
