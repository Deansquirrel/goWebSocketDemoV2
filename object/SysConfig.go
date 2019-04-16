package object

import (
	"encoding/json"
	"fmt"
	"strings"
)

import log "github.com/Deansquirrel/goToolLog"

//系统配置（Server|Client）
type SystemConfig struct {
	Total   systemConfigTotal   `toml:"total"`
	DB      systemConfigDB      `toml:"configDb"`
	Iris    systemConfigIris    `toml:"iris"`
	Service systemConfigService `toml:"service"`
	Server  systemConfigServer  `toml:"server"`
}

func (sc *SystemConfig) FormatConfig() {
	sc.Total.FormatConfig()
	sc.DB.FormatConfig()
	sc.Iris.FormatConfig()
	sc.Service.FormatConfig()
	sc.Server.FormatConfig()
}

func (sc *SystemConfig) ToString() string {
	d, err := json.Marshal(sc)
	if err != nil {
		log.Warn(fmt.Sprintf("SystemConfig转换为字符串时遇到错误：%s", err.Error()))
		return ""
	}
	return string(d)
}

//通用配置（Server|Client）
type systemConfigTotal struct {
	StdOut   bool   `toml:"stdOut"`
	LogLevel string `toml:"logLevel"`
}

func (t *systemConfigTotal) FormatConfig() {
	//去除首尾空格
	t.LogLevel = strings.Trim(t.LogLevel, " ")
	//设置默认日志级别
	if t.LogLevel == "" {
		t.LogLevel = "warn"
	}
	//设置字符串转换为小写
	t.LogLevel = strings.ToLower(t.LogLevel)
	t.LogLevel = t.checkLogLevel(t.LogLevel)
}

//校验SysConfig中iris日志级别设置（Server|Client）
func (t *systemConfigTotal) checkLogLevel(level string) string {
	switch level {
	case "debug", "info", "warn", "error":
		return level
	default:
		return "warn"
	}
}

//配置库（Server）
type systemConfigDB struct {
	Server string `toml:"server"`
	Port   int    `toml:"port"`
	DbName string `toml:"dbName"`
	User   string `toml:"user"`
	Pwd    string `toml:"pwd"`
}

func (c *systemConfigDB) FormatConfig() {
	c.Server = strings.Trim(c.Server, " ")
	if c.Port == 0 {
		c.Port = 1433
	}
	c.DbName = strings.Trim(c.DbName, " ")
	c.User = strings.Trim(c.User, " ")
	c.Pwd = strings.Trim(c.Pwd, " ")
}

//Iris（Server|Client）
type systemConfigIris struct {
	Port     int    `toml:"port"`
	LogLevel string `toml:"logLevel"`
}

//格式化
func (i *systemConfigIris) FormatConfig() {
	//设置默认端口 8000
	if i.Port == 0 {
		i.Port = 8000
	}
	//去除首尾空格
	i.LogLevel = strings.Trim(i.LogLevel, " ")
	//设置Iris默认日志级别
	if i.LogLevel == "" {
		i.LogLevel = "warn"
	}
	//设置字符串转换为小写
	i.LogLevel = strings.ToLower(i.LogLevel)
}

//服务配置（Server|Client）
type systemConfigService struct {
	Name        string `toml:"name"`
	DisplayName string `toml:"displayName"`
	Description string `toml:"description"`
}

//格式化
func (sc *systemConfigService) FormatConfig() {
	sc.Name = strings.Trim(sc.Name, " ")
	sc.DisplayName = strings.Trim(sc.DisplayName, " ")
	sc.Description = strings.Trim(sc.Description, " ")
	if sc.Name == "" {
		sc.Name = "goWebSocketDemo"
	}
	if sc.DisplayName == "" {
		sc.DisplayName = "goWebSocketDemo"
	}
	if sc.Description == "" {
		sc.Description = sc.Name
	}
}

//服务端地址（Client）
type systemConfigServer struct {
	Address string `toml:"address"`
}

func (sc *systemConfigServer) FormatConfig() {
	sc.Address = strings.Trim(sc.Address, " ")
}
