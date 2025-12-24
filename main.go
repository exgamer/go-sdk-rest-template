package main

import (
	"github.com/exgamer/go-sdk-rest-template/internal/app"
	"log"
)

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
