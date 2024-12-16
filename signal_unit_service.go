package main

import (
	"github.com/squadracorsepolito/acmelib"
)

type SignalUnit struct {
	base

	Symbol string `json:"symbol"`

	ReferenceCount int               `json:"referenceCount"`
	References     []SignalReference `json:"references"`
}

type SignalUnitService struct {
	*service0[*acmelib.SignalUnit, SignalUnit, *signalUnitHandler]
}

func newSignalUnitService(sidebar *sidebarController) *SignalUnitService {
	return &SignalUnitService{
		service0: newService0(serviceKindSignalUnit, newSignalUnitHandler(sidebar), sidebar),
	}
}

func (s *SignalUnitService) UpdateName(entityID string, req UpdateNameReq) (SignalUnit, error) {
	return s.handle(entityID, &req, s.handler.updateName)
}

func (s *SignalUnitService) UpdateDesc(entityID string, req UpdateDescReq) (SignalUnit, error) {
	return s.handle(entityID, &req, s.handler.updateDesc)
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

//
//
// //

// type SignalUnitService0 struct {
// 	*service[*acmelib.SignalUnit, SignalUnit]
// }

// func signalUnitConverter(sigType *acmelib.SignalUnit) SignalUnit {
// 	return SignalUnit{
// 		base: getBase(sigType),

// 		Symbol: sigType.Symbol(),

// 		ReferenceCount: sigType.ReferenceCount(),
// 		References:     getSignalReferences(sigType),
// 	}
// }

// func newSignalUnitService() *SignalUnitService0 {
// 	return &SignalUnitService0{
// 		service: newService(proxy.sigUnitCh, signalUnitConverter),
// 	}
// }

// func (s *SignalUnitService0) GetInvalidNames(entityID string) []string {
// 	s.mux.RLock()
// 	defer s.mux.RUnlock()

// 	names := []string{}
// 	for _, tmpSigUnit := range s.pool {
// 		if tmpSigUnit.EntityID() == acmelib.EntityID(entityID) {
// 			continue
// 		}

// 		names = append(names, tmpSigUnit.Name())
// 	}

// 	return names
// }

// func (s *SignalUnitService0) UpdateName(entityID string, name string) (SignalUnit, error) {
// 	// retrieve the signal unit
// 	sigUnit, err := s.getEntity(entityID)
// 	if err != nil {
// 		return SignalUnit{}, err
// 	}

// 	// Lock to guarantee unique access
// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	// Actual signal unit name
// 	oldName := sigUnit.Name()

// 	if name == oldName {
// 		return s.converterFn(sigUnit), nil
// 	}

// 	// Update name
// 	sigUnit.SetName(name)

// 	// Push the new name to the sideBar
// 	proxy.pushSidebarUpdate(sigUnit.EntityID(), name)

// 	// Add history operation
// 	proxy.pushHistoryOperation(
// 		operationDomainSignalUnit,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			// Rollback of the name
// 			sigUnit.SetName(oldName)
// 			proxy.pushSidebarUpdate(sigUnit.EntityID(), oldName)

// 			return s.converterFn(sigUnit), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			// Final update
// 			sigUnit.SetName(name)
// 			proxy.pushSidebarUpdate(sigUnit.EntityID(), name)

// 			return s.converterFn(sigUnit), nil
// 		},
// 	)

// 	// return the updated entity
// 	return s.converterFn(sigUnit), nil
// }

// func (s *SignalUnitService0) UpdateDesc(entityID string, desc string) (SignalUnit, error) {
// 	sigUnit, err := s.getEntity(entityID)
// 	if err != nil {
// 		return SignalUnit{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	oldDesc := sigUnit.Desc()
// 	if desc == oldDesc {
// 		return s.converterFn(sigUnit), nil
// 	}

// 	sigUnit.SetDesc(desc)

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalUnit,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			sigUnit.SetDesc(oldDesc)

// 			return s.converterFn(sigUnit), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			sigUnit.SetDesc(desc)

// 			return s.converterFn(sigUnit), nil
// 		},
// 	)

// 	return s.converterFn(sigUnit), nil
// }

// func (s *SignalUnitService0) UpdateSymbol(entityID string, symbol string) (SignalUnit, error) {
// 	sigUnit, err := s.getEntity(entityID)
// 	if err != nil {
// 		return SignalUnit{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	oldSymbol := sigUnit.Symbol()
// 	if symbol == oldSymbol {
// 		return s.converterFn(sigUnit), nil
// 	}

// 	sigUnit.SetSymbol(symbol)

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalUnit,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			sigUnit.SetSymbol(oldSymbol)
// 			return s.converterFn(sigUnit), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			sigUnit.SetSymbol(symbol)
// 			return s.converterFn(sigUnit), nil
// 		},
// 	)

// 	return s.converterFn(sigUnit), nil
// }
