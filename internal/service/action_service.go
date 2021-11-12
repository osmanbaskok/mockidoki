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
	eventChannel, messageBytes, halt := getChannelAndMessageAsBytes(writer, request, service)

	if halt {
		return
	}

	eventMessage := string(messageBytes)

	err := service.kafkaProducer.Produce(eventMessage, *eventChannel)

	response := service.httpResponse
	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Event eventMessage could not be sent")
		return
	}

	writer.WriteHeader(http.StatusOK)
	return
}

func (service *ActionService) ProcessList(writer http.ResponseWriter, request *http.Request) {
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

	for _, message := range list {
		eventMessage, _ := json.Marshal(message)
		err := service.kafkaProducer.Produce(string(eventMessage), *eventChannel)

		if err != nil {
			response.RespondWithError(writer, http.StatusBadRequest, "Event message could not be sent")
			return
		}
	}

	writer.WriteHeader(http.StatusOK)
	return
}

func getChannelAndMessageAsBytes(writer http.ResponseWriter, request *http.Request, service *ActionService) (*string, []byte, bool) {
	response := service.httpResponse

	vars := mux.Vars(request)
	key, ok := vars["key"]
	if !ok {
		response.RespondWithError(writer, http.StatusNotFound, "Missing key parameter")
		return nil, nil, true
	}

	eventChannel, err := service.repo.FindEventChannelByKey(key)

	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Channel not found")
		return nil, nil, true
	}

	messageBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.RespondWithError(writer, http.StatusBadRequest, "Request could not be read")
		return nil, nil, true
	}
	return eventChannel, messageBytes, false
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
		response.RespondWithError(writer, http.StatusInternalServerError, "An error occurred when saving request. ["+err.Error()+"]")
		return
	}

	response.RespondWithJSON(writer, http.StatusCreated, nil)
}
