package dbquery

import (
	"database/sql"
	"log"
)

func GetResults(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	values := make([]interface{}, len(columns))
	valuePointers := make([]interface{}, len(columns))

	var results []map[string]interface{}
	for rows.Next() {
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			log.Fatal(err)
			return nil, err
		}

		row := make(map[string]interface{})
		for i, column := range columns {
			row[column] = values[i]
		}

		results = append(results, row)
	}

	return results, nil
}
