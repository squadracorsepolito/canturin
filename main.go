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

//go:embed frontend/build
var assets embed.FS

var app *application.App

var proxy *appProxy

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	proxy = newAppProxy()

	// Initialize the services
	sidebarSrv := newSidebarService()

	msgServ := newMessageService()
	sigTypeServ := newSignalTypeService()
	sigUnitServ := newSignalUnitService()

	// Create a new Wails application by providing the necesvar (sary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app = application.New(application.Options{
		Name:        "canturin",
		Description: "",
		Services: []application.Service{

			application.NewService(sidebarSrv),

			application.NewService(msgServ),
			application.NewService(sigTypeServ),
			application.NewService(sigUnitServ),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	initMenus()

	// Add the menu bar to perform tasks like Edit, File, Help
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.HelpMenu)

	openMenu := menu.AddSubmenu("File")
	openMenu.Add("Open File").OnClick(func(ctx *application.Context) {
		result, _ := application.OpenFileDialog().
			CanChooseFiles(true).
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		} else {
			application.InfoDialog().SetMessage("No file selected").Show()
		}
	})

	app.SetMenu(menu)

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
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// // Create a goroutine that emits an event containing the current time every second.
	// // The frontend can listen to this event and update the UI accordingly.
	// go func() {
	// 	for {
	// 		now := time.Now().Format(time.RFC1123)
	// 		app.Events.Emit(&application.WailsEvent{
	// 			Name: "time",
	// 			Data: now,
	// 		})
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	app.OnApplicationEvent(events.Common.ApplicationStarted, func(_ *application.ApplicationEvent) {
		loadNetwork()
	})

	// Run the application. This blocks until the application has been exited.
	// If an error occurred while running the application, log it and exit.
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
