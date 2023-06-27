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
		repo.GET("", control.List)                            // 列表
		repo.GET("/:id", control.GetById)                     // 通过id 查找
		repo.POST("", control.Create)                         // 创建
		repo.PUT("", control.Modify)                          // 修改
		repo.DELETE("/:id", control.Delete)                   // 删除
		repo.GET("branch/:id", control.BranchList)            // 分支
		repo.PUT("code", control.ModifyCode)                  // 修改代码
		repo.PUT("testRun", control.TestRun)                  // 测试运行
		repo.GET("ws/:id", controller.WsScriptTestRunConsole) // 测试运行输出监控
	}

	// 版本管理
	version := repo.Group("version")
	{
		version.DELETE("/:id", control.RemoveVersion)            // 移除版本
		version.POST("", control.CreateVersion)                  // 创建版本
		version.GET("", control.GetVersionList)                  // 获取版本列表
		version.GET("lastVerNum/:repoId", control.GetLastVerNum) // 获取最大版本号
		version.GET("build/:id", control.Build)                  // 编译
		version.GET("ws/:id", controller.WsScriptCompileConsole) // 编译输出监控
	}

	release := version.Group("release")
	{
		release.GET("", control.GetReleaseList)       // 获取发布列表
		release.DELETE("/:id", control.RemoveRelease) // 移除发布程序
	}
}
