package main

import (
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type keybindingsHandler struct {
	m map[string]func(window *application.WebviewWindow)
}

func newKeybindingsHandler() *keybindingsHandler {
	return &keybindingsHandler{
		m: make(map[string]func(window *application.WebviewWindow)),
	}
}

func (kh *keybindingsHandler) init() {
	kh.registerWindow(kh.appendKeys(kh.metaKey(), "s"), kh.saveHandler)
	kh.registerWindow(kh.appendKeys(kh.metaKey(), "z"), kh.undoHandler)
	kh.registerWindow(kh.appendKeys(kh.metaKey(), "y"), kh.redoHandler)
}

func (kh *keybindingsHandler) getWindowKeybindings() map[string]func(window *application.WebviewWindow) {
	return kh.m
}

func (kh *keybindingsHandler) registerWindow(key string, cb func(*application.WebviewWindow)) {
	kh.m[key] = cb
}

func (kh *keybindingsHandler) appendKeys(keys ...string) string {
	if len(keys) == 0 {
		return ""
	}

	res := keys[0]
	for _, key := range keys {
		res += "+" + key
	}

	return res
}

func (kh *keybindingsHandler) metaKey() string {
	if runtime.GOOS == "darwin" {
		return "CMD"
	}
	return "ctrl"
}

func (kh *keybindingsHandler) saveHandler(_ *application.WebviewWindow) {
	if err := manager.trySaveNetwork(); err != nil {
		log.Print(err)
	}
}

func (kh *keybindingsHandler) undoHandler(_ *application.WebviewWindow) {
	manager.historySrv.Undo()
	manager.historySrv.emitHistoryChange()
}

func (kh *keybindingsHandler) redoHandler(_ *application.WebviewWindow) {
	manager.historySrv.Redo()
	manager.historySrv.emitHistoryChange()
}
