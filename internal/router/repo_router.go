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

func (repo RepoRouter) Register(group *gin.RouterGroup) {
	f := group.Group(repo.Name)
	f.GET("", repo.control.List)
	f.POST("", repo.control.Create)
	f.PUT("", repo.control.Modify)
	f.DELETE("/:id", repo.control.Delete)
}
