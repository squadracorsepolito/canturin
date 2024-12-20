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
	menu.AddRole(application.EditMenu)

	fileMenu := menu.AddSubmenu("File")

	h.register(fileMenu, "Open File", h.openFile)
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
	path, err := application.OpenFileDialog().
		AddFilter("protobuf binary file", "*.binpb").
		AddFilter("protobuf JSON file", "*.json").
		AddFilter("protobuf text file", "*.txtpb").
		PromptForSingleSelection()

	if err != nil {
		return err
	}

	if path == "" {
		return nil
	}

	fileEncoding := acmelib.SaveEncodingWire
	switch filepath.Ext(path) {
	case ".binpb":
		fileEncoding = acmelib.SaveEncodingWire
	case ".json":
		fileEncoding = acmelib.SaveEncodingJSON
	case ".txtpb":
		fileEncoding = acmelib.SaveEncodingText
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	net, err := acmelib.LoadNetwork(file, fileEncoding)
	if err != nil {
		return err
	}

	manager.loadNetwork(net)

	return nil
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
