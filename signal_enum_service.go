package main

import (
	"fmt"

	"github.com/squadracorsepolito/acmelib"
)

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

type SignalEnum struct {
	base

	Size     int               `json:"size"`
	MinSize  int               `json:"minSize"`
	MaxIndex int               `json:"maxIndex"`
	Values   []SignalEnumValue `json:"values"`

	References []SignalReference `json:"references"`
}

func newSignalEnum(sigEnum *acmelib.SignalEnum) SignalEnum {
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

type SignalEnumService struct {
	*service[*acmelib.SignalEnum, SignalEnum, *signalEnumHandler]
}

func newSignalEnumService(sidebar *sidebarController) *SignalEnumService {
	return &SignalEnumService{
		service: newService(serviceKindSignalEnum, newSignalEnumHandler(sidebar), sidebar),
	}
}

func (s *SignalEnumService) Create(req CreateSignalEnumReq) (SignalEnum, error) {
	sigEnum := acmelib.NewSignalEnum(req.Name)

	if len(req.Desc) > 0 {
		sigEnum.SetDesc(req.Desc)
	}

	if req.MinSize > 0 {
		sigEnum.SetMinSize(req.MinSize)
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	s.addEntity(sigEnum)
	s.sidebar.sendAdd(sigEnum)

	s.sendHistoryOp(
		func() (*acmelib.SignalEnum, error) {
			s.removeEntity(sigEnum.EntityID().String())
			s.sidebar.sendDelete(sigEnum)
			return sigEnum, nil
		},
		func() (*acmelib.SignalEnum, error) {
			s.addEntity(sigEnum)
			s.sidebar.sendAdd(sigEnum)
			return sigEnum, nil
		},
	)

	return s.handler.toResponse(sigEnum), nil
}

func (s *SignalEnumService) Delete(entityID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	sigEnum, err := s.getEntity(entityID)
	if err != nil {
		return err
	}

	if sigEnum.ReferenceCount() > 0 {
		return fmt.Errorf("signal enum %s is referenced %d times", sigEnum.Name(), sigEnum.ReferenceCount())
	}

	s.removeEntity(entityID)
	s.sidebar.sendDelete(sigEnum)

	s.sendHistoryOp(
		func() (*acmelib.SignalEnum, error) {
			s.addEntity(sigEnum)
			s.sidebar.sendAdd(sigEnum)
			return sigEnum, nil
		},
		func() (*acmelib.SignalEnum, error) {
			s.removeEntity(sigEnum.EntityID().String())
			s.sidebar.sendDelete(sigEnum)
			return sigEnum, nil
		},
	)

	return nil
}

func (s *SignalEnumService) UpdateName(entityID string, req UpdateNameReq) (SignalEnum, error) {
	return s.handle(entityID, &req, s.handler.updateName)
}

func (s *SignalEnumService) UpdateDesc(entityID string, req UpdateDescReq) (SignalEnum, error) {
	return s.handle(entityID, &req, s.handler.updateDesc)
}

func (s *SignalEnumService) AddValue(entityID string) (SignalEnum, error) {
	return s.handle(entityID, nil, s.handler.addValue)
}

func (s *SignalEnumService) RemoveValues(entityID string, req RemoveValuesReq) (SignalEnum, error) {
	return s.handle(entityID, &req, s.handler.removeValues)
}

func (s *SignalEnumService) ReorderValue(entityID string, req ReorderValueReq) (SignalEnum, error) {
	return s.handle(entityID, &req, s.handler.reorderValueHandler)
}

func (s *SignalEnumService) UpdateValueName(entityID string, req UpdateValueNameReq) (SignalEnum, error) {
	return s.handle(entityID, &req, s.handler.updateValueName)
}

func (s *SignalEnumService) UpdateValueDesc(entityID string, req UpdateValueDescReq) (SignalEnum, error) {
	return s.handle(entityID, &req, s.handler.updateValueDesc)
}

func (s *SignalEnumService) UpdateValueIndex(entityID string, req UpdateValueIndexReq) (SignalEnum, error) {
	return s.handle(entityID, &req, s.handler.updateValueIndex)
}

type signalEnumRes = response[*acmelib.SignalEnum]

type signalEnumHandler struct {
	*commonServiceHandler
}

func newSignalEnumHandler(sidebar *sidebarController) *signalEnumHandler {
	return &signalEnumHandler{
		commonServiceHandler: newCommonServiceHandler(sidebar),
	}
}

func (h *signalEnumHandler) toResponse(sigEnum *acmelib.SignalEnum) SignalEnum {
	return newSignalEnum(sigEnum)
}

func (h *signalEnumHandler) updateName(sigEnum *acmelib.SignalEnum, req *request, res *signalEnumRes) error {
	parsedReq := req.toUpdateName()

	name := parsedReq.Name

	oldName := sigEnum.Name()
	if oldName == name {
		return nil
	}

	sigEnum.UpdateName(name)
	h.sidebar.sendUpdateName(sigEnum)

	res.setUndo(
		func() (*acmelib.SignalEnum, error) {
			sigEnum.UpdateName(oldName)
			h.sidebar.sendUpdateName(sigEnum)
			return sigEnum, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalEnum, error) {
			sigEnum.UpdateName(name)
			h.sidebar.sendUpdateName(sigEnum)
			return sigEnum, nil
		},
	)

	return nil
}

func (h *signalEnumHandler) updateDesc(sigEnum *acmelib.SignalEnum, req *request, res *signalEnumRes) error {
	parsedReq := req.toUpdateDesc()

	desc := parsedReq.Desc

	oldDesc := sigEnum.Desc()
	if oldDesc == desc {
		return nil
	}

	sigEnum.SetDesc(desc)

	res.setUndo(
		func() (*acmelib.SignalEnum, error) {
			sigEnum.SetDesc(oldDesc)
			return sigEnum, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalEnum, error) {
			sigEnum.SetDesc(desc)
			return sigEnum, nil
		},
	)

	return nil
}

func (h *signalEnumHandler) addValue(sigEnum *acmelib.SignalEnum, _ *request, res *signalEnumRes) error {
	valNames := make(map[string]struct{})
	valIndexes := make(map[int]struct{})
	for _, tmpVal := range sigEnum.Values() {
		valNames[tmpVal.Name()] = struct{}{}
		valIndexes[tmpVal.Index()] = struct{}{}
	}

	valNameIdx := len(valNames) + 1
	valName := ""
	for {
		valName = fmt.Sprintf("NEW_VALUE_%d", valNameIdx)
		if _, ok := valNames[valName]; !ok {
			break
		}
		valNameIdx++
	}

	valIndex := 0
	for {
		if _, ok := valIndexes[valIndex]; !ok {
			break
		}
		valIndex++
	}

	sigEnumVal := acmelib.NewSignalEnumValue(valName, valIndex)
	if err := sigEnum.AddValue(sigEnumVal); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.SignalEnum, error) {
			if err := sigEnum.RemoveValue(sigEnumVal.EntityID()); err != nil {
				return nil, err
			}
			return sigEnum, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalEnum, error) {
			if err := sigEnum.AddValue(sigEnumVal); err != nil {
				return nil, err
			}
			return sigEnum, nil
		},
	)

	return nil
}

