package main

import (
	"fmt"

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
		service: newService(signalTypeConverter),
	}
}

func (s *SignalTypeService) sendSidebarAdd(sigType *acmelib.SignalType) {
	item := newSignalTypeSidebarItem(sigType)
	manager.sidebar.sendAdd(newSidebarAddReq(item, SidebarSignalTypesPrefix))
}

func (s *SignalTypeService) sendSidebarUpdateName(sigType *acmelib.SignalType) {
	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(sigType.EntityID().String(), sigType.Name()))
}

func (s *SignalTypeService) sendSidebarDelete(sigType *acmelib.SignalType) {
	manager.sidebar.sendDelete(newSidebarDeleteReq(sigType.EntityID().String()))
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

	s.addEntity(sigType)
	s.sendSidebarAdd(sigType)

	proxy.pushHistoryOperation(
		operationDomainSignalType,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			s.deleteEntity(sigType.EntityID().String())
			s.sendSidebarDelete(sigType)

			return s.converterFn(sigType), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			s.addEntity(sigType)
			s.sendSidebarAdd(sigType)

			return s.converterFn(sigType), nil
		},
	)

	return s.converterFn(sigType), nil
}

func (s *SignalTypeService) Delete(entityID string) error {
	sigType, err := s.getEntity(entityID)
	if err != nil {
		return err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	if sigType.ReferenceCount() > 0 {
		return fmt.Errorf("signal type %s is referenced %d times", sigType.Name(), sigType.ReferenceCount())
	}

	s.deleteEntity(entityID)
	s.sendSidebarDelete(sigType)

	proxy.pushHistoryOperation(
		operationDomainSignalType,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			s.addEntity(sigType)
			s.sendSidebarAdd(sigType)

			return s.converterFn(sigType), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			s.deleteEntity(sigType.EntityID().String())
			s.sendSidebarDelete(sigType)

			return s.converterFn(sigType), nil
		},
	)

	return nil
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
	s.sendSidebarUpdateName(sigType)

	// add to the history
	proxy.pushHistoryOperation(
		operationDomainSignalType,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetName(oldName)
			s.sendSidebarUpdateName(sigType)

			return s.converterFn(sigType), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigType.SetName(name)
			s.sendSidebarUpdateName(sigType)

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
