package route

import (
	"git.pmx.cn/hci/microservice-app/api/handler/socialite"

	"github.com/gin-gonic/gin"
)

func RegisterSocialite(router *gin.RouterGroup) {
	r := router.Group("/socialite/wx")
	r.POST("/js_login", socialite.JsLogin)
	r.POST("/js_config", socialite.JsConfig)
}
