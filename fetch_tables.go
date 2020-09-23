package xlsx4db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func dbType(db *sql.DB) string {
	return reflect.ValueOf(db.Driver()).Type().String()
}

func isMySQL(db *sql.DB) bool {
	return dbType(db) == "*mysql.MySQLDriver"
}

func isPostgreSQL(db *sql.DB) bool {
	t := dbType(db)
	return t == "*pq.drv" || t == "*pq.Driver"
}

func fetchTableRows(rows *sql.Rows) ([]string, error) {
	var tables []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	err := rows.Err()
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func fetchTablesMySQL(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return fetchTableRows(rows)
}

func fetchTablesPostgreSQL(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT relname FROM pg_stat_user_tables")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return fetchTableRows(rows)
}

// FetchTables fetches all accessible tables from database.
func FetchTables(db *sql.DB) ([]string, error) {
	if isMySQL(db) {
		return fetchTablesMySQL(db)
	}
	if isPostgreSQL(db) {
		return fetchTablesPostgreSQL(db)
	}
	return nil, fmt.Errorf("not supported DB: %#v", dbType(db))
}

func buildInsertQueryMySQL(db *sql.DB, table string, columns []string) (string, error) {
	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = "?"
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table,
		strings.Join(columns, ", "), strings.Join(placeholders, ", ")), nil
}

func buildInsertQueryPostgreSQL(db *sql.DB, table string, columns []string) (string, error) {
	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table,
		strings.Join(columns, ", "), strings.Join(placeholders, ", ")), nil
}

// BuildInsertQuery builds na insert statement with given columns.
func BuildInsertQuery(db *sql.DB, table string, columns []string) (string, error) {
	if isMySQL(db) {
		return buildInsertQueryMySQL(db, table, columns)
	}
	if isPostgreSQL(db) {
		return buildInsertQueryPostgreSQL(db, table, columns)
	}
	return "", fmt.Errorf("not supported DB: %#v", dbType(db))
}

func buildUpsertQueryMySQL(db *sql.DB, table string, columns []string) (string, error) {
	placeholders := make([]string, len(columns))
	updates := make([]string, len(columns))
	for i, cname := range columns {
		placeholders[i] = "?"
		updates[i] = fmt.Sprintf("%s=?", cname)
	}
	q := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updates, ", "))
	return q, nil
}

//func buildUpsertQueryPostgreSQL(db *sql.DB, table string, columns []string) (string, error) {
//	placeholders := make([]string, len(columns))
//	updates := make([]string, len(columns))
//	for i, cname := range columns {
//		placeholders[i] = fmt.Sprintf("$%d", i+1)
//		updates[i] = fmt.Sprintf("%s=$%d", cname, i+len(columns)+1)
//	}
//	q := fmt.Sprintf(
//		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO UPDATE SET %s",
//		table,
//		strings.Join(columns, ", "),
//		strings.Join(placeholders, ", "),
//		strings.Join(updates, ", "))
//	return q, nil
//}

// BuildUpsertQuery builds "insert or update" statement with given params.
func BuildUpsertQuery(db *sql.DB, table string, columns []string) (string, error) {
	if isMySQL(db) {
		return buildUpsertQueryMySQL(db, table, columns)
	}
	//if isPostgreSQL(db) {
	//	return buildUpsertQueryPostgreSQL(db, table, columns)
	//}
	return "", fmt.Errorf("BuildUpsertQuery don't supported DB: %#v", dbType(db))
}
