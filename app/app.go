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
	router         *mux.Router
	db             models.Datastore
	userID         int
	currentBoardID int
	token          string
	refreshToken   string
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
	app.router.HandleFunc("/board/{board_id}/task/{task_id}", validate(app.getTask)).Methods("Get")
	app.router.HandleFunc("/board/{board_id:[0-9]+}/task/create", validate(app.addTask)).Methods("Post")
	app.router.HandleFunc("/board/all", validate(app.listBoards)).Methods("Get")
	app.router.HandleFunc("/board/{board_id:[0-9]+}", validate(app.getBoard)).Methods("Get")
	app.router.HandleFunc("/board/create", validate(app.addBoard)).Methods("Post")
	app.router.HandleFunc("/board/newStatus", validate(app.newStatus)).Methods("Post")
	app.router.HandleFunc("/register", app.register).Methods("Post")
	app.router.HandleFunc("/login", app.login).Methods("Post")
	app.router.HandleFunc("/refresh", app.refresh).Methods("Get")
}

func (app *App) run(addr string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, app.router)
	http.ListenAndServe(addr, loggedRouter)
}

func (app *App) listBoards(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unathorized", 401)
		return
	}

	_ = claims
	tasks, err := app.db.AllBoards(app.userID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, tasks)
}

func (app *App) newStatus(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unathorized", 401)
		// return nil
	}

	_ = claims

	status := &models.Status{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&status); err != nil {
		utils.ServerError(w, err)
	}
	status, err := app.db.AddStatus(status.Status, app.currentBoardID)
	if err != nil {
		utils.ServerError(w, err)
	}

	utils.RespondJSON(w, http.StatusOK, status)
}

func (app *App) listTasks(w http.ResponseWriter, r *http.Request, boardID int) []*models.Task {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unathorized", 401)
		return nil
	}

	_ = claims
	tasks, err := app.db.GetBoardTasks(boardID)
	if err != nil {
		utils.ServerError(w, err)
		return nil
	}

	// utils.RespondJSON(w, http.StatusOK, tasks)
	return tasks
}

func (app *App) addBoard(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unauthorized", 401)
		return
	}

	_ = claims
	board := &models.Board{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&board); err != nil {
		utils.BadRequest(w, "payload is required ")
		return
	}
	defer r.Body.Close()

	if board.Title == "" {
		utils.BadRequest(w, "title is required")
		return
	}

	board, err := app.db.AddBoard(board.Title, board.Description, app.userID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, board)
}

func (app *App) getBoard(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unauthorized", 401)
		return
	}

	_ = claims
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["board_id"])
	if err != nil {
		utils.BadRequest(w, "ID must be an int")
	}

	// board, err := app.db.GetBoard(id, app.userID)
	// if err != nil {
	// utils.ServerError(w, err)
	// return
	// }

	utils.RespondJSON(w, http.StatusOK, app.listTasks(w, r, id))
	//utils.RespondJSON(w, http.StatusOK, board)
}

func (app *App) addTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unauthorized", 401)
		return
	}

	_ = claims
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["board_id"])
	if err != nil {
		utils.BadRequest(w, "ID of board must be an int")
	}

	task := &models.Task{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		utils.BadRequest(w, "payload is required ")
		return
	}
	defer r.Body.Close()

	if task.Title == "" {
		utils.BadRequest(w, "title is required")
		return
	}

	task, err = app.db.AddTask(task.Title, task.Description, task.Status, boardID, app.userID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, task)
}

func (app *App) getTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unauthorized", 401)
		return
	}

	_ = claims
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["task_id"])
	if err != nil {
		utils.BadRequest(w, "ID of task must be an int")
	}

	task, err := app.db.GetTask(taskID)
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

	user, err := app.db.Register(user.Email, user.Name, user.Surname, user.Password)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	app.userID = user.ID

	utils.RespondJSON(w, http.StatusCreated, user)
}

func (app *App) login(w http.ResponseWriter, r *http.Request) {
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

	user, err := app.db.Login(user.Email, user.Password)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	app.userID = user.ID
	app.token = user.Token
	app.refreshToken = user.RefreshToken

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

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && token.Raw == app.refreshToken {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
		return
	}

	cookie := http.Cookie{Name: "Auth", Value: app.token, Expires: models.GenerateCookie(), HttpOnly: true}
	http.SetCookie(w, &cookie)

	utils.RespondJSON(w, http.StatusCreated, reqToken)
	return
}
