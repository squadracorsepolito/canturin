package main

import (
	"os"
	"path/filepath"

	"github.com/squadracorsepolito/acmelib"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type menuHandler struct{}

func newMenuHandler() *menuHandler {
	return &menuHandler{}
}

func (h *menuHandler) init() {
	app := application.Get()

	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)

	fileMenu := menu.AddSubmenu("File")

	h.register(fileMenu, "Open File", h.openFile)
	h.register(fileMenu, "Save File", h.saveFile)
	h.register(fileMenu, "Save File As", h.saveFileAs)

	h.register(fileMenu, "Import DBC", h.importDBC)

	h.register(fileMenu, "Reload", h.reload)

	app.SetMenu(menu)
}

func (h *menuHandler) register(menu *application.Menu, name string, cb func(*application.Context) error) {
	menu.Add(name).OnClick(func(ctx *application.Context) {
		if err := cb(ctx); err != nil {
			application.ErrorDialog().SetMessage(err.Error()).Show()
		}
	})
}

func (h *menuHandler) openFile(_ *application.Context) error {
	dialog := application.OpenFileDialog()

	dialog.AddFilter("protobuf binary file", "*.binpb")
	dialog.AddFilter("protobuf JSON file", "*.json")
	dialog.AddFilter("protobuf text file", "*.txtpb")

	filename, err := dialog.PromptForSingleSelection()
	if err != nil {
		return err
	}

	return manager.openNetwork(filename)
}

func (h *menuHandler) saveFile(_ *application.Context) error {
	return manager.trySaveNetwork()
}

func (h *menuHandler) saveFileAs(_ *application.Context) error {
	dialog := application.SaveFileDialog()

	dialog.AddFilter("protobuf binary file", "*.binpb")
	dialog.AddFilter("protobuf JSON file", "*.json")
	dialog.AddFilter("protobuf text file", "*.txtpb")

	filename, err := dialog.PromptForSingleSelection()
	if err != nil {
		return err
	}

	return manager.saveNetworkAs(filename)
}

func (h *menuHandler) importDBC(_ *application.Context) error {
	path, err := application.OpenFileDialog().
		AddFilter("DBC file", "*.dbc").
		PromptForSingleSelection()

	if err != nil {
		return err
	}

	if path == "" {
		return nil
	}

	dbcFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dbcFile.Close()

	fileName := filepath.Base(path)
	busName := fileName[:len(fileName)-len(filepath.Ext(path))]

	bus, err := acmelib.ImportDBCFile(busName, dbcFile)
	if err != nil {
		return err
	}

	if err := manager.network.AddBus(bus); err != nil {
		return err
	}

	manager.reloadNetwork()

	return nil
}

func (h *menuHandler) reload(_ *application.Context) error {
	manager.reloadNetwork()

	return nil
}
