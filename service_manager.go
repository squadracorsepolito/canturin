package main

import "github.com/wailsapp/wails/v3/pkg/application"

type serviceManager struct {
	sidebar *SidebarService
	history *HistoryService

	bus        *BusService
	node       *NodeService
	message    *MessageService
	signalType *SignalTypeService
	signalUnit *SignalUnitService
	signalEnum *SignalEnumService
}

func newServiceManager() *serviceManager {
	sidebar := newSidebarService()
	sidebarController := sidebar.getController()

	bus := newBusService(sidebarController)

	return &serviceManager{
		sidebar: sidebar,
		history: newHistoryService(),

		bus:        bus,
		node:       newNodeService(sidebarController, bus),
		message:    newMessageService(),
		signalType: newSignalTypeService(sidebarController),
		signalUnit: newSignalUnitService(sidebarController),
		signalEnum: newSignalEnumService(sidebarController),
	}
}

func (m *serviceManager) getServices() []application.Service {
	return []application.Service{
		application.NewService(m.sidebar),
		application.NewService(manager.history),

		application.NewService(manager.bus),
		application.NewService(manager.node),
		application.NewService(manager.message),
		application.NewService(manager.signalType),
		application.NewService(manager.signalUnit),
		application.NewService(manager.signalEnum),
	}
}
