package main

import (
	"github.com/atqwerty/choresBackend/pkg/controllers"
	"github.com/atqwerty/choresBackend/internal/config"
)

func main() {
	config := config.GetConf()
	app := &initController.Session{}
	app.Start(config)
}
