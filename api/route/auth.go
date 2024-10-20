package route

import (
	"github.com/gin-gonic/gin"

	"git.pmx.cn/hci/microservice-app/api/handler/auth"
)

func RegisterAuth(router *gin.RouterGroup) {
	r := router.Group("/v1/auth")
	r.POST("/login", auth.Login)
	r.GET("/auth", auth.Auth)
	r.POST("/token", auth.Auth)
	r.POST("/singin", auth.Auth)
	r.POST("/sms", auth.SendSms)
	r.POST("/register", auth.Register)


}
