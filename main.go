package main

import (
	"exgamer.kz/example-service/cmd/app/src"
	"exgamer.kz/example-service/cmd/app/src/config"
	_ "exgamer.kz/example-service/docs"
	"log"
)

//	@title			Example Service Go API
//	@version		1.0
//	@description	Example Service Go API

//	@host		localhost:8090
//	@BasePath	/
func main() {
	// Configuration
	appConfig, dbConfig, redisConfig, err := config.InitConfig()

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app := src.NewApp(appConfig, dbConfig, redisConfig)
	app.Run()
}
