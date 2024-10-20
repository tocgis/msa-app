package model

// Users 用户基础表
type User struct {
	//ID 			  string  `gorm:"unique;column:id;type:int(11)" json:"id"`
	UUID          string  `gorm:"unique;column:uuid;type:char(36)" json:"uuid"`
	Name          string  `gorm:"column:name;type:varchar(100)" json:"-"`
	Email         string  `gorm:"column:email;type:varchar(150)" json:"-"`
	Phone         string  `gorm:"unique;column:phone;type:varchar(50)" json:"phone"`
	Password      string  `gorm:"column:password;type:varchar(128)" json:"-"`
	PasswordSalt  string  `gorm:"column:password_salt;type:char(6)" json:"-"`
	RememberToken string  `gorm:"column:remember_token;type:varchar(256)" json:"-"`
	InviteCode    string  `gorm:"column:invite_code;type:varchar(45)" json:"invite_code"`       // 唯一邀请码
	InviteID      *uint64 `gorm:"column:invite_id;type:bigint(20) unsigned;default:0" json:"-"` // 邀请人id

	BaseModel
}

// TableName get sql table name.获取数据库表名
func (m *User) TableName() string {
	return "users"
}

func NewUser() *User {
	return &User{}
}
