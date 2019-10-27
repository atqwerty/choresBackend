package models

import (
	"database/sql"
	"fmt"

	"github.com/atqwerty/choresBackend/app/config"
)

type Datastore interface {
	AllTodos() ([]*Todo, error)
	AddTodo(string, string) (*Todo, error)
	GetTodo(int) (*Todo, error)
}

// DB ...
type DB struct {
	*sql.DB
}

// InitDB ...
func InitDB(dbConfig *config.DBConfig) (*DB, error) {
	dbURL := fmt.Sprintf(dbConfig.Username + ":" + dbConfig.Password + "@tcp(172.17.0.2:3306)/" + dbConfig.Name)

	db, err := sql.Open(dbConfig.Dialect, dbURL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// ps, err := purse.New(filepath.Join(".", "/sql"))
	// contents, ok := ps.Get("init.sql")
	// if !ok {
	// 	fmt.Println("SQL file not loaded")
	// }

	db.Query("CREATE DATABASE choresdb;")

	return &DB{db}, nil
}
