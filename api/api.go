package api

import (
	"git.pmx.cn/hci/microservice-app/api/route"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	r := router.Group("/api")
	route.RegisterAuth(r) // 登录授权
	route.RegisterUser(r) // 用户
	route.RegisterSocialite(r) // 微信等社交平台对接
	route.RegisterCaptcha(r) // 验证码

	route.RegisterPublic(r) // 公共模块
	route.RegisterAssess(r) // 人力资本身价评估

}
