package dao

import (
	"git.pmx.cn/hci/microservice-app/pkg/utils/gormutil"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDao() {
	Db = db()
}

func db() *gorm.DB {
	return gormutil.DB()
}