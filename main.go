package main

import (
	"github.com/atqwerty/choresBackend/app"
	"github.com/atqwerty/choresBackend/app/config"
)

func main() {
	config := config.GetConf()
	app := &app.App{}
	app.Start(config)
}
