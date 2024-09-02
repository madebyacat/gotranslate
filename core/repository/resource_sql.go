package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gotranslate/core/contracts"
	"gotranslate/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// This implementation was mostly to experiment
type ResourceSql struct {
	Pool *pgxpool.Pool
}

func NewResourceSql(pool *pgxpool.Pool) *ResourceSql {
	return &ResourceSql{Pool: pool}
}

var _ contracts.ResoureRepository = (*ResourceSql)(nil)

func (repo *ResourceSql) Init() error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS resources (
			Key TEXT NOT NULL,
			LanguageCode TEXT NOT NULL,
			Text TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`

	_, err := repo.Pool.Exec(context.Background(), createTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ResourceSql) GetResourcesByLanguageCode(languageCode string) ([]models.Resource, error) {
	return repo.getResources(resourceFilter{LanguageCode: languageCode})
}

func (repo *ResourceSql) GetResourcesByKey(key string) ([]models.Resource, error) {
	return repo.getResources(resourceFilter{Key: key})
}

func (repo *ResourceSql) AddResources(resources ...models.Resource) error {
	sqlStatement := `INSERT INTO resources (Key, LanguageCode, Text) VALUES `
	columns := 3
	params := []interface{}{}
	totalResources := len(resources)

	for i, resource := range resources {
		sqlStatement += fmt.Sprintf("($%d, $%d, $%d)", i*columns+1, i*columns+2, i*columns+3)
		if i+1 < totalResources {
			sqlStatement += ", "
		} else {
			sqlStatement += ";"
		}
		param := []interface{}{resource.Key, resource.LanguageCode, resource.Text}
		params = append(params, param...)
	}

	_, err := repo.Pool.Exec(context.Background(), sqlStatement, params...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ResourceSql) UpdateResourceValues(resources ...models.Resource) (rowsAffected int64, err error) {
	type sqlWithParams struct {
		sqlStatement string
		params       []interface{}
	}

	generatedUpdates := []sqlWithParams{}
	for _, resource := range resources {
		update := sqlWithParams{
			"UPDATE resources SET Text = $1 WHERE Key = $2 AND LanguageCode = $3;",
			[]interface{}{resource.Text, resource.Key, resource.LanguageCode},
		}
		generatedUpdates = append(generatedUpdates, update)
	}

	for i, update := range generatedUpdates {
		cmd, err := repo.Pool.Exec(context.Background(), update.sqlStatement, update.params...)
		if err != nil {
			shouldUpdateCount := len(resources)
			if i > 0 {
				return rowsAffected, fmt.Errorf("partially updated %d/%d rows. Error: %v", i, shouldUpdateCount, err.Error())
			} else {
				return rowsAffected, fmt.Errorf("no entries were updated. Error: %v", err.Error())
			}
		}
		rowsAffected += cmd.RowsAffected()
	}

	return rowsAffected, nil
}

func (repo *ResourceSql) RemoveResources(key, languageCode string) (rowsAffected int64, err error) {
	sqlStatement := "DELETE FROM resources WHERE Key = $1 AND LanguageCode = $2;"
	cmd, err := repo.Pool.Exec(context.Background(), sqlStatement, key, languageCode)
	if err != nil {
		return 0, err
	}
	return cmd.RowsAffected(), nil
}

func (repo *ResourceSql) ExistingLanguageCodes() (results []models.LanguageResult, err error) {
	query := `SELECT languagecode as "LanguageCode", COUNT(*) as "Count" FROM resources GROUP BY languagecode`
	rows, err := repo.Pool.Query(context.Background(), query)
	if err != nil {
		return []models.LanguageResult{}, err
	}

	for rows.Next() {
		item := models.LanguageResult{}
		rows.Scan(&item.LanguageCode, &item.Count)
		results = append(results, item)
	}

	return results, nil
}

func (repo *ResourceSql) getResources(filters ...resourceFilter) ([]models.Resource, error) {
	results, emptyResult := []models.Resource{}, []models.Resource{}

	query, params, err := generateQueryAndParameters(filters...)
	if err != nil {
		return emptyResult, err
	}

	rows, err := repo.Pool.Query(context.Background(), query, params...)
	if err != nil {
		return emptyResult, err
	}
	defer rows.Close()

	for rows.Next() {
		var resource models.Resource
		err := rows.Scan(&resource.Key, &resource.LanguageCode, &resource.Text)
		if err != nil {
			return nil, err
		}
		results = append(results, resource)
	}

	return results, nil
}

func generateQueryAndParameters(filters ...resourceFilter) (query string, params []any, err error) {
	filtersCount := len(filters)

	if filtersCount == 0 {
		return "", nil, errors.New("no filters defined")
	}

	if filtersCount == 1 {
		query = `
			SELECT
				Key,
				LanguageCode,
				Text
			FROM resources
			WHERE 1=1
		`
		filter := filters[0]

		if len(filter.LanguageCode) > 0 {
			params = append(params, filter.LanguageCode)
			query += fmt.Sprintf(" AND LanguageCode = $%d", len(params))
		}

		if len(filter.Key) > 0 {
			params = append(params, filter.Key)
			query += fmt.Sprintf(" AND Key = $%d", len(params))
		}
	} else if filtersCount > 1 {
		query = `
			WITH filter_data AS (
				SELECT * FROM jsonb_to_recordset($1::jsonb) AS x(k TEXT, l TEXT)
			)
			SELECT r.*
			FROM resources r
			INNER JOIN filter_data f ON r.Key = f.Key AND r.LanguageCode = f.LanguageCode
		`
		filterData, err := json.Marshal(filters)
		if err != nil {
			return "", nil, errors.New(err.Error())
		}
		params = append(params, filterData)
	}

	return query, params, nil
}
