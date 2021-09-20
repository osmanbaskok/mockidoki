package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
	"mockidoki/config"
)

type ActionRepository struct {
	connection string
}

func (repo *ActionRepository) Save(dao ActionDao) error {
	db, err := sql.Open("postgres", repo.connection)
	defer db.Close()

	if err != nil {
		log.Fatalf("Error when connecting db : %s", err.Error())
	}

	query := fmt.Sprintf("insert into action (key, channel, description, is_deleted) "+
		"values('%s','%s','%s','%t')", dao.Key, dao.Channel, dao.Description, dao.IsDeleted)

	_, err = db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (repo *ActionRepository) FindEventChannelByKey(key string) (*string, error) {
	db, err := sql.Open("postgres", repo.connection)
	defer db.Close()

	if err != nil {
		log.Fatalf("Error when connecting db : %s", err.Error())
	}

	query := fmt.Sprintf("select channel from action where is_deleted = false and key = '%s'", key)

	data, err := db.Query(query)

	if err != nil {
		log.Fatalf("Error when running query : %s", err.Error())
	}

	if data.Next() {

		var channel string
		err = data.Scan(&channel)
		if err != nil {
			log.Fatalf("Error when scanning data : %s", err.Error())
		}

		return &channel, nil
	}

	return nil, errors.New("No record found!")
}

func NewActionRepository(config config.PostgresConfig) *ActionRepository {
	connection := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)
	return &ActionRepository{connection: connection}
}
