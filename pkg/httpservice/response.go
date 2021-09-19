package httpservice

import (
	"encoding/json"
	"net/http"
)

type Response struct {
}

func (response *Response) RespondWithError(w http.ResponseWriter, code int, message string) {
	response.RespondWithJSON(w, code, map[string]string{"error": message})
}

func (response *Response) RespondWithJSON(writer http.ResponseWriter, code int, payload interface{}) {

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	if payload != nil {
		jsonPayload, _ := json.Marshal(payload)
		_, _ = writer.Write(jsonPayload)
	}
}
