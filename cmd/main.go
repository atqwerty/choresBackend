package main

import (
	"github.com/atqwerty/choresBackend/pkg"
	"github.com/atqwerty/choresBackend/internal/config"
)

func main() {
	config := config.GetConf()
	app := &controllers.Session{}
	app.Start(config)
}
