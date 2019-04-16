package object

type IClientManager interface {
	GetChRegister() chan<- IClient
	GetChUnregister() chan<- string
	GetChBroadcast() chan<- *SocketMessage
	GetClient(id string) IClient
}
