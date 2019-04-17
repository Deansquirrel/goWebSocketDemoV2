package object

type OprMessage struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Data string `json:"data"`
}

//File
type File struct {
	Name string `json:"name"`
}

type ReFile struct {
	ErrCode int
	ErrMsg  string
	Data    []byte
}
