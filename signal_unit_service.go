package main

import (
	"fmt"

	"github.com/squadracorsepolito/acmelib"
)

type SignalUnitKind string

const (
	SignalUnitKindCustom      SignalUnitKind = "custom"
	SignalUnitKindTemperature SignalUnitKind = "temperature"
	SignalUnitKindElectrical  SignalUnitKind = "electrical"
	SignalUnitKindPower       SignalUnitKind = "power"
)

func newSignalUnitKind(kind acmelib.SignalUnitKind) SignalUnitKind {
	switch kind {
	case acmelib.SignalUnitKindCustom:
		return SignalUnitKindCustom
	case acmelib.SignalUnitKindTemperature:
		return SignalUnitKindTemperature
	case acmelib.SignalUnitKindElectrical:
		return SignalUnitKindElectrical
	case acmelib.SignalUnitKindPower:
		return SignalUnitKindPower
	default:
		return SignalUnitKindCustom
	}
}

func (k SignalUnitKind) parse() acmelib.SignalUnitKind {
	switch k {
	case SignalUnitKindCustom:
		return acmelib.SignalUnitKindCustom
	case SignalUnitKindTemperature:
		return acmelib.SignalUnitKindTemperature
	case SignalUnitKindElectrical:
		return acmelib.SignalUnitKindElectrical
	case SignalUnitKindPower:
		return acmelib.SignalUnitKindPower
	default:
		return acmelib.SignalUnitKindCustom
	}
}

type SignalUnit struct {
	base

	Kind   SignalUnitKind `json:"kind"`
	Symbol string         `json:"symbol"`

	ReferenceCount int               `json:"referenceCount"`
	References     []SignalReference `json:"references"`
}

type SignalUnitService struct {
	*service[*acmelib.SignalUnit, SignalUnit, *signalUnitHandler]
}

func newSignalUnitService(sidebar *sidebarController) *SignalUnitService {
	return &SignalUnitService{
		service: newService(serviceKindSignalUnit, newSignalUnitHandler(sidebar), sidebar),
	}
}

func (s *SignalUnitService) Create(req CreateSignalUnitReq) (SignalUnit, error) {
	sigUnit := acmelib.NewSignalUnit(req.Name, req.Kind.parse(), req.Symbol)
	sigUnit.SetDesc(req.Desc)

	s.mux.Lock()
	defer s.mux.Unlock()

	s.addEntity(sigUnit)
	s.sidebar.sendAdd(sigUnit)

	s.sendHistoryOp(
		func() (*acmelib.SignalUnit, error) {
			s.removeEntity(sigUnit.EntityID().String())
			s.sidebar.sendDelete(sigUnit)
			return sigUnit, nil
		},
		func() (*acmelib.SignalUnit, error) {
			s.addEntity(sigUnit)
			s.sidebar.sendAdd(sigUnit)
			return sigUnit, nil
		},
	)

	return s.handler.toResponse(sigUnit), nil
}

func (s *SignalUnitService) Delete(entityID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	sigUnit, err := s.getEntity(entityID)
	if err != nil {
		return err
	}

	if sigUnit.ReferenceCount() > 0 {
		return fmt.Errorf("signal unit %s is referenced %d times", sigUnit.Name(), sigUnit.ReferenceCount())
	}

	s.removeEntity(entityID)
	s.sidebar.sendDelete(sigUnit)

	s.sendHistoryOp(
		func() (*acmelib.SignalUnit, error) {
			s.addEntity(sigUnit)
			s.sidebar.sendAdd(sigUnit)
			return sigUnit, nil
		},
		func() (*acmelib.SignalUnit, error) {
			s.removeEntity(sigUnit.EntityID().String())
			s.sidebar.sendDelete(sigUnit)
			return sigUnit, nil
		},
	)

	return nil
}

func (s *SignalUnitService) UpdateName(entityID string, req UpdateNameReq) (SignalUnit, error) {
	return s.handle(entityID, &req, s.handler.updateName)
}

