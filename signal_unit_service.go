package main

import (
	"github.com/squadracorsepolito/acmelib"
)

type SignalUnit struct {
	base

	Symbol     string            `json:"symbol"`
	References []SignalReference `json:"references"`
}

type SignalReference struct {
	entityStub

	ParentMessage entityStub `json:"parentMessage"`
}

type SignalUnitService struct {
	*service[*acmelib.SignalUnit, SignalUnit]
}

func newSignalUnitService(sigUnitCh chan *acmelib.SignalUnit) *SignalUnitService {
	return &SignalUnitService{
		service: newService(sigUnitCh, func(su *acmelib.SignalUnit) SignalUnit {
			res := SignalUnit{
				base: getBase(su),

				Symbol:     su.Symbol(),
				References: []SignalReference{},
			}

			for _, tmpStdSig := range su.References() {
				res.References = append(res.References, SignalReference{
					entityStub: getEntityStub(tmpStdSig),

					ParentMessage: getEntityStub(tmpStdSig.ParentMessage()),
				})
			}

			return res
		}),
	}
}
