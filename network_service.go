package main

import (
	"errors"
	"sync"

	"github.com/squadracorsepolito/acmelib"
)

type Network struct {
	BaseEntity
}

func newNetwork(net *acmelib.Network) Network {
	if net == nil {
		return Network{}
	}

	return Network{
		BaseEntity: newBaseEntity(net),
	}
}

type NetworkService struct {
	handler *networkHandler

	mux     *sync.RWMutex
	network *acmelib.Network

	sidebarCtr *sidebarController
	historyCtr *historyController
}

func newNetworkService(handler *networkHandler, mux *sync.RWMutex, sidebarCtr *sidebarController, historyCtr *historyController) *NetworkService {
	return &NetworkService{
		handler: handler,

		mux:     mux,
		network: nil,

		sidebarCtr: sidebarCtr,
		historyCtr: historyCtr,
	}
}

func (s *NetworkService) load(net *acmelib.Network) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.network = net
	s.sidebarCtr.sendLoad(net)
}

func (s *NetworkService) clear() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.network = nil
}

func (s *NetworkService) handle(reqDataPtr any, handlerFn func(*acmelib.Network, *request, *networkRes) error) (dummyRes Network, _ error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.network == nil {
		return dummyRes, errors.New("network not loaded")
	}

	req := newRequest(reqDataPtr)
	res := newNetworkResponse()

	if err := handlerFn(s.network, req, res); err != nil {
		return dummyRes, err
	}

	s.historyCtr.sendOperation(
		serviceKindNetwork,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := res.undo(); err != nil {
				return nil, err
			}

			return s.handler.toResponse(s.network), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := res.redo(); err != nil {
				return nil, err
			}

			return s.handler.toResponse(s.network), nil
		},
	)

	return s.handler.toResponse(s.network), nil
}

func (s *NetworkService) AddBus() (Network, error) {
	return s.handle(nil, s.handler.addBus)
}

func (s *NetworkService) DeleteBus(req DeleteBusReq) (Network, error) {
	return s.handle(&req, s.handler.deleteBus)
}

type networkHandler struct {
	busCtr *busController
}

func newNetworkHandler(busCtr *busController) *networkHandler {
	return &networkHandler{
		busCtr: busCtr,
	}
}

func (h *networkHandler) toResponse(net *acmelib.Network) Network {
	return newNetwork(net)
}

func (h *networkHandler) addBus(net *acmelib.Network, _ *request, res *networkRes) error {
	takenName := make(map[string]struct{})
	for _, bus := range net.Buses() {
		takenName[bus.Name()] = struct{}{}
	}

	bus := acmelib.NewBus(getNewName("bus", takenName))

	if err := net.AddBus(bus); err != nil {
		return err
	}

	h.busCtr.sendAdd(bus)

	res.setUndo(
		func() error {
			if err := net.RemoveBus(bus.EntityID()); err != nil {
				return err
			}

			h.busCtr.sendDelete(bus)

			return nil
		},
	)

	res.setRedo(
		func() error {
			if err := net.AddBus(bus); err != nil {
				return err
			}

			h.busCtr.sendAdd(bus)

			return nil
		},
	)

	return nil
}

func (h *networkHandler) deleteBus(net *acmelib.Network, req *request, res *networkRes) error {
	parsedReq := req.toDeleteBus()

	busEntID := acmelib.EntityID(parsedReq.BusEntityID)

	var bus *acmelib.Bus
	found := false
	for _, tmpBus := range net.Buses() {
		if tmpBus.EntityID() == busEntID {
			bus = tmpBus
			found = true
			break
		}
	}

	if !found {
		return errors.New("bus not found")
	}

	if err := net.RemoveBus(busEntID); err != nil {
		return err
	}

	h.busCtr.sendDelete(bus)

	res.setUndo(
		func() error {
			if err := net.AddBus(bus); err != nil {
				return err
			}

			h.busCtr.sendAdd(bus)

			return nil
		},
	)

	res.setRedo(
		func() error {
			if err := net.RemoveBus(busEntID); err != nil {
				return err
			}

			h.busCtr.sendDelete(bus)

			return nil
		},
	)

	return nil
}
