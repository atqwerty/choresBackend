package models

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User ...
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Password     string    `json:"password"`
	Token        string    `json:"token"`
	ExpireCookie time.Time `json:"expire_cookie"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GetUser ...
func (db *DB) GetUser(id int) (*User, error) {
	user := User{}
	row := db.QueryRow("SELECT * FROM user WHERE id=" + strconv.Itoa(id) + ";")
	if err := row.Scan(&user); err != nil {
		return nil, err
	}

	token, expireCookie, err := generateToken()
	if err != nil {
		return nil, err
	}
	user.Token = token
	user.ExpireCookie = expireCookie
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
	query := `SELECT * WHERE email=? AND password=?;`
	stmt := db.QueryRow(query, email, password)
	if err := stmt.Scan(&user); err != nil {
		return nil, err
	}

	token, expireCookie, err := generateToken()
	if err != nil {
		return nil, err
	}

	user.Token = token
	user.ExpireCookie = expireCookie
	return &user, nil
}

func generateToken() (string, time.Time, error) {
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	expireCookie := time.Now().Add(time.Second * 15)

	claims := Claims{
		"atqwerty",
		jwt.StandardClaims{
			ExpiresAt: expireToken,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := token.SignedString([]byte("secret"))

	return signedToken, expireCookie, nil
}

// GenerateCookie ...
func GenerateCookie() time.Time {
	expireCookie := time.Now().Add(time.Second * 15)

	return expireCookie
}
