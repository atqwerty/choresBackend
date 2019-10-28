package app

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/atqwerty/choresBackend/app/config"
	"github.com/atqwerty/choresBackend/app/models"
	"github.com/atqwerty/choresBackend/app/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type App struct {
	router *mux.Router
	db     models.Datastore
	// userDb models.UserStore
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
	app.router.HandleFunc("/todo", app.listTodos).Methods("Get")
	app.router.HandleFunc("/todo/{id:[0-9]+}", app.getTodo).Methods("Get")
	app.router.HandleFunc("/todo/create", app.addTodo).Methods("Post")
	app.router.HandleFunc("/register", app.register).Methods("Post")
}

func (app *App) run(addr string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, app.router)
	http.ListenAndServe(addr, loggedRouter)
}

func (app *App) listTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := app.db.AllTodos()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}

func (app *App) addTodo(w http.ResponseWriter, r *http.Request) {
	todo := &models.Todo{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&todo); err != nil {
		utils.BadRequest(w, "payload is required "+err.Error())
		return
	}
	defer r.Body.Close()

	if todo.Title == "" {
		utils.BadRequest(w, "title is required")
		return
	}

	todo, err := app.db.AddTodo(todo.Title, todo.Content)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, todo)
}

func (app *App) getTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "ID must be an int")
	}

	todo, err := app.db.GetTodo(id)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
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

	// utils.BadRequest(w, user.Email+" "+user.Password)
	// return

	user, err := app.db.Register(user.Email, user.Name, user.Surname, user.Password)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, user)
}
