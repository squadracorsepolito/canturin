package main

import "github.com/squadracorsepolito/acmelib"

type BusBase struct {
	base
}

func getBusBase(bus *acmelib.Bus) BusBase {
	if bus == nil {
		return BusBase{}
	}

	return BusBase{
		base: getBase(bus),
	}
}

type Bus struct {
	base

	NodeInterfaces []NodeInterface `json:"nodeInterfaces"`
}

func getBus(bus *acmelib.Bus) Bus {
	nodeInts := []NodeInterface{}

	for _, nodeInt := range bus.NodeInterfaces() {
		nodeInts = append(nodeInts, getNodeInterface(nodeInt))
	}

	return Bus{
		base: getBase(bus),

		NodeInterfaces: nodeInts,
	}
}

type BusService struct {
	*service[*acmelib.Bus, Bus]
}

func newBusService() *BusService {
	return &BusService{
		service: newService(proxy.busCh, getBus),
	}
}

func (s *BusService) GetInvalidNames(entityID string) []string {
	s.mux.Lock()
	defer s.mux.Unlock()

	names := []string{}
	for _, tmpBus := range s.pool {
		if tmpBus.EntityID() != acmelib.EntityID(entityID) {
			names = append(names, tmpBus.Name())
		}
	}

	return names
}

func (s *BusService) ListBase() []BusBase {
	s.mux.Lock()
	defer s.mux.Unlock()

	briefs := []BusBase{}
	for _, bus := range s.pool {
		briefs = append(briefs, getBusBase(bus))
	}

	return briefs
}

func (s *BusService) sendSidebarUpdateName(bus *acmelib.Bus) {
	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(bus.EntityID().String(), bus.Name()))

	msgBusGroupKey := manager.sidebar.getMessageBusGroupKey(bus)
	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(msgBusGroupKey, bus.Name()))
}

func (s *BusService) UpdateName(entityID string, name string) (Bus, error) {
	bus, err := s.getEntity(entityID)
	if err != nil {
		return Bus{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldName := bus.Name()
	if name == oldName {
		return s.converterFn(bus), nil
	}

	if err := bus.UpdateName(name); err != nil {
		return Bus{}, err
	}

	s.sendSidebarUpdateName(bus)

	proxy.pushHistoryOperation(
		operationDomainBus,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := bus.UpdateName(oldName); err != nil {
				return Bus{}, err
			}

			s.sendSidebarUpdateName(bus)

			return s.converterFn(bus), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := bus.UpdateName(name); err != nil {
				return Bus{}, err
			}

			s.sendSidebarUpdateName(bus)

			return s.converterFn(bus), nil
		},
	)

	return s.converterFn(bus), nil
}

func (s *BusService) UpdateDesc(entityID string, desc string) (Bus, error) {
	bus, err := s.getEntity(entityID)
	if err != nil {
		return Bus{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldDesc := bus.Desc()
	if desc == oldDesc {
		return s.converterFn(bus), nil
	}

	bus.SetDesc(desc)

	proxy.pushHistoryOperation(
		operationDomainBus,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			bus.SetDesc(oldDesc)

			return s.converterFn(bus), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			bus.SetDesc(desc)

			return s.converterFn(bus), nil
		},
	)

	return s.converterFn(bus), nil
}
