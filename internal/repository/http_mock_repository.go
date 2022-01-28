package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
	"mockidoki/config"
)

type HttpMockRepository struct {
	connection string
}

func (repo *HttpMockRepository) Find(method string, fullUrl string, header string, body string) (*HttpMockDao, error) {
	db, err := sql.Open("postgres", repo.connection)
	defer db.Close()

	if err != nil {
		log.Fatalf("Error when connecting db : %s", err.Error())
	}

	query := fmt.Sprintf("select id, method, matching_url, matching_body, matching_header, response_status, response_body, response_header from" +
		" (select id, method, coalesce(matching_url, '') matching_url, coalesce(matching_body, '') matching_body, coalesce(matching_header, '') matching_header, "+
		"response_status, coalesce(response_body, '') response_body, coalesce(response_header, '[]') response_header,"+
		" REGEXP_MATCHES('%s', coalesce(matching_url, '')),"+
		" REGEXP_MATCHES('%s', coalesce(matching_body, '')),"+
		" REGEXP_MATCHES('%s', coalesce(matching_header,''))"+
		" from http_mock where is_deleted = false and method = '%s') mock", fullUrl, body, header, method)

	data, err := db.Query(query)

	if err != nil {
		log.Fatalf("Error when running query : %s", err.Error())
	}

	if data.Next() {
		var httpMockDao HttpMockDao
		err = data.Scan(&httpMockDao.Id, &httpMockDao.Method, &httpMockDao.MatchingUrl, &httpMockDao.MatchingBody, &httpMockDao.MatchingHeader,
			&httpMockDao.ResponseStatus, &httpMockDao.ResponseBody, &httpMockDao.ResponseHeader)
		if err != nil {
			log.Fatalf("Error when scanning data : %s", err.Error())
		}

		return &httpMockDao, nil
	}

	return nil, errors.New("No record found!")
}

func NewHttpMockRepository(config config.PostgresConfig) *HttpMockRepository {
	connection := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)
	return &HttpMockRepository{connection: connection}
}
