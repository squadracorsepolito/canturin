package main

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type History struct {
	OperationCount int `json:"operationCount"`
	CurrentIndex   int `json:"currentIndex"`
}

type operationDomain int

const (
	operationDomainSignalType operationDomain = iota
	operationDomainSignalUnit
	operationDomainSignalEnum
)

type operationFunc func() (any, error)

type operation struct {
	domain operationDomain
	undo   operationFunc
	redo   operationFunc
}

type HistoryService struct {
	operations []*operation
	currOpIdx  int

	mux sync.RWMutex

	stopCh chan struct{}
}

func newHistoryService() *HistoryService {
	return &HistoryService{
		operations: []*operation{},
		currOpIdx:  -1,

		stopCh: make(chan struct{}),
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
		case op := <-proxy.historyOperationCh:
			s.addOperation(op)

		case <-s.stopCh:
			return
		}
	}
}

func (s *HistoryService) getState() History {
	return History{
		OperationCount: len(s.operations),
		CurrentIndex:   s.currOpIdx,
	}
}

func (s *HistoryService) emitHistoryChange() {
	app.EmitEvent(HistoryChange, s.getState())
}

func (s *HistoryService) addOperation(op *operation) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.currOpIdx == -1 {
		s.operations = []*operation{}
	} else if s.currOpIdx < len(s.operations)-1 {
		s.operations = s.operations[:s.currOpIdx]
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

	op := s.operations[s.currOpIdx]

	res, err := op.undo()
	if err != nil {
		return s.getState(), err
	}
	s.sendModifyEvent(op.domain, res)

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
	s.sendModifyEvent(op.domain, res)

	return s.getState(), nil
}

func (s *HistoryService) sendModifyEvent(opDomain operationDomain, res any) {
	eventName := ""
	switch opDomain {
	case operationDomainSignalType:
		eventName = HistorySignalTypeModify
	case operationDomainSignalUnit:
		eventName = HistorySignalUnitModify
	case operationDomainSignalEnum:
		eventName = HistorySignalEnumModify
	}

	app.EmitEvent(eventName, res)
}
