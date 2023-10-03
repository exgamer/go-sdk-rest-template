package config

import (
	"github.com/davecgh/go-spew/spew"
	structures2 "github.com/exgamer/go-rest-sdk/pkg/config/structures"
	"github.com/exgamer/go-rest-sdk/pkg/helpers/config"
	"log"
)

func InitConfig() (*structures2.AppConfig, *structures2.DbConfig, *structures2.RedisConfig, error) {
	appConfig, dbConfig, redisConfig, err := config.InitBaseConfig()

	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	if appConfig.AppEnv != "prod" {
		spew.Dump(appConfig)
		spew.Dump(dbConfig)
	}

	return appConfig, dbConfig, redisConfig, nil
}
