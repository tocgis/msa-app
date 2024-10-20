package model

import (
	"database/sql/driver"
	"fmt"
	"git.pmx.cn/hci/microservice-app/pkg/utils/gormutil"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

type BaseModel struct {
	CreatedAt DateTime       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt DateTime       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ID        uint           `gorm:"primarykey;autoIncrement" json:"id"`
}

func DB() *gorm.DB {
	db = gormutil.DB()
	return db
}

type DateTime struct {
	time.Time
}

// Scan implements the Scanner interface.
func (t *DateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// MarshalJSON 实现它的json序列化方法
func (t DateTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (t DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

type Date struct {
	time.Time
}

// Scan implements the Scanner interface.
func (t *Date) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Date{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// MarshalJSON 实现它的json序列化方法
func (t Date) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))
	return []byte(stamp), nil
}

func (t Date) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}
