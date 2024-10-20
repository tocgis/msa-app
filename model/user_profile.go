package model

import (
	"gorm.io/datatypes"
	"time"
)

// UserProfile 用户信息表
type UserProfile struct {
	User           User           `gorm:"ForeignKey:UserID" json:"-"`              // 用户基础表
	UserID         uint64         `gorm:"index:fk_user_profiles_users;column:user_id;type:bigint(20) unsigned;not null" json:"user_id"`
	Nickname       string         `gorm:"column:nickname;type:varchar(128)" json:"nickname"`             // 昵称
	NicknameStatus *int8          `gorm:"column:nickname_status;type:tinyint(4)" json:"nickname_status"` // 昵称状态
	Sex            *int8          `gorm:"column:sex;type:tinyint(4)" json:"sex"`                         // 性别
	Mobile         string         `gorm:"column:mobile;type:varchar(32)" json:"mobile"`                  // 手机号
	MobileStatus   *int8          `gorm:"column:mobile_status;type:tinyint(4)" json:"mobile_status"`     // 手机状态
	Realname       string         `gorm:"column:realname;type:varchar(45)" json:"realname"`              // 实名
	RealnameStatus *int8          `gorm:"column:realname_status;type:tinyint(4)" json:"realname_status"` // 实名状态
	IDcardno       string         `gorm:"column:idcardno;type:varchar(32)" json:"idcardno"`              // 身份证号码
	Province       *int           `gorm:"column:province;type:int(11)" json:"province"`                  // 省份id
	City           *int           `gorm:"column:city;type:int(11)" json:"city"`                          // 城市id
	Area           *int           `gorm:"column:area;type:int(11)" json:"area"`                          // 区域id
	Avatar         *int           `gorm:"column:avatar;type:int(11)" json:"avatar"`                      // 头像id
	Autograph      string         `gorm:"column:autograph;type:varchar(256)" json:"autograph"`           // 个人签名,
	Introduce      string         `gorm:"column:introduce;type:varchar(512)" json:"introduce"`           // 个人简介
	Tags           string         `gorm:"column:tags;type:varchar(256)" json:"tags"`                     // 个人标签
	Birthday       datatypes.Date `gorm:"column:birthday;type:date" json:"birthday"`                     // 生日
	Degree         string         `gorm:"column:degree;type:varchar(128)" json:"degree"`                 // 学位
	Major          string         `gorm:"column:major;type:varchar(128)" json:"major"`                   // 专业
	Company        string         `gorm:"column:company;type:varchar(128)" json:"company"`               // 公司名称
	CompanyID      *int64         `gorm:"column:company_id;type:bigint(20)" json:"company_id"`           // 公司id
	Post           string         `gorm:"column:post;type:varchar(45)" json:"post"`                      // 职位
	CreatedAt      time.Time      `gorm:"column:created_at;type:timestamp" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;type:timestamp" json:"updated_at"`


}

// TableName get sql table name.获取数据库表名
func (m *UserProfile) TableName() string {
	return "user_profiles"
}

func NewUserProfile() *UserProfile {
	return &UserProfile{}
}
