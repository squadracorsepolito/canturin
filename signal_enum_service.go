package main

import (
	"fmt"

	"github.com/squadracorsepolito/acmelib"
)

type SignalEnum struct {
	base

	Size     int               `json:"size"`
	MinSize  int               `json:"minSize"`
	MaxIndex int               `json:"maxIndex"`
	Values   []SignalEnumValue `json:"values"`

	References []SignalReference `json:"references"`
}

func signalEnumConverter(sigEnum *acmelib.SignalEnum) SignalEnum {
	values := []SignalEnumValue{}
	for _, val := range sigEnum.Values() {
		values = append(values, signalEnumValueConverter(val))
	}

	references := []SignalReference{}
	for _, ref := range sigEnum.References() {
		parMsg := ref.ParentMessage()
		parNode := parMsg.SenderNodeInterface().Node()
		parBus := parMsg.SenderNodeInterface().ParentBus()

		references = append(references, SignalReference{
			Bus:     getEntityStub(parBus),
			Node:    getEntityStub(parNode),
			Message: getEntityStub(parMsg),
			Signal:  getEntityStub(ref),
		})
	}

	return SignalEnum{
		base: getBase(sigEnum),

		Size:     sigEnum.GetSize(),
		MinSize:  sigEnum.MinSize(),
		MaxIndex: sigEnum.MaxIndex(),
		Values:   values,

		References: references,
	}
}

type SignalEnumValue struct {
	base

	Index int `json:"index"`
}

func signalEnumValueConverter(sigEnumValue *acmelib.SignalEnumValue) SignalEnumValue {
	return SignalEnumValue{
		base: getBase(sigEnumValue),

		Index: sigEnumValue.Index(),
	}
}

type SignalEnumService struct {
	*service[*acmelib.SignalEnum, SignalEnum]
}

func newSignalEnumService() *SignalEnumService {
	return &SignalEnumService{
		service: newService(proxy.signalEnumCh, signalEnumConverter),
	}
}

