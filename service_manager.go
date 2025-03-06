package main

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/squadracorsepolito/acmelib"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/exp/maps"
)

type serviceManager struct {
	filePath string

	configSrv *ConfigService

	mux     *sync.RWMutex
	network *acmelib.Network

	sidebarSrv *SidebarService

	historySrv *HistoryService
	historyCtr *historyController

	networkSrv *NetworkService

	busSrv *BusService
	busCtr *busController

	nodeSrv *NodeService
	nodeCtr *nodeController

	messageSrv *MessageService
	messageCtr *messageController

	signalSrv *SignalService
	signalCtr *signalController

	signalTypeSrv *SignalTypeService
	signalTypeCtr *signalTypeController

	signalUnitSrv *SignalUnitService
	signalUnitCtr *signalUnitController

	signalEnumSrv *SignalEnumService
	signalEnumCtr *signalEnumController
}

func newServiceManager() *serviceManager {
	mux := &sync.RWMutex{}

	sidebarSrv := newSidebarService()
	sidebarCtr := sidebarSrv.getController()

	historySrv := newHistoryService()
	historyCtr := historySrv.getController()

	signalTypeSrv := newSignalTypeService(mux, sidebarCtr)
	signalTypeSrv.setHistoryController(historyCtr)
	signalTypeCtr := signalTypeSrv.getController()

	signalUnitSrv := newSignalUnitService(mux, sidebarCtr)
	signalUnitSrv.setHistoryController(historyCtr)
	signalUnitCtr := signalUnitSrv.getController()

	signalEnumSrv := newSignalEnumService(mux, sidebarCtr)
	signalEnumSrv.setHistoryController(historyCtr)
	signalEnumCtr := signalEnumSrv.getController()

	signalSrv := newSignalService(mux, sidebarCtr, signalTypeCtr, signalUnitCtr, signalEnumCtr)
	signalSrv.setHistoryController(historyCtr)
	signalCtr := signalSrv.getController()

	messageSrv := newMessageService(mux, sidebarCtr, signalCtr)
	messageSrv.setHistoryController(historyCtr)
	messageCtr := messageSrv.getController()

	busSrv := newBusService(mux, sidebarCtr)
	busSrv.setHistoryController(historyCtr)
	busCtr := busSrv.getController()

	nodeSrv := newNodeService(mux, sidebarCtr, busSrv, messageCtr)
	nodeSrv.setHistoryController(historyCtr)
	nodeCtr := nodeSrv.getController()

	networkSrv := newNetworkService(newNetworkHandler(sidebarCtr, busCtr), mux, sidebarCtr, historyCtr)

	return &serviceManager{
		filePath: "",

		configSrv: newConfigService(),

		mux:     mux,
		network: nil,

		sidebarSrv: sidebarSrv,

		historySrv: historySrv,
		historyCtr: historyCtr,

		networkSrv: networkSrv,

		busSrv: busSrv,
		busCtr: busCtr,

		nodeSrv: nodeSrv,
		nodeCtr: nodeCtr,

		messageSrv: messageSrv,
		messageCtr: messageCtr,

		signalSrv: signalSrv,
		signalCtr: signalCtr,

		signalTypeSrv: signalTypeSrv,
		signalTypeCtr: signalTypeCtr,

		signalUnitSrv: signalUnitSrv,
		signalUnitCtr: signalUnitCtr,

		signalEnumSrv: signalEnumSrv,
		signalEnumCtr: signalEnumCtr,
	}
}

func (m *serviceManager) getServices() []application.Service {
	return []application.Service{
		application.NewService(m.configSrv),

		application.NewService(m.sidebarSrv),
		application.NewService(manager.historySrv),

		application.NewService(manager.networkSrv),
		application.NewService(manager.busSrv),
		application.NewService(manager.nodeSrv),
		application.NewService(manager.messageSrv),
		application.NewService(manager.signalSrv),
		application.NewService(manager.signalTypeSrv),
		application.NewService(manager.signalUnitSrv),
		application.NewService(manager.signalEnumSrv),
	}
}

func (m *serviceManager) createNetwork() {
	net := acmelib.NewNetwork("new_network")

	m.filePath = ""
	m.clearServices()
	m.initNetwork(net)
	m.historySrv.setSaved(false)
}

