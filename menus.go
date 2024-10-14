package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func initMenus() {
	app := application.Get()

	// open file acmelib.

	log.Print(app.GetPID())
}
