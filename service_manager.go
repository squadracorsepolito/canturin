package main

import (
	"github.com/squadracorsepolito/acmelib"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/exp/maps"
)

type serviceManager struct {
	network *acmelib.Network

	sidebar           *SidebarService
	sidebarController *sidebarController
	history           *HistoryService

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
	sidebar := newSidebarService()
	sidebarCtr := sidebar.getController()

	signalTypeSrv := newSignalTypeService(sidebarCtr)
	signalTypeCtr := signalTypeSrv.getController()

	signalUnitSrv := newSignalUnitService(sidebarCtr)
	signalUnitCtr := signalUnitSrv.getController()

	signalEnumSrv := newSignalEnumService(sidebarCtr)
	signalEnumCtr := signalEnumSrv.getController()

	signalSrv := newSignalService(sidebarCtr, signalTypeCtr, signalUnitCtr, signalEnumCtr)
	signalCtr := signalSrv.getController()

	messageSrv := newMessageService(sidebarCtr, signalCtr)
	messageCtr := messageSrv.getController()

	busSrv := newBusService(sidebarCtr)
	busCtr := busSrv.getController()

	nodeSrv := newNodeService(sidebarCtr, busSrv, messageCtr)
	nodeCtr := nodeSrv.getController()

	return &serviceManager{
		network: acmelib.NewNetwork("Unnamed Network"),

		sidebar:           sidebar,
		sidebarController: sidebarCtr,
		history:           newHistoryService(),

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
		application.NewService(m.sidebar),
		application.NewService(manager.history),

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

	m.sidebar.sendLoad(newSidebarLoadReq(net))

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
	m.sidebar.clear()

	m.busCtr.sendClear()
	m.nodeCtr.sendClear()
	m.messageCtr.sendClear()
	m.signalCtr.sendClear()
	m.signalTypeCtr.sendClear()
	m.signalUnitCtr.sendClear()
	m.signalEnumCtr.sendClear()
}
