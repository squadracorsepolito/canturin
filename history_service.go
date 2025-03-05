package main

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type History struct {
	OperationCount int  `json:"operationCount"`
	CurrentIndex   int  `json:"currentIndex"`
	Saved          bool `json:"saved"`
}

type operationFunc func() (any, error)

type operation struct {
	serviceKind serviceKind
	undo        operationFunc
	redo        operationFunc
}

type HistoryService struct {
	operations []*operation
	currOpIdx  int

	saved bool

	mux sync.RWMutex

	operationCh chan *operation
	stopCh      chan struct{}
}

func newHistoryService() *HistoryService {
	return &HistoryService{
		operations: []*operation{},
		currOpIdx:  -1,

		saved: true,

		operationCh: make(chan *operation),
		stopCh:      make(chan struct{}),
	}
}

func (s *HistoryService) OnStartup(_ context.Context, _ application.ServiceOptions) error {
	go s.run()
	return nil
}

func (s *HistoryService) OnShutdown() error {
	s.stopCh <- struct{}{}
	return nil
}

func (s *HistoryService) run() {
	for {
		select {
		case op := <-s.operationCh:
			s.handleOperation(op)

		case <-s.stopCh:
			return
		}
	}
}

func (s *HistoryService) getState() History {
	return History{
		OperationCount: len(s.operations),
		CurrentIndex:   s.currOpIdx,
		Saved:          s.saved,
	}
}

func (s *HistoryService) emitHistoryChange() {
	app.EmitEvent(HistoryChange, s.getState())
}

func (s *HistoryService) handleOperation(op *operation) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.saved = false

	if s.currOpIdx == -1 {
		s.operations = []*operation{}
	} else if s.currOpIdx < len(s.operations)-1 {
		s.operations = s.operations[:s.currOpIdx+1]
	}

	s.operations = append(s.operations, op)
	s.currOpIdx++

	s.emitHistoryChange()
}

func (s *HistoryService) Undo() (History, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.currOpIdx == -1 {
		return s.getState(), nil
	}

	if s.currOpIdx > len(s.operations)-1 {
		s.currOpIdx = len(s.operations) - 1
	}

	op := s.operations[s.currOpIdx]

	res, err := op.undo()
	if err != nil {
		return s.getState(), err
	}

	s.saved = false
	s.sendModifyEvent(op.serviceKind, res)

	s.currOpIdx--

	return s.getState(), nil
}

func (s *HistoryService) Redo() (History, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.currOpIdx == len(s.operations)-1 {
		return s.getState(), nil
	}

	s.currOpIdx++

	op := s.operations[s.currOpIdx]
	res, err := s.operations[s.currOpIdx].redo()
	if err != nil {
		return s.getState(), err
	}

	s.saved = false
	s.sendModifyEvent(op.serviceKind, res)

	return s.getState(), nil
}

func (s *HistoryService) sendModifyEvent(opDomain serviceKind, res any) {
	eventName := ""
	switch opDomain {
	case serviceKindNetwork:
		eventName = HistoryNetworkModify
	case serviceKindBus:
		eventName = HistoryBusModify
	case serviceKindNode:
		eventName = HistoryNodeModify
	case serviceKindMessage:
		eventName = HistoryMessageModify
	case serviceKindSignal:
		eventName = HistorySignalModify
	case serviceKindSignalType:
		eventName = HistorySignalTypeModify
	case serviceKindSignalUnit:
		eventName = HistorySignalUnitModify
	case serviceKindSignalEnum:
		eventName = HistorySignalEnumModify
	}

	application.Get().EmitEvent(eventName, res)
}

func (s *HistoryService) save() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.saved = true
	s.emitHistoryChange()
}

func (s *HistoryService) isSaved() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.saved
}

func (s *HistoryService) getController() *historyController {
	return &historyController{
		operationCh: s.operationCh,
	}
}

type historyController struct {
	operationCh chan<- *operation
}

func (hc *historyController) sendOperation(serviceKind serviceKind, undo, redo operationFunc) {
	hc.operationCh <- &operation{
		serviceKind: serviceKind,
		undo:        undo,
		redo:        redo,
	}
}
