package repository

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)

type HttpMockDao struct {
	Id             int                     `json:"id"`
	Method         string                  `json:"method"`
	MatchingHeader string                  `json:"matching_header"`
	MatchingUrl    string                  `json:"matching_url"`
	MatchingBody   string                  `json:"matching_body"`
	ResponseStatus int                     `json:"response_status"`
	ResponseBody   string                  `json:"response_body"`
	ResponseHeader ResponseHeaderItemArray `json:"response_header"`
	Description    string                  `json:"description"`
	IsDeleted      bool                    `json:"is_deleted"`
}

type ResponseHeaderItem struct {
	Header string `json:"header"`
	Value  string `json:"value"`
}

type ResponseHeaderItemArray []ResponseHeaderItem

func (a *ResponseHeaderItemArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	return json.Unmarshal(b, &a)
}

func (p ResponseHeaderItemArray) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}