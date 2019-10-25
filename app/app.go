package app

import "github.com/gorilla/mux"

type App struct {
	router *mux.Router
	db     models.Datastore
}
