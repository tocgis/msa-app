package model

import "time"

// UserSocialite 用户三方信息绑定
type UserSocialite struct {
	AppID     string    `gorm:"index:app_id;column:app_id;type:varchar(32)" json:"app_id"`
	UserID    uint64    `gorm:"primaryKey;column:user_id;type:bigint(20) unsigned;not null" json:"-"` // User ID.
	User      User      `gorm:"joinForeignKey:user_id;foreignKey:id" json:"users_list"`               // 用户基础表
	UnionID   string    `gorm:"primaryKey;column:union_id;type:varchar(191);not null" json:"-"`       // Provider union ID.
	Type      string    `gorm:"column:type;type:varchar(100);not null" json:"type"`                   // Provider Type.
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *UserSocialite) TableName() string {
	return "user_socialites"
}

func NewUserSocialite() *UserSocialite {
	return &UserSocialite{}
}
