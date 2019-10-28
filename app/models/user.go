package models

import "strconv"

// User ...
type User struct {
	id       int
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

// GetUser ...
func (db *DB) GetUser(id int) (*User, error) {
	user := User{}
	row := db.QueryRow("SELECT email, name, surname FROM user WHERE id=" + strconv.Itoa(id) + ";")
	if err := row.Scan(&user.Email, &user.Name, &user.Surname); err != nil {
		return nil, err
	}

	return &user, nil
}

// Register ...
func (db *DB) Register(email, name, surname, password string) (*User, error) {
	stmt, err := db.Prepare("INSERT INTO user (email, name, surname, password) VALUES(?, ?, ?, ?);")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	idQuery, err := stmt.Exec(email, name, surname, password)
	if err != nil {
		return nil, err
	}

	id64, err := idQuery.LastInsertId()
	if err != nil {
		return nil, err
	}

	id := int(id64)
	return db.GetUser(id)
}

// Login ...
func (db *DB) Login(email, password string) (*User, error) {
	user := User{}
	row, err := db.Query("SELECT email, name, surname FROM user WHERE email=", email, " AND password=", password, ";")
	if err != nil {
		return nil, err
	}

	row.Scan(&user.Email, &user.Name, &user.Surname)

	return &user, nil
}
