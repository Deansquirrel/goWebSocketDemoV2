package object

type IClientManager interface {
	GetChRegister() chan<- IClient
	GetChUnregister() chan<- string
	GetChBroadcast() chan<- SocketMessage
	GetIdList() []string
	GetClient(id string) IClient
}
