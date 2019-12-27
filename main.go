package main

import (
	"github.com/atqwerty/choresBackend/cmd"
	"github.com/atqwerty/choresBackend/app/config"
)

func main() {
	config := config.GetConf()
	app := &app.App{}
	app.Start(config)
}
