package main

import "github.com/squadracorsepolito/acmelib"

type SignalTypeService struct {
	*service[*acmelib.SignalType, SignalType]
}

func newSignalTypeService(sigTypeCh chan *acmelib.SignalType) *SignalTypeService {
	return &SignalTypeService{
		service: newService(sigTypeCh, func(st *acmelib.SignalType) SignalType {
			return SignalType{
				base: getBase(st),
			}
		}),
	}
}
