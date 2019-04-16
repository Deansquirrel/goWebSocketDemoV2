package object

type IClient interface {
	GetId() string
	SetId(id string)
	GetChReceive() <-chan *SocketMessage
	GetChSend() chan<- *SocketMessage
	GetChClose() <-chan struct{}
	Close()
}
