package main

import "github.com/squadracorsepolito/acmelib"

type BusType string

const (
	BusTypeCAN2A BusType = "CAN_2.0A"
)

func (bt BusType) parse() acmelib.BusType {
	switch bt {
	case BusTypeCAN2A:
		return acmelib.BusTypeCAN2A
	default:
		return acmelib.BusTypeCAN2A
	}
}

func getBusType(typ acmelib.BusType) BusType {
	return BusType(typ.String())
}

type AttachedNode struct {
	base

	ID              uint `json:"id"`
	InterfaceNumber int  `json:"interfaceNumber"`
}

func getAttachedNode(nodeInt *acmelib.NodeInterface) AttachedNode {
	node := nodeInt.Node()

	return AttachedNode{
		base: getBase(node),

		ID:              uint(node.ID()),
		InterfaceNumber: nodeInt.Number(),
	}
}

type Bus struct {
	base

	Type     BusType `json:"type"`
	Baudrate int     `json:"baudrate"`

	AttachedNodes []AttachedNode `json:"attachedNodes"`
}

type BusService struct {
	*service0[*acmelib.Bus, Bus, *busHandlers]
}

func newBusService(sidebar *sidebarController) *BusService {
	handler := &busHandlers{}

	return &BusService{
		service0: newService0(serviceKindBus, handler, sidebar),
	}
}

// func (s *BusService) sendSidebarUpdateName(bus *acmelib.Bus) {
// 	s.service0.sendSidebarUpdateName(bus)

// 	msgBusKey := manager.sidebar.getMessageBusGroupKey(bus)
// 	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(msgBusKey, bus.Name()))
// }

func (s *BusService) Create(req CreateBusReq) (Bus, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	bus := acmelib.NewBus(req.Name)
	bus.SetDesc(req.Desc)
	bus.SetType(req.BusType.parse())
	bus.SetBaudrate(req.Baudrate)

	s.addEntity(bus)

	s.sidebar.sendAdd(bus)

	// manager.sidebar.sendAdd(newSidebarAddReq(newBusSidebarItem(bus), SidebarBusesPrefix))
	// manager.sidebar.sendAdd(newSidebarAddReq(newMessageBusGroupSidebarItem(bus), SidebarMessagesPrefix))

	s.sendHistoryOp(
		func() (*acmelib.Bus, error) {
			s.removeEntity(bus.EntityID().String())

			// manager.sidebar.sendDelete(newSidebarDeleteReq(bus.EntityID().String()))
			// manager.sidebar.sendDelete(newSidebarDeleteReq(manager.sidebar.getMessageBusGroupKey(bus)))

			s.sidebar.sendDelete(bus)

			return bus, nil
		},
		func() (*acmelib.Bus, error) {
			s.addEntity(bus)

			// manager.sidebar.sendAdd(newSidebarAddReq(newBusSidebarItem(bus), SidebarBusesPrefix))
			// manager.sidebar.sendAdd(newSidebarAddReq(newMessageBusGroupSidebarItem(bus), SidebarMessagesPrefix))

			s.sidebar.sendAdd(bus)

			return bus, nil
		},
	)

	return s.hanlders.toResponse(bus), nil
}

func (s *BusService) UpdateName(entityID string, name string) (Bus, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	bus, err := s.getEntity(entityID)
	if err != nil {
		return Bus{}, err
	}

	oldName := bus.Name()
	if name == oldName {
		return s.hanlders.toResponse(bus), nil
	}

	if err := s.hanlders.updateName(bus, name); err != nil {
		return Bus{}, err
	}

	// s.sendSidebarUpdateName(bus)
	s.sidebar.sendUpdateName(bus)

	s.sendHistoryOp(
		func() (*acmelib.Bus, error) {
			s.hanlders.updateName(bus, oldName)
			// s.sendSidebarUpdateName(bus)
			s.sidebar.sendUpdateName(bus)

			return bus, nil
		},
		func() (*acmelib.Bus, error) {
			s.hanlders.updateName(bus, name)
			// s.sendSidebarUpdateName(bus)
			s.sidebar.sendUpdateName(bus)

			return bus, nil
		},
	)

	return s.hanlders.toResponse(bus), nil
}

func (s *BusService) UpdateDesc(entityID string, req UpdateDescReq) (Bus, error) {
	return s.handle(entityID, &req, s.hanlders.updateDesc)
}

func (s *BusService) UpdateBusType(entityID string, req UpdateBusTypeReq) (Bus, error) {
	return s.handle(entityID, &req, s.hanlders.updateBusType)
}

func (s *BusService) UpdateBaudrate(entityID string, req UpdateBaudrateReq) (Bus, error) {
	return s.handle(entityID, &req, s.hanlders.updateBaudrate)
}

type busHandlers struct{}

func (h *busHandlers) toResponse(bus *acmelib.Bus) Bus {
	attNodes := []AttachedNode{}

	for _, nodeInt := range bus.NodeInterfaces() {
		attNodes = append(attNodes, getAttachedNode(nodeInt))
	}

	return Bus{
		base: getBase(bus),

		Type:     getBusType(bus.Type()),
		Baudrate: bus.Baudrate(),

		AttachedNodes: attNodes,
	}
}

type busRes = response[*acmelib.Bus]

func (h *busHandlers) updateName(bus *acmelib.Bus, name string) error {
	return bus.UpdateName(name)
}

func (h *busHandlers) updateDesc(bus *acmelib.Bus, req *request, res *busRes) error {
	parsedReq := req.toUpdateDesc()

	desc := parsedReq.Desc

	oldDesc := bus.Desc()
	if desc == oldDesc {
		return nil
	}

	bus.SetDesc(desc)

	res.setUndo(
		func() (*acmelib.Bus, error) {
			bus.SetDesc(oldDesc)
			return bus, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Bus, error) {
			bus.SetDesc(desc)
			return bus, nil
		},
	)

	return nil
}

func (h *busHandlers) updateBusType(bus *acmelib.Bus, req *request, res *busRes) error {
	parsedReq := req.toUpdateBusType()

	typ := parsedReq.BusType

	var busType acmelib.BusType
	switch typ {
	case BusTypeCAN2A:
		busType = acmelib.BusTypeCAN2A
	}

	oldBusType := bus.Type()
	if oldBusType == busType {
		return nil
	}

	bus.SetType(busType)

	res.setUndo(
		func() (*acmelib.Bus, error) {
			bus.SetType(oldBusType)
			return bus, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Bus, error) {
			bus.SetType(busType)
			return bus, nil
		},
	)

	return nil
}

func (h *busHandlers) updateBaudrate(bus *acmelib.Bus, req *request, res *busRes) error {
	parsedReq := req.toUpdateBaudrate()

	baudrate := parsedReq.Baudrate

	oldBaudrate := bus.Baudrate()
	if oldBaudrate == baudrate {
		return nil
	}

	bus.SetBaudrate(baudrate)

	res.setUndo(
		func() (*acmelib.Bus, error) {
			bus.SetBaudrate(oldBaudrate)
			return bus, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Bus, error) {
			bus.SetBaudrate(baudrate)
			return bus, nil
		},
	)

	return nil
}
