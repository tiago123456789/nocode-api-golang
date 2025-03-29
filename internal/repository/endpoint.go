package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tiago123456789/nocode-api-golang/internal/types"
	"go.opentelemetry.io/otel/trace"
)

type EndpointRepositoryInterface interface {
	Create(ctx context.Context, endpoint types.Endpoint) (int, error)
	GetByPath(ctx context.Context, path string) (int, error)
	GetAllCreated(ctx context.Context) (map[string]types.Endpoint, error)
	Setup() error
	Delete(ctx context.Context, id int64) (string, error)
}

type EndpointRepository struct {
	db     *sql.DB
	tracer trace.Tracer
}

func EndpointRepositoryNew(
	db *sql.DB, tracer trace.Tracer) *EndpointRepository {
	return &EndpointRepository{
		db:     db,
		tracer: tracer,
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

func (e *EndpointRepository) GetAllCreated(ctx context.Context) (map[string]types.Endpoint, error) {
	ctx, span := e.tracer.Start(ctx, "get-all-repository")
	defer span.End()

	rows, err := e.db.Query("select id, path, data from endpoints ORDER BY id ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var registers []types.EndpointFromDatabase
	for rows.Next() {
		var item types.EndpointFromDatabase
		if err := rows.Scan(&item.ID, &item.Path, &item.Data); err != nil {
			return nil, err
		}
		registers = append(registers, item)
	}

	endpoints := map[string]types.Endpoint{}
	for _, value := range registers {
		var item types.Endpoint
		json.Unmarshal([]byte(value.Data), &item)
		item.ID = value.ID
		endpoints[value.Path] = item
	}

	return endpoints, nil
}

func (e *EndpointRepository) GetByPath(ctx context.Context, path string) (int, error) {
	ctx, span := e.tracer.Start(ctx, "get-by-path-repository")
	defer span.End()

	id := 0
	err := e.db.QueryRow("SELECT id FROM endpoints where path = $1", path).Scan(&id)
	return id, err
}

func (e *EndpointRepository) Create(ctx context.Context, endpoint types.Endpoint) (int, error) {
	ctx, span := e.tracer.Start(ctx, "create-endpoint-repository")
	defer span.End()

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

func (e *EndpointRepository) Delete(ctx context.Context, id int64) (string, error) {
	ctx, span := e.tracer.Start(ctx, "delete-endpoint-repository")
	defer span.End()
	sql := `DELETE FROM endpoints WHERE id=$1 RETURNING path;`
	var path string
	err := e.db.QueryRow(sql, id).Scan(&path)
	if err != nil {
		return "", err
	}

	return path, nil
}
