package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	connStr := os.Getenv("DATABASE_URL")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("could not open database: %v", err)
	}

	sqlFilePath := "postgres/init.sql"

	// Execute the SQL file
	err = executeSQLFile(db, sqlFilePath)
	if err != nil {
		return fmt.Errorf("failed to execute SQL file: %v", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func executeSQLFile(db *sql.DB, filePath string) error {
	sqlContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	_, err = db.Exec(string(sqlContent))
	if err != nil {
		return fmt.Errorf("failed to execute SQL statements: %v", err)
	}

	log.Println("SQL file executed successfully.")
	return nil
}
