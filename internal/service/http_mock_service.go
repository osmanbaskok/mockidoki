package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mockidoki/internal/repository"
	"mockidoki/pkg/httpservice"
	"net/http"
)

type HttpMockService struct {
	repo         repository.HttpMockRepository
	httpResponse httpservice.Response
}

func NewHttpMockService(repository repository.HttpMockRepository, httpResponse httpservice.Response) *HttpMockService {
	return &HttpMockService{repo: repository, httpResponse: httpResponse}
}

func (service *HttpMockService) Process(writer http.ResponseWriter, request *http.Request) {
	httpMock, halt := getHttpMock(writer, request, service)

	if halt {
		return
	}

	response := service.httpResponse
	var responseBody interface{}
	_ = json.Unmarshal([]byte(httpMock.ResponseBody), &responseBody)
	for _, header := range httpMock.ResponseHeader {
		writer.Header().Set(header.Header, header.Value)
	}
	response.RespondWithJSON(writer, httpMock.ResponseStatus, responseBody)
}

func getHttpMock(writer http.ResponseWriter, request *http.Request, service *HttpMockService) (*repository.HttpMockDao, bool) {
	response := service.httpResponse

	fullUrl := request.URL.RequestURI()
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		response.RespondWithError(writer, http.StatusInternalServerError, "Http mock body cannot be parsed")
		return nil, true
	}

	httpMock, err := service.repo.Find(request.Method, fullUrl, getAllRequestHeaderValues(request), string(body))

	if err != nil {
		response.RespondWithError(writer, http.StatusNotFound, "Http mock not found")
		return nil, true
	}

	return httpMock, false
}

func getAllRequestHeaderValues(request *http.Request) string {
	var strBuffer bytes.Buffer
	for _, values := range request.Header {
		for _, value := range values {
			strBuffer.WriteString(value + ",")
		}
	}
	return strBuffer.String()
}
