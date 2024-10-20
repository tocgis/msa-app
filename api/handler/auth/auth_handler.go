package auth

import (
	"net/http"

	"git.pmx.cn/hci/microservice-app/dto"
	"git.pmx.cn/hci/microservice-app/services/public"
	"git.pmx.cn/hci/microservice-app/services/user"
	"github.com/gin-gonic/gin"
)

// Login 登录
func Login(c *gin.Context) {

	//reqDTO := &auth.LoginRequest{}
	req := &dto.LoginDto{}

	c.BindJSON(req)

	service := user.NewService()

	//service.Login(req)
	resp, err := service.Login(req)

	//resp, err := authClient.GetClient().Login(context.Background(), reqDTO) // 调用客户端
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, resp)
}

// Auth 用户认证
func Auth(c *gin.Context) {
	//reqDTO := &auth.AuthRequest{}
	//if err := c.BindJSON(reqDTO); err != nil {
	//	c.String(http.StatusBadRequest, err.Error())
	//	return
	//}
	req := &dto.AuthDto{}

	token := c.Query("token")

	req.Token = token

	//resp, err := authClient.GetClient().Auth(context.Background(), req)

	resp, err := user.NewService().Auth(req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func Register(c *gin.Context) {
	req := &dto.RegisterDto{}

	resp, err := user.NewService().Register(req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func SendSms(c *gin.Context) {
	req := &dto.SmsDto{}

	c.BindJSON(req)

	resp, err := public.NewService().SendSms(req, c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
