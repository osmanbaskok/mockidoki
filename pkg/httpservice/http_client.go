package httpservice

import (
	"bytes"
	"net/http"
	"time"
)

type Client interface {
	Post(url string, requestSerialized []byte) (*http.Response, error)
}

type service struct {
	client  *http.Client
	request request
}

type request struct {
}

func Create(timeout int) Client {
	client := http.Client{
		Timeout: time.Duration(timeout * int(time.Second)),
	}
	request := new(request)

	return &service{client: &client, request: *request}
}

func (srv *service) Post(url string, requestSerialized []byte) (*http.Response, error) {
	request, err := srv.request.NewRequest(url, requestSerialized)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")

	resp, err := srv.client.Do(request)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}

func (request *request) NewRequest(url string, requestSerialized []byte) (*http.Request, error) {
	return http.NewRequest("POST", url, bytes.NewBuffer(requestSerialized))
}
