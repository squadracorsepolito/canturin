package main

import (
	"github.com/squadracorsepolito/acmelib"
)

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

	proxy.pushSidebarAdd(SidebarNodeKindSignalType, sigType.EntityID(), proxy.network.EntityID(), sigType.Name())

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

	s.mux.Lock()
	defer s.mux.Unlock()

	// verify that the new name is not equal to the old one
	oldName := sigType.Name()
	if name == oldName {
		return s.converterFn(sigType), nil
	}

	sigType.SetName(name)

	// push the new name in the sidebar
	proxy.pushSidebarUpdate(sigType.EntityID(), name)

	// add to the history
	proxy.pushHistoryOperation(
		operationDomainSignalType,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetName(oldName)
			proxy.pushSidebarUpdate(sigType.EntityID(), oldName)

			return s.converterFn(sigType), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetName(name)
			proxy.pushSidebarUpdate(sigType.EntityID(), name)

			return s.converterFn(sigType), nil
		},
	)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateDesc(entityID string, desc string) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldDesc := sigType.Desc()
	if desc == oldDesc {
		return s.converterFn(sigType), nil
	}

	sigType.SetDesc(desc)

	proxy.pushHistoryOperation(
		operationDomainSignalType,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetDesc(oldDesc)

			return s.converterFn(sigType), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetDesc(desc)

			return s.converterFn(sigType), nil
		},
	)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateMin(entityID string, min float64) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldMin := sigType.Min()
	if min == oldMin {
		return s.converterFn(sigType), nil
	}

	sigType.SetMin(min)

	proxy.pushHistoryOperation(
		operationDomainSignalType,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetMin(oldMin)

			return s.converterFn(sigType), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetMin(min)

			return s.converterFn(sigType), nil
		},
	)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateMax(entityID string, max float64) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldMax := sigType.Max()
	if max == oldMax {
		return s.converterFn(sigType), nil
	}

	sigType.SetMax(max)

	proxy.pushHistoryOperation(
		operationDomainSignalType,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetMax(oldMax)

			return s.converterFn(sigType), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetMax(max)

			return s.converterFn(sigType), nil
		},
	)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateScale(entityID string, scale float64) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldScale := sigType.Scale()
	if scale == oldScale {
		return s.converterFn(sigType), nil
	}

	sigType.SetScale(scale)

	proxy.pushHistoryOperation(
		operationDomainSignalType,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetScale(oldScale)

			return s.converterFn(sigType), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetScale(scale)

			return s.converterFn(sigType), nil
		},
	)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) UpdateOffset(entityID string, offset float64) (SignalType, error) {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return SignalType{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldOffset := sigType.Offset()
	if offset == oldOffset {
		return s.converterFn(sigType), nil
	}

	sigType.SetOffset(offset)

	proxy.pushHistoryOperation(
		operationDomainSignalType,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetOffset(oldOffset)

			return s.converterFn(sigType), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetOffset(offset)

			return s.converterFn(sigType), nil
		},
	)

	return s.converterFn(sigType), nil
}
