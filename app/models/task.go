package models

import (
	"strconv"
)

// Task ...
type Task struct {
	id       int
	Title    string `json:"title"`
	Content  string `json:"content"`
	Finished bool   `json:"finished"`
}

// MarkFinished ...
func (t *Task) MarkFinished() {
	t.Finished = true
}

// MarkUnfinished ...
func (t *Task) MarkUnfinished() {
	t.Finished = false
}

// AllTasks ...
func (db *DB) AllTasks() ([]*Task, error) {
	rows, err := db.Query("SELECT title, description FROM task;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		task := &Task{}
		rows.Scan(&task.Title, &task.Content)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// AddTask ...
func (db *DB) AddTask(title, content string) (*Task, error) {
	stmt, err := db.Prepare("INSERT INTO task (title, content) VALUES(?, ?);")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	idQuery, err := stmt.Exec(title, content)
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
	row := db.QueryRow("SELECT * FROM task WHERE id=" + strconv.Itoa(id) + ";")
	if err := row.Scan(&task.id, &task.Title, &task.Content); err != nil {
		return nil, err
	}

	return &task, nil
}
