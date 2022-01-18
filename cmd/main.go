package main

import (
	"github.com/gorilla/mux"
	"log"
	"mockidoki/config"
	"mockidoki/internal/repository"
	"mockidoki/internal/service"
	"mockidoki/pkg/httpservice"
	"mockidoki/pkg/messagebus"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	configurationManager := config.NewConfigurationManager("config/config.yml", "local")
	postgresConfig := configurationManager.GetPostgresConfig()
	kafkaConfig := configurationManager.GetKafkaConfig()

	eventMockRepository := repository.NewEventMockRepository(postgresConfig)
	httpActionRepository := repository.NewHttpMockRepository(postgresConfig)
	kafkaProducer := messagebus.NewKafkaProducer(kafkaConfig)
	response := new(httpservice.Response)
	actionService := service.NewEventMockService(*eventMockRepository, *kafkaProducer, *response)
	httpActionService := service.NewHttpMockService(*httpActionRepository, *response)

	router.HandleFunc("/http-mocks", httpActionService.Process).Methods("POST")
	router.HandleFunc("/http-mocks", httpActionService.Process).Methods("GET")
	router.HandleFunc("/event-mocks/{key}/process", actionService.Process).Methods("POST")
	router.HandleFunc("/event-mocks/{key}/process-list", actionService.ProcessList).Methods("POST")
	router.HandleFunc("/mocks", actionService.Create).Methods("POST")
	router.HandleFunc("/management/health", func(writer http.ResponseWriter, request *http.Request) {
		payload := map[string]interface{}{"status": "ok"}
		response.RespondWithJSON(writer, http.StatusOK, payload)
	}).Methods("GET")

	port := configurationManager.GetServerConfig().Port
	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		log.Fatalf("Error when running the application : %s", err.Error())
	}
}
