package main

import (
	"errors"
	"sync"

	"github.com/squadracorsepolito/acmelib"
)

type Network struct {
	BaseEntity

	Buses []BusBase `json:"buses"`
}

func newNetwork(net *acmelib.Network) Network {
	if net == nil {
		return Network{}
	}

	res := Network{
		BaseEntity: newBaseEntity(net),

		Buses: []BusBase{},
	}

	for _, bus := range net.Buses() {
		res.Buses = append(res.Buses, newBusBase(bus))
	}

	return res
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

func (s *NetworkService) Get() Network {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.handler.toResponse(s.network)
}

func (s *NetworkService) UpdateName(req UpdateNameReq) (Network, error) {
	return s.handle(&req, s.handler.updateName)
}

func (s *NetworkService) UpdateDesc(req UpdateDescReq) (Network, error) {
	return s.handle(&req, s.handler.updateDesc)
}

func (s *NetworkService) AddBus() (Network, error) {
	return s.handle(nil, s.handler.addBus)
}

func (s *NetworkService) DeleteBuses(req DeleteBusesReq) (Network, error) {
	return s.handle(&req, s.handler.deleteBuses)
}

type networkHandler struct {
	sidebarCtr *sidebarController
	busCtr     *busController
}

func newNetworkHandler(sidebarCtr *sidebarController, busCtr *busController) *networkHandler {
	return &networkHandler{
		sidebarCtr: sidebarCtr,
		busCtr:     busCtr,
	}
}

func (h *networkHandler) toResponse(net *acmelib.Network) Network {
	return newNetwork(net)
}

func (h *networkHandler) updateName(net *acmelib.Network, req *request, res *networkRes) error {
	parsedReq := req.toUpdateName()

	name := parsedReq.Name
	oldName := net.Name()
	if oldName == name {
		return nil
	}

	net.UpdateName(name)

	h.sidebarCtr.sendUpdateName(net)

	res.setUndo(
		func() error {
			net.UpdateName(oldName)
			h.sidebarCtr.sendUpdateName(net)
			return nil
		},
	)

	res.setRedo(
		func() error {
			net.UpdateName(name)
			h.sidebarCtr.sendUpdateName(net)
			return nil
		},
	)

	return nil
}

func (h *networkHandler) updateDesc(net *acmelib.Network, req *request, res *networkRes) error {
	parsedReq := req.toUpdateDesc()

	desc := parsedReq.Desc
	oldDesc := net.Desc()
	if desc == oldDesc {
		return nil
	}

	net.SetDesc(desc)

	res.setUndo(
		func() error {
			net.SetDesc(oldDesc)
			return nil
		},
	)

	res.setRedo(
		func() error {
			net.SetDesc(desc)
			return nil
		},
	)

	return nil
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

func (h *networkHandler) deleteBuses(net *acmelib.Network, req *request, res *networkRes) error {
	parsedReq := req.toDeleteBuses()

	if len(parsedReq.BusEntityIDs) == 0 {
		return nil
	}

	remBusIDs := make(map[string]struct{})
	for _, busID := range parsedReq.BusEntityIDs {
		remBusIDs[busID] = struct{}{}
	}

	remBuses := []*acmelib.Bus{}
	for _, bus := range net.Buses() {
		if _, ok := remBusIDs[bus.EntityID().String()]; ok {
			remBuses = append(remBuses, bus)
		}
	}

	for _, bus := range remBuses {
		if err := net.RemoveBus(bus.EntityID()); err != nil {
			return err
		}

		h.sidebarCtr.sendDelete(bus)
	}

	res.setUndo(
		func() error {
			for _, bus := range remBuses {
				if err := net.AddBus(bus); err != nil {
					return err
				}

				h.sidebarCtr.sendAdd(bus)
			}

			return nil
		},
	)

	res.setRedo(
		func() error {
			for _, bus := range remBuses {
				if err := net.RemoveBus(bus.EntityID()); err != nil {
					return err
				}

				h.sidebarCtr.sendDelete(bus)
			}

			return nil
		},
	)

	return nil
}
