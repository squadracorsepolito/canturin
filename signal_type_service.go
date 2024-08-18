package main

import "github.com/squadracorsepolito/acmelib"

type SignalType struct {
	base

	Kind   acmelib.SignalTypeKind `json:"kind"`
	Min    float64                `json:"min"`
	Max    float64                `json:"max"`
	Scale  float64                `json:"scale"`
	Offset float64                `json:"offset"`

	ReferenceCount int               `json:"referenceCount"`
	References     []SignalReference `json:"references"`
}

type SignalTypeService struct {
	*service[*acmelib.SignalType, SignalType]
}

func newSignalTypeService(sigTypeCh chan *acmelib.SignalType) *SignalTypeService {
	return &SignalTypeService{
		service: newService(sigTypeCh, func(st *acmelib.SignalType) SignalType {
			return SignalType{
				base: getBase(st),

				Kind:   st.Kind(),
				Min:    st.Min(),
				Max:    st.Max(),
				Scale:  st.Scale(),
				Offset: st.Scale(),

				ReferenceCount: st.ReferenceCount(),
				References:     getSignalReferences(st),
			}
		}),
	}
}
