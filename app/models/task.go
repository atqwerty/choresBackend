package models

import (
	"strconv"
)

// Task ...
type Task struct {
	ID          int
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Finished    bool   `json:"finished"`
}

// MarkFinished ...
func (t *Task) MarkFinished() {
	t.Finished = true
}

// MarkUnfinished ...
func (t *Task) MarkUnfinished() {
	t.Finished = false
}

// GetBoardTasks ...
func (db *DB) GetBoardTasks(boardID int) ([]*Task, error) {
	rows, err := db.Query("SELECT id, title, description, status FROM task WHERE board_id=" + strconv.Itoa(boardID) + ";")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		task := &Task{}
		rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// AddTask ...
func (db *DB) AddTask(title, description string, status, boardID, hostID int) (*Task, error) {
	stmtStatus, err := db.Prepare("INSERT INTO statuses (status, board_id) VALUES (?, ?);")
	if err != nil {
		return nil, err
	}
	defer stmtStatus.Close()

	statusIDQuery, err := stmtStatus.Exec(status, boardID)
	if err != nil {
		return nil, err
	}

	statusID64, err := statusIDQuery.LastInsertId()
	if err != nil {
		return nil, err
	}

	statusID := int(statusID64)

	stmt, err := db.Prepare("INSERT INTO task (title, description, status, board_id, host_id) VALUES(?, ?, ?, ?, ?);")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	idQuery, err := stmt.Exec(title, description, statusID, boardID, hostID)
	if err != nil {
		return nil, err
	}

	id64, err := idQuery.LastInsertId()
	if err != nil {
		return nil, err
	}

	id := int(id64)
	return db.GetTask(id)
}

// GetTask ...
func (db *DB) GetTask(id int) (*Task, error) {
	task := Task{}
	row := db.QueryRow("SELECT id, title, description, status FROM task WHERE id=" + strconv.Itoa(id) + ";")
	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status); err != nil {
		return nil, err
	}

	return &task, nil
}
