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

func newSignalUnitService(sigUnitCh chan *acmelib.SignalUnit) *SignalUnitService {
	return &SignalUnitService{
		service: newService(sigUnitCh, func(su *acmelib.SignalUnit) SignalUnit {
			return SignalUnit{
				base: getBase(su),

				Symbol: su.Symbol(),

				ReferenceCount: su.ReferenceCount(),
				References:     getSignalReferences(su),
			}
		}),
	}
}

func (s *SignalUnitService) UpdateName(entityID string, newName string) (SignalUnit, error) {
	sigUnit, err := s.getEntity(entityID)
	if err != nil {
		return SignalUnit{}, err
	}

	sigUnit.SetName(newName)

	return s.converterFn(sigUnit), nil
}

func (s *SignalUnitService) UpdateDesc(entityID string, newDesc string) (SignalUnit, error) {
	sigUnit, err := s.getEntity(entityID)
	if err != nil {
		return SignalUnit{}, err
	}

	sigUnit.SetDesc(newDesc)

	return s.converterFn(sigUnit), nil
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
