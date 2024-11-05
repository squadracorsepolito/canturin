package main

import (
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

	sigEnum.UpdateName(name)

	proxy.pushSidebarUpdate(sigEnum.EntityID(), name)

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) UpdateDesc(entityID string, desc string) (SignalEnum, error) {
	sigEnum, err := s.getEntity(entityID)
	if err != nil {
		return SignalEnum{}, err
	}

	sigEnum.SetDesc(desc)

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) ReorderValue(enumEntID, valueEntID string, from, to int) (SignalEnum, error) {
	if from == to {
		return s.converterFn(nil), nil
	}

	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	targetEnumVal := sigEnum.Values()[from]
	if err := sigEnum.RemoveValue(targetEnumVal.EntityID()); err != nil {
		return SignalEnum{}, err
	}
	lastValIdx := targetEnumVal.Index()

	restValues := sigEnum.Values()
	if to < from {
		// move up
		for i := from - 1; i >= to; i-- {
			tmpVal := restValues[i]
			tmpValIdx := tmpVal.Index()
			if err := tmpVal.UpdateIndex(lastValIdx); err != nil {
				return SignalEnum{}, err
			}
			lastValIdx = tmpValIdx
		}
	} else {
		// move down
		for i := from; i < to; i++ {
			tmpVal := restValues[i]
			tmpValIdx := tmpVal.Index()
			if err := tmpVal.UpdateIndex(lastValIdx); err != nil {
				return SignalEnum{}, err
			}
			lastValIdx = tmpValIdx
		}
	}

	if err := targetEnumVal.UpdateIndex(lastValIdx); err != nil {
		return SignalEnum{}, err
	}

	if err := sigEnum.AddValue(targetEnumVal); err != nil {
		return SignalEnum{}, err
	}

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) RemoveValues(enumEntID string, valueEntIDs ...string) (SignalEnum, error) {
	if len(valueEntIDs) == 0 {
		return s.converterFn(nil), nil
	}

	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	targetIDs := make(map[string]struct{})
	for _, valueEntID := range valueEntIDs {
		targetIDs[valueEntID] = struct{}{}
	}

	for _, tmpValue := range sigEnum.Values() {
		tmpEntId := tmpValue.EntityID()

		_, ok := targetIDs[tmpEntId.String()]
		if !ok {
			continue
		}

		if err := sigEnum.RemoveValue(tmpValue.EntityID()); err != nil {
			return SignalEnum{}, err
		}
	}

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) UpdateValueName(enumEntID, valueEntID, name string) (SignalEnum, error) {
	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	for _, tmpValue := range sigEnum.Values() {
		if tmpValue.EntityID().String() == valueEntID {
			if err := tmpValue.UpdateName(name); err != nil {
				return SignalEnum{}, err
			}
			break
		}
	}

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) UpdateValueIndex(enumEntID, valueEntID string, index int) (SignalEnum, error) {
	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	for _, tmpValue := range sigEnum.Values() {
		if tmpValue.EntityID().String() == valueEntID {
			if err := tmpValue.UpdateIndex(index); err != nil {
				return SignalEnum{}, err
			}
			break
		}
	}

	return s.converterFn(sigEnum), nil
}

func (s *SignalEnumService) UpdateValueDesc(enumEntID, valueEntID, desc string) (SignalEnum, error) {
	sigEnum, err := s.getEntity(enumEntID)
	if err != nil {
		return SignalEnum{}, err
	}

	for _, tmpValue := range sigEnum.Values() {
		if tmpValue.EntityID().String() == valueEntID {
			tmpValue.SetDesc(desc)
			break
		}
	}

	return s.converterFn(sigEnum), nil
}
