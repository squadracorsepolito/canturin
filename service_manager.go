package main

import "github.com/wailsapp/wails/v3/pkg/application"

type serviceManager struct {
	sidebar        *SidebarService0
	sidebarService *SidebarService
	historyService *HistoryService

	busService        *BusService
	nodeService       *NodeService
	messageService    *MessageService
	signalTypeService *SignalTypeService
	signalUnitService *SignalUnitService
	signalEnumService *SignalEnumService
}

func newServiceManager() *serviceManager {
	return &serviceManager{
		sidebar:        newSidebarService0(),
		sidebarService: newSidebarService(),
		historyService: newHistoryService(),

		busService:        newBusService(),
		nodeService:       newNodeService(),
		messageService:    newMessageService(),
		signalTypeService: newSignalTypeService(),
		signalUnitService: newSignalUnitService(),
		signalEnumService: newSignalEnumService(),
	}
}

func (m *serviceManager) getServices() []application.Service {
	return []application.Service{
		application.NewService(m.sidebar),
		application.NewService(manager.sidebarService),
		application.NewService(manager.historyService),

		application.NewService(manager.busService),
		application.NewService(manager.nodeService),
		application.NewService(manager.messageService),
		application.NewService(manager.signalTypeService),
		application.NewService(manager.signalUnitService),
		application.NewService(manager.signalEnumService),
	}
}
