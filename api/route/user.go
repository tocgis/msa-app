package route

import (
	"git.pmx.cn/hci/microservice-app/api/handler/user"
	"github.com/gin-gonic/gin"
)

func RegisterUser(router *gin.RouterGroup) {
	r := router.Group("/user")
	r.GET("/profile", user.GetProfile)
	r.POST("/register", user.Register)
}
