package main

import (
	"github.com/atqwerty/choresBackend/app"
	"choresBackend/app/config"
)

func main() {
	config := config.GetConf()
	app := &app.App{}
	app.Start(config)
}
