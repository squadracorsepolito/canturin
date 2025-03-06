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
	configDirName = "canturin"
)

type ConfigNetwork struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func newConfigNetwork(name, path string) ConfigNetwork {
	return ConfigNetwork{Name: name, Path: path}
}

type Config struct {
	Version        int             `json:"version"`
	OpenedNetworks []ConfigNetwork `json:"openedNetworks"`
}

func newDefaultConfig() *Config {
	return &Config{
		Version:        1,
		OpenedNetworks: []ConfigNetwork{},
	}
}

type ConfigService struct {
	mux *sync.RWMutex

	dir         string
	cfgFilePath string

	cfg       *Config
	cfgTained bool

	saveCh chan struct{}
}

func newConfigService() *ConfigService {
	return &ConfigService{
		mux: new(sync.RWMutex),

		dir:         "",
		cfgFilePath: "",

		cfg:       newDefaultConfig(),
		cfgTained: true,

		saveCh: make(chan struct{}),
	}
}

func (cs *ConfigService) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (cs *ConfigService) initFolder() error {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(baseDir, configDirName)
	cs.dir = dir
	if !cs.fileExists(dir) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}

	cfgFilePath := filepath.Join(dir, "config.json")
	cs.cfgFilePath = cfgFilePath
	if cs.fileExists(cfgFilePath) {
		return cs.load(cfgFilePath)
	}

	return cs.handleSave()
}

func (cs *ConfigService) load(path string) error {
	fileBuf, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cfg := new(Config)
	if err := json.Unmarshal(fileBuf, cfg); err != nil {
		return err
	}
	cs.cfg = cfg
	cs.cfgTained = false

	cs.filterOpenedNetworks()

	return nil
}

func (cs *ConfigService) handleSave() error {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	if !cs.cfgTained {
		return nil
	}

	file, err := os.Create(cs.cfgFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileBuf, err := json.Marshal(cs.cfg)
	if err != nil {
		return err
	}

	_, err = file.Write(fileBuf)
	if err != nil {
		return err
	}

	cs.cfgTained = false

	return nil
}

func (cs *ConfigService) run(ctx context.Context) {
	for {
		select {
		case <-cs.saveCh:
			cs.handleSave()

		case <-ctx.Done():
			return
		}
	}
}

func (cs *ConfigService) OnStartup(ctx context.Context, _ application.ServiceOptions) error {
	if err := cs.initFolder(); err != nil {
		return err
	}

	go cs.run(ctx)

	return nil
}

func (cs *ConfigService) OnShutdown() {
	if err := cs.handleSave(); err != nil {
		printError(err)
	}
}

func (cs *ConfigService) sendSave() {
	cs.saveCh <- struct{}{}
}

func (cs *ConfigService) filterOpenedNetworks() {
	tmpOpenedNets := []ConfigNetwork{}

	for _, openNet := range cs.cfg.OpenedNetworks {
		if cs.fileExists(openNet.Path) {
			tmpOpenedNets = append(tmpOpenedNets, openNet)
		}
	}

	cs.cfg.OpenedNetworks = tmpOpenedNets
}

func (cs *ConfigService) addOpenedNetwork(name, path string) {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	tmpOpenNets := []ConfigNetwork{newConfigNetwork(name, path)}

	for _, openNet := range cs.cfg.OpenedNetworks {
		if openNet.Path == path {
			continue
		}

		tmpOpenNets = append(tmpOpenNets, openNet)
	}

	cs.cfg.OpenedNetworks = tmpOpenNets
	cs.cfgTained = true

	cs.sendSave()
}

func (cs *ConfigService) renameOpenedNetwork(name, path string) {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	found := false
	for idx, openNet := range cs.cfg.OpenedNetworks {
		if openNet.Path == path {
			cs.cfg.OpenedNetworks[idx].Name = name
			found = true
			break
		}
	}

	if !found {
		return
	}

	cs.cfgTained = true

	cs.sendSave()
}

func (cs *ConfigService) Get() Config {
	cs.mux.RLock()
	defer cs.mux.RUnlock()

	return *cs.cfg
}
