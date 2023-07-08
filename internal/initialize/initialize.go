package initialize

import (
	"fmt"

	"github.com/issueye/version-mana/internal/global"
	"github.com/issueye/version-mana/internal/gojs"
)

func Initialize() {
	// 初始化运行文件
	InitRuntime()
	// 配置参数
	InitConfig()
	// 日志
	InitLog()
	// 数据
	InitData()
	// http服务
	InitServer()
	// 启动服务
	ShowInfo()
	// 初始化JS
	gojs.InitGojs()
	// 监听服务
	_ = global.HttpServer.ListenAndServe()
}

var (
	AppName string
	Branch  string
	Commit  string
	Date    string
	Version string
)

func ShowInfo() {
	bannerStr := `
                                            ██                                                                          
                                            ▀▀                                                                          
 ██▄  ▄██   ▄████▄    ██▄████  ▄▄█████▄   ████      ▄████▄   ██▄████▄            ████▄██▄   ▄█████▄  ██▄████▄   ▄█████▄ 
  ██  ██   ██▄▄▄▄██   ██▀      ██▄▄▄▄ ▀     ██     ██▀  ▀██  ██▀   ██            ██ ██ ██   ▀ ▄▄▄██  ██▀   ██   ▀ ▄▄▄██ 
  ▀█▄▄█▀   ██▀▀▀▀▀▀   ██        ▀▀▀▀██▄     ██     ██    ██  ██    ██            ██ ██ ██  ▄██▀▀▀██  ██    ██  ▄██▀▀▀██ 
   ████    ▀██▄▄▄▄█   ██       █▄▄▄▄▄██  ▄▄▄██▄▄▄  ▀██▄▄██▀  ██    ██            ██ ██ ██  ██▄▄▄███  ██    ██  ██▄▄▄███ 
    ▀▀       ▀▀▀▀▀    ▀▀        ▀▀▀▀▀▀   ▀▀▀▀▀▀▀▀    ▀▀▀▀    ▀▀    ▀▀            ▀▀ ▀▀ ▀▀   ▀▀▀▀ ▀▀  ▀▀    ▀▀   ▀▀▀▀ ▀▀ 

	版本管理 version-mana
	`
	fmt.Println(bannerStr) // mona12 风格

	info := `
	AppName: %s
	Branch : %s
	Commit : %s
	Date   : %s
	Version: %s
	
	`
	fmt.Printf(info+"\n", AppName, Branch, Commit, Date, Version)
}
