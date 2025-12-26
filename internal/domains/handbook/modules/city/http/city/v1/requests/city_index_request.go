package requests

import (
	"github.com/exgamer/gosdk-http-core/pkg/gin/validation"
)

// CityIndexRequest Request для списка
type CityIndexRequest struct {
	validation.Request
	Id   uint   `form:"id" description:"ID"`
	Name string `form:"id" description:"Название города"`
}
