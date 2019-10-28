package models

// User ...
type User struct {
	id       int
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

// Register ...
func (db *DB) Register(email, name, surname, password string) (*User, error) {
	// TODO
	stmt, err := db.Prepare("INSERT INTO user (email, name, surname, password) VALUES()")
}
