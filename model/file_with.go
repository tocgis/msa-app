package model

import "gorm.io/gorm"

// FileWiths 文件使用关联
type FileWiths struct {
	gorm.Model
	FileID  int    `gorm:"index:fk_file_withs_files1_idx;column:file_id;type:int(11);not null" json:"file_id"`
	Files   Files  `gorm:"joinForeignKey:file_id;foreignKey:id" json:"files_list"` // 文件存储表
	UserID  *int   `gorm:"column:user_id;type:int(11)" json:"user_id"`
	Channel string `gorm:"column:channel;type:varchar(100)" json:"channel"` // 记录频道 (task:image)
	Raw     string `gorm:"column:raw;type:varchar(100)" json:"raw"`         // 原始频道关联信息 (如关联id)
	Size    string `gorm:"column:size;type:varchar(50)" json:"size"`        // 图片尺寸，目标文件如果是图片的话则存在。便于客户端提前预设盒子
}

// TableName get sql table name.获取数据库表名
func (m *FileWiths) TableName() string {
	return "file_withs"
}
