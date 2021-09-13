package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"mockidoki/config"
	"mockidoki/internal/repository"
	"mockidoki/internal/service"
	message2 "mockidoki/pkg/message"
)

func main() {

	e := echo.New()
	configurationManager := config.NewConfigurationManager("config/config.yml", "local")
	postgresConfig := configurationManager.GetPostgresConfig()
	kafkaConfig := configurationManager.GetKafkaConfig()

	actionRepository := repository.NewActionRepository(postgresConfig)
	kafkaProducer := message2.NewKafkaProducer(kafkaConfig)
	actionService := service.NewActionService(*actionRepository, *kafkaProducer)

	e.POST("/actions", actionService.Process)
	e.GET("/management/health", func(c echo.Context) error {
		c.Response().WriteHeader(200)
		c.Response().Header().Set("Content-Type", "application/json")
		_, _ = c.Response().Write([]byte(`{"status": "ok"}`))
		return nil
	})

	port := configurationManager.GetServerConfig().Port
	err := e.Start(":" + port)
	if err != nil {
		log.Fatalf("Error when connecting db : %s", err.Error())
	}
}
