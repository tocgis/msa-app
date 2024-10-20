package aliyun

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/spf13/viper"
)

func AliSms(d map[string]string) (err error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", viper.GetString("aliyun.sms.accessKeyId"), viper.GetString("aliyun.sms.accessKeySecret"))

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = viper.GetString("aliyun.sms.signName")
	request.PhoneNumbers = d["phone"]

	if d["templateCodeType"] == "changePass" {
		request.TemplateCode = viper.GetString("aliyun.sms.template.changPassCode")
	} else if d["templateCodeType"] == "login" {
		request.TemplateCode = viper.GetString("aliyun.sms.template.loginCode")
	} else if d["templateCodeType"] == "authentication" {
		request.TemplateCode = viper.GetString("aliyun.sms.template.authenticationCode")
	} else if d["templateCodeType"] == "remoteLogin" {
		request.TemplateCode = viper.GetString("aliyun.sms.template.remoteLoginCode")
	} else if d["templateCodeType"] == "importantInfo" {
		request.TemplateCode = viper.GetString("aliyun.sms.template.updateInfoCode")
	} else if d["templateCodeType"] == "register" {
		request.TemplateCode = viper.GetString("aliyun.sms.template.registerCode")
	}

	request.TemplateParam = "{\"code\":"+d["code"]+"}"

	response, err := client.SendSms(request)
	if err != nil {
		return
	}
	if response.Code != "OK" {
		errData ,_:= json.Marshal(response)
		//log.Info("sms", "发送短信错误", string(errData))
		fmt.Println(errData)
		if response.Code == "isv.DAY_LIMIT_CONTROL" {
			return errors.New("每日发送次数超出上限")
		} else if response.Code == "isv.BUSINESS_LIMIT_CONTROL" {
			return errors.New("发送过于频繁，请稍后再试")
		} else {
			return errors.New("发送失败")
		}

	}
	return

}