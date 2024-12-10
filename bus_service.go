package main

import "github.com/squadracorsepolito/acmelib"

// type BusBase struct {
// 	base
// }

// func getBusBase(bus *acmelib.Bus) BusBase {
// 	if bus == nil {
// 		return BusBase{}
// 	}

// 	return BusBase{
// 		base: getBase(bus),
// 	}
// }

type BusType string

const (
	BusTypeCAN2A BusType = "CAN_2.0A"
)

func getBusType(typ acmelib.BusType) BusType {
	return BusType(typ.String())
}

type Bus struct {
	base

	Type     BusType `json:"type"`
	Baudrate int     `json:"baudrate"`

	AttachedInterfaces []NodeInterface `json:"attachedInterfaces"`
}

// func getBus(bus *acmelib.Bus) Bus {
// 	nodeInts := []NodeInterface{}

// 	for _, nodeInt := range bus.NodeInterfaces() {
// 		nodeInts = append(nodeInts, getNodeInterface(nodeInt))
// 	}

// 	return Bus{
// 		base: getBase(bus),

// 		NodeInterfaces: nodeInts,
// 	}
// }

type BusService struct {
	*service0[*acmelib.Bus, Bus, *busHandlers]
}

func newBusService() *BusService {
	handler := &busHandlers{}

	return &BusService{
		service0: newService0(serviceKindBus, handler),
	}
}

func (s *BusService) sendSidebarUpdateName(bus *acmelib.Bus) {
	s.service0.sendSidebarUpdateName(bus)

	msgBusKey := manager.sidebar.getMessageBusGroupKey(bus)
	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(msgBusKey, bus.Name()))
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

	s.sendSidebarUpdateName(bus)

	s.sendHistoryOp(
		func() (*acmelib.Bus, error) {
			s.hanlders.updateName(bus, oldName)
			s.sendSidebarUpdateName(bus)

			return bus, nil
		},
		func() (*acmelib.Bus, error) {
			s.hanlders.updateName(bus, name)
			s.sendSidebarUpdateName(bus)

			return bus, nil
		},
	)

	return s.hanlders.toResponse(bus), nil
}

func (s *BusService) UpdateDesc(entityID string, req UpdateDescReq) (Bus, error) {
	return s.process(entityID, newRequest(&req), s.hanlders.updateDesc)
}

func (s *BusService) UpdateBusType(entityID string, req UpdateBusTypeReq) (Bus, error) {
	return s.process(entityID, newRequest(&req), s.hanlders.updateType)
}

func (s *BusService) UpdateBaudrate(entityID string, req UpdateBaudrateReq) (Bus, error) {
	return s.process(entityID, newRequest(&req), s.hanlders.updateBaudrate)
}

type busHandlers struct{}

func (h *busHandlers) toResponse(bus *acmelib.Bus) Bus {
	nodeInts := []NodeInterface{}

	for _, nodeInt := range bus.NodeInterfaces() {
		nodeInts = append(nodeInts, getNodeInterface(nodeInt))
	}

	return Bus{
		base: getBase(bus),

		Type:     getBusType(bus.Type()),
		Baudrate: bus.Baudrate(),

		AttachedInterfaces: nodeInts,
	}
}

type busRes = *response[*acmelib.Bus]

func (h *busHandlers) updateName(bus *acmelib.Bus, name string) error {
	return bus.UpdateName(name)
}

func (h *busHandlers) updateDesc(bus *acmelib.Bus, req *request) (busRes, error) {
	parsedReq := req.toUpdateDesc()

	desc := parsedReq.Desc

	oldDesc := bus.Desc()
	if desc == oldDesc {
		return newUnchangedResponse[*acmelib.Bus](), nil
	}

	bus.SetDesc(desc)

	return newResponse(
		func() (*acmelib.Bus, error) {
			bus.SetDesc(oldDesc)
			return bus, nil
		},
		func() (*acmelib.Bus, error) {
			bus.SetDesc(desc)
			return bus, nil
		},
	), nil
}

