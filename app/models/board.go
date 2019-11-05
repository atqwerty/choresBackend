package models

import (
	"log"
	"strconv"
)

// Board ...
type Board struct {
	id          int
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AllBoards ...
func (db *DB) AllBoards(userID int) ([]*Board, error) {
	query := "SELECT title, description FROM board WHERE host_id=" + strconv.Itoa(userID) + ";"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := make([]*Board, 0)
	for rows.Next() {
		board := &Board{}
		rows.Scan(&board.Title, &board.Description)
		boards = append(boards, board)
	}

	return boards, nil
}

// AddBoard ...
func (db *DB) AddBoard(title, description string, hostID int) (*Board, error) {
	stmt, err := db.Prepare("INSERT INTO board (title, description, host_id) VALUES(?, ?, ?);")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	idQuery, err := stmt.Exec(title, description, hostID)
	if err != nil {
		return nil, err
	}

	id64, err := idQuery.LastInsertId()
	if err != nil {
		return nil, err
	}

	id := int(id64)
	return db.GetBoard(id)
}

// LinkWithUser ...
func (db *DB) LinkWithUser(boardID, userID int) error {
	stmt, err := db.Prepare("INSERT INTO user_board VALUES (?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	idQuery, err := stmt.Exec(boardID, userID)
	if err != nil {
		return err
	}

	log.Print(idQuery)
	return nil
}

// GetBoard ...
func (db *DB) GetBoard(id int) (*Board, error) {
	board := Board{}
	row := db.QueryRow("SELECT * FROM board WHERE id=" + strconv.Itoa(id) + ";")
	if err := row.Scan(&board.id, &board.Title, &board.Description); err != nil {
		return nil, err
	}

	return &board, nil
}
