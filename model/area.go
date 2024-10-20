package model


// Area 地址表
type Area struct {
	ID           int    `gorm:"primaryKey;column:id;type:int(11);not null" json:"-"`
	AreaID       int    `gorm:"unique;column:area_id;type:int(11);not null" json:"area_id"`        // 地址id
	AreaName     string `gorm:"column:area_name;type:varchar(100);not null" json:"area_name"`      // 地址名称
	ParentAreaID int    `gorm:"column:parent_area_id;type:int(11);not null" json:"parent_area_id"` // 上级地址id
	Level        bool   `gorm:"column:level;type:tinyint(1);not null" json:"level"`                // 1省或直辖市 2市 3区县 4乡镇
	AreaCode     int    `gorm:"column:area_code;type:int(11);not null;default:0" json:"area_code"` // 行政编码
}

// TableName get sql table name.获取数据库表名
func (m *Area) TableName() string {
	return "area"
}

func NewArea() *Area {
	return &Area{}
}