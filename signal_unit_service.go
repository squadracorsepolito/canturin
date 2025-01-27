package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

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

type SignalUnitBrief struct {
	BaseEntity

	Kind SignalUnitKind `json:"kind"`
}

func newSignalUnitBrief(sigUnit *acmelib.SignalUnit) SignalUnitBrief {
	return SignalUnitBrief{
		BaseEntity: newBaseEntity(sigUnit),

		Kind: newSignalUnitKind(sigUnit.Kind()),
	}
}

type SignalUnit struct {
	base

	Kind   SignalUnitKind `json:"kind"`
	Symbol string         `json:"symbol"`

	ReferenceCount int         `json:"referenceCount"`
	References     []Reference `json:"references"`
}

func newSignalUnit(sigUnit *acmelib.SignalUnit) SignalUnit {
	refCount := sigUnit.ReferenceCount()

	res := SignalUnit{
		base: getBase(sigUnit),

		Kind:   newSignalUnitKind(sigUnit.Kind()),
		Symbol: sigUnit.Symbol(),

		ReferenceCount: refCount,
	}

	if refCount == 0 {
		return res
	}

	rootRefs := []*reference{}
	refs := make(map[acmelib.EntityID]*reference)
	for _, stdSig := range sigUnit.References() {
		sigRef := newReference(stdSig)
		refs[stdSig.EntityID()] = sigRef

		var msgRef *reference
		msg := stdSig.ParentMessage()
		msgRef, ok := refs[msg.EntityID()]
		if !ok {
			msgRef = newReference(msg)
			refs[msg.EntityID()] = msgRef
		}
		msgRef.addChild(sigRef)

		var nodeRef *reference
		node := msg.SenderNodeInterface().Node()
		nodeRef, ok = refs[node.EntityID()]
		if !ok {
			nodeRef = newReference(node)
			refs[node.EntityID()] = nodeRef
		}
		nodeRef.addChild(msgRef)

		var busRef *reference
		bus := msg.SenderNodeInterface().ParentBus()
		busRef, ok = refs[bus.EntityID()]
		if !ok {
			busRef = newReference(bus)
			refs[bus.EntityID()] = busRef
			rootRefs = append(rootRefs, busRef)
		}
		busRef.addChild(nodeRef)
	}

	for _, tmpRef := range rootRefs {
		res.References = append(res.References, tmpRef.toResponse())
	}

	return res
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

func (s *SignalUnitService) ListBrief() []SignalUnitBrief {
	s.mux.RLock()
	defer s.mux.RUnlock()

	res := []SignalUnitBrief{}
	for _, sigUnit := range s.entities {
		res = append(res, newSignalUnitBrief(sigUnit))
	}

	slices.SortFunc(res, func(a, b SignalUnitBrief) int {
		if a.Kind == b.Kind {
			return strings.Compare(a.Name, b.Name)
		}

		return int(a.Kind.parse()) - int(b.Kind.parse())
	})

	return res
}

func (s *SignalUnitService) UpdateName(entityID string, req UpdateNameReq) (SignalUnit, error) {
	return s.handle(entityID, &req, s.handler.updateName)
}

func (s *SignalUnitService) UpdateDesc(entityID string, req UpdateDescReq) (SignalUnit, error) {
	return s.handle(entityID, &req, s.handler.updateDesc)
}

func (s *SignalUnitService) UpdateKind(entityID string, req UpdateSignalUnitKindReq) (SignalUnit, error) {
	return s.handle(entityID, &req, s.handler.updateKind)
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
	return newSignalUnit(sigUnit)
}

func (h *signalUnitHandler) updateName(sigUnit *acmelib.SignalUnit, req *request, res *signalUnitRes) error {
	parsedReq := req.toUpdateName()

	name := strings.TrimSpace(parsedReq.Name)

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

func (h *signalUnitHandler) updateKind(sigUnit *acmelib.SignalUnit, req *request, res *signalUnitRes) error {
	parsedReq := req.toUpdateSignalUnitKind()

	kind := parsedReq.Kind.parse()

	log.Print(kind.String())

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
