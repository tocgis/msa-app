package user

import (
	"errors"
	"fmt"
	"time"

	"git.pmx.cn/hci/microservice-app/dao"
	"git.pmx.cn/hci/microservice-app/dto"
	"git.pmx.cn/hci/microservice-app/model"
	"git.pmx.cn/hci/microservice-app/pkg/utils"
	"git.pmx.cn/hci/microservice-app/pkg/utils/jwtutil"
)


type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) Login(req *dto.LoginDto) (interface{}, error) {
	var user *model.User
	ds := dao.NewUserDao()
	user = model.NewUser()


	verify := false

	// req.Type 默认为0，使用验证码登录，为1时，使用密码登录
	// 使用验证码登录时，如果验证码正确，无用户时创建新用户
	if req.Type == 0 {
		// 验证码正确
		if req.Code == "888888" {
			verify = true
		} else {
			return nil, errors.New("验证码错误")
		}
	} else {
		return nil, errors.New("暂时不支持密码登录")
	}

	if verify == true {
		u , _ := ds.GetUserByPhone(req.Name)
		if u.ID >1 {
			user = u
		} else {
			user.Phone = req.Name
			//user.Name = req.Name
			user.UUID = utils.CreateUUID()

			dao.Db.Create(&user)
		}
	}

	// 颁发 JWT
	//newSalt := random.GenValidateCode(6)
	//newToken := req.Name + req.Password + newSalt

	//password := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	newToken , err := jwtutil.CreateJwtToken(user.Phone, int(user.ID))
	if (err != nil ) {
		return nil, err
	}
	//fmt.Println(newToken)


	resp := map[string]interface{}{
		"token" : newToken,
		"exp" : time.Now().Unix() + 3600 * 24 * 14,
		"refresh" : time.Now().Unix() + 3600 * 24 *7,
		//"user" : user,
	}

	return resp, nil
}

// Auth 获取 登录认证 数据
func (s *service) Auth(req *dto.AuthDto) (interface{}, error ){

	token := req.Token
	jwtInfo, err := jwtutil.ParseToken(token)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(jwtInfo)
	exp , _ := utils.Interface2int64(jwtInfo["exp"])
	fmt.Println(exp)
	timeout := exp - time.Now().Unix()

	resp := map[string]interface{}{
		"jwt" : jwtInfo,
		"now": time.Now().Unix(),
		"sec" : timeout,
	}

	return resp, nil
}

func (s *service) UserInfo(req *dto.UserDto) (interface{}, error) {
	var userinfo *model.UserProfile
	ds := dao.NewUserDao()
	userinfo = model.NewUserProfile()

	if req.Id > 0 {
		userinfo, _ = ds.GetUserInfoById(req.Id)
	}

	resp := map[string]interface{}{
		"user": userinfo,
	}

	return resp, nil
}

func (s *service) Register(req *dto.RegisterDto) (interface{}, error) {

	resp := map[string]interface{}{
		"user": "",
	}

	return resp, nil
}