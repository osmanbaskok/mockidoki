package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"mockidoki/internal/repository"
	"mockidoki/pkg/httpservice"
	"mockidoki/pkg/messagebus"
	"net/http"
)

type ActionService struct {
	repo          repository.ActionRepository
	kafkaProducer messagebus.KafkaProducer
	httpResponse  httpservice.Response
}

func NewActionService(repository repository.ActionRepository, kafkaProducer messagebus.KafkaProducer, httpResponse httpservice.Response) *ActionService {
	return &ActionService{repo: repository, kafkaProducer: kafkaProducer, httpResponse: httpResponse}
}

func (service *ActionService) Process(writer http.ResponseWriter, request *http.Request) {
	response := service.httpResponse

	vars := mux.Vars(request)
	key, ok := vars["key"]
	if !ok {
		response.RespondWithError(writer, http.StatusNotFound, "Missing key parameter")
		return
	}

	eventChannel, err := service.repo.FindEventChannelByKey(key)

	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Channel not found")
		return
	}

	messageBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Request could not be read")
		return
	}

	eventMessage := string(messageBytes)

	err = service.kafkaProducer.Produce(eventMessage, *eventChannel)

	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Event messagebus could not be sent")
		return
	}

	writer.WriteHeader(http.StatusOK)
	return
}

func (service *ActionService) Create(writer http.ResponseWriter, request *http.Request) {

	response := service.httpResponse

	var actionRequest ActionRequest
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&actionRequest); err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Request could not be read")
		return
	}
	defer request.Body.Close()

	err := service.repo.Save(actionRequest.ToDao())
	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "An error occurred when saving request")
		return
	}

	response.RespondWithJSON(writer, http.StatusCreated, nil)
}
