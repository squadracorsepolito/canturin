package main

import (
	"sync"

	"github.com/squadracorsepolito/acmelib"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/exp/maps"
)

type serviceManager struct {
	network *acmelib.Network

	sidebarSrv *SidebarService
	historySrv *HistoryService

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

	return &serviceManager{
		network: acmelib.NewNetwork("Unnamed Network"),

		sidebarSrv: sidebarSrv,
		historySrv: historySrv,

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
		application.NewService(m.sidebarSrv),
		application.NewService(manager.historySrv),

		application.NewService(manager.busSrv),
		application.NewService(manager.nodeSrv),
		application.NewService(manager.messageSrv),
		application.NewService(manager.signalSrv),
		application.NewService(manager.signalTypeSrv),
		application.NewService(manager.signalUnitSrv),
		application.NewService(manager.signalEnumSrv),
	}
}

func (m *serviceManager) initNetwork(net *acmelib.Network) {
	m.network = net

	m.sidebarSrv.sendLoad(newSidebarLoadReq(net))

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

	m.busCtr.sendLoad(buses)
	m.nodeCtr.sendLoad(maps.Values(nodes))
	m.messageCtr.sendLoad(messages)
	m.signalCtr.sendLoad(signals)
	m.signalTypeCtr.sendLoad(maps.Values(sigTypes))
	m.signalUnitCtr.sendLoad(maps.Values(sigUnits))
	m.signalEnumCtr.sendLoad(maps.Values(sigEnums))
}

func (m *serviceManager) loadNetwork(net *acmelib.Network) {
	m.clearServices()
	m.initNetwork(net)
}

func (m *serviceManager) reloadNetwork() {
	m.initNetwork(m.network)
}

func (m *serviceManager) clearServices() {
	m.sidebarSrv.clear()

	m.busCtr.sendClear()
	m.nodeCtr.sendClear()
	m.messageCtr.sendClear()
	m.signalCtr.sendClear()
	m.signalTypeCtr.sendClear()
	m.signalUnitCtr.sendClear()
	m.signalEnumCtr.sendClear()
}