func (h *signalEnumHandler) removeValues(sigEnum *acmelib.SignalEnum, req *request, res *signalEnumRes) error {
	parsedReq := req.toRemoveValues()

	if len(parsedReq.ValueEntityIDs) == 0 {
		return nil
	}

	remValIDs := make(map[string]struct{})
	for _, tmpEntID := range parsedReq.ValueEntityIDs {
		remValIDs[tmpEntID] = struct{}{}
	}

	remValues := []*acmelib.SignalEnumValue{}
	for _, tmpVal := range sigEnum.Values() {
		if _, ok := remValIDs[tmpVal.EntityID().String()]; ok {
			remValues = append(remValues, tmpVal)
		}
	}

	for _, tmpVal := range remValues {
		if err := sigEnum.RemoveValue(tmpVal.EntityID()); err != nil {
			return err
		}
	}

	res.setUndo(
		func() (*acmelib.SignalEnum, error) {
			for _, tmpVal := range remValues {
				if err := sigEnum.AddValue(tmpVal); err != nil {
					return nil, err
				}
			}
			return sigEnum, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalEnum, error) {
			for _, tmpVal := range remValues {
				if err := sigEnum.RemoveValue(tmpVal.EntityID()); err != nil {
					return nil, err
				}
			}
			return sigEnum, nil
		},
	)

	return nil
}

func (h *signalEnumHandler) reorderValue(sigEnum *acmelib.SignalEnum, sigEnumVal *acmelib.SignalEnumValue, from, to int) error {
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

func (h *signalEnumHandler) reorderValueHandler(sigEnum *acmelib.SignalEnum, req *request, res *signalEnumRes) error {
	parsedReq := req.toReorderValue()

	from := parsedReq.From
	to := parsedReq.To

	if from == to {
		return nil
	}

	sigEnumVal, err := sigEnum.GetValue(acmelib.EntityID(parsedReq.ValueEntityID))
	if err != nil {
		return err
	}

	if err := h.reorderValue(sigEnum, sigEnumVal, from, to); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.SignalEnum, error) {
			if err := h.reorderValue(sigEnum, sigEnumVal, to, from); err != nil {
				return nil, err
			}
			return sigEnum, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalEnum, error) {
			if err := h.reorderValue(sigEnum, sigEnumVal, from, to); err != nil {
				return nil, err
			}
			return sigEnum, nil
		},
	)

	return nil
}

func (h *signalEnumHandler) updateValueName(sigEnum *acmelib.SignalEnum, req *request, res *signalEnumRes) error {
	parsedReq := req.toUpdateValueName()

	name := parsedReq.Name

	sigEnumVal, err := sigEnum.GetValue(acmelib.EntityID(parsedReq.ValueEntityID))
	if err != nil {
		return err
	}

	oldName := sigEnumVal.Name()
	if oldName == name {
		return nil
	}

	if err := sigEnumVal.UpdateName(name); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.SignalEnum, error) {
			if err := sigEnumVal.UpdateName(oldName); err != nil {
				return nil, err
			}
			return sigEnum, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalEnum, error) {
			if err := sigEnumVal.UpdateName(name); err != nil {
				return nil, err
			}
			return sigEnum, nil
		},
	)

	return nil
}

func (h *signalEnumHandler) updateValueDesc(sigEnum *acmelib.SignalEnum, req *request, res *signalEnumRes) error {
	parsedReq := req.toUpdateValueDesc()

	desc := parsedReq.Desc

	sigEnumVal, err := sigEnum.GetValue(acmelib.EntityID(parsedReq.ValueEntityID))
	if err != nil {
		return err
	}

	oldDesc := sigEnumVal.Desc()
	if oldDesc == desc {
		return nil
	}

	sigEnumVal.SetDesc(desc)

	res.setUndo(
		func() (*acmelib.SignalEnum, error) {
			sigEnumVal.SetDesc(oldDesc)
			return sigEnum, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalEnum, error) {
			sigEnumVal.SetDesc(desc)
			return sigEnum, nil
		},
	)

	return nil
}

func (h *signalEnumHandler) updateValueIndex(sigEnum *acmelib.SignalEnum, req *request, res *signalEnumRes) error {
	parsedReq := req.toUpdateValueIndex()

	index := parsedReq.Index

	sigEnumVal, err := sigEnum.GetValue(acmelib.EntityID(parsedReq.ValueEntityID))
	if err != nil {
		return err
	}

	oldIndex := sigEnumVal.Index()
	if oldIndex == index {
		return nil
	}

	if err := sigEnumVal.UpdateIndex(index); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.SignalEnum, error) {
			if err := sigEnumVal.UpdateIndex(oldIndex); err != nil {
				return nil, err
			}
			return sigEnum, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalEnum, error) {
			if err := sigEnumVal.UpdateIndex(index); err != nil {
				return nil, err
			}
			return sigEnum, nil
		},
	)

	return nil
}

//
//
//

// type SignalEnumService0 struct {
// 	*service[*acmelib.SignalEnum, SignalEnum]
// }

// func newSignalEnumService() *SignalEnumService0 {
// 	return &SignalEnumService0{
// 		service: newService(newSignalEnum),
// 	}
// }

// func (s *SignalEnumService0) sendSidebarAdd(sigEnum *acmelib.SignalEnum) {
// 	item := newSignalEnumSidebarItem(sigEnum)
// 	manager.sidebar.sendAdd(newSidebarAddReq(item, SidebarSignalEnumsPrefix))
// }

// func (s *SignalEnumService0) sendSidebarUpdateName(sigEnum *acmelib.SignalEnum) {
// 	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(sigEnum.EntityID().String(), sigEnum.Name()))
// }

// func (s *SignalEnumService0) sendSidebarDelete(sigEnum *acmelib.SignalEnum) {
// 	manager.sidebar.sendDelete(newSidebarDeleteReq(sigEnum.EntityID().String()))
// }

// func (s *SignalEnumService0) Create(name, desc string, minSize int) (SignalEnum, error) {
// 	sigEnum := acmelib.NewSignalEnum(name)

// 	if minSize > 0 {
// 		sigEnum.SetMinSize(minSize)
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	s.addEntity(sigEnum)
// 	s.sendSidebarAdd(sigEnum)

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalEnum,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			s.deleteEntity(sigEnum.EntityID().String())
// 			s.sendSidebarDelete(sigEnum)

// 			return s.converterFn(sigEnum), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			s.addEntity(sigEnum)
// 			s.sendSidebarAdd(sigEnum)

// 			return s.converterFn(sigEnum), nil
// 		},
// 	)

// 	return s.converterFn(sigEnum), nil
// }

// func (s *SignalEnumService0) Delete(entityID string) error {
// 	sigEnum, err := s.getEntity(entityID)
// 	if err != nil {
// 		return err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	if sigEnum.ReferenceCount() > 0 {
// 		return fmt.Errorf("signal enum %s is referenced %d times", sigEnum.Name(), sigEnum.ReferenceCount())
// 	}

// 	s.deleteEntity(entityID)
// 	s.sendSidebarDelete(sigEnum)

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalEnum,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			s.addEntity(sigEnum)
// 			s.sendSidebarAdd(sigEnum)

// 			return s.converterFn(sigEnum), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			s.deleteEntity(sigEnum.EntityID().String())
// 			s.sendSidebarDelete(sigEnum)

// 			return s.converterFn(sigEnum), nil
// 		},
// 	)

// 	return nil
// }

// func (s *SignalEnumService0) UpdateName(entityID string, name string) (SignalEnum, error) {
// 	sigEnum, err := s.getEntity(entityID)
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	oldName := sigEnum.Name()
// 	if name == oldName {
// 		return s.converterFn(sigEnum), nil
// 	}

// 	sigEnum.UpdateName(name)

// 	s.sendSidebarUpdateName(sigEnum)

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalEnum,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			sigEnum.UpdateName(oldName)
// 			s.sendSidebarUpdateName(sigEnum)

// 			return s.converterFn(sigEnum), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			sigEnum.UpdateName(name)
// 			s.sendSidebarUpdateName(sigEnum)

// 			return s.converterFn(sigEnum), nil
// 		},
// 	)

// 	return s.converterFn(sigEnum), nil
// }

// func (s *SignalEnumService0) UpdateDesc(entityID string, desc string) (SignalEnum, error) {
// 	sigEnum, err := s.getEntity(entityID)
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	oldDesc := sigEnum.Desc()
// 	if desc == oldDesc {
// 		return s.converterFn(sigEnum), nil
// 	}

// 	sigEnum.SetDesc(desc)

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalEnum,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			sigEnum.SetDesc(oldDesc)

// 			return s.converterFn(sigEnum), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			sigEnum.SetDesc(desc)

// 			return s.converterFn(sigEnum), nil
// 		},
// 	)

// 	return s.converterFn(sigEnum), nil
// }

// func (s *SignalEnumService0) reorderValue(sigEnum *acmelib.SignalEnum, sigEnumVal *acmelib.SignalEnumValue, from, to int) error {
// 	if from == to {
// 		return nil
// 	}

