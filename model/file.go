package model

import "time"

// Files 文件存储表
type Files struct {
	ID             int       `gorm:"primaryKey;column:id;type:int(11);not null" json:"-"`
	Hash           string    `gorm:"column:hash;type:varchar(150)" json:"hash"`                       // 文件hash
	OriginFilename string    `gorm:"column:origin_filename;type:varchar(191)" json:"origin_filename"` // 原始文件名
	Filename       string    `gorm:"column:filename;type:varchar(191)" json:"filename"`
	Mine           string    `gorm:"column:mine;type:varchar(100)" json:"mine"`                // mine type
	Width          *float64  `gorm:"column:width;type:double(8,2)" json:"width"`               // 图片宽度
	Height         *float64  `gorm:"column:height;type:double(8,2)" json:"height"`             // 图片高度
	Filesize       int64     `gorm:"column:filesize;type:bigint(20);not null" json:"filesize"` // 文件大小
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Files) TableName() string {
	return "files"
}
