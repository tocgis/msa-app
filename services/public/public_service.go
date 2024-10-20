package public

import (
	"fmt"
	"git.pmx.cn/hci/microservice-app/dto"
	"git.pmx.cn/hci/microservice-app/pkg/redis"
	"git.pmx.cn/hci/microservice-app/pkg/utils"
	"github.com/spf13/viper"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type service struct{}

func NewService() *service {
	return &service{}
}

// SendSms 发送短信
func (s *service) SendSms(req *dto.SmsDto, c *gin.Context) (rs interface{}, err error) {
	var (
		smsType string
	)

	switch req.Type {
	case "login":
		smsType = "login"
	case "register":
		smsType = "register"
	case "change_pass":
		smsType = "changePass"
	case "update_info":
		smsType = "importantInfo"
	case "auth":
		smsType = "authentication"

	default:
		return nil, nil
	}

	fmt.Println("type", smsType)

	smsString := utils.CreateCaptcha() // 验证码6位
	fmt.Println(smsString)

	// 验证手机号码
	fmt.Println(req.Phone)

	// 使用 session 记录客户端
	session := sessions.Default(c)
	cid := session.Get("client_id")
	if cid == nil {
		cid = utils.CreateUUID()

		session.Set("client_id", cid)
		session.Save()
	}

	// 发送短信
	//发送短信验证码
	sendData := make(map[string]string)
	sendData["code"] = smsString
	sendData["phone"] = req.Phone
	sendData["templateCodeType"] = smsType

	//err = aliyun.AliSms(sendData)
	if err != nil {
		return
	}

	rdb := redis.GetRedis(
		viper.GetString("redis.host"),
		viper.GetString("redis.pass"),
		viper.GetInt("redis.systemDb"),
	) // 连接redis 服务

	ctx := c.Request.Context()

	err = rdb.HSet(ctx, req.Phone, "cid", cid, "sms", smsString, "time", time.Now().Unix()).Err()
	rdb.HMSet(ctx, req.Phone+"-1234", cid, smsString)
	//rdb.ZScan(ctx, "sorted-hash-key", 0, "prefix:*").Iterator()
	if err != nil {
		panic(err)
	}

	rs = map[string]interface{}{
		"exp":     time.Now().Unix() + 3600*24*14,
		"refresh": time.Now().Unix() + 3600*24*7,
		"type":    smsType,
		//"user" : user,
		"code": smsString,
	}
	return rs, nil
}

// Captcha 验证码
func (s *service) Captcha(d *dto.CaptchaDto) (rs interface{}, err error) {

	return rs, nil
}
