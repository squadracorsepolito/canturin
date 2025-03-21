package main

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"sync"

	"github.com/squadracorsepolito/acmelib"
)

type SignalUnitService struct {
	*service[*acmelib.SignalUnit, SignalUnit, *signalUnitHandler]
}

func newSignalUnitService(mux *sync.RWMutex, sidebar *sidebarController) *SignalUnitService {
	return &SignalUnitService{
		service: newService(serviceKindSignalUnit, newSignalUnitHandler(sidebar), mux, sidebar),
	}
}

func (s *SignalUnitService) Create() (SignalUnit, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	takenNames := make(map[string]struct{})
	for _, sigUnit := range s.entities {
		takenNames[sigUnit.Name()] = struct{}{}
	}

	sigUnit := acmelib.NewSignalUnit(getNewName("signal_unit", takenNames), acmelib.SignalUnitKindCustom, "")

	s.addEntity(sigUnit)
	s.sidebarCtr.sendAdd(sigUnit)

	s.sendHistoryOp(
		func() (*acmelib.SignalUnit, error) {
			s.removeEntity(sigUnit.EntityID().String())
			s.sidebarCtr.sendDelete(sigUnit)
			return sigUnit, nil
		},
		func() (*acmelib.SignalUnit, error) {
			s.addEntity(sigUnit)
			s.sidebarCtr.sendAdd(sigUnit)
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
	s.sidebarCtr.sendDelete(sigUnit)

	s.sendHistoryOp(
		func() (*acmelib.SignalUnit, error) {
			s.addEntity(sigUnit)
			s.sidebarCtr.sendAdd(sigUnit)
			return sigUnit, nil
		},
		func() (*acmelib.SignalUnit, error) {
			s.removeEntity(sigUnit.EntityID().String())
			s.sidebarCtr.sendDelete(sigUnit)
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
	h.sidebarCtr.sendUpdateName(sigUnit)

	res.setUndo(
		func() (*acmelib.SignalUnit, error) {
			sigUnit.SetName(oldName)
			h.sidebarCtr.sendUpdateName(sigUnit)
			return sigUnit, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalUnit, error) {
			sigUnit.SetName(name)
			h.sidebarCtr.sendUpdateName(sigUnit)
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
