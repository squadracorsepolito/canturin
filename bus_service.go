package main

import (
	"github.com/squadracorsepolito/acmelib"
)

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
	*service0[*acmelib.Bus, Bus, *busHandler]
}

func newBusService(sidebar *sidebarController) *BusService {
	return &BusService{
		service0: newService0(serviceKindBus, newBusHandler(sidebar), sidebar),
	}
}

func (s *BusService) Create(req CreateBusReq) (Bus, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	bus := acmelib.NewBus(req.Name)
	bus.SetDesc(req.Desc)
	bus.SetType(req.BusType.parse())
	bus.SetBaudrate(req.Baudrate)

	s.addEntity(bus)
	s.sidebar.sendAdd(bus)

	s.sendHistoryOp(
		func() (*acmelib.Bus, error) {
			s.removeEntity(bus.EntityID().String())
			s.sidebar.sendDelete(bus)

			return bus, nil
		},
		func() (*acmelib.Bus, error) {
			s.addEntity(bus)
			s.sidebar.sendAdd(bus)

			return bus, nil
		},
	)

	return s.handler.toResponse(bus), nil
}

func (s *BusService) Delete(entityID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	bus, err := s.getEntity(entityID)
	if err != nil {
		return err
	}

	nodeInts := bus.NodeInterfaces()
	bus.RemoveAllNodeInterfaces()

	s.removeEntity(entityID)
	s.sidebar.sendDelete(bus)

	s.sendHistoryOp(
		func() (*acmelib.Bus, error) {
			for _, nodeInt := range nodeInts {
				if err := bus.AddNodeInterface(nodeInt); err != nil {
					return nil, err
				}
			}

			s.addEntity(bus)
			s.sidebar.sendAdd(bus)

			return bus, nil
		},
		func() (*acmelib.Bus, error) {
			bus.RemoveAllNodeInterfaces()

			s.removeEntity(entityID)
			s.sidebar.sendDelete(bus)

			return bus, nil
		},
	)

	return nil
}

func (s *BusService) UpdateName(entityID string, req UpdateNameReq) (Bus, error) {
	return s.handle(entityID, &req, s.handler.updateName)
}

func (s *BusService) UpdateDesc(entityID string, req UpdateDescReq) (Bus, error) {
	return s.handle(entityID, &req, s.handler.updateDesc)
}

func (s *BusService) UpdateBusType(entityID string, req UpdateBusTypeReq) (Bus, error) {
	return s.handle(entityID, &req, s.handler.updateBusType)
}

func (s *BusService) UpdateBaudrate(entityID string, req UpdateBaudrateReq) (Bus, error) {
	return s.handle(entityID, &req, s.handler.updateBaudrate)
}

type busRes = response[*acmelib.Bus]

type busHandler struct {
	*commonServiceHandler
}

func newBusHandler(sidebar *sidebarController) *busHandler {
	return &busHandler{
		commonServiceHandler: newCommonServiceHandler(sidebar),
	}
}

func (h *busHandler) toResponse(bus *acmelib.Bus) Bus {
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

func (h *busHandler) updateName(bus *acmelib.Bus, req *request, res *busRes) error {
	parsedReq := req.toUpdateName()

	name := parsedReq.Name

	oldName := bus.Name()
	if name == oldName {
		return nil
	}

	if err := bus.UpdateName(name); err != nil {
		return err
	}

	h.sidebar.sendUpdateName(bus)

	res.setUndo(
		func() (*acmelib.Bus, error) {
			if err := bus.UpdateName(oldName); err != nil {
				return nil, err
			}

			h.sidebar.sendUpdateName(bus)

			return bus, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Bus, error) {
			if err := bus.UpdateName(name); err != nil {
				return nil, err
			}

			h.sidebar.sendUpdateName(bus)

			return bus, nil
		},
	)

	return nil
}

func (h *busHandler) updateDesc(bus *acmelib.Bus, req *request, res *busRes) error {
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

func (h *busHandler) updateBusType(bus *acmelib.Bus, req *request, res *busRes) error {
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

func (h *busHandler) updateBaudrate(bus *acmelib.Bus, req *request, res *busRes) error {
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
