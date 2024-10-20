package socialite

import (
	"context"
	"git.pmx.cn/hci/microservice-app/proto/socialite"
	"time"
)

func init() {

}

func NewSocialiteService() socialite.SocialiteServiceServer {
	return service{}
}

type service struct {}

func (s service) WxJsConfig(ctx context.Context, param *socialite.NoParam) (*socialite.WxJsConfigResponse, error) {
	var res socialite.WxJsConfigResponse
	res.Appid = "wxb2a8cfe55d404daf"
	res.Signature = "fe71f68535d059ffeada7bd17bb28646e834aea5"
	res.Noncestr = "1598599463341104900"
	res.Timestamp = time.Now().Unix()

	return &res, nil
}

func (s service) WxJsLogin(ctx context.Context, request *socialite.WxJsLoginRequest) (*socialite.WxJsLoginResponse, error) {
	var res socialite.WxJsLoginResponse
	res.IsLogin = 1;
	res.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo1MSwibmFtZSI6IumCu-WutuWls-WtqSIsInBob25lIjoiMTg4MjM3MDYxNjIiLCJleHAiOjE1OTA0OTk3NTl9.ZSmTYD3xwlCE6axC-sG_Wk-qwP3Mp26JvOk2liR4L-4"
	
	return &res, nil
}

