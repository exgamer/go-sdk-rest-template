package city

// City — таблица городов
type city struct {
	ID     uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name   string `gorm:"column:name"`
	Status int    `gorm:"column:status"`
}

func (city) TableName() string {
	return "city"
}
