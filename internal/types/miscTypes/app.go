package miscTypes

import (
	"github.com/gorilla/mux"
	"github.com/atqwerty/choresBackend/internal/models"
)

type PrototypeSession struct {
	Router         *mux.Router
	Db             models.Datastore
	CurrentUserID  int
	CurrentBoardID int
	Token          string
	RefreshToken   string
}

type KeyPrototype int

type TokenPrototype struct {
	Token string `json:"token"`
}
