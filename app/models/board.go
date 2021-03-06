package models

import (
	"log"
	"strconv"
)

// Board ...
type Board struct {
	ID          int
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Tasks       []*Task `json:"tasks"`
}

type Status struct {
	ID      int
	Status  string `json:"status"`
	BoardID int    `json:"board_id"`
}

type ReturnStatus struct {
	ID int `json:"id"`
	Status string `json:"status"`
}

type IncomingStatus struct {
	Status string `json:"status"`
}

// AllBoards ...
func (db *DB) AllBoards(userID int) ([]*Board, error) {
	// query :=
	rows, err := db.Query("SELECT id, title, description FROM board WHERE id IN (SELECT board_id FROM user_board WHERE user_id=" + strconv.Itoa(userID) + ");")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := make([]*Board, 0)
	for rows.Next() {
		board := &Board{}
		rows.Scan(&board.ID, &board.Title, &board.Description)
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
	query := "INSERT INTO user_board (user_id, board_id) VALUES (" + strconv.Itoa(hostID) + ", " + strconv.Itoa(id) + ");"
	db.Query(query)
	return db.GetBoard(id, hostID)
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
func (db *DB) GetBoard(id, userID int) (*Board, error) {
	board := Board{}
	row := db.QueryRow("SELECT id, title, description FROM board WHERE id=(SELECT board_id FROM user_board WHERE user_id=" + strconv.Itoa(userID) + " AND board_id=" + strconv.Itoa(id) + ");")
	if err := row.Scan(&board.ID, &board.Title, &board.Description); err != nil {
		return nil, err
	}

	return &board, nil
}

// AddStatus ...
func (db *DB) AddStatus(status string, boardID int) (*ReturnStatus, error) {
	stmt, err := db.Prepare("INSERT INTO statuses (status, board_id) VALUES (?, ?);")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	idQuery, err := stmt.Exec(status, strconv.Itoa(boardID))
	if err != nil {
		return nil, err
	}

	id64, err := idQuery.LastInsertId()
	if err != nil {
		return nil, err
	}

	returningStatus := ReturnStatus{}
	row := db.QueryRow("SELECT id, status FROM statuses WHERE id=" + strconv.Itoa(int(id64)) +";")
	if err := row.Scan(&returningStatus.ID, &returningStatus.Status); err != nil {
		return nil, err
	}
	return &returningStatus, nil
}

func (db *DB) GetStatuses(boardID int) ([]*ReturnStatus, error) {
	rows, err := db.Query("SELECT id, status FROM statuses WHERE board_id=" + strconv.Itoa(boardID) + ";")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := make([]*ReturnStatus, 0)
	for rows.Next() {
		status := &ReturnStatus{}
		rows.Scan(&status.ID, &status.Status)
		statuses = append(statuses, status)
	}
	return statuses, nil
}
