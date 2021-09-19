package service

import "mockidoki/internal/repository"

type ActionRequest struct {
	Key         string `json:"key"`
	Channel     string `json:"channel"`
	Description string `json:"description"`
}

func (request *ActionRequest) ToDao() repository.ActionDao {

	return repository.ActionDao{
		Key:         request.Key,
		Channel:     request.Channel,
		Description: request.Description,
		IsDeleted:   false,
	}
}
