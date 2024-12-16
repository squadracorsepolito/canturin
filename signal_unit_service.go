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
	*service[*acmelib.SignalUnit, SignalUnit]
}

func signalUnitConverter(sigType *acmelib.SignalUnit) SignalUnit {
	return SignalUnit{
		base: getBase(sigType),

		Symbol: sigType.Symbol(),

		ReferenceCount: sigType.ReferenceCount(),
		References:     getSignalReferences(sigType),
	}
}

func newSignalUnitService() *SignalUnitService {
	return &SignalUnitService{
		service: newService(proxy.sigUnitCh, signalUnitConverter),
	}
}

func (s *SignalUnitService) GetInvalidNames(entityID string) []string {
	s.mux.RLock()
	defer s.mux.RUnlock()

	names := []string{}
	for _, tmpSigUnit := range s.pool {
		if tmpSigUnit.EntityID() == acmelib.EntityID(entityID) {
			continue
		}

		names = append(names, tmpSigUnit.Name())
	}

	return names
}

func (s *SignalUnitService) UpdateName(entityID string, name string) (SignalUnit, error) {
	// retrieve the signal unit
	sigUnit, err := s.getEntity(entityID)
	if err != nil {
		return SignalUnit{}, err
	}

	// Lock to guarantee unique access
	s.mux.Lock()
	defer s.mux.Unlock()

	// Actual signal unit name
	oldName := sigUnit.Name()

	if name == oldName {
		return s.converterFn(sigUnit), nil
	}

	// Update name
	sigUnit.SetName(name)

	// Push the new name to the sideBar
	proxy.pushSidebarUpdate(sigUnit.EntityID(), name)

	// Add history operation
	proxy.pushHistoryOperation(
		operationDomainSignalUnit,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			// Rollback of the name
			sigUnit.SetName(oldName)
			proxy.pushSidebarUpdate(sigUnit.EntityID(), oldName)

			return s.converterFn(sigUnit), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			// Final update
			sigUnit.SetName(name)
			proxy.pushSidebarUpdate(sigUnit.EntityID(), name)

			return s.converterFn(sigUnit), nil
		},
	)

	// return the updated entity
	return s.converterFn(sigUnit), nil
}

func (s *SignalUnitService) UpdateDesc(entityID string, desc string) (SignalUnit, error) {
	sigUnit, err := s.getEntity(entityID)
	if err != nil {
		return SignalUnit{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldDesc := sigUnit.Desc()
	if desc == oldDesc {
		return s.converterFn(sigUnit), nil
	}

	sigUnit.SetDesc(desc)

	proxy.pushHistoryOperation(
		operationDomainSignalUnit,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigUnit.SetDesc(oldDesc)

			return s.converterFn(sigUnit), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigUnit.SetDesc(desc)

			return s.converterFn(sigUnit), nil
		},
	)

	return s.converterFn(sigUnit), nil
}

func (s *SignalUnitService) UpdateSymbol(entityID string, symbol string) (SignalUnit, error) {
	sigUnit, err := s.getEntity(entityID)
	if err != nil {
		return SignalUnit{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldSymbol := sigUnit.Symbol()
	if symbol == oldSymbol {
		return s.converterFn(sigUnit), nil
	}

	sigUnit.SetSymbol(symbol)

	proxy.pushHistoryOperation(
		operationDomainSignalUnit,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigUnit.SetSymbol(oldSymbol)
			return s.converterFn(sigUnit), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigUnit.SetSymbol(symbol)
			return s.converterFn(sigUnit), nil
		},
	)

	return s.converterFn(sigUnit), nil
}
