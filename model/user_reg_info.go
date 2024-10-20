package model

// UserRegInfo 用户注册记录
type UserRegInfo struct {
	UserID        uint64 `gorm:"index:fk_user_reg_logs_users1;column:user_id;type:bigint(20) unsigned;not null" json:"user_id"`
	User          User   `gorm:"joinForeignKey:user_id;foreignKey:id" json:"users_list"` // 用户基础表
	RegIP         string `gorm:"column:reg_ip;type:varchar(45)" json:"reg_ip"`           // 注册时ip
	RegCountry    string `gorm:"column:reg_country;type:varchar(45)" json:"reg_country"` // 注册时国家
	RegRegion     string `gorm:"column:reg_region;type:varchar(45)" json:"reg_region"`   // 注册时省份
	RegRegionName string `gorm:"column:reg_region_name;type:varchar(45)" json:"reg_region_name"`
	RegCity       string `gorm:"column:reg_city;type:varchar(45)" json:"reg_city"`     // 注册时城市
	RegLati       string `gorm:"column:reg_lati;type:varchar(45)" json:"reg_lati"`     // 纬度
	RegLong       string `gorm:"column:reg_long;type:varchar(45)" json:"reg_long"`     // 经度
	RegSource     string `gorm:"column:reg_source;type:varchar(45)" json:"reg_source"` // 注册来源
	LastIP        string `gorm:"column:last_ip;type:varchar(45)" json:"last_ip"`       // 最后登IP
	LastTime      *int   `gorm:"column:last_time;type:int(11)" json:"last_time"`       // 最后登录时间
}

// TableName get sql table name.获取数据库表名
func (m *UserRegInfo) TableName() string {
	return "user_reg_info"
}

func NewUserRegInfo() *UserRegInfo {
	return &UserRegInfo{}
}
