package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"mockidoki/config"
)

type ActionRepository struct {
	connection string
}

func (p *ActionRepository) FindEventChannelByKey(key string) *string {
	db, err := sql.Open("postgres", p.connection)
	defer db.Close()

	if err != nil {
		log.Fatalf("Error when connecting db : %s", err.Error())
		return nil
	}

	query := fmt.Sprintf("select channel from events where is_deleted = false and key = '%s'", key)

	data, err := db.Query(query)

	if err != nil {
		log.Fatalf("Error when running query : %s", err.Error())
		return nil
	}

	if data.Next() {

		var channel string
		err = data.Scan(&channel)
		if err != nil {
			log.Fatalf("Error when scanning data : %s", err.Error())
			return nil
		}

		return &channel
	}

	return nil
}

func NewActionRepository(config config.PostgresConfig) *ActionRepository {
	connection := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)
	return &ActionRepository{connection: connection}
}
