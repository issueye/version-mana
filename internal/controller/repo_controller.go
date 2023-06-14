package controller

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/issueye/version-mana/internal/global"
	"github.com/issueye/version-mana/internal/logic"
	"github.com/issueye/version-mana/internal/model"
	"github.com/issueye/version-mana/internal/service"
)

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

	fmt.Println("req", req)

	list, err := logic.NewRepo().GetVersionList(req)
	if err != nil {
		control.FailByMsg(err.Error())
		return
	}

	control.SuccessAutoData(req, list)
}