func (s *SignalEnumService) Create(name, desc string, minSize int) (SignalEnum, error) {
	sigEnum := acmelib.NewSignalEnum(name)

	if minSize > 0 {
		sigEnum.SetMinSize(minSize)
	}

	s.mux.Lock()
	defer s.mux.Unlock()
	s.pool[sigEnum.EntityID()] = sigEnum

	proxy.pushSidebarAdd(SidebarNodeKindSignalEnum, sigEnum.EntityID(), proxy.network.EntityID(), name)

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) GetInvalidNames(entityID string) []string {
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

func (s *SignalEnumService) UpdateName(entityID string, name string) (SignalEnum, error) {
	sigEnum, err := s.getEntity(entityID)
	if err != nil {
		return SignalEnum{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldName := sigEnum.Name()
	if name == oldName {
		return s.converterFn(sigEnum), nil
	}

	sigEnum.UpdateName(name)

	proxy.pushSidebarUpdate(sigEnum.EntityID(), name)

	proxy.pushHistoryOperation(
		operationDomainSignalEnum,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigEnum.UpdateName(oldName)
			proxy.pushSidebarUpdate(sigEnum.EntityID(), oldName)

			return s.converterFn(sigEnum), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigEnum.UpdateName(name)
			proxy.pushSidebarUpdate(sigEnum.EntityID(), name)

			return s.converterFn(sigEnum), nil
		},
	)

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) UpdateDesc(entityID string, desc string) (SignalEnum, error) {
	sigEnum, err := s.getEntity(entityID)
	if err != nil {
		return SignalEnum{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldDesc := sigEnum.Desc()
	if desc == oldDesc {
		return s.converterFn(sigEnum), nil
	}

	sigEnum.SetDesc(desc)

	proxy.pushHistoryOperation(
		operationDomainSignalEnum,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigEnum.SetDesc(oldDesc)

			return s.converterFn(sigEnum), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigEnum.SetDesc(desc)

			return s.converterFn(sigEnum), nil
		},
	)

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) reorderValue(sigEnum *acmelib.SignalEnum, sigEnumVal *acmelib.SignalEnumValue, from, to int) error {
	if from == to {
		return nil
	}

	if err := sigEnum.RemoveValue(sigEnumVal.EntityID()); err != nil {
		return err
	}
	lastValIdx := sigEnumVal.Index()

	restValues := sigEnum.Values()
	if to < from {
		// move up
		for i := from - 1; i >= to; i-- {
			tmpVal := restValues[i]
			tmpValIdx := tmpVal.Index()
			if err := tmpVal.UpdateIndex(lastValIdx); err != nil {
				return err
			}
			lastValIdx = tmpValIdx
		}
	} else {
		// move down
		for i := from; i < to; i++ {
			tmpVal := restValues[i]
			tmpValIdx := tmpVal.Index()
			if err := tmpVal.UpdateIndex(lastValIdx); err != nil {
				return err
			}
			lastValIdx = tmpValIdx
		}
	}

	if err := sigEnumVal.UpdateIndex(lastValIdx); err != nil {
		return err
	}

	if err := sigEnum.AddValue(sigEnumVal); err != nil {
		return err
	}

	return nil
}

func (s *SignalEnumService) ReorderValue(enumEntID, valueEntID string, from, to int) (SignalEnum, error) {
	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	if from == to {
		return s.converterFn(sigEnum), nil
	}

	sigEnumVal, err := sigEnum.GetValue(acmelib.EntityID(valueEntID))
	if err != nil {
		return SignalEnum{}, err
	}

	if err := s.reorderValue(sigEnum, sigEnumVal, from, to); err != nil {
		return SignalEnum{}, err
	}

	proxy.pushHistoryOperation(operationDomainSignalEnum,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := s.reorderValue(sigEnum, sigEnumVal, to, from); err != nil {
				return SignalEnum{}, err
			}

			return s.converterFn(sigEnum), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := s.reorderValue(sigEnum, sigEnumVal, from, to); err != nil {
				return SignalEnum{}, err
			}

			return s.converterFn(sigEnum), nil
		},
	)

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) AddValue(enumEntID string) (SignalEnum, error) {
	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	valNames := make(map[string]struct{})
	valIndexes := make(map[int]struct{})
	for _, tmpVal := range sigEnum.Values() {
		valNames[tmpVal.Name()] = struct{}{}
		valIndexes[tmpVal.Index()] = struct{}{}
	}

	valNameIdx := len(valNames) + 1
	newValName := ""
	for {
		newValName = fmt.Sprintf("NEW_VALUE_%d", valNameIdx)
		if _, ok := valNames[newValName]; !ok {
			break
		}
		valNameIdx++
	}

	newValIndex := 0
	for {
		if _, ok := valIndexes[newValIndex]; !ok {
			break
		}
		newValIndex++
	}

	newSigEnumVal := acmelib.NewSignalEnumValue(newValName, newValIndex)
	if err := sigEnum.AddValue(newSigEnumVal); err != nil {
		return SignalEnum{}, err
	}

	proxy.pushHistoryOperation(
		operationDomainSignalEnum,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := sigEnum.RemoveValue(newSigEnumVal.EntityID()); err != nil {
				return SignalEnum{}, err
			}

			return s.converterFn(sigEnum), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := sigEnum.AddValue(newSigEnumVal); err != nil {
				return SignalEnum{}, err
			}

			return s.converterFn(sigEnum), nil
		},
	)

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) RemoveValues(enumEntID string, valueEntIDs ...string) (SignalEnum, error) {
	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	if len(valueEntIDs) == 0 {
		return s.converterFn(sigEnum), nil
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	targetIDs := make(map[string]struct{})
	for _, valueEntID := range valueEntIDs {
		targetIDs[valueEntID] = struct{}{}
	}

	sigEnumValues := []*acmelib.SignalEnumValue{}
	for _, tmpValue := range sigEnum.Values() {
		tmpEntID := tmpValue.EntityID()

		_, ok := targetIDs[tmpEntID.String()]
		if !ok {
			continue
		}

		sigEnumValues = append(sigEnumValues, tmpValue)
	}

	for _, tmpValue := range sigEnumValues {
		if err := sigEnum.RemoveValue(tmpValue.EntityID()); err != nil {
			return SignalEnum{}, err
		}
	}

	proxy.pushHistoryOperation(
		operationDomainSignalEnum,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			for _, tmpValue := range sigEnumValues {
				if err := sigEnum.AddValue(tmpValue); err != nil {
					return SignalEnum{}, err
				}
			}

			return s.converterFn(sigEnum), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			for _, tmpValue := range sigEnumValues {
				if err := sigEnum.RemoveValue(tmpValue.EntityID()); err != nil {
					return SignalEnum{}, err
				}
			}

			return s.converterFn(sigEnum), nil
		},
	)

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) UpdateValueName(enumEntID, valueEntID, name string) (SignalEnum, error) {
	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	sigEnumVal, err := sigEnum.GetValue(acmelib.EntityID(valueEntID))
	if err != nil {
		return SignalEnum{}, err
	}

	oldName := sigEnumVal.Name()
	if name == oldName {
		return s.converterFn(sigEnum), nil
	}

	if err := sigEnumVal.UpdateName(name); err != nil {
		return SignalEnum{}, err
	}

	proxy.pushHistoryOperation(
		operationDomainSignalEnum,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := sigEnumVal.UpdateName(oldName); err != nil {
				return SignalEnum{}, err
			}

			return s.converterFn(sigEnum), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := sigEnumVal.UpdateName(name); err != nil {
				return SignalEnum{}, err
			}

			return s.converterFn(sigEnum), nil
		},
	)

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) UpdateValueIndex(enumEntID, valueEntID string, index int) (SignalEnum, error) {
	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	sigEnumValue, err := sigEnum.GetValue(acmelib.EntityID(valueEntID))
	if err != nil {
		return SignalEnum{}, err
	}

	oldIndex := sigEnumValue.Index()
	if index == oldIndex {
		return s.converterFn(sigEnum), nil
	}

	if err := sigEnumValue.UpdateIndex(index); err != nil {
		return SignalEnum{}, err
	}

	proxy.pushHistoryOperation(
		operationDomainSignalEnum,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := sigEnumValue.UpdateIndex(oldIndex); err != nil {
				return SignalEnum{}, err
			}

			return s.converterFn(sigEnum), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := sigEnumValue.UpdateIndex(index); err != nil {
				return SignalEnum{}, err
			}

			return s.converterFn(sigEnum), nil
		},
	)

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) UpdateValueDesc(enumEntID, valueEntID, desc string) (SignalEnum, error) {
	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	sigEnumValue, err := sigEnum.GetValue(acmelib.EntityID(valueEntID))
	if err != nil {
		return SignalEnum{}, err
	}

	oldDesc := sigEnumValue.Desc()
	if desc == oldDesc {
		return s.converterFn(sigEnum), nil
	}

	sigEnumValue.SetDesc(desc)

	proxy.pushHistoryOperation(
		operationDomainSignalEnum,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigEnumValue.SetDesc(oldDesc)

			return s.converterFn(sigEnum), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			sigEnumValue.SetDesc(desc)

			return s.converterFn(sigEnum), nil
		},
	)

	return s.converterFn(sigEnum), nil
}
