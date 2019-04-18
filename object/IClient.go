package object

type IClient interface {
	GetId() string
	SetId(id string)
	GetChReceive() <-chan CtrlMessage
	GetChSend() chan<- CtrlMessage
	GetChClose() <-chan struct{}
	Close()
}
