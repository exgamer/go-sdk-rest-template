// Отдельная точка входа для ops-команд проекта, в обход полного приложения
// (main.go в корне) — без HTTP/Rabbit kernels. Вся сборка — в
// internal/console, по аналогии с internal/app для основного сервиса.
package main

import (
	"log"

	"github.com/exgamer/go-sdk-rest-template/internal/console"
)

func main() {
	appInstance, err := console.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	defer appInstance.Close()

	if err := appInstance.RunAll(); err != nil {
		log.Fatal(err)
	}
}
