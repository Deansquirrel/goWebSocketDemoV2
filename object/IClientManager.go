package object

type IClientManager interface {
	GetChRegister() chan<- IClient
	GetChUnregister() chan<- string
	GetChBroadcast() chan<- CtrlMessage
	GetIdList() []string
	GetClient(id string) IClient
}
