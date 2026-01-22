package main

import (
	"github.com/exgamer/go-sdk-rest-template/internal/app"
	"log"
)

// @title        Rest-Template (API)
// @version      1.0
// @description  <b>Описание:</b> Пример приложения
// @host         0.0.0.0:8090
// @BasePath     /
func main() {
	appInstance, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = appInstance.RunAll()
	if err != nil {
		log.Fatal(err)
	}

	appInstance.WaitForShutdown()
}
