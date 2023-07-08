package router

import (
	"fmt"
	"io/fs"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/issueye/version-mana/internal/config"
	"github.com/issueye/version-mana/internal/global"
	"github.com/issueye/version-mana/internal/middleware"
)

type IRouters interface {
	Register(group *gin.RouterGroup)
}

func InitRouter(r *gin.Engine) {
	name := config.GetParam("api-name", "api").String()
	v := config.GetParam("api-version", "v1").String()

	apiName := r.Group(name)

	if strings.ToUpper(global.TagName) == "UI" {
		r.Use(func(ctx *gin.Context) {
			path := ctx.Request.URL.Path
			fmt.Println("Handle ->", path)
			reg := "/admin/api/v1"
			matched, err := regexp.MatchString(reg, ctx.Request.URL.Path)
			if err != nil {
				ctx.Next()
				return
			}

			if !matched {
				ctx.Next()
				return
			}

			re := regexp.MustCompile(reg)
			ctx.Request.URL.Path = re.ReplaceAllString(ctx.Request.URL.Path, "/api/v1")
			r.HandleContext(ctx)
		})

		f, _ := fs.Sub(global.Assets, "assets/admin")
		r.StaticFS("/admin", http.FS(f))
	}

	version := apiName.Group(v)
	global.Auth = middleware.NewAuth()

	// 用户鉴权
	version.POST("login", global.Auth.LoginHandler)
	version.GET("logout", global.Auth.LogoutHandler)
	version.GET("refreshToken", global.Auth.RefreshHandler)

	// 鉴权
	// version.Use(global.Auth.MiddlewareFunc())
	registerVersionRouter(version,
		NewParamsRouter(), // 参数配置
		NewJobRouter(),    // 定时任务
		NewRepoRouter(),   // 代码仓库
	)
}

// registerRouter 注册路由
func registerVersionRouter(r *gin.RouterGroup, iRouters ...IRouters) {
	for _, iRouter := range iRouters {
		iRouter.Register(r)
	}
}
