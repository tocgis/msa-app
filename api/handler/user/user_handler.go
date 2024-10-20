package user

import (
	client "git.pmx.cn/hci/microservice-app/client/user"
	"git.pmx.cn/hci/microservice-app/dto"
	"git.pmx.cn/hci/microservice-app/proto/user"
	userService "git.pmx.cn/hci/microservice-app/services/user"

	"github.com/gin-gonic/gin"

	"context"
	"net/http"
	"strconv"
)

func GetProfiles(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	req := &user.GetProfileRequest{UserId: userId}
	resp, err := client.GetClient().GetProfile(context.Background(), req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func GetProfile(c *gin.Context) {
	id , err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	req := &dto.UserDto{
		Id: id,
	}


	resp, err := userService.NewService().UserInfo(req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, resp)

}

// Register 注册用户
func Register(c *gin.Context) {

}