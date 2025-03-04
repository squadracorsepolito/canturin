package main

import (
	"embed"
	"log"
	"log/slog"

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

var manager *serviceManager

const testdataPath = "./testdata/SC24.binpb"

// main function serves as the application's entry point.
// Main initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second.
// It subsequently runs the application and logs any error that might occur.
func main() {
	manager = newServiceManager()
	menuHandler := newMenuHandler()

	// Create a new Wails application by providing the necesvar (sary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app = application.New(application.Options{
		Name:        "canturin",
		Description: "",

		Services: manager.getServices(),

		// Key bindings for undo/redo actions, triggering functions on specific key combinations.
		KeyBindings: map[string]func(window *application.WebviewWindow){
			"ctrl+s": func(_ *application.WebviewWindow) {
				if !manager.canSave() {
					return
				}

				if err := manager.saveNetwork(); err != nil {
					log.Print(err)
				}
			},

			"ctrl+z": func(_ *application.WebviewWindow) {
				manager.historySrv.Undo()
				manager.historySrv.emitHistoryChange()
			},
			"ctrl+y": func(_ *application.WebviewWindow) {
				manager.historySrv.Redo()
				manager.historySrv.emitHistoryChange()
			},
		},

		// Configure the asset handler to serve embedded frontend files.
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},

		// macOS-specific options: close the app when the last window is closed.
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},

		LogLevel: slog.LevelError,
	})

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
		BackgroundColour:       application.NewRGB(27, 38, 54),
		URL:                    "/",
		OpenInspectorOnStartup: true,
	})

	// TODO! remove this event in production
	app.OnApplicationEvent(events.Common.ApplicationStarted, func(_ *application.ApplicationEvent) {
		manager.openNetwork(testdataPath)
	})

	menuHandler.init()

	// Run the application. This blocks until the application has been exited.
	// If an error occurred while running the application, log it and exit.
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
