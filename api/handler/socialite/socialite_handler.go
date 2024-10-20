package socialite

import (
	"context"
	client "git.pmx.cn/hci/microservice-app/client/socialite"
	"git.pmx.cn/hci/microservice-app/proto/socialite"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JsLogin(c *gin.Context) {
	reqDTO := &socialite.WxJsLoginRequest{}
	c.BindJSON(reqDTO)

	resp, err := client.GetClient().WxJsLogin(context.Background(), reqDTO)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, resp)


}

func JsConfig(c *gin.Context) {
	reqDTO := &socialite.NoParam{}
	c.BindJSON(reqDTO)

	resp, err := client.GetClient().WxJsConfig(context.Background(), reqDTO)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, resp)
}