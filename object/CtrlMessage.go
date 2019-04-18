package object

const (
	//===============================================================
	//Server + Client
	CtrlMessageReturn = "return"
	//===============================================================
	//Server
	CtrlMessageHello            = "hello"
	CtrlMessageUpdateId         = "updateId"
	CtrlMessageDownloadFileList = "downloadFileList"
	//===============================================================
	//Client
	CtrlMessageTest               = "test"
	CtrlMessageReDownloadFileList = "reDownloadFileList"
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

//downloadFile
type DownloadFileList struct {
	SubPath string `json:"subpath"`
}

//reDownloadFile
type ReDownloadFileList struct {
	FileList []DownloadFile `json:"filelist"`
}

//downloadFileObject
type DownloadFile struct {
	Name    string `json:"name"`
	SubPath string `json:"subpath"`
	MD5     string `json:"md5"`
}

//reDownloadFileObject
type ReDownloadFile struct {
	ErrCode int
	ErrMsg  string
	Data    []byte
}
