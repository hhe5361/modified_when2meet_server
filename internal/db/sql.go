package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func InitDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	return db
}

func ScanStruct[T any](rows *sql.Rows) (T, error) {
	var dest T
	v := reflect.ValueOf(&dest).Elem()
	t := v.Type()

	columns, err := rows.Columns()
	if err != nil {
		var zero T
		return zero, err
	}

	fieldMap := map[string]int{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		fieldMap[tag] = i
	}

	scanTargets := make([]any, len(columns))
	columnPointers := make([]any, len(columns))
	for i, col := range columns {
		if fieldIndex, ok := fieldMap[col]; ok {
			field := v.Field(fieldIndex)
			columnPointers[i] = field.Addr().Interface()
		} else {
			var dummy any
			columnPointers[i] = &dummy
		}
		scanTargets[i] = columnPointers[i]
	}

	if err := rows.Scan(scanTargets...); err != nil {
		var zero T
		return zero, err
	}

	return dest, nil
}

// get row
func QueryOnlyRow[T any](db *sql.DB, query string, scanFunc func(*sql.Rows) (T, error), args ...any) (T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		var zero T
		return zero, errors.New("query exec is failed")
	}
	defer rows.Close()

	if !rows.Next() {
		var zero T
		return zero, sql.ErrNoRows
	}

	return scanFunc(rows)
}

// get rows
func QueryRows[T any](db *sql.DB, query string, scanFunc func(*sql.Rows) (T, error), args ...any) ([]T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		result, err := scanFunc(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func QueryExec(db *sql.DB, query string, args ...any) (int64, error) {
	res, err := db.Exec(query, args...)
	if err != nil {
		return -1, fmt.Errorf("query execution failed: %v, query: %s, args: %v", err, query, args)
	}
	id, _ := res.LastInsertId()

	return id, nil
}
