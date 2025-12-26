package requests

import (
	"github.com/exgamer/gosdk-http-core/pkg/gin/validation"
)

// CityCreateRequest Request для создания города
type CityCreateRequest struct {
	validation.Request
	Name string `json:"name" binding:"required" validate:"required"`
}
