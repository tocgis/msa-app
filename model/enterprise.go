package model

import "gorm.io/gorm"

// Enterprises 企业
type Enterprises struct {
	gorm.Model
	UUID          string `gorm:"column:uuid;type:char(36)" json:"uuid"`
	EnName        string `gorm:"column:en_name;type:varchar(256)" json:"en_name"`                // 企业名称
	EnCreditCode  string `gorm:"column:en_credit_code;type:varchar(64)" json:"en_credit_code"`   // 信息识别码
	EnLegalPerson string `gorm:"column:en_legal_person;type:varchar(45)" json:"en_legal_person"` // 法人
	EnMobile      string `gorm:"column:en_mobile;type:varchar(45)" json:"en_mobile"`             // 企业电话
	EnEmail       string `gorm:"column:en_email;type:varchar(128)" json:"en_email"`              // 企业email
	EnAddress     string `gorm:"column:en_address;type:varchar(256)" json:"en_address"`          // 企业所在地址
	EnWebsite     string `gorm:"column:en_website;type:varchar(256)" json:"en_website"`          // 企业官网
	RefUserID     *int64 `gorm:"column:ref_user_id;type:bigint(20)" json:"ref_user_id"`          // 企业注册人id
}

// TableName get sql table name.获取数据库表名
func (m *Enterprises) TableName() string {
	return "enterprises"
}
