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
