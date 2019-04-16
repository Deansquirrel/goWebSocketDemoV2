package object

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

type ReturnMessage struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
