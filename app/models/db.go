package models

import (
	"database/sql"
	"fmt"

	"github.com/atqwerty/choresBackend/app/config"
)

// Datastore ...
type Datastore interface {
	GetBoardTasks(boardID int) ([]*Task, error)
	AddTask(string, string, string, int, int) (*Task, error)
	GetTask(int) (*Task, error)
	GetUser(int) (*User, error)
	Register(string, string, string, string) (*User, error)
	Login(string, string) (*User, error)
	AllBoards(userID int) ([]*Board, error)
	AddBoard(title, description string, hostID int) (*Board, error)
	GetBoard(id, userID int) (*Board, error)
	LinkWithUser(boardID, userID int) error
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

	return &DB{db}, nil
}
