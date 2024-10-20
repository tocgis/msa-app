package assess

import (
	"context"
	"net/http"

	client "git.pmx.cn/hci/microservice-app/client/assess"
	"git.pmx.cn/hci/microservice-app/proto/assess"

	"github.com/gin-gonic/gin"
)

//func RegisterAssess(router *gin.RouterGroup) {
//	r := router.Group("/assess")
//	r.POST("/init_score", initScore)
//	r.GET("/score_info", scoreInfo)
//	r.POST("/basic", basicSave)
//	r.POST("/education", educationSave)
//	r.POST("/work", workSave)
//}

func InitScore(c *gin.Context) {
	reqDTO := &assess.ScoreRequest{}
	c.BindJSON(reqDTO)

	resp, err := client.GetClient().InitScore(context.Background(), reqDTO)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, resp)

}

func ScoreInfo(c *gin.Context) {
	reqDTO := &assess.ScoreRequest{}

	phone := c.Query("phone")
	userId := c.Query("user_id")

	reqDTO.Phone = phone
	reqDTO.UserId = userId

	resp, err := client.GetClient().ScoreInfo(context.Background(), reqDTO)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)



}

func BasicSave(c *gin.Context) {
	reqDTO := &assess.BasicInfoRequest{}
	c.BindJSON(reqDTO)

	resp, err := client.GetClient().BasicSave(context.Background(), reqDTO)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)

}

func EducationSave(c *gin.Context) {
	reqDTO := &assess.EducationRequest{}
	c.BindJSON(reqDTO)

	resp, err := client.GetClient().EducationSave(context.Background(), reqDTO)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)

}

func WorkSave(c *gin.Context) {
	reqDTO := &assess.WorkinfoRequest{}
	c.BindJSON(reqDTO)

	resp, err := client.GetClient().WorkSave(context.Background(), reqDTO)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)

}