package main

import "github.com/squadracorsepolito/acmelib"

type SignalTypeKind string

const (
	SignalTypeKindCustom  SignalTypeKind = "custom"
	SignalTypeKindFlag    SignalTypeKind = "flag"
	SignalTypeKindInteger SignalTypeKind = "integer"
	SignalTypeKindDecimal SignalTypeKind = "decimal"
)

func newSignalTypeKind(kind acmelib.SignalTypeKind) SignalTypeKind {
	return SignalTypeKind(kind.String())
}

type SignalType struct {
	base

	Kind   SignalTypeKind `json:"kind"`
	Size   int            `json:"size"`
	Signed bool           `json:"signed"`
	Min    float64        `json:"min"`
	Max    float64        `json:"max"`
	Scale  float64        `json:"scale"`
	Offset float64        `json:"offset"`

	ReferenceCount int               `json:"referenceCount"`
	References     []SignalReference `json:"references"`
}

type SignalTypeService struct {
	*service[*acmelib.SignalType, SignalType]
}

func signalTypeConverter(sigType *acmelib.SignalType) SignalType {
	return SignalType{
		base: getBase(sigType),

		Kind:   newSignalTypeKind(sigType.Kind()),
		Size:   int(sigType.Size()),
		Signed: sigType.Signed(),
		Min:    sigType.Min(),
		Max:    sigType.Max(),
		Scale:  sigType.Scale(),
		Offset: sigType.Offset(),

		ReferenceCount: sigType.ReferenceCount(),
		References:     getSignalReferences(sigType),
	}
}

func newSignalTypeService() *SignalTypeService {
	return &SignalTypeService{
		service: newService(proxy.sigTypeCh, signalTypeConverter),
	}
}

func (s *SignalTypeService) Create(kind SignalTypeKind, name, desc string, size int, signed bool, min, max, scale, offset float64) (SignalType, error) {
	sigType := &acmelib.SignalType{}
	switch kind {
	case SignalTypeKindCustom:
		tmpSigType, err := acmelib.NewCustomSignalType(name, size, signed, min, max, scale, offset)
		if err != nil {
			return SignalType{}, err
		}
		sigType = tmpSigType

	case SignalTypeKindFlag:
		sigType = acmelib.NewFlagSignalType(name)

	case SignalTypeKindInteger:
		tmpSigType, err := acmelib.NewIntegerSignalType(name, size, signed)
		if err != nil {
			return SignalType{}, err
		}
		sigType = tmpSigType

	case SignalTypeKindDecimal:
		tmpSigType, err := acmelib.NewDecimalSignalType(name, size, signed)
		if err != nil {
			return SignalType{}, err
		}
		sigType = tmpSigType
	}

	if len(desc) > 0 {
		sigType.SetDesc(desc)
	}

	sigType.SetMin(min)
	sigType.SetMax(max)
	sigType.SetScale(scale)
	sigType.SetOffset(offset)

	s.mux.Lock()
	defer s.mux.Unlock()
	s.pool[sigType.EntityID()] = sigType

	return s.converterFn(sigType), nil
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

	proxy.pushSidebarUpdate(sigType.EntityID(), name)

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
