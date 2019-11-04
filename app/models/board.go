package models

// Board ...
type Board struct {
	id          int
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AddBoard ...
func (db *DB) AddBoard(title, description string, hostID int) ([]*Board, error) {

}
