package main

import (
	"github.com/squadracorsepolito/acmelib"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type serviceManager struct {
	network *acmelib.Network

	sidebar           *SidebarService
	sidebarController *sidebarController
	history           *HistoryService

	bus        *BusService
	node       *NodeService
	message    *MessageService
	signal     *SignalService
	signalType *SignalTypeService
	signalUnit *SignalUnitService
	signalEnum *SignalEnumService
}

func newServiceManager() *serviceManager {
	sidebar := newSidebarService()
	sidebarController := sidebar.getController()

	bus := newBusService(sidebarController)

	signalType := newSignalTypeService(sidebarController)
	signalUnit := newSignalUnitService(sidebarController)
	signalEnum := newSignalEnumService(sidebarController)

	return &serviceManager{
		network: acmelib.NewNetwork("Unnamed Network"),

		sidebar:           sidebar,
		sidebarController: sidebarController,
		history:           newHistoryService(),

		bus:     bus,
		node:    newNodeService(sidebarController, bus),
		message: newMessageService(sidebarController),

		signal:     newSignalService(sidebarController, signalType, signalUnit, signalEnum),
		signalType: signalType,
		signalUnit: signalUnit,
		signalEnum: signalEnum,
	}
}

func (m *serviceManager) getServices() []application.Service {
	return []application.Service{
		application.NewService(m.sidebar),
		application.NewService(manager.history),

		application.NewService(manager.bus),
		application.NewService(manager.node),
		application.NewService(manager.message),
		application.NewService(manager.signal),
		application.NewService(manager.signalType),
		application.NewService(manager.signalUnit),
		application.NewService(manager.signalEnum),
	}
}

func (m *serviceManager) initNetwork(net *acmelib.Network) {
	m.network = net

	m.sidebar.sendLoad(newSidebarLoadReq(net))

	nodes := make(map[acmelib.EntityID]*acmelib.Node)

	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	// Iterate over the buses, node interfaces and messages
	for _, bus := range net.Buses() {
		manager.bus.sendLoad(bus)

		for _, nodeInt := range bus.NodeInterfaces() {
			tmpNode := nodeInt.Node()
			nodes[tmpNode.EntityID()] = tmpNode

			for _, msg := range nodeInt.SentMessages() {
				manager.message.sendLoad(msg)

				// Iterate over the signals
				for _, sig := range msg.Signals() {
					manager.signal.sendLoad(sig)

					switch sig.Kind() {
					// Standard Signal
					case acmelib.SignalKindStandard:
						stdSig, err := sig.ToStandard()
						if err != nil {
							panic(err)
						}
						sigTypes[stdSig.Type().EntityID()] = stdSig.Type()

						if stdSig.Unit() != nil {
							sigUnits[stdSig.Unit().EntityID()] = stdSig.Unit()
						}

					// Enum Signal
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

	for _, node := range nodes {
		m.node.sendLoad(node)
	}

	for _, sigType := range sigTypes {
		m.signalType.sendLoad(sigType)
	}

	for _, sigUnit := range sigUnits {
		m.signalUnit.sendLoad(sigUnit)
	}

	for _, sigEnum := range sigEnums {
		m.signalEnum.sendLoad(sigEnum)
	}
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

	m.bus.sendClear()
	m.node.sendClear()
	m.message.sendClear()
	m.signal.sendClear()
	m.signalType.sendClear()
	m.signalUnit.sendClear()
	m.signalEnum.sendClear()
}
