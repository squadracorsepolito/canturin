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
