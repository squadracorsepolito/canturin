package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func initMenus() {
	app := application.Get()

	/*open file acmelib*/
	
	// Add the menu bar to perform tasks like Edit, File, Help
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.HelpMenu)

	// Create "File" submenu with an "Open File" option.
	// When clicked, it opens a file dialog for the user to select a file.
	openMenu := menu.AddSubmenu("File")
	openMenu.Add("Open File").OnClick(func(ctx *application.Context) {
		result, _ := application.OpenFileDialog().
			CanChooseFiles(true).
			PromptForSingleSelection()
		if result != "" {
			loadNetwork(result)
		} else {
			application.InfoDialog().SetMessage("No file selected").Show()
		}
	})

	app.SetMenu(menu)

	log.Print(app.GetPID())
}
