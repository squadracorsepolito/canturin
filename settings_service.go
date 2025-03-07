package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	settingsDirName = "canturin"
)

type RecentNetwork struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func newRecentNetwork(name, path string) RecentNetwork {
	return RecentNetwork{Name: name, Path: path}
}

type Settings struct {
	Version        int             `json:"version"`
	RecentNetworks []RecentNetwork `json:"recentNetworks"`
}

func newDefaultSettings() *Settings {
	return &Settings{
		Version:        1,
		RecentNetworks: []RecentNetwork{},
	}
}

type SettingsService struct {
	mux *sync.RWMutex

	dir              string
	settingsFilePath string

	settings       *Settings
	settingsTained bool

	saveCh chan struct{}
}

func newConfigService() *SettingsService {
	return &SettingsService{
		mux: new(sync.RWMutex),

		dir:              "",
		settingsFilePath: "",

		settings:       newDefaultSettings(),
		settingsTained: true,

		saveCh: make(chan struct{}),
	}
}

func (cs *SettingsService) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (cs *SettingsService) initFolder() error {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(baseDir, settingsDirName)
	cs.dir = dir
	if !cs.fileExists(dir) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}

	settingsFilePath := filepath.Join(dir, "settings.json")
	cs.settingsFilePath = settingsFilePath
	if cs.fileExists(settingsFilePath) {
		return cs.load(settingsFilePath)
	}

	return cs.handleSave()
}

func (cs *SettingsService) load(path string) error {
	fileBuf, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cfg := new(Settings)
	if err := json.Unmarshal(fileBuf, cfg); err != nil {
		return err
	}
	cs.settings = cfg
	cs.settingsTained = false

	cs.filterRecentNetworks()

	return nil
}

func (cs *SettingsService) handleSave() error {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	if !cs.settingsTained {
		return nil
	}

	file, err := os.Create(cs.settingsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileBuf, err := json.Marshal(cs.settings)
	if err != nil {
		return err
	}

	_, err = file.Write(fileBuf)
	if err != nil {
		return err
	}

	cs.settingsTained = false

	return nil
}

func (cs *SettingsService) run(ctx context.Context) {
	for {
		select {
		case <-cs.saveCh:
			cs.handleSave()

		case <-ctx.Done():
			return
		}
	}
}

func (cs *SettingsService) OnStartup(ctx context.Context, _ application.ServiceOptions) error {
	if err := cs.initFolder(); err != nil {
		return err
	}

	go cs.run(ctx)

	return nil
}

func (cs *SettingsService) OnShutdown() {
	if err := cs.handleSave(); err != nil {
		printError(err)
	}
}

func (cs *SettingsService) sendSave() {
	cs.saveCh <- struct{}{}
}

func (cs *SettingsService) filterRecentNetworks() {
	tmpRecentNets := []RecentNetwork{}

	for _, recNet := range cs.settings.RecentNetworks {
		if cs.fileExists(recNet.Path) {
			tmpRecentNets = append(tmpRecentNets, recNet)
		}
	}

	cs.settings.RecentNetworks = tmpRecentNets
}

func (cs *SettingsService) addRecentNetwork(name, path string) {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	tmpRecentNets := []RecentNetwork{newRecentNetwork(name, path)}

	for _, recNet := range cs.settings.RecentNetworks {
		if recNet.Path == path {
			continue
		}

		tmpRecentNets = append(tmpRecentNets, recNet)
	}

	cs.settings.RecentNetworks = tmpRecentNets
	cs.settingsTained = true

	cs.sendSave()
}

func (cs *SettingsService) renameRecentNetwork(name, path string) {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	found := false
	for idx, recNet := range cs.settings.RecentNetworks {
		if recNet.Path == path {
			cs.settings.RecentNetworks[idx].Name = name
			found = true
			break
		}
	}

	if !found {
		return
	}

	cs.settingsTained = true

	cs.sendSave()
}

func (cs *SettingsService) Get() Settings {
	cs.mux.RLock()
	defer cs.mux.RUnlock()

	return *cs.settings
}