func (m *serviceManager) initNetwork(net *acmelib.Network) {
	m.network = net

	buses := net.Buses()
	nodes := make(map[acmelib.EntityID]*acmelib.Node)
	messages := []*acmelib.Message{}
	signals := []acmelib.Signal{}
	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	for _, bus := range buses {
		for _, nodeInt := range bus.NodeInterfaces() {
			tmpNode := nodeInt.Node()
			nodes[tmpNode.EntityID()] = tmpNode

			tmpMessages := nodeInt.SentMessages()
			messages = append(messages, tmpMessages...)

			for _, msg := range tmpMessages {
				tmpSignals := msg.Signals()
				signals = append(signals, tmpSignals...)

				for _, sig := range msg.Signals() {

					switch sig.Kind() {
					case acmelib.SignalKindStandard:
						stdSig, err := sig.ToStandard()
						if err != nil {
							panic(err)
						}
						sigTypes[stdSig.Type().EntityID()] = stdSig.Type()

						if stdSig.Unit() != nil {
							sigUnits[stdSig.Unit().EntityID()] = stdSig.Unit()
						}

					case acmelib.SignalKindEnum:
						enumSig, err := sig.ToEnum()
						if err != nil {
							panic(err)
						}
						sigEnums[enumSig.Enum().EntityID()] = enumSig.Enum()
					}
				}
			}
		}
	}

	m.networkSrv.load(net)
	m.busCtr.sendLoad(buses)
	m.nodeCtr.sendLoad(maps.Values(nodes))
	m.messageCtr.sendLoad(messages)
	m.signalCtr.sendLoad(signals)
	m.signalTypeCtr.sendLoad(maps.Values(sigTypes))
	m.signalUnitCtr.sendLoad(maps.Values(sigUnits))
	m.signalEnumCtr.sendLoad(maps.Values(sigEnums))
}

func (m *serviceManager) getEncoding(path string) acmelib.SaveEncoding {
	switch filepath.Ext(path) {
	case ".binpb":
		return acmelib.SaveEncodingWire
	case ".json":
		return acmelib.SaveEncodingJSON
	case ".txtpb":
		return acmelib.SaveEncodingText
	}
	return acmelib.SaveEncodingWire
}

func (m *serviceManager) openNetwork(path string) error {
	if path == "" {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	net, err := acmelib.LoadNetwork(file, m.getEncoding(path))
	if err != nil {
		return err
	}

	m.clearServices()
	m.initNetwork(net)
	m.historySrv.setSaved(true)

	m.filePath = path

	m.configSrv.addOpenedNetwork(net.Name(), m.filePath)

	return nil
}

func (m *serviceManager) saveNetwork() error {
	if m.filePath == "" {
		dialog := newSaveNetworkDialog()
		filename, err := dialog.PromptForSingleSelection()
		if err != nil {
			application.Get().Logger.Error(err.Error())
			return nil
		}

		return m.saveNetworkAs(filename)
	}

	file, err := os.Create(m.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	m.mux.Lock()
	defer m.mux.Unlock()

	fileEnc := m.getEncoding(m.filePath)
	switch fileEnc {
	case acmelib.SaveEncodingWire:
		err = acmelib.SaveNetwork(m.network, fileEnc, file, nil, nil)

	case acmelib.SaveEncodingJSON:
		err = acmelib.SaveNetwork(m.network, fileEnc, nil, file, nil)

	case acmelib.SaveEncodingText:
		err = acmelib.SaveNetwork(m.network, fileEnc, nil, nil, file)
	}

	if err != nil {
		return err
	}

	m.historySrv.setSaved(true)

	m.configSrv.addOpenedNetwork(m.network.Name(), m.filePath)

	printInfo("NETWORK SAVED")

	return nil
}

func (m *serviceManager) trySaveNetwork() error {
	if m.historySrv.isSaved() {
		return nil
	}

	return m.saveNetwork()
}

func (m *serviceManager) saveNetworkAs(filename string) error {
	m.filePath = filename
	return m.saveNetwork()
}

func (m *serviceManager) reloadNetwork() {
	m.clearServices()
	m.initNetwork(m.network)
}

func (m *serviceManager) importDBC(path string) error {
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
		printError(err)
		return err
	}

	m.mux.Lock()
	if err := m.network.AddBus(bus); err != nil {
		m.mux.Unlock()
		printError(err)
		return err
	}
	m.mux.Unlock()

	m.initNetwork(m.network)

	m.historyCtr.sendOperation(
		serviceKindNetwork,
		func() (any, error) {
			m.mux.Lock()
			defer m.mux.Unlock()

			if err := m.network.RemoveBus(bus.EntityID()); err != nil {
				return nil, err
			}

			m.busCtr.sendDelete(bus)

			return newNetwork(m.network), nil
		},
		func() (any, error) {
			m.mux.Lock()
			defer m.mux.Unlock()

			if err := m.network.AddBus(bus); err != nil {
				return nil, err
			}

			m.busCtr.sendAdd(bus)

			return newNetwork(m.network), nil
		},
	)

	return nil
}

func (m *serviceManager) exportDBC(path string) error {
	if path == "" {
		return nil
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	return acmelib.ExportNetwork(manager.network, path)
}

func (m *serviceManager) clearServices() {
	m.sidebarSrv.clear()
	m.historySrv.clear()

	m.networkSrv.clear()
	m.busCtr.sendClear()
	m.nodeCtr.sendClear()
	m.messageCtr.sendClear()
	m.signalCtr.sendClear()
	m.signalTypeCtr.sendClear()
	m.signalUnitCtr.sendClear()
	m.signalEnumCtr.sendClear()
}
