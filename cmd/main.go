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

	actionRepository := repository.NewActionRepository(postgresConfig)
	kafkaProducer := messagebus.NewKafkaProducer(kafkaConfig)
	response := new(httpservice.Response)
	actionService := service.NewActionService(*actionRepository, *kafkaProducer, *response)

	router.HandleFunc("/actions/{key}/process", actionService.Process).Methods("POST")
	router.HandleFunc("/actions/{key}/process-list", actionService.ProcessList).Methods("POST")
	router.HandleFunc("/actions", actionService.Create).Methods("POST")
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
