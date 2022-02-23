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

	query := fmt.Sprintf("select id, method, coalesce(matching_url, '') as matching_url, coalesce(matching_body, '') as matching_body, "+
		"coalesce(matching_header, '') as matching_header, response_status, coalesce(response_body, '') as response_body, "+
		"coalesce(response_header, '[]') as response_header from http_mock "+
		"where (matching_url is null or '%s' ~ matching_url) "+
		"and (matching_body is null or '%s' ~ matching_body) "+
		"and (matching_header is null or '%s' ~ matching_header) "+
		"and is_deleted = false and method = '%s'", fullUrl, body, header, method)

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
