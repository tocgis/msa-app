package route

import (
	"git.pmx.cn/hci/microservice-app/api/handler/assess"
	"github.com/gin-gonic/gin"
)

func RegisterAssess(router *gin.RouterGroup) {
	r := router.Group("/assess")
	r.POST("/init_score", assess.InitScore)
	r.GET("/score_info", assess.ScoreInfo)
	r.POST("/basic", assess.BasicSave)
	r.POST("/education", assess.EducationSave)
	r.POST("/work", assess.WorkSave)
}
