package main

import (
	"github.com/latolukasz/orm"
	log "github.com/sirupsen/logrus"
	"os"
)

var engine *orm.Engine

func InitDatabase(mysqlUri string, redisUri string) (error, *orm.Engine) {

	if engine != nil {
		return nil, engine
	}

	registry := &orm.Registry{}
	registry.RegisterMySQLPool(mysqlUri)
	registry.RegisterRedis(redisUri, 0)
	registry.RegisterRedis(redisUri, 2, "lockers_pool")

	var image image
	var imageCategory imageCategory
	registry.RegisterEntity(&image, &imageCategory)

	registry.RegisterEnumSlice("status", []string{"visible", "hidden"})

	validatedRegistry, err := registry.Validate()
	if err != nil {
		return err, nil
	}
	engine = validatedRegistry.CreateEngine()

	imageCategories := validatedRegistry.GetTableSchemaForEntity(&imageCategory)
	imageCategories.UpdateSchema(engine)
	images := validatedRegistry.GetTableSchemaForEntity(&image)
	images.UpdateSchema(engine)

	return nil, engine
}

func SetupTestEngine() {
	if engine != nil {
		return
	}

	testMysql := os.Getenv("TEST_MYSQL_URI")
	testRedis := os.Getenv("TEST_REDIS_URI")

	if testRedis == "" || testMysql == "" {
		log.Panic("No mysql or redis uri provided, panicking...")
	}

	err, _ := InitDatabase(testMysql, testRedis)
	if err != nil {
		log.Panic("failed setting database up, panicking...", err)
	}
}