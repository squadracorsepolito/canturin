package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func initMenus() {
	app := application.Get()

	log.Print(app.GetPID())
}
