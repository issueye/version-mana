package router

import (
	"github.com/gin-gonic/gin"
	"github.com/issueye/version-mana/internal/controller"
)

type RepoRouter struct {
	Name    string
	control *controller.RepoController
}

func NewRepoRouter() *RepoRouter {
	return &RepoRouter{
		Name:    "repo",
		control: controller.NewRepoController(),
	}
}

func (r RepoRouter) Register(group *gin.RouterGroup) {
	control := controller.NewRepoController()

	// 仓库管理
	repo := group.Group(r.Name)
	{
		repo.GET("", control.List)
		repo.POST("", control.Create)
		repo.PUT("", control.Modify)
		repo.DELETE("/:id", control.Delete)
		repo.GET("branch/:id", control.BranchList)
		repo.PUT("code", control.ModifyCode)
		repo.PUT("testRun", control.TestRun)
	}

	// 版本管理
	version := repo.Group("version")
	{
		version.DELETE("/:id", control.RemoveVersion)            // 移除版本
		version.POST("", control.CreateVersion)                  // 创建版本
		version.GET("", control.GetVersionList)                  // 获取版本列表
		version.GET("lastVerNum/:repoId", control.GetLastVerNum) // 获取最大版本号
		version.GET("build/:id", control.Build)                  // 编译
	}
}
