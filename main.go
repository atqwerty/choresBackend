package main

import (
	"github.com/atqwerty/choresBackend/cmd"
	"github.com/atqwerty/choresBackend/internal/config"
)

func main() {
	config := config.GetConf()
	app := &cmd.Session{}
	app.Start(config)
}
