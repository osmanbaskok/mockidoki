package service

import "mockidoki/internal/repository"

type EventMockRequest struct {
	Key         string `json:"key"`
	Channel     string `json:"channel"`
	Description string `json:"description"`
}

func (request *EventMockRequest) ToDao() repository.EventMockDao {

	return repository.EventMockDao{
		Key:         request.Key,
		Channel:     request.Channel,
		Description: request.Description,
		IsDeleted:   false,
	}
}
