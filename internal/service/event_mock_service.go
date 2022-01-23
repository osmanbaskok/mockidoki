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

type EventMockService struct {
	repo          repository.EventMockRepository
	kafkaProducer messagebus.KafkaProducer
	httpResponse  httpservice.Response
}

func NewEventMockService(repository repository.EventMockRepository, kafkaProducer messagebus.KafkaProducer, httpResponse httpservice.Response) *EventMockService {
	return &EventMockService{repo: repository, kafkaProducer: kafkaProducer, httpResponse: httpResponse}
}

func (service *EventMockService) Process(writer http.ResponseWriter, request *http.Request) {
	eventChannel, messageBytes, halt := getChannelAndMessageAsBytes(writer, request, service)

	if halt {
		return
	}

	err := service.kafkaProducer.Produce(messageBytes, *eventChannel)

	response := service.httpResponse
	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Event eventMessage could not be sent")
		return
	}

	writer.WriteHeader(http.StatusOK)
	return
}

func (service *EventMockService) ProcessList(writer http.ResponseWriter, request *http.Request) {
	eventChannel, messageBytes, halt := getChannelAndMessageAsBytes(writer, request, service)

	if halt {
		return
	}

	response := service.httpResponse
	var list []interface{}
	if err := json.Unmarshal(messageBytes, &list); err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Request could not be parsed")
		return
	}

	err := service.kafkaProducer.ProduceList(list, *eventChannel)

	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Event message could not be sent")
		return
	}

	writer.WriteHeader(http.StatusOK)
	return
}

func getChannelAndMessageAsBytes(writer http.ResponseWriter, request *http.Request, service *EventMockService) (*string, []byte, bool) {
	response := service.httpResponse

	vars := mux.Vars(request)
	key, ok := vars["key"]
	if !ok {
		response.RespondWithError(writer, http.StatusNotFound, "Missing key parameter")
		return nil, nil, true
	}

	eventChannel, err := service.repo.FindEventChannelByKey(key)

	if err != nil {
		response.RespondWithError(writer, http.StatusNotFound, "Event mock not found")
		return nil, nil, true
	}

	messageBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Request could not be read")
		return nil, nil, true
	}
	return eventChannel, messageBytes, false
}

func (service *EventMockService) Create(writer http.ResponseWriter, request *http.Request) {

	response := service.httpResponse

	var eventMockRequest EventMockRequest
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&eventMockRequest); err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Request could not be read")
		return
	}
	defer request.Body.Close()

	err := service.repo.Save(eventMockRequest.ToDao())
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, "An error occurred when saving request. ["+err.Error()+"]")
		return
	}

	response.RespondWithJSON(writer, http.StatusCreated, nil)
}
