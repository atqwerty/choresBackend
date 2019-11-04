package models

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User ...
type User struct {
	id       int
	Email    string            `json:"email"`
	Name     string            `json:"name"`
	Surname  string            `json:"surname"`
	Password string            `json:"password"`
	Token    map[string]string `json:"token"`
}

// GetUser ...
func (db *DB) GetUser(id int) (*User, error) {
	user := User{}
	row := db.QueryRow("SELECT email, name, surname FROM user WHERE id=" + strconv.Itoa(id) + ";")
	if err := row.Scan(&user.Email, &user.Name, &user.Surname); err != nil {
		return nil, err
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	user.Token = token
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
	query := `SELECT email, name, surname FROM user WHERE email=? AND password=?;`
	stmt := db.QueryRow(query, email, password)
	if err := stmt.Scan(&user.Email, &user.Name, &user.Surname); err != nil {
		return nil, err
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	user.Token = token
	return &user, nil
}

func generateToken() (map[string]string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Second * 10).Unix()

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	return map[string]string{"token": tokenString}, nil
}
