package service

import (
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"mockidoki/internal/message"
	"mockidoki/internal/repository"
	"net/http"
)

type ActionService struct {
	Repo          repository.ActionRepository
	KafkaProducer message.KafkaProducer
}

func NewActionService(repository repository.ActionRepository, kafkaProducer message.KafkaProducer) *ActionService {
	return &ActionService{Repo: repository, KafkaProducer: kafkaProducer}
}

func (service *ActionService) Process(c echo.Context) error {

	actionKey := c.Request().Header.Get("Action-Key")

	eventChannel := service.Repo.FindEventChannelByKey(actionKey)

	if eventChannel == nil {
		return c.JSON(http.StatusBadRequest, "Channel not found")
	}

	messageBytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Request body could not be read")
	}

	eventMessage := string(messageBytes)

	err = service.KafkaProducer.Produce(eventMessage, *eventChannel)

	if err != nil {
		return c.JSON(http.StatusBadRequest, "Event message could not be sent")
	}

	return c.NoContent(http.StatusOK)
}
