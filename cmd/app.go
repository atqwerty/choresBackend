package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/atqwerty/choresBackend/internal/config"
	"github.com/atqwerty/choresBackend/internal/models"
	"github.com/atqwerty/choresBackend/internal/utils"
	"github.com/atqwerty/choresBackend/internal/types/miscTypes"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const MyKey miscTypes.KeyPrototype = 0

type Session miscTypes.PrototypeSession

type Token miscTypes.TokenPrototype

type TaskStatus struct {
	StatusID int `json:"status_id"`
	TaskID int `json:"task_id"`
}

func (app *Session) Start(conf *config.Config) {
	db, err := models.InitDB(conf.DBConfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	app.Db = db
	app.Router = mux.NewRouter()
	app.initRouters()
	app.run()
}

func (app *Session) initRouters() {
	app.Router.HandleFunc("/", app.status).Methods("Get")
	app.Router.HandleFunc("/board/{board_id}/task/{task_id}", validate(app.getTask)).Methods("Get")
	app.Router.HandleFunc("/board/{board_id:[0-9]+}/task/create", validate(app.addTask)).Methods("Post")
	app.Router.HandleFunc("/board/all", validate(app.listBoards)).Methods("Get")
	app.Router.HandleFunc("/board/{board_id:[0-9]+}", validate(app.getBoard)).Methods("Get")
	app.Router.HandleFunc("/board/create", validate(app.addBoard)).Methods("Post")
	app.Router.HandleFunc("/board/newStatus", validate(app.newStatus)).Methods("Post")
	app.Router.HandleFunc("/board/newStatusMobile", validate(app.newStatusMobile)).Methods("Post")
	app.Router.HandleFunc("/register", app.register).Methods("Post")
	app.Router.HandleFunc("/login", app.login).Methods("Post")
	app.Router.HandleFunc("/refresh", app.refresh).Methods("Get")
	app.Router.HandleFunc("/board/{board_id:[0-9]+}/getStatuses", validate(app.getStatuses)).Methods("Get")
	app.Router.HandleFunc("/board/{board_id:[0-9]+}/updateStatus", validate(app.updateStatus)).Methods("Post")
}

func (app *Session) run() {
	port := os.Getenv("PORT")
	loggedRouter := handlers.LoggingHandler(os.Stdout, app.Router)
	http.ListenAndServe(":" + port, loggedRouter)
}

func (app *Session) listBoards(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unathorized", 401)
		return
	}

	_ = claims
	boards, err := app.Db.AllBoards(app.CurrentUserID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, boards)
}

func (app *Session) newStatus(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unathorized", 401)
		// return nil
	}

	_ = claims

	status := &models.IncomingStatus{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&status); err != nil {
		utils.ServerError(w, err)
	}
	returnedStatus, err := app.Db.AddStatus(status.Status, app.CurrentBoardID)
	if err != nil {
		utils.ServerError(w, err)
	}

	utils.RespondJSON(w, http.StatusOK, returnedStatus)
}

func (app *Session) newStatusMobile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unathorized", 401)
	}

	_ = claims

	status := &models.ReturnStatus{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&status); err != nil {
		utils.ServerError(w, err)
	}
	returnedStatus, err := app.Db.AddStatus(status.Status, status.ID)
	if err != nil {
		utils.ServerError(w, err)
	}

	utils.RespondJSON(w, http.StatusOK, returnedStatus)
}

func (app *Session) updateStatus(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unathorized", 401)
	}

	_ = claims

	var myStoredVariable = &TaskStatus{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&myStoredVariable); err != nil {
		utils.BadRequest(w, err.Error())
	}
	err := app.Db.UpdateTaskStatus(myStoredVariable.StatusID, myStoredVariable.TaskID)
	if err != nil {
		utils.ServerError(w, err)
	}

	utils.RespondJSON(w, http.StatusOK, myStoredVariable)
}

func (app *Session) getStatuses(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unathorized", 401)
	}

	_ = claims
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["board_id"])
	if err != nil {
		utils.BadRequest(w, "ID must be an int")
	}

	statuses, err := app.Db.GetStatuses(id)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, statuses)
}

func (app *Session) listTasks(w http.ResponseWriter, r *http.Request, boardID int) []*models.Task {
	claims, ok := r.Context().Value(MyKey).(models.Claims)
	if !ok {
		http.Error(w, "Unathorized", 401)
		return nil
	}

	_ = claims
	tasks, err := app.Db.GetBoardTasks(boardID)
	if err != nil {
		utils.ServerError(w, err)
		return nil
	}

	// utils.RespondJSON(w, http.StatusOK, tasks)
	return tasks
}

func (app *Session) addBoard(w http.ResponseWriter, r *http.Request) {
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

	board, err := app.Db.AddBoard(board.Title, board.Description, app.CurrentUserID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, board)
}

func (app *Session) getBoard(w http.ResponseWriter, r *http.Request) {
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

	// board, err := app.db.GetBoard(id, app.CurrentUserID)
	// if err != nil {
	// utils.ServerError(w, err)
	// return
	// }

	app.CurrentBoardID = id
	utils.RespondJSON(w, http.StatusOK, app.listTasks(w, r, id))
	//utils.RespondJSON(w, http.StatusOK, board)
}

func (app *Session) addTask(w http.ResponseWriter, r *http.Request) {
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

	task := &models.IncomingTask{}
	returnTask := &models.Task{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		utils.BadRequest(w, "asdfasfasf")
		return
	}
	defer r.Body.Close()

	if task.Title == "" {
		utils.BadRequest(w, "title is required")
		return
	}

	returnTask, err = app.Db.AddTask(task.Title, task.Description, task.Status, boardID, app.CurrentUserID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, returnTask)
}

func (app *Session) getTask(w http.ResponseWriter, r *http.Request) {
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

	task, err := app.Db.GetTask(taskID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, task)
}

func (app *Session) status(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, "API is up and working!")
}

func (app *Session) register(w http.ResponseWriter, r *http.Request) {
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

func (app *Session) login(w http.ResponseWriter, r *http.Request) {
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

func (app *Session) refresh(w http.ResponseWriter, r *http.Request) {
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
