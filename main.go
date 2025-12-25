package main

import (
	"github.com/exgamer/go-sdk-rest-template/internal/app"
	"log"
)

// @title        Rest-Template (API)
// @version      1.0
// @description  <b>Описание:</b> Сервис доставки.<br/><b>Ответственная команда:</b> Checkout Team;<br/><b>Git:</b> <a href="https://gitlab.almanit.kz/jmart/shipping-service-go">https://gitlab.almanit.kz/jmart/shipping-service-go</a>;<br/><b>Golang SDK:</b> <a href="https://gitlab.almanit.kz/jmart/gosdk">https://gitlab.almanit.kz/jmart/gosdk</a>;<br/><b>Документация:</b> <a href="https://conf.almanit.kz/display/CHE/shipping-service-go">https://conf.almanit.kz/display/CHE/shipping-service-go</a>;<br/><b>Readme file:</b> <a href="https://gitlab.almanit.kz/jmart/shipping-service-go/-/blob/master/README.md">README.md</a>
// @host         0.0.0.0:8090
// @BasePath     /
func main() {
	appInstance, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = appInstance.RunHttpKernel()
	if err != nil {
		log.Fatal(err)
	}

	appInstance.WaitForShutdown()
}
