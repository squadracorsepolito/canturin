package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed frontend/build/*
var assets embed.FS

var app *application.App

var proxy *appProxy

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	proxy = newAppProxy()

	// Initialize the services
	sidebarService := newSidebarService()
	historyService := newHistoryService()

	busService := newBusService()
	messageService := newMessageService()
	signalTypeService := newSignalTypeService()
	signalUnitService := newSignalUnitService()
	signalEnumService := newSignalEnumService()

	// Create a new Wails application by providing the necesvar (sary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app = application.New(application.Options{
		Name:        "canturin",
		Description: "",

		Services: []application.Service{

			application.NewService(sidebarService),
			application.NewService(historyService),

			application.NewService(busService),
			application.NewService(messageService),
			application.NewService(signalTypeService),
			application.NewService(signalUnitService),
			application.NewService(signalEnumService),
		},

		KeyBindings: map[string]func(window *application.WebviewWindow){
			"ctrl+z": func(w *application.WebviewWindow) {
				historyService.Undo()
				historyService.emitHistoryChange()
			},
			"ctrl+y": func(w *application.WebviewWindow) {
				historyService.Redo()
				historyService.emitHistoryChange()
			},
		},

		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},

		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	initMenus()

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "canturin",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	app.OnApplicationEvent(events.Common.ApplicationStarted, func(_ *application.ApplicationEvent) {
		loadNetwork("./testdata/SC24.binpb")
	})

	// Run the application. This blocks until the application has been exited.
	// If an error occurred while running the application, log it and exit.
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
