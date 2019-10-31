package config

import (
	"log"
	"runtime/debug"

	"github.com/lytics/confl"
)

const (
	// json格式配置文件
	CONF_FILE = "./confhttp.ini"
)

var (
	Conf = new(Config)
)

type Config struct {
	Debug bool // 是否调试模式
	Mail  Mail
	Log   Log
	Redis Redis
	App   App
}

// 服务邮件配置
type Mail struct {
	Title     string
	Alias     string
	Host      string // smtp.163.com
	Port      int    // 25
	User      string
	Pwd       string
	Receivers []string // 错误邮件接收者列表
}

// 扩展部分配置
type Log struct {
	// 日志文件目录
	LogDir string
	// 错误日志配置
	LogErrMbytes     int
	LogErrMaxDays    int
	LogErrMaxBackups int
	// 通用日志配置
	LogNormMbytes     int
	LogNormMaxDays    int
	LogNormMaxBackups int
}

// Redis数据库配置
type Redis struct {
	MaxActive int
	MaxIdle   int
}

// Http相关
type App struct {
	HttpCharset          string
	DefaultMapSize       int
	FormatSignSaltClient string
	FormatSignSaltServ   string
}

func init() {
	if _, err := confl.DecodeFile(CONF_FILE, Conf); err != nil {
		log.Fatalf("[Config] read conf file error. info=%s trace=%s\n", err.Error(), string(debug.Stack()))
	} else {
		// fmt.Println("[Config] read conf ok.", Conf)
	}
}
