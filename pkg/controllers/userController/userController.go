package userController

import (
	"encoding/json"
	"net/http"

	"github.com/atqwerty/choresBackend/internal/utils"
	"github.com/atqwerty/choresBackend/internal/models"
	"github.com/atqwerty/choreshBackend/pkg/controllers/appController"
)

func (app *appController.Session) Register(w http.ResponseWriter, r *http.Request) {
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

func (app *appController.Session) Login(w http.ResponseWriter, r *http.Request) {
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
