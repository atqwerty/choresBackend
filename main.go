package main

import (
	"github.com/atqwerty/choresBackend/cmd"
	"github.com/atqwerty/choresBackend/internal/config"
)

func main() {
	config := config.GetConf()
	app := &app.App{}
	app.Start(config)
}
