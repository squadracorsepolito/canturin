package main

import (
	"embed"
	_ "embed"

	// "fmt"
	"log"
	// "os"

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

var manager *serviceManager

// main function serves as the application's entry point.
// Main initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second.
// It subsequently runs the application and logs any error that might occur.
func main() {
	proxy = newAppProxy()

	manager = newServiceManager()

	// Path to the file used for loading network data.
	// infilepath := "./testdata/SC24.binpb"

	// Initialize the services
	// sidebarService := newSidebarService()
	// historyService := newHistoryService()

	// busService := newBusService()
	// nodeService := newNodeService()
	// messageService := newMessageService()
	// signalTypeService := newSignalTypeService()
	// signalUnitService := newSignalUnitService()
	// signalEnumService := newSignalEnumService()

	// Create a new Wails application by providing the necesvar (sary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app = application.New(application.Options{
		Name:        "canturin",
		Description: "",

		Services: manager.getServices(),

		// Services: []application.Service{

		// 	application.NewService(sidebarService),
		// 	application.NewService(historyService),

		// 	application.NewService(busService),
		// 	application.NewService(nodeService),
		// 	application.NewService(messageService),
		// 	application.NewService(signalTypeService),
		// 	application.NewService(signalUnitService),
		// 	application.NewService(signalEnumService),
		// },

		// Key bindings for undo/redo actions, triggering functions on specific key combinations.
		KeyBindings: map[string]func(window *application.WebviewWindow){
			"ctrl+z": func(w *application.WebviewWindow) {
				// historyService.Undo()
				// historyService.emitHistoryChange()
				manager.history.Undo()
				manager.history.emitHistoryChange()
			},
			"ctrl+y": func(w *application.WebviewWindow) {
				// historyService.Redo()
				// historyService.emitHistoryChange()
				manager.history.Redo()
				manager.history.emitHistoryChange()
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

	// Emit an application event to load network data when the application starts.
	app.OnApplicationEvent(events.Common.ApplicationStarted, func(_ *application.ApplicationEvent) {
		// loadNetwork(infilepath)
	})

	// Run the application. This blocks until the application has been exited.
	// If an error occurred while running the application, log it and exit.
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// function to process and read files
// func processFile(filePath string) {
// 	// read selected file
// 	data, err := os.ReadFile(filePath)
// 	if err != nil {
// 		fmt.Println("Error reading file:", err)
// 		return
// 	}

// 	// Show file content on the terminal
// 	fmt.Printf("File %s loaded successfully! Content:\n%s\n", filePath, string(data))
// }
