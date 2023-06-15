package gojs

import (
	"fmt"
	"path/filepath"

	"github.com/dop251/goja"
	licheeJs "github.com/issueye/lichee-js"
	"github.com/issueye/version-mana/internal/global"
	"github.com/issueye/version-mana/internal/model"
	"github.com/issueye/version-mana/internal/service"
	"github.com/issueye/version-mana/pkg/utils"
)

func InitGojs() {
	go gojs()
}

var (
	VersId  = make(chan string, 10)
	RunCode = make(chan *model.RepoCode, 10)
)

func gojs() {
	for {
		select {
		case id := <-VersId:
			{
				runCompileScript(id)
			}
		case code := <-RunCode:
			{
				runCompileScript(code.ID, code.Code)
			}
		}
	}
}

func runCompileScript(id string, args ...any) {
	// 获取id信息
	info, err := service.NewRepo(global.DB).GetVersionById(id)
	if err != nil {
		return
	}

	// 获取仓库信息
	r, err := service.NewRepo(global.DB).GetById(info.RepoID)
	if err != nil {
		return
	}

	c := licheeJs.NewCore()
	c.SetLogOutMode(licheeJs.LOM_DEBUG)
	workDir := utils.GetWorkDir()
	repoDir := filepath.Join(workDir, "runtime", "git_repo", r.ServerName, info.AppName)
	workDir = filepath.Join(workDir, "runtime", "logs")
	c.SetLogPath(workDir)

	// 添加 js 变量
	c.SetGlobalProperty("repo_info", r)
	c.SetGlobalProperty("workDir", repoDir)

	code := r.Code
	if len(args) > 0 {
		code = args[0].(string)
	} else {
		c.SetGlobalProperty("version_info", info)
	}

	// 注册模块
	InitVm(c)

	err = c.RunString("test", code)
	if err != nil {
		fmt.Println("运行错误，错误原因：", err.Error())
		return
	}
}

func InitVm(c *licheeJs.Core) {
	// 运行 command 命令
	vm := c.GetRts()

	vm.Set("convToString", func(data []byte, t int) string {
		switch t {
		case 0:
			return ConvertByte2String(data, UTF8)
		case 1:
			return ConvertByte2String(data, GB18030)
		default:
			return string(data)
		}
	})

	c.RegisterModule("vmm/exec", func(vm *goja.Runtime, module *goja.Object) {
		module.Set("new", func() goja.Value {
			return NewExecJs(vm, NewExec())
		})
	})
}
