package models

// City — таблица городов
type City struct {
	ID   uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"column:name"`
}

func (City) TableName() string {
	return "city"
}
