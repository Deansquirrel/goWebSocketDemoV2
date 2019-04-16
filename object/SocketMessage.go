package object

type SocketMessage struct {
	ClientId    string
	MessageType int
	Data        []byte
}
