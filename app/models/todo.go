package models

import "strconv"

// Todo ...
type Todo struct {
	id       int
	Title    string `json:"title"`
	Content  string `json:"content"`
	Finished bool   `json:"finished"`
}

// MarkFinished ...
func (t *Todo) MarkFinished() {
	t.Finished = true
}

// MarkUnfinished
func (t *Todo) MarkUnfinished() {
	t.Finished = false
}

// AllTodos ...
func (db *DB) AllTodos() ([]*Todo, error) {
	rows, err := db.Query("SELECT * FROM todo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]*Todo, 0)
	for rows.Next() {
		todo := &Todo{}
		rows.Scan(&todo.id, &todo.Title, &todo.Content)
		todos = append(todos, todo)
	}

	return todos, nil
}

// AddTodo ...
func (db *DB) AddTodo(title, content string) (*Todo, error) {
	stmt, err := db.Prepare("INSERT INTO todo (title, content) VALUES(" + title + ", " + content + ");")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var id int
	if err = stmt.QueryRow(title, content).Scan(&id); err != nil {
		return nil, err
	}

	return db.GetTodo(id)
}

// GetTodo ...
func (db *DB) GetTodo(id int) (*Todo, error) {
	todo := Todo{}
	row := db.QueryRow("SELECT * FROM todo WHERE id=" + strconv.Itoa(id) + ";")
	if err := row.Scan(&todo.id, &todo.Title, &todo.Content); err != nil {
		return nil, nil
	}

	return &todo, nil
}
