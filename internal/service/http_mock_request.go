package service

type HttpMockRequest struct {
	MatchingHeader string `json:"matchingHeader"`
	MatchingUrl    string `json:"matchingUrl"`
	MatchingBody   string `json:"matchingBody"`
	ResponseBody   string `json:"responseBody"`
	ResponseHeader string `json:"responseHeader"`
	ResponseStatus int    `json:"responseStatus"`
	Description    string `json:"description"`
}

//func (request *HttpMockRequest) ToDao() repository.HttpMockDao {
//
//	responseHeader, err := json.Marshal(request.ResponseHeader)
//
//	if err != nil {
//		log.Fatalf("Error when converting response header into json: %s", err.Error())
//	}
//
//	return repository.HttpMockDao{
//		MatchingHeader: request.MatchingHeader,
//		MatchingUrl:    request.MatchingUrl,
//		ResponseBody:   request.ResponseBody,
//		ResponseHeader: responseHeader,
//		ResponseStatus: request.ResponseStatus,
//		Description:    request.Description,
//		IsDeleted:      false,
//	}
//}
