package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/atqwerty/choresBackend/app/config"
)

// Datastore ...
type Datastore interface {
	GetBoardTasks(boardID int) ([]*Task, error)
	AddTask(string, string, int, int, int) (*Task, error)
	GetTask(int) (*Task, error)
	GetUser(int) (*User, error)
	Register(string, string, string, string) (*User, error)
	Login(string, string) (*User, error)
	AllBoards(userID int) ([]*Board, error)
	AddBoard(title, description string, hostID int) (*Board, error)
	GetBoard(id, userID int) (*Board, error)
	LinkWithUser(boardID, userID int) error
	AddStatus(string, int) (*ReturnStatus, error)
	GetStatuses(int) ([]*ReturnStatus, error)
	UpdateTaskStatus(int, int) error
	GenerateCookie() time.Time
}

// DB ...
type DB struct {
	*sql.DB
}

// InitDB ...
func InitDB(dbConfig *config.DBConfig) (*DB, error) {
	dbURL := fmt.Sprintf(dbConfig.Username + ":" + dbConfig.Password + "@tcp(remotemysql.com:3306)/" + dbConfig.Name)

	db, err := sql.Open(dbConfig.Dialect, dbURL)
	db.SetMaxOpenConns(700)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		fmt.Print("asdf")
		return nil, err
	}

	// ps, err := purse.New(filepath.Join(".", "/sql"))
	// contents, ok := ps.Get("init.sql")
	// if !ok {
	// 	fmt.Println("SQL file not loaded")
	// }

	return &DB{db}, nil
}
