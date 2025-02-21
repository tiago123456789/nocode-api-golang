package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type EndpointRepositoryInterface interface {
	Create(endpoint types.Endpoint) (int, error)
	GetByPath(path string) (int, error)
	GetAllCreated() (map[string]types.Endpoint, error)
	Setup() error
}

type EndpointRepository struct {
	db *sql.DB
}

func EndpointRepositoryNew(db *sql.DB) *EndpointRepository {
	return &EndpointRepository{
		db: db,
	}
}

func (e *EndpointRepository) Setup() error {
	_, err := e.db.Exec(`
		CREATE TABLE IF NOT EXISTS endpoints(
			id SERIAL PRIMARY KEY,
			path VARCHAR(255) NOT NULL,
			data JSONB
		);
		CREATE TABLE IF NOT EXISTS auth(
			id SERIAL PRIMARY KEY,
			name VARCHAR(70) NOT NULL,
			email VARCHAR(150) NOT NULL,
			password VARCHAR(255) NOT NULL
		);
		`,
	)

	if err != nil {
		return err
	}

	return nil
}

func (e *EndpointRepository) GetAllCreated() (map[string]types.Endpoint, error) {
	rows, err := e.db.Query("select path, data from endpoints")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var registers []types.EndpointFromDatabase
	for rows.Next() {
		var item types.EndpointFromDatabase
		if err := rows.Scan(&item.Path, &item.Data); err != nil {
			return nil, err
		}
		registers = append(registers, item)
	}

	endpoints := map[string]types.Endpoint{}
	for _, value := range registers {
		var item types.Endpoint
		json.Unmarshal([]byte(value.Data), &item)
		endpoints[value.Path] = item
	}

	return endpoints, nil
}

func (e *EndpointRepository) GetByPath(path string) (int, error) {
	id := 0
	err := e.db.QueryRow("SELECT id FROM endpoints where path = $1", path).Scan(&id)
	return id, err
}

func (e *EndpointRepository) Create(endpoint types.Endpoint) (int, error) {
	id := 0
	sql := fmt.Sprintf(
		`INSERT INTO endpoints(path, data) VALUES ($1, $2) RETURNING id;`,
	)
	data, _ := json.Marshal(endpoint)
	var params []interface{}
	params = append(params, endpoint.Path)
	params = append(params, data)

	err := e.db.QueryRow(sql, params...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil

}