// 	if err := sigEnum.RemoveValue(sigEnumVal.EntityID()); err != nil {
// 		return err
// 	}
// 	lastValIdx := sigEnumVal.Index()

// 	restValues := sigEnum.Values()
// 	if to < from {
// 		// move up
// 		for i := from - 1; i >= to; i-- {
// 			tmpVal := restValues[i]
// 			tmpValIdx := tmpVal.Index()
// 			if err := tmpVal.UpdateIndex(lastValIdx); err != nil {
// 				return err
// 			}
// 			lastValIdx = tmpValIdx
// 		}
// 	} else {
// 		// move down
// 		for i := from; i < to; i++ {
// 			tmpVal := restValues[i]
// 			tmpValIdx := tmpVal.Index()
// 			if err := tmpVal.UpdateIndex(lastValIdx); err != nil {
// 				return err
// 			}
// 			lastValIdx = tmpValIdx
// 		}
// 	}

// 	if err := sigEnumVal.UpdateIndex(lastValIdx); err != nil {
// 		return err
// 	}

// 	if err := sigEnum.AddValue(sigEnumVal); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *SignalEnumService0) ReorderValue(enumEntID, valueEntID string, from, to int) (SignalEnum, error) {
// 	sigEnum, err := s.getEntity(enumEntID)
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	if from == to {
// 		return s.converterFn(sigEnum), nil
// 	}

// 	sigEnumVal, err := sigEnum.GetValue(acmelib.EntityID(valueEntID))
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	if err := s.reorderValue(sigEnum, sigEnumVal, from, to); err != nil {
// 		return SignalEnum{}, err
// 	}

// 	proxy.pushHistoryOperation(operationDomainSignalEnum,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := s.reorderValue(sigEnum, sigEnumVal, to, from); err != nil {
// 				return SignalEnum{}, err
// 			}

// 			return s.converterFn(sigEnum), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := s.reorderValue(sigEnum, sigEnumVal, from, to); err != nil {
// 				return SignalEnum{}, err
// 			}

// 			return s.converterFn(sigEnum), nil
// 		},
// 	)

// 	return s.converterFn(sigEnum), nil
// }

// func (s *SignalEnumService0) AddValue(enumEntID string) (SignalEnum, error) {
// 	sigEnum, err := s.getEntity(enumEntID)
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	valNames := make(map[string]struct{})
// 	valIndexes := make(map[int]struct{})
// 	for _, tmpVal := range sigEnum.Values() {
// 		valNames[tmpVal.Name()] = struct{}{}
// 		valIndexes[tmpVal.Index()] = struct{}{}
// 	}

// 	valNameIdx := len(valNames) + 1
// 	newValName := ""
// 	for {
// 		newValName = fmt.Sprintf("NEW_VALUE_%d", valNameIdx)
// 		if _, ok := valNames[newValName]; !ok {
// 			break
// 		}
// 		valNameIdx++
// 	}

// 	newValIndex := 0
// 	for {
// 		if _, ok := valIndexes[newValIndex]; !ok {
// 			break
// 		}
// 		newValIndex++
// 	}

// 	newSigEnumVal := acmelib.NewSignalEnumValue(newValName, newValIndex)
// 	if err := sigEnum.AddValue(newSigEnumVal); err != nil {
// 		return SignalEnum{}, err
// 	}

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalEnum,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := sigEnum.RemoveValue(newSigEnumVal.EntityID()); err != nil {
// 				return SignalEnum{}, err
// 			}

// 			return s.converterFn(sigEnum), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := sigEnum.AddValue(newSigEnumVal); err != nil {
// 				return SignalEnum{}, err
// 			}

// 			return s.converterFn(sigEnum), nil
// 		},
// 	)

// 	return s.converterFn(sigEnum), nil
// }

// func (s *SignalEnumService0) RemoveValues(enumEntID string, valueEntIDs ...string) (SignalEnum, error) {
// 	sigEnum, err := s.getEntity(enumEntID)
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	if len(valueEntIDs) == 0 {
// 		return s.converterFn(sigEnum), nil
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	targetIDs := make(map[string]struct{})
// 	for _, valueEntID := range valueEntIDs {
// 		targetIDs[valueEntID] = struct{}{}
// 	}

// 	sigEnumValues := []*acmelib.SignalEnumValue{}
// 	for _, tmpValue := range sigEnum.Values() {
// 		tmpEntID := tmpValue.EntityID()

// 		_, ok := targetIDs[tmpEntID.String()]
// 		if !ok {
// 			continue
// 		}

// 		sigEnumValues = append(sigEnumValues, tmpValue)
// 	}

// 	for _, tmpValue := range sigEnumValues {
// 		if err := sigEnum.RemoveValue(tmpValue.EntityID()); err != nil {
// 			return SignalEnum{}, err
// 		}
// 	}

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalEnum,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			for _, tmpValue := range sigEnumValues {
// 				if err := sigEnum.AddValue(tmpValue); err != nil {
// 					return SignalEnum{}, err
// 				}
// 			}

// 			return s.converterFn(sigEnum), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			for _, tmpValue := range sigEnumValues {
// 				if err := sigEnum.RemoveValue(tmpValue.EntityID()); err != nil {
// 					return SignalEnum{}, err
// 				}
// 			}

// 			return s.converterFn(sigEnum), nil
// 		},
// 	)

// 	return s.converterFn(sigEnum), nil
// }

// func (s *SignalEnumService0) UpdateValueName(enumEntID, valueEntID, name string) (SignalEnum, error) {
// 	sigEnum, err := s.getEntity(enumEntID)
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	sigEnumVal, err := sigEnum.GetValue(acmelib.EntityID(valueEntID))
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	oldName := sigEnumVal.Name()
// 	if name == oldName {
// 		return s.converterFn(sigEnum), nil
// 	}

// 	if err := sigEnumVal.UpdateName(name); err != nil {
// 		return SignalEnum{}, err
// 	}

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalEnum,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := sigEnumVal.UpdateName(oldName); err != nil {
// 				return SignalEnum{}, err
// 			}

// 			return s.converterFn(sigEnum), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := sigEnumVal.UpdateName(name); err != nil {
// 				return SignalEnum{}, err
// 			}

// 			return s.converterFn(sigEnum), nil
// 		},
// 	)

// 	return s.converterFn(sigEnum), nil
// }

// func (s *SignalEnumService0) UpdateValueIndex(enumEntID, valueEntID string, index int) (SignalEnum, error) {
// 	sigEnum, err := s.getEntity(enumEntID)
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	sigEnumValue, err := sigEnum.GetValue(acmelib.EntityID(valueEntID))
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	oldIndex := sigEnumValue.Index()
// 	if index == oldIndex {
// 		return s.converterFn(sigEnum), nil
// 	}

// 	if err := sigEnumValue.UpdateIndex(index); err != nil {
// 		return SignalEnum{}, err
// 	}

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalEnum,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := sigEnumValue.UpdateIndex(oldIndex); err != nil {
// 				return SignalEnum{}, err
// 			}

// 			return s.converterFn(sigEnum), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := sigEnumValue.UpdateIndex(index); err != nil {
// 				return SignalEnum{}, err
// 			}

// 			return s.converterFn(sigEnum), nil
// 		},
// 	)

// 	return s.converterFn(sigEnum), nil
// }

// func (s *SignalEnumService0) UpdateValueDesc(enumEntID, valueEntID, desc string) (SignalEnum, error) {
// 	sigEnum, err := s.getEntity(enumEntID)
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	sigEnumValue, err := sigEnum.GetValue(acmelib.EntityID(valueEntID))
// 	if err != nil {
// 		return SignalEnum{}, err
// 	}

// 	oldDesc := sigEnumValue.Desc()
// 	if desc == oldDesc {
// 		return s.converterFn(sigEnum), nil
// 	}

// 	sigEnumValue.SetDesc(desc)

// 	proxy.pushHistoryOperation(
// 		operationDomainSignalEnum,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			sigEnumValue.SetDesc(oldDesc)

// 			return s.converterFn(sigEnum), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			sigEnumValue.SetDesc(desc)

// 			return s.converterFn(sigEnum), nil
// 		},
// 	)

// 	return s.converterFn(sigEnum), nil
// }
