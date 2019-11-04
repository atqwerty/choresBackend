package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/atqwerty/choresBackend/app/config"
	"github.com/atqwerty/choresBackend/app/models"
	"github.com/atqwerty/choresBackend/app/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Key int

const MyKey Key = 0

type App struct {
	router *mux.Router
	db     models.Datastore
	// userDb models.UserStore
}

type Token struct {
	token string `json:"token"`
}

func (app *App) Start(conf *config.Config) {
	db, err := models.InitDB(conf.DBConfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	app.db = db
	app.router = mux.NewRouter()
	app.initRouters()
	app.run(":8080")
}

func (app *App) initRouters() {
	app.router.HandleFunc("/", app.status).Methods("Get")
	app.router.HandleFunc("/todo", validate(app.listTasks)).Methods("Get")
	app.router.HandleFunc("/todo/{id:[0-9]+}", validate(app.getTask)).Methods("Get")
	app.router.HandleFunc("/todo/create", validate(app.addTask)).Methods("Post")
	app.router.HandleFunc("/register", app.register).Methods("Post")
	app.router.HandleFunc("/login", app.login).Methods("Post")
	app.router.HandleFunc("/refresh", app.refresh).Methods("Post")
}

func (app *App) run(addr string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, app.router)
	http.ListenAndServe(addr, loggedRouter)
}

func (app *App) listTasks(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unathorized", 401)
		return
	}

	fmt.Fprintf(w, "Hello %s", claims.Username)
	tasks, err := app.db.AllTasks()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, tasks)
}

func (app *App) addTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Hello %s", claims.Username)

	task := &models.Task{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		utils.BadRequest(w, "payload is required "+err.Error())
		return
	}
	defer r.Body.Close()

	if task.Title == "" {
		utils.BadRequest(w, "title is required")
		return
	}

	task, err := app.db.AddTask(task.Title, task.Content)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, task)
}

func (app *App) getTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Hello %s", claims.Username)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "ID must be an int")
	}

	task, err := app.db.GetTask(id)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, task)
}

func (app *App) status(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, "API is up and working!")
}

func (app *App) register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		utils.BadRequest(w, "payload is required "+err.Error())
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

	user, err := app.db.Register(user.Email, user.Name, user.Surname, user.Password)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, user)
}

func (app *App) login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		utils.BadRequest(w, "payload is required"+err.Error())
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

	user, err := app.db.Login(user.Email, user.Password)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	cookie := http.Cookie{Name: "Auth", Value: user.Token, Expires: user.ExpireCookie, HttpOnly: true}
	http.SetCookie(w, &cookie)
	utils.RespondJSON(w, http.StatusOK, user)
}

func validate(page http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("Auth")
		if err != nil {
			http.Error(res, "Unauthorized", 401)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			http.Error(res, "Unauthorized", 401)
			return
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			ctx := context.WithValue(req.Context(), MyKey, *claims)
			page(res, req.WithContext(ctx))
		} else {
			http.Error(res, "Unauthorized", 401)
			return
		}
	})
}

func (app *App) refresh(w http.ResponseWriter, r *http.Request) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	reqToken = splitToken[1]

	fmt.Fprintf(w, "%s", reqToken)
	return
}
