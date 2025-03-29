package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/tiago123456789/nocode-api-golang/internal/types"
	dbquery "github.com/tiago123456789/nocode-api-golang/pkg/dbQuery"
	"go.opentelemetry.io/otel/trace"
)

type CustomEndpointInterface interface {
	GetById(
		ctx context.Context,
		table string,
		id string,
	) ([]map[string]interface{}, error)
	GetAll(
		ctx context.Context,
		endpoint types.Endpoint,
	) ([]map[string]interface{}, error)
	GetAllByCustomQuery(
		ctx context.Context,
		endpoint types.Endpoint,
		params []interface{},
	) ([]map[string]interface{}, error)
	Delete(
		ctx context.Context,
		table string,
		id string,
	) error
	Create(
		ctx context.Context,
		newRegister map[string]interface{},
		table string,
	) (int64, error)
	Update(
		ctx context.Context,
		newRegister map[string]interface{},
		table string,
		idRegister string,
	) error
}

type CustomEndpointRepository struct {
	db     *sql.DB
	tracer trace.Tracer
}

func CustomEndpointRepositoryNew(db *sql.DB, tracer trace.Tracer,
) *CustomEndpointRepository {
	return &CustomEndpointRepository{
		db:     db,
		tracer: tracer,
	}
}

func (c *CustomEndpointRepository) Update(
	ctx context.Context,
	newRegister map[string]interface{},
	table string,
	idRegister string,
) error {
	ctx, span := c.tracer.Start(ctx, "update-register-repository")
	defer span.End()
	count := 1
	fieldToUpdate := make([]string, 0, len(newRegister))
	values := make([]interface{}, 0, len(newRegister))
	for key, value := range newRegister {
		fieldToUpdate = append(
			fieldToUpdate,
			fmt.Sprintf("%s = $%d", key, count),
		)
		values = append(values, value)
		count += 1
	}

	sql := fmt.Sprintf(
		`UPDATE "%s" SET %s WHERE id = $%d RETURNING id;`,
		table,
		strings.Join(fieldToUpdate, ","),
		count,
	)

	values = append(values, idRegister)

	var id int64
	err := c.db.QueryRow(sql, values...).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (c *CustomEndpointRepository) Create(
	ctx context.Context,
	newRegister map[string]interface{},
	table string,
) (int64, error) {
	ctx, span := c.tracer.Start(ctx, "create-register-repository")
	defer span.End()

	count := 1
	keys := make([]string, 0, len(newRegister))
	values := make([]interface{}, 0, len(newRegister))
	keysToReplace := make([]string, 0, len(newRegister))
	for key, value := range newRegister {
		keys = append(keys, key)
		values = append(values, value)
		keysToReplace = append(
			keysToReplace,
			fmt.Sprintf("$%d", count),
		)
		count += 1
	}

	sql := fmt.Sprintf(
		`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id;`,
		table,
		strings.Join(keys, ","),
		strings.Join(keysToReplace, ","),
	)

	var id int64
	err := c.db.QueryRow(sql, values...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *CustomEndpointRepository) GetAllByCustomQuery(
	ctx context.Context,
	endpoint types.Endpoint,
	params []interface{},
) ([]map[string]interface{}, error) {
	ctx, span := c.tracer.Start(ctx, "get-all-by-custom-query-register-repository")
	defer span.End()

	sql := endpoint.Query
	if len(endpoint.QueryParams) > 0 {
		if len(params) != len(endpoint.QueryParams) {
			return nil, errors.New(fmt.Sprintf(
				"You need to provide the params: %s",
				strings.Join(endpoint.QueryParams, ","),
			))
		}

		rows, err := c.db.Query(sql, params...)
		if err != nil {
			return nil, err
		}

		defer rows.Close()
		results, _ := dbquery.GetResults(rows)
		return results, nil
	}

	rows, err := c.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results, _ := dbquery.GetResults(rows)
	return results, nil
}

func (c *CustomEndpointRepository) Delete(
	ctx context.Context,
	table string,
	id string,
) error {
	ctx, span := c.tracer.Start(ctx, "delete-register-repository")
	defer span.End()
	sql := fmt.Sprintf(`DELETE FROM "%s" WHERE id=$1;`, table)
	_, err := c.db.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *CustomEndpointRepository) GetAll(
	ctx context.Context,
	endpoint types.Endpoint,
) ([]map[string]interface{}, error) {
	ctx, span := c.tracer.Start(ctx, "get-all-register-repository")
	defer span.End()

	sql := fmt.Sprintf("SELECT * FROM \"%s\"", endpoint.Table)
	rows, err := c.db.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	results, _ := dbquery.GetResults(rows)
	return results, nil
}

func (c *CustomEndpointRepository) GetById(
	ctx context.Context,
	table string,
	id string,
) ([]map[string]interface{}, error) {
	ctx, span := c.tracer.Start(ctx, "get-by-register-repository")
	defer span.End()

	sql := fmt.Sprintf(`SELECT * FROM "%s" WHERE id=$1;`, table)
	rows, err := c.db.Query(sql, id)
	if err != nil {
		return nil, err
	}

	results, err := dbquery.GetResults(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}
