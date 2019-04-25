package global

import (
	"context"
	"github.com/Deansquirrel/goWebSocketDemoV2/object"
)

const (
	//PreVersion = "0.0.0 Build20190101"
	//TestVersion = "0.0.0 Build20190101"
	Version = "0.0.0 Build20190101"
)

const (
	//http连接超时时间
	HttpConnectTimeout = 30
	//重连时间
	ReConnectDuration = 30
	//心跳时间间隔
	ClientHeartBeatDuration = 60
	//客户端下载文件夹名称
	ClientFileFolder = "clientFile"
	//最大线程数
	//MaxThread = 10
)

var Ctx context.Context
var Cancel func()

//程序启动参数
var Args *object.ProgramArgs

//配置文件是否存在
//var IsConfigExists bool
//系统参数
var SysConfig *object.SystemConfig
