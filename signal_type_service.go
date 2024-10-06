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

func (s *SignalTypeService) GetInvalidNames(entityID string) []string {
	s.mux.RLock()
	defer s.mux.RUnlock()

	names := []string{}
	for _, tmpSigType := range s.pool {
		if tmpSigType.EntityID() == acmelib.EntityID(entityID) {
			continue
		}

		names = append(names, tmpSigType.Name())
	}

	return names
}

func (s *SignalTypeService) UpdateName(entityID string, name string) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	sigType.SetName(name)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateDesc(entityID string, desc string) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	sigType.SetDesc(desc)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateKind(entityID string, kind acmelib.SignalTypeKind) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	// TODO: updateKind in acmelib

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateSize(entityID string, size int) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	// TODO: updateSize in acmelib

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateMin(entityID string, min float64) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	sigType.SetMin(min)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateMax(entityID string, max float64) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	sigType.SetMax(max)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateScale(entityID string, scale float64) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	sigType.SetScale(scale)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateOffset(entityID string, offset float64) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	sigType.SetOffset(offset)

	return s.converterFn(sigType), nil
}
