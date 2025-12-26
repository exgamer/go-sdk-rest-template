package requests

import (
	"github.com/exgamer/gosdk-http-core/pkg/gin/validation"
)

// CityUpdateRequest Request для редактирования города
type CityUpdateRequest struct {
	validation.Request
	Name string `json:"name" binding:"required" validate:"required"`
}
