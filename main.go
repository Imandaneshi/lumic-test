package main

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/service/registry"
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/orm"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"lumic/graph/generated"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "lumic",
		Usage: "test project for coretrix",
		Before: func(ctx *cli.Context) error {
			err, _ := InitDatabase(config.mysqlUri, config.redisUri)
			if err != nil {
				log.Panic(err)
			}
			err = setupLogging(config.debug)
			if err != nil {
				log.Panic(err)
			}
			indexer := orm.NewRedisSearchIndexer(engine)
			go indexer.Run(context.Background())
			altersSearch := engine.GetRedisSearchIndexAlters()
			for _, alter := range altersSearch {
				alter.Execute()
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				EnvVars:     []string{"DEBUG"},
				Destination: &config.debug,
			},
			&cli.StringFlag{
				Name:        "mysql-uri",
				EnvVars:     []string{"MYSQL_URI"},
				Destination: &config.mysqlUri,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "test-mysql-uri",
				EnvVars:     []string{"TEST_MYSQL_URI"},
				Destination: &config.testMysqlUri,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "redis-uri",
				EnvVars:     []string{"REDIS_URI"},
				Destination: &config.redisUri,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "secret-key",
				EnvVars:     []string{"SECRET_KEY"},
				Destination: &config.secret,
				Required:    true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "api",
				Description: "Starts the web server",
				Action: func(context *cli.Context) error {

					s, deferFunc := hitrix.New(
						"lumic", config.secret,
					).RegisterDIService(
						registry.ServiceProviderErrorLogger(), //register redis error logger
						registry.ServiceDefinitionOrmRegistry(entity.Init),
						registry.ServiceDefinitionOrmEngine(),
						registry.ServiceDefinitionOrmEngineForContext(false),
						registry.ServiceProviderConfigDirectory("config"),
					).Build()
					defer deferFunc()
					s.RunServer(8068, generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}),
						func(ginEngine *gin.Engine) {

						//here you can register all your middlewares
					}, func(server *handler.Server) {

						})
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Please select one of the commands !")
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
