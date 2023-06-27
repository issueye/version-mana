package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/issueye/version-mana/internal/global"
	"github.com/issueye/version-mana/internal/gojs"
	"github.com/issueye/version-mana/internal/logic"
	"github.com/issueye/version-mana/internal/model"
	"github.com/issueye/version-mana/internal/service"
	"github.com/issueye/version-mana/pkg/ws"
)

// websocket 升级并跨域
var (
	upgrade = &websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// 日志查看 ws
func WsScriptTestRunConsole(ctx *gin.Context) {
	control := New(ctx)
	id := control.Param("id")
	if id == "" {
		control.FailBind(errors.New("[id]不能为空"))
		return
	}

	// 升级为 websocket
	conn, err := upgrade.Upgrade(control.Writer, control.Request, nil)
	if err != nil {
		control.FailByMsgf("升级协议失败，失败原因：%s", err.Error())
	}

	fmt.Println("id => websocket", id)

	ws.NewConn(id, conn)
	control.Success()
}

// 日志查看 ws
func WsScriptCompileConsole(ctx *gin.Context) {
	control := New(ctx)
	id := control.Param("id")
	if id == "" {
		control.FailBind(errors.New("[id]不能为空"))
		return
	}

	// 升级为 websocket
	conn, err := upgrade.Upgrade(control.Writer, control.Request, nil)
	if err != nil {
		control.FailByMsgf("升级协议失败，失败原因：%s", err.Error())
	}

	ws.NewConn(id, conn)
	control.Success()
}

type RepoController struct {
	Controller
}

func NewRepoController() *RepoController {
	return new(RepoController)
}

// 创建一个代码仓库信息
func (RepoController) Create(ctx *gin.Context) {
	control := New(ctx)

	// 绑定参数
	req := new(model.CreateRepository)
	err := control.Bind(req)
	if err != nil {
		control.FailBind(err)
		return
	}

	err = service.NewRepo(global.DB).Create(req)
	if err != nil {
		control.FailByMsgf("创建代码仓库信息失败，失败原因：%s", err.Error())
		return
	}

	control.Success()
}

// 查询代码仓库列表
func (RepoController) GetById(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("id")
	if id == "" {
		control.FailBind(errors.New("[id]不能为空"))
		return
	}

	data, err := service.NewRepo(global.DB).GetById(id)
	if err != nil {
		global.Log.Errorf("查询代码仓库信息列表失败，失败原因： %s", err.Error())
		control.FailByMsg("查询代码仓库信息列表失败")
		return
	}

	control.SuccessData(data)
}

// 查询代码仓库列表
func (RepoController) List(ctx *gin.Context) {
	control := New(ctx)

	req := new(model.QueryRepository)
	err := control.Bind(req)
	if err != nil {
		global.Log.Errorf("绑定请求内容失败，失败原因：%s", err.Error())
		control.FailBind(err)
		return
	}

	list, err := service.NewRepo(global.DB).Query(req)
	if err != nil {
		global.Log.Errorf("查询代码仓库信息列表失败，失败原因： %s", err.Error())
		control.FailByMsg("查询代码仓库信息列表失败")
		return
	}

	control.SuccessAutoData(req, list)
}

// 修改代码仓库信息
func (RepoController) Modify(ctx *gin.Context) {
	control := New(ctx)

	req := new(model.ModifyRepository)
	err := ctx.Bind(req)
	if err != nil {
		global.Log.Errorf("绑定参数失败，失败原因：%s", err.Error())
		control.FailBind(err)
		return
	}

	err = service.NewRepo(global.DB).Modify(req)
	if err != nil {
		control.FailByMsgf("修改代码仓库信息失败，失败原因：%s", err.Error())
		return
	}

	control.Success()
}

// 删除代码仓库信息
func (RepoController) Delete(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("id")
	if id == "" {
		control.FailBind(errors.New("[id]不能为空"))
		return
	}

	// 查询代码仓库信息
	info, err := service.NewRepo(global.DB).GetById(id)
	if err != nil {
		control.FailByMsgf("查询代码仓库信息失败，失败原因：%s", err.Error())
		return
	}

	if info.ID == "" {
		control.FailByMsg("未查询到代码仓库信息")
		return
	}

	err = service.NewRepo(global.DB).DelById(id)
	if err != nil {
		control.FailByMsgf("删除代码仓库信息失败，失败原因：%s", err.Error())
		return
	}

	control.Success()
}

// 获取分支列表
func (RepoController) BranchList(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("id")
	if id == "" {
		control.FailBind(errors.New("[id]不能为空"))
		return
	}

	list, err := logic.NewRepo().BranchList(id)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	control.SuccessData(list)
}

// 获取分支的发布类型的最大版号
func (RepoController) GetLastVerNum(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("repoId")
	if id == "" {
		control.FailBind(errors.New("[repoId]不能为空"))
		return
	}

	req := new(model.QryLastVer)
	err := control.Bind(req)
	if err != nil {
		control.FailBind(err)
		return
	}

	data, err := service.NewRepo(global.DB).GetLastVerNum(id, req)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	control.SuccessData(data)
}

// 删除代码仓库信息
func (RepoController) RemoveVersion(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("id")
	if id == "" {
		control.FailBind(errors.New("[id]不能为空"))
		return
	}

	// 移除版本
	err := logic.NewRepo().RemoveVersion(id)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	control.Success()
}

// 删除代码仓库信息
func (RepoController) ModifyCode(ctx *gin.Context) {
	control := New(ctx)

	req := new(model.RepoCode)
	err := control.Bind(req)
	if err != nil {
		control.FailBind(err)
		return
	}

	err = service.NewRepo(global.DB).ModifyCode(req)
	if err != nil {
		control.FailByMsgf("修改代码失败，失败原因：%s", err.Error())
		return
	}

	control.Success()
}

// 删除代码仓库信息
func (RepoController) TestRun(ctx *gin.Context) {
	control := New(ctx)

	req := new(model.RepoCode)
	err := control.Bind(req)
	if err != nil {
		control.FailBind(err)
		return
	}

	gojs.RunCode <- req
	control.Success()
}

// 删除代码仓库信息
func (RepoController) Build(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("id")
	if id == "" {
		control.FailBind(errors.New("[id]不能为空"))
		return
	}
	gojs.VersId <- id
	control.Success()
}

// 创建版本信息
func (RepoController) CreateVersion(ctx *gin.Context) {
	control := New(ctx)

	// 绑定参数
	req := new(model.CreateVersion)
	err := control.Bind(req)
	if err != nil {
		control.FailBind(err)
		return
	}

	err = logic.NewRepo().CreateVersion(req)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	control.Success()
}

// 获取版本列表
func (RepoController) GetVersionList(ctx *gin.Context) {
	control := New(ctx)

	// 绑定参数
	req := new(model.QueryVersion)
	err := control.Bind(req)
	if err != nil {
		control.FailBind(err)
		return
	}

	list, err := logic.NewRepo().GetVersionList(req)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	control.SuccessAutoData(req, list)
}

// 获取发布列表
func (RepoController) GetReleaseList(ctx *gin.Context) {
	control := New(ctx)

	// 绑定参数
	req := new(model.QueryRelease)
	err := control.Bind(req)
	if err != nil {
		control.FailBind(err)
		return
	}

	list, err := service.NewRepo(global.DB).GetReleaseList(req)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	control.SuccessAutoData(req, list)
}

// 移除程序
func (RepoController) RemoveRelease(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("id")
	if id == "" {
		control.FailBind(errors.New("[id]不能为空"))
		return
	}

	data, err := service.NewRepo(global.DB).GetReleaseById(id)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	err = service.NewRepo(global.DB).DelReleaseById(id)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	// 移除文件目录中的程序
	name := data.AppName
	if data.Platform == 0 {
		name += ".exe"
	}

	path := filepath.Join("runtime", "static", "app", data.AppName, name)
	err = os.Remove(path)
	if err != nil {
		control.FailByMsgf("移除文件失败，失败原因：%s", err.Error())
		return
	}

	control.Success()
}

func (RepoController) HandleDownloadFile(ctx *gin.Context) {
	control := New(ctx)
	id := control.Param("id")
	if id == "" {
		control.FailBind(errors.New("[id]不能为空"))
		return
	}

	ri, err := service.NewRepo(global.DB).GetReleaseById(id)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	// 更新下载次数
	err = service.NewRepo(global.DB).DownCountInc(id)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	file := filepath.Join("runtime", "static", "app", ri.AppName, ri.AppName)
	if ri.Platform == 0 {
		file += ".exe"
	}

	fmt.Println("file", file)

	control.File(file)
}
