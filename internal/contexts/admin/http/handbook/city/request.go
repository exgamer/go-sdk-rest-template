package city

import (
	"github.com/exgamer/gosdk-http-core/pkg/gin/validation"
)

// indexRequest Request для списка
type indexRequest struct {
	validation.Request
	Id   uint   `form:"id" description:"ID"`
	Name string `form:"id" description:"Название города"`
}

// createRequest Request для создания города
type createRequest struct {
	validation.Request
	Name string `json:"name" binding:"required" validate:"required"`
}

// updateRequest Request для редактирования города
type updateRequest struct {
	validation.Request
	Name string `json:"name" binding:"required" validate:"required"`
}
