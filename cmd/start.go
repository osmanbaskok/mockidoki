package cmd

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

func Start() {
	router := mux.NewRouter()

	configurationManager := config.NewConfigurationManager("config/config.yml", "local")
	postgresConfig := configurationManager.GetPostgresConfig()
	kafkaConfig := configurationManager.GetKafkaConfig()

	eventMockRepository := repository.NewEventMockRepository(postgresConfig)
	httpActionRepository := repository.NewHttpMockRepository(postgresConfig)
	kafkaProducer := messagebus.NewKafkaProducer(kafkaConfig)
	response := new(httpservice.Response)
	eventMockService := service.NewEventMockService(*eventMockRepository, *kafkaProducer, *response)
	httpMockService := service.NewHttpMockService(*httpActionRepository, *response)

	router.HandleFunc("/http-mocks", httpMockService.Process).Methods("GET", "POST", "PUT")
	router.HandleFunc("/http-mocks/{key}", httpMockService.Process).Methods("GET", "POST", "PUT")
	router.HandleFunc("/http-mocks/{key}/{key2}", httpMockService.Process).Methods("GET", "POST", "PUT")
	router.HandleFunc("/http-mocks/{key}/{key2}/{key3}", httpMockService.Process).Methods("GET", "POST", "PUT")
	router.HandleFunc("/http-mocks/{key}/{key2}/{key3}/{key4}", httpMockService.Process).Methods("GET", "POST", "PUT")
	router.HandleFunc("/http-mocks/{key}/{key2}/{key3}/{key4}/{key5}", httpMockService.Process).Methods("GET", "POST", "PUT")
	router.HandleFunc("/http-mocks/{key}/{key2}/{key3}/{key4}/{key5}/{key6}", httpMockService.Process).Methods("GET", "POST", "PUT")
	router.HandleFunc("/http-mocks/{key}/{key2}/{key3}/{key4}/{key5}/{key6}/{key7}", httpMockService.Process).Methods("GET", "POST", "PUT")
	router.HandleFunc("/http-mocks/{key}/{key2}/{key3}/{key4}/{key5}/{key6}/{key7}/{key8}", httpMockService.Process).Methods("GET", "POST", "PUT")
	router.HandleFunc("/http-mocks/{key}/{key2}/{key3}/{key4}/{key5}/{key6}/{key7}/{key8}/{key9}", httpMockService.Process).Methods("GET", "POST", "PUT")
	router.HandleFunc("/http-mocks/{key}/{key2}/{key3}/{key4}/{key5}/{key6}/{key7}/{key8}/{key9}/{key10}", httpMockService.Process).Methods("GET", "POST", "PUT")

	router.HandleFunc("/event-mocks/{key}/process", eventMockService.Process).Methods("POST")
	router.HandleFunc("/event-mocks/{key}/process-list", eventMockService.ProcessList).Methods("POST")
	router.HandleFunc("/event-mocks", eventMockService.Create).Methods("POST")
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
