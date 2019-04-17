package object

const (
	//===============================================================
	//Server + Client
	CtrlMessageReturn = "return"
	//===============================================================
	//Server
	CtrlMessageHello    = "hello"
	CtrlMessageUpdateId = "updateId"
	//===============================================================
	//Client
	CtrlMessageTest = "test"
	//===============================================================
)

type CtrlMessage struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Data string `json:"data"`
}

//hello
type Hello struct {
	Msg string `json:"msg"`
}

//updateId
type UpdateId struct {
	Id string `json:"id"`
}

//return
type ReturnMessage struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
