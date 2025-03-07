package main

import (
	"strings"
	"sync"

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

func newBusType(typ acmelib.BusType) BusType {
	return BusType(typ.String())
}

type AttachedNode struct {
	BaseEntity

	ID              uint `json:"id"`
	InterfaceNumber int  `json:"interfaceNumber"`
}

func newAttachedNode(nodeInt *acmelib.NodeInterface) AttachedNode {
	node := nodeInt.Node()

	return AttachedNode{
		BaseEntity: newBaseEntity(node),

		ID:              uint(node.ID()),
		InterfaceNumber: nodeInt.Number(),
	}
}

type BusBase struct {
	BaseEntity

	Baudrate int `json:"baudrate"`
}

func newBusBase(bus *acmelib.Bus) BusBase {
	if bus == nil {
		return BusBase{}
	}

	return BusBase{
		BaseEntity: newBaseEntity(bus),
		Baudrate:   bus.Baudrate(),
	}
}

type Bus struct {
	BaseEntity

	Type     BusType `json:"type"`
	Baudrate int     `json:"baudrate"`

	AttachedNodes []AttachedNode `json:"attachedNodes"`
}

func newBus(bus *acmelib.Bus) Bus {
	attNodes := []AttachedNode{}

	for _, nodeInt := range bus.NodeInterfaces() {
		attNodes = append(attNodes, newAttachedNode(nodeInt))
	}

	return Bus{
		BaseEntity: newBaseEntity(bus),

		Type:     newBusType(bus.Type()),
		Baudrate: bus.Baudrate(),

		AttachedNodes: attNodes,
	}
}

type BusLoadMessage struct {
	BaseEntity

	BitsPerSec float64 `json:"bitsPerSec"`
	Percentage float64 `json:"percentage"`
}

func newBusLoadMessage(msg *acmelib.MessageLoad) BusLoadMessage {
	return BusLoadMessage{
		BaseEntity: newBaseEntity(msg.Message),
		BitsPerSec: msg.BitsPerSec,
		Percentage: msg.Percentage,
	}
}

type BusLoad struct {
	Percentage float64          `json:"percentage"`
	Messages   []BusLoadMessage `json:"messages"`
}

func newBusLoad(load float64, msgLoads []*acmelib.MessageLoad) BusLoad {
	messages := []BusLoadMessage{}
	for _, tmpMsgLoad := range msgLoads {
		messages = append(messages, newBusLoadMessage(tmpMsgLoad))
	}

	return BusLoad{
		Percentage: load,
		Messages:   messages,
	}
}

type BusService struct {
	*service[*acmelib.Bus, Bus, *busHandler]
}

func newBusService(mux *sync.RWMutex, sidebar *sidebarController) *BusService {
	return &BusService{
		service: newService(serviceKindBus, newBusHandler(sidebar), mux, sidebar),
	}
}

func (s *BusService) GetLoad(entityID string) (BusLoad, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	bus, err := s.getEntity(entityID)
	if err != nil {
		return BusLoad{}, err
	}

	load, msgLoads, err := acmelib.CalculateBusLoad(bus, 1000)
	if err != nil {
		return BusLoad{}, err
	}

	return newBusLoad(load, msgLoads), nil
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
	return newBus(bus)
}

func (h *busHandler) updateName(bus *acmelib.Bus, req *request, res *busRes) error {
	parsedReq := req.toUpdateName()

	name := strings.TrimSpace(parsedReq.Name)

	oldName := bus.Name()
	if name == oldName {
		return nil
	}

	if err := bus.UpdateName(name); err != nil {
		return err
	}

	h.sidebarCtr.sendUpdateName(bus)

	res.setUndo(
		func() (*acmelib.Bus, error) {
			if err := bus.UpdateName(oldName); err != nil {
				return nil, err
			}

			h.sidebarCtr.sendUpdateName(bus)

			return bus, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Bus, error) {
			if err := bus.UpdateName(name); err != nil {
				return nil, err
			}

			h.sidebarCtr.sendUpdateName(bus)

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

func (h *busHandler) addNodeInterface(bus *acmelib.Bus, req *request, res *busRes) error {
	return nil
}