func (h *busHandlers) updateType(bus *acmelib.Bus, req *request) (busRes, error) {
	parsedReq := req.toUpdateBusType()

	typ := parsedReq.BusType

	var busType acmelib.BusType
	switch typ {
	case BusTypeCAN2A:
		busType = acmelib.BusTypeCAN2A
	}

	oldBusType := bus.Type()
	if oldBusType == busType {
		return newUnchangedResponse[*acmelib.Bus](), nil
	}

	bus.SetType(busType)

	return newResponse(
		func() (*acmelib.Bus, error) {
			bus.SetType(oldBusType)
			return bus, nil
		},
		func() (*acmelib.Bus, error) {
			bus.SetType(busType)
			return bus, nil
		},
	), nil
}

func (h *busHandlers) updateBaudrate(bus *acmelib.Bus, req *request) (busRes, error) {
	parsedReq := req.toUpdateBaudrate()

	baudrate := parsedReq.Baudrate

	oldBaudrate := bus.Baudrate()
	if oldBaudrate == baudrate {
		return newUnchangedResponse[*acmelib.Bus](), nil
	}

	bus.SetBaudrate(baudrate)

	return newResponse(
		func() (*acmelib.Bus, error) {
			bus.SetBaudrate(oldBaudrate)
			return bus, nil
		},
		func() (*acmelib.Bus, error) {
			bus.SetBaudrate(baudrate)
			return bus, nil
		},
	), nil
}

// type BusService struct {
// 	*service[*acmelib.Bus, Bus]
// }

// func newBusService() *BusService {
// 	return &BusService{
// 		service: newService(getBus),
// 	}
// }

// func (s *BusService) ListBase() []BusBase {
// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	briefs := []BusBase{}
// 	for _, bus := range s.pool {
// 		briefs = append(briefs, getBusBase(bus))
// 	}

// 	return briefs
// }

// func (s *BusService) sendSidebarUpdateName(bus *acmelib.Bus) {
// 	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(bus.EntityID().String(), bus.Name()))

// 	msgBusGroupKey := manager.sidebar.getMessageBusGroupKey(bus)
// 	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(msgBusGroupKey, bus.Name()))
// }

// func (s *BusService) UpdateName(entityID string, name string) (Bus, error) {
// 	bus, err := s.getEntity(entityID)
// 	if err != nil {
// 		return Bus{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	oldName := bus.Name()
// 	if name == oldName {
// 		return s.converterFn(bus), nil
// 	}

// 	if err := bus.UpdateName(name); err != nil {
// 		return Bus{}, err
// 	}

// 	s.sendSidebarUpdateName(bus)

// 	proxy.pushHistoryOperation(
// 		operationDomainBus,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := bus.UpdateName(oldName); err != nil {
// 				return Bus{}, err
// 			}

// 			s.sendSidebarUpdateName(bus)

// 			return s.converterFn(bus), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := bus.UpdateName(name); err != nil {
// 				return Bus{}, err
// 			}

// 			s.sendSidebarUpdateName(bus)

// 			return s.converterFn(bus), nil
// 		},
// 	)

// 	return s.converterFn(bus), nil
// }

// func (s *BusService) UpdateDesc(entityID string, desc string) (Bus, error) {
// 	bus, err := s.getEntity(entityID)
// 	if err != nil {
// 		return Bus{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	oldDesc := bus.Desc()
// 	if desc == oldDesc {
// 		return s.converterFn(bus), nil
// 	}

// 	bus.SetDesc(desc)

// 	proxy.pushHistoryOperation(
// 		operationDomainBus,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			bus.SetDesc(oldDesc)

// 			return s.converterFn(bus), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			bus.SetDesc(desc)

// 			return s.converterFn(bus), nil
// 		},
// 	)

// 	return s.converterFn(bus), nil
// }
