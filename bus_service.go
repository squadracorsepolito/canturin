package main

import "github.com/squadracorsepolito/acmelib"

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

	proxy.pushSidebarUpdate(bus.EntityID(), name)

	proxy.pushHistoryOperation(
		operationDomainBus,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := bus.UpdateName(oldName); err != nil {
				return Bus{}, err
			}

			proxy.pushSidebarUpdate(bus.EntityID(), oldName)

			return s.converterFn(bus), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := bus.UpdateName(name); err != nil {
				return Bus{}, err
			}

			proxy.pushSidebarUpdate(bus.EntityID(), name)

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
