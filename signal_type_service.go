package main

import "github.com/squadracorsepolito/acmelib"

type SignalType struct {
	base

	Kind   acmelib.SignalTypeKind `json:"kind"`
	Size   int                    `json:"size"`
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
				Size:   st.Size(),
				Min:    st.Min(),
				Max:    st.Max(),
				Scale:  st.Scale(),
				Offset: st.Offset(),

				ReferenceCount: st.ReferenceCount(),
				References:     getSignalReferences(st),
			}
		}),
	}
}

func (s *SignalTypeService) UpdateName(entityID string, newName string) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	sigType.SetName(newName)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateDesc(entityID string, newDesc string) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	sigType.SetDesc(newDesc)

	return s.converterFn(sigType), nil
}
