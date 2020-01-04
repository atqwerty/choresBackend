package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"

	"github.com/atqwerty/choresBackend/internal/utils"
	"github.com/atqwerty/choresBackend/internal/models"
	"github.com/dgrijalva/jwt-go"
)

type UserController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
}

func (app *Session) Register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		utils.BadRequest(w, "payload is required ")
		return
	}
	defer r.Body.Close()

	if user.Email == "" {
		utils.BadRequest(w, "email is required")
		return
	}
	if user.Name == "" {
		utils.BadRequest(w, "name is required")
		return
	}
	if user.Surname == "" {
		utils.BadRequest(w, "surname is required")
		return
	}
	if user.Password == "" {
		utils.BadRequest(w, "password is required")
		return
	}

	user, err := app.Db.Register(user.Email, user.Name, user.Surname, user.Password)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	app.CurrentUserID = user.ID

	utils.RespondJSON(w, http.StatusOK, user)
}

func (app *Session) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		utils.BadRequest(w, "payload is required")
		return
	}
	defer r.Body.Close()

	if user.Email == "" {
		utils.BadRequest(w, "email is missing")
		return
	}
	if user.Password == "" {
		utils.BadRequest(w, "password is missing")
		return
	}

	user, err := app.Db.Login(user.Email, user.Password)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	app.CurrentUserID = user.ID
	app.Token = user.Token
	app.RefreshToken = user.RefreshToken

	cookie := http.Cookie{Name: "Auth", Value: user.Token, Expires: user.ExpireCookie, HttpOnly: true}
	http.SetCookie(w, &cookie)
	utils.RespondJSON(w, http.StatusOK, user)
}

func (app *Session) Refresh(w http.ResponseWriter, r *http.Request) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = strings.Replace(splitToken[1], "\n", "", -1)

	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && token.Raw == app.RefreshToken {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
		return
	}

	cookie := http.Cookie{Name: "Auth", Value: app.Token, Expires: app.Db.GenerateCookie(), HttpOnly: true}
	http.SetCookie(w, &cookie)

	utils.RespondJSON(w, http.StatusOK, reqToken)
	return
}
