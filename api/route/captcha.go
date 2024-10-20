package route

import (
	"git.pmx.cn/hci/microservice-app/api/handler/captcha"
	"github.com/gin-gonic/gin"
)

func RegisterCaptcha(router *gin.RouterGroup) {

	g := router.Group("/captcha") // 验证码
	g.GET("/block_puzzle", captcha.Dev) // 滑块验证码
	g.POST("/puzzle_check", captcha.Dev) // 校验验证码是否滑动成功
	g.POST("/system_verify", captcha.Dev) // 业务层调用校验封装
	g.GET("/image", captcha.Dev) // 图形数字

	g1 := router.Group("/captcha_sms") // 短信验证码
	g1.GET("/config", captcha.Dev) // 获取短信验证码配置状态，如果需要滑块，则由前端调用滑块功能
	g1.POST("/code", captcha.Dev) // 发送短信验证码
	g1.POST("/verify", captcha.Dev) // 校验短信验证码

}
