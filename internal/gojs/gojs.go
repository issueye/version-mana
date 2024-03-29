package gojs

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dop251/goja"
	licheeJs "github.com/issueye/lichee-js"
	"github.com/issueye/version-mana/internal/global"
	"github.com/issueye/version-mana/internal/logic"
	"github.com/issueye/version-mana/internal/model"
	"github.com/issueye/version-mana/internal/service"
	"github.com/issueye/version-mana/pkg/utils"
	"github.com/issueye/version-mana/pkg/ws"
)

func InitGojs() {
	go gojs()
}

var (
	VersId  = make(chan string, 10)
	RunCode = make(chan *model.RepoCode, 10)
)

type WriteFileCallBack func(*FileOperation)

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

	c := licheeJs.NewCore()
	c.SetLogOutMode(licheeJs.LOM_DEBUG)
	workDir := utils.GetWorkDir()
	logDir := filepath.Join(workDir, "runtime", "logs")
	c.SetLogPath(logDir)

	// 如果是测试运行，则之传入仓库信息

	var (
		repo *model.Repository
		vers *model.AppVersionInfo
		err  error
		code string
		t    int // 类型 0 测试编译 1 编译
	)
	if len(args) > 0 {
		repo, vers, err = injectRepoInfo(id, workDir, c)
		if err != nil {
			sendMessage(id, err.Error())
			return
		}

		code = args[0].(string)
		t = 0
	} else {
		repo, vers, err = injectVersionInfo(id, workDir, c)
		if err != nil {
			sendMessage(id, err.Error())
			return
		}

		code = repo.Code
		t = 1
	}

	repoDir := filepath.Join(workDir, "runtime", "git_repo", repo.ServerName, vers.AppName)
	c.SetGlobalProperty("vers", vers)
	c.SetGlobalProperty("repo", repo)
	c.SetGlobalProperty("workDir", repoDir)
	c.SetGlobalProperty("runType", t)

	cacheDir := filepath.Join(workDir, "runtime", "git_repo", "cache")
	c.SetGlobalProperty("cacheDir", cacheDir)

	// 注册模块
	InitVm(c, id)

	err = c.RunString("test", code)
	if err != nil {
		fmt.Println("运行错误，错误原因：", err.Error())
		sendMessage(id, err.Error())
		return
	}
}

// 保存文件信息
func Save(id string, t int) {
	logic.NewRepo().CreateRelease(id, t)
}

func sendMessage(id, msg string) {
	value, ok := ws.SMap.Load(id)
	if ok {
		wc := value.(*ws.WsConn)
		err := wc.OutChanWrite([]byte(msg))
		if err != nil {
			global.Log.Errorf("传送输出信息失败，失败原因：%s", err.Error())
			return
		}
	}
}

// 注入仓库信息
func injectRepoInfo(id string, workDir string, c *licheeJs.Core) (*model.Repository, *model.AppVersionInfo, error) {
	// 获取仓库信息
	r, err := service.NewRepo(global.DB).GetById(id)
	if err != nil {
		return nil, nil, err
	}

	// 随机选择当前仓库的任意版本
	vers, err := service.NewRepo(global.DB).GetVerByRepoId(id)
	if err != nil {
		return nil, nil, err
	}

	return r, vers, nil
}

func injectVersionInfo(id string, workDir string, c *licheeJs.Core) (*model.Repository, *model.AppVersionInfo, error) {
	// 获取版本信息
	version, err := service.NewRepo(global.DB).GetVersionById(id)
	if err != nil {
		return nil, nil, err
	}

	// 注入仓库信息
	r, err := service.NewRepo(global.DB).GetById(version.RepoID)
	if err != nil {
		return nil, nil, err
	}

	return r, version, nil
}

func InitVm(c *licheeJs.Core, id string) {
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

	// 文件夹不存在时创建
	vm.Set("createNotExists", func(path string) {
		_, err := utils.PathExists(path)
		if err != nil {
			sendMessage(id, fmt.Sprintf("创建文件夹失败，失败原因：%s", err.Error()))
			return
		}
	})

	// 文件夹存在时，删除文件夹
	vm.Set("removeExists", func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			sendMessage(id, fmt.Sprintf("删除文件夹失败，失败原因：%s", err.Error()))
			return
		}
	})

	// 移动文件
	vm.Set("moveFile", func(old, new string) bool {
		err := os.MkdirAll(filepath.Dir(new), os.ModePerm)
		if err != nil {
			sendMessage(id, fmt.Sprintf("移动文件失败，失败原因：%s", err.Error()))
			return false
		}

		err = os.Rename(old, new)
		if err != nil {
			sendMessage(id, fmt.Sprintf("移动文件失败，失败原因：%s", err.Error()))
			return false
		}

		return true
	})

	vm.Set("openFile", func(path string, back WriteFileCallBack) string {
		f2, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err.Error()
		}
		defer f2.Close()
		fo := CreateFileOp(f2)
		back(fo)
		fo.Flush()
		return ""
	})

	// md5
	vm.Set("fileMD5", func(path string) string {
		f, err := os.Open(path)
		if err != nil {
			fmt.Errorf("打开文件失败，filename=%v, err=%v", path, err)
			return ""
		}

		defer f.Close()

		md5h := md5.New()
		io.Copy(md5h, f)
		return hex.EncodeToString(md5h.Sum(nil))
	})

	// sha256
	vm.Set("fileSha256", func(path string) string {
		hash := sha256.New()
		f, err := os.Open(path)
		if err != nil {
			fmt.Errorf("打开文件失败，filename=%v, err=%v", path, err)
			return ""
		}
		defer f.Close()
		if _, err := io.Copy(hash, f); err != nil {
			fmt.Errorf("io 拷贝失败，filename=%v, err=%v", path, err)
		}
		bytes := hash.Sum(nil)
		//将字符串编码为16进制格式,返回字符串
		hashCode := hex.EncodeToString(bytes)
		//返回哈希值
		return hashCode
	})

	// 获取当前程序的根目录
	vm.Set("rootPath", func() string {
		return utils.GetWorkDir()
	})

	vm.Set("createRelease", func(id string, t int) {
		Save(id, t)
	})

	// 回调
	c.ConsoleCallBack = func(args ...any) {
		if len(args) > 0 {
			data := args[0].(string)
			sendMessage(id, data)
		}
	}

	c.RegisterModule("vmm/exec", func(vm *goja.Runtime, module *goja.Object) {
		module.Set("new", func() goja.Value {
			return NewExecJs(vm, NewExec())
		})
	})
}
