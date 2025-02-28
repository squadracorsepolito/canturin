package main

import (
	"sync"

	"github.com/squadracorsepolito/acmelib"
)

type Network struct {
	BaseEntity
}

func newNetwork(net *acmelib.Network) Network {
	return Network{
		BaseEntity: newBaseEntity(net),
	}
}

type NetworkService struct {
	network *acmelib.Network
	mux     *sync.RWMutex
	factory *entityFactory

	sidebarCtr *sidebarController
	historyCtr *HistoryService
	busCtr     *busController
}

func newNetworkService(mux *sync.RWMutex, sidebarCtr *sidebarController) *NetworkService {
	return &NetworkService{
		network: nil,
		mux:     mux,
		factory: newEntityFactory(),

		sidebarCtr: sidebarCtr,
	}
}

func (s *NetworkService) AddBus() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.network == nil {
		return nil
	}

	bus, err := s.factory.createBus(s.network)
	if err != nil {
		return err
	}

	if err := s.network.AddBus(bus); err != nil {
		return err
	}

	s.busCtr.sendAdd(bus)

	return nil
}
