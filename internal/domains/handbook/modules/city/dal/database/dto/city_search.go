package dto

// CitySearch DTO для поиска по городам
type CitySearch struct {
	Id      uint   `description:"ID"`
	Name    string `description:"name"`
	Page    uint   `form:"page"`
	PerPage uint   `form:"per_page"`
}