func (s *SignalUnitService) UpdateDesc(entityID string, req UpdateDescReq) (SignalUnit, error) {
	return s.handle(entityID, &req, s.handler.updateDesc)
}

func (s *SignalUnitService) UpdateSignalUnitKind(entityID string, req UpdateSignalUnitKindReq) (SignalUnit, error) {
	return s.handle(entityID, &req, s.handler.updateSignalUnitKind)
}

func (s *SignalUnitService) UpdateSymbol(entityID string, req UpdateSymbolReq) (SignalUnit, error) {
	return s.handle(entityID, &req, s.handler.updateSymbol)
}

type signalUnitRes = response[*acmelib.SignalUnit]

type signalUnitHandler struct {
	*commonServiceHandler
}

func newSignalUnitHandler(sidebar *sidebarController) *signalUnitHandler {
	return &signalUnitHandler{
		commonServiceHandler: newCommonServiceHandler(sidebar),
	}
}

func (h *signalUnitHandler) toResponse(sigUnit *acmelib.SignalUnit) SignalUnit {
	return SignalUnit{
		base: getBase(sigUnit),

		Kind:   newSignalUnitKind(sigUnit.Kind()),
		Symbol: sigUnit.Symbol(),

		ReferenceCount: sigUnit.ReferenceCount(),
		References:     getSignalReferences(sigUnit),
	}
}

func (h *signalUnitHandler) updateName(sigUnit *acmelib.SignalUnit, req *request, res *signalUnitRes) error {
	parsedReq := req.toUpdateName()

	name := parsedReq.Name

	oldName := sigUnit.Name()
	if oldName == name {
		return nil
	}

	sigUnit.SetName(name)
	h.sidebar.sendUpdateName(sigUnit)

	res.setUndo(
		func() (*acmelib.SignalUnit, error) {
			sigUnit.SetName(oldName)
			h.sidebar.sendUpdateName(sigUnit)
			return sigUnit, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalUnit, error) {
			sigUnit.SetName(name)
			h.sidebar.sendUpdateName(sigUnit)
			return sigUnit, nil
		},
	)

	return nil
}

func (h *signalUnitHandler) updateDesc(sigUnit *acmelib.SignalUnit, req *request, res *signalUnitRes) error {
	parsedReq := req.toUpdateDesc()

	desc := parsedReq.Desc

	oldDesc := sigUnit.Desc()
	if oldDesc == desc {
		return nil
	}

	sigUnit.SetDesc(desc)

	res.setUndo(
		func() (*acmelib.SignalUnit, error) {
			sigUnit.SetDesc(oldDesc)
			return sigUnit, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalUnit, error) {
			sigUnit.SetDesc(desc)
			return sigUnit, nil
		},
	)

	return nil
}

func (h *signalUnitHandler) updateSignalUnitKind(sigUnit *acmelib.SignalUnit, req *request, res *signalUnitRes) error {
	parsedReq := req.toUpdateSignalUnitKind()

	kind := parsedReq.Kind.parse()

	oldKind := sigUnit.Kind()
	if oldKind == kind {
		return nil
	}

	sigUnit.SetKind(kind)

	res.setUndo(
		func() (*acmelib.SignalUnit, error) {
			sigUnit.SetKind(oldKind)
			return sigUnit, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalUnit, error) {
			sigUnit.SetKind(kind)
			return sigUnit, nil
		},
	)

	return nil
}

func (h *signalUnitHandler) updateSymbol(sigUnit *acmelib.SignalUnit, req *request, res *signalUnitRes) error {
	parsedReq := req.toUpdateSymbol()

	symbol := parsedReq.Symbol

	oldSymbol := sigUnit.Symbol()
	if oldSymbol == symbol {
		return nil
	}

	sigUnit.SetSymbol(symbol)

	res.setUndo(
		func() (*acmelib.SignalUnit, error) {
			sigUnit.SetSymbol(oldSymbol)
			return sigUnit, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalUnit, error) {
			sigUnit.SetSymbol(symbol)
			return sigUnit, nil
		},
	)

	return nil
}
