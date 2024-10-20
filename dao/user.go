package dao

import (
	"fmt"
	"git.pmx.cn/hci/microservice-app/model"
)

type UserDao struct {}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (m *UserDao) GetUserByPhone(phone string) (*model.User, error){
	 user := new(model.User)

	//user := model.Users{
	//	Name: phone,
	//}

	//result := map[string]interface{}{}
	Db.Model(&model.User{}).First(&user, "phone = ?", phone)

	//result := db().First(&user)
	//affected := result.RowsAffected // 返回找到的记录数
	//err := result.Error             // returns error or nil

	fmt.Println(user.ID)
	return user, nil
}

func (m *UserDao) GetUserInfoById(id int64) (*model.UserProfile, error) {
	userinfo := new(model.UserProfile)

	result := userinfo
	//res := Db.Model(model.UserProfile{}).Association("User").Find(&result, "user_profiles.user_id = ? ", id)
	res := Db.Where("user_id = ?", id).First(&userinfo)

	fmt.Println(res)
	return result, nil
}

func (m *UserDao) CreateUser(name string, phone string, password string) (*model.User, error) {

	user := new(model.User)

	result := user

	return result, nil
}