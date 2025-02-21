package repository

import (
	"database/sql"

	dbquery "github.com/tiago123456789/nocode-api-golang/pkg/dbQuery"
)

type TableRepositoryInterface interface {
	GetByName(name string) (map[string]interface{}, error)
	GetAll() ([]string, error)
	GetColumnsFromTable(table string) ([]string, error)
}

type TableRepository struct {
	db *sql.DB
}

func TableRepositoryNew(db *sql.DB) *TableRepository {
	return &TableRepository{
		db: db,
	}
}

func (t *TableRepository) GetColumnsFromTable(table string) ([]string, error) {
	sql := "SELECT column_name as column FROM information_schema.columns WHERE table_name = $1"
	rows, err := t.db.Query(sql, table)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var results []string
	for rows.Next() {
		var column string
		rows.Scan(&column)

		results = append(results, column)
	}

	return results, nil
}

func (t *TableRepository) GetAll() ([]string, error) {
	sql := "SELECT table_name as table FROM information_schema.tables WHERE table_schema = 'public'"
	rows, err := t.db.Query(sql)
	defer rows.Close()

	var results []string
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)

		results = append(results, tableName)
	}

	return results, err
}

func (t *TableRepository) GetByName(name string) (map[string]interface{}, error) {
	sql := "SELECT table_name as table FROM information_schema.tables WHERE table_name = $1 and table_schema = 'public'"
	rows, err := t.db.Query(sql, name)
	defer rows.Close()

	results, err := dbquery.GetResults(rows)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results[0], err
}
