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

	h.register(fileMenu, "New Network", h.newNetwork)

	fileMenu.AddSeparator()

	h.register(fileMenu, "Open Network", h.openNetwork)

	fileMenu.AddSeparator()

	h.register(fileMenu, "Save Network", h.saveNetwork)
	h.register(fileMenu, "Save Network As", h.saveNetworkAs)

	fileMenu.AddSeparator()

	h.register(fileMenu, "Import DBC", h.importDBC)
	h.register(fileMenu, "Export DBC", h.exportDBC)

	fileMenu.AddSeparator()

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

func (h *menuHandler) newNetwork(_ *application.Context) error {
	manager.createNetwork()
	return nil
}

func (h *menuHandler) openNetwork(_ *application.Context) error {
	dialog := newOpenNetworkDialog()
	filename, err := dialog.PromptForSingleSelection()
	if err != nil {
		printError(err)
		return nil
	}

	return manager.openNetwork(filename)
}

func (h *menuHandler) saveNetwork(_ *application.Context) error {
	return manager.trySaveNetwork()
}

func (h *menuHandler) saveNetworkAs(_ *application.Context) error {
	dialog := newSaveNetworkDialog()
	filename, err := dialog.PromptForSingleSelection()
	if err != nil {
		printError(err)
		return nil
	}

	return manager.saveNetworkAs(filename)
}

func (h *menuHandler) importDBC(_ *application.Context) error {
	dialog := application.OpenFileDialog()

	dialog.AddFilter("DBC file", "*.dbc")

	path, err := dialog.PromptForSingleSelection()
	if err != nil {
		printError(err)
		return nil
	}

	if path == "" {
		return nil
	}

	dbcFile, err := os.Open(path)
	if err != nil {
		printError(err)
		return err
	}
	defer dbcFile.Close()

	fileName := filepath.Base(path)
	busName := fileName[:len(fileName)-len(filepath.Ext(path))]

	bus, err := acmelib.ImportDBCFile(busName, dbcFile)
	if err != nil {
		printError(err)
		return err
	}

	if err := manager.network.AddBus(bus); err != nil {
		printError(err)
		return err
	}

	manager.reloadNetwork()

	return nil
}

func (h *menuHandler) exportDBC(_ *application.Context) error {
	dialog := application.OpenFileDialog()
	dialog.CanChooseFiles(false)
	dialog.CanChooseDirectories(true)

	path, err := dialog.PromptForSingleSelection()
	if err != nil {
		printError(err)
		return nil
	}

	if path == "" {
		return nil
	}

	return acmelib.ExportNetwork(manager.network, path)
}

func (h *menuHandler) reload(_ *application.Context) error {
	manager.reloadNetwork()

	return nil
}
