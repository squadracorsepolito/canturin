package main

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/squadracorsepolito/acmelib"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type BaseEntity struct {
	EntityID   string    `json:"entityId"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	CreateTime time.Time `json:"createTime"`
}

func newBaseEntity(e entity) BaseEntity {
	return BaseEntity{
		EntityID:   e.EntityID().String(),
		Name:       e.Name(),
		Desc:       e.Desc(),
		CreateTime: e.CreateTime(),
	}
}

type serviceKind int

const (
	serviceKindNetwork serviceKind = iota
	serviceKindBus
	serviceKindNode
	serviceKindMessage
	serviceKindSignal
	serviceKindSignalType
	serviceKindSignalUnit
	serviceKindSignalEnum
)

type serviceHandler[E entity, R any] interface {
	toResponse(entity E) R
}

type commonServiceHandler struct {
	sidebarCtr *sidebarController
}

func newCommonServiceHandler(sidebarCtr *sidebarController) *commonServiceHandler {
	return &commonServiceHandler{
		sidebarCtr: sidebarCtr,
	}
}

type service[E entity, R any, H serviceHandler[E, R]] struct {
	kind serviceKind

	handler H

	mux      *sync.RWMutex
	entities map[acmelib.EntityID]E

	loadCh   chan []E
	addCh    chan E
	deleteCh chan E
	clearCh  chan struct{}

	sidebarCtr *sidebarController
	historyCtr *historyController
}

func newService[E entity, R any, H serviceHandler[E, R]](kind serviceKind, handlers H, mux *sync.RWMutex, sidebarCtr *sidebarController) *service[E, R, H] {
	return &service[E, R, H]{
		kind: kind,

		handler: handlers,

		mux:      mux,
		entities: make(map[acmelib.EntityID]E),

		loadCh:   make(chan []E),
		addCh:    make(chan E),
		deleteCh: make(chan E),
		clearCh:  make(chan struct{}),

		sidebarCtr: sidebarCtr,
	}
}

func (s *service[E, R, H]) setHistoryController(historyCtr *historyController) {
	s.historyCtr = historyCtr
}

func (s *service[E, R, H]) OnStartup(ctx context.Context, _ application.ServiceOptions) error {
	go s.run(ctx)
	return nil
}

func (s *service[E, R, H]) OnShutdown() {}

func (s *service[E, R, H]) run(ctx context.Context) {
	for {
		select {
		case entities := <-s.loadCh:
			s.handleLoad(entities)

		case ent := <-s.addCh:
			s.handleAdd(ent)

		case ent := <-s.deleteCh:
			s.handleDelete(ent)

		case <-s.clearCh:
			s.handleClear()

		case <-ctx.Done():
			return
		}
	}
}

func (s *service[E, R, H]) handleLoad(entities []E) {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, ent := range entities {
		s.addEntity(ent)
	}
}

func (s *service[E, R, H]) handleAdd(ent E) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.addEntity(ent)
	s.sidebarCtr.sendAdd(ent)

	addEventName := ""
	switch s.kind {
	case serviceKindBus:
		addEventName = BusAdded
	case serviceKindNode:
		addEventName = NodeAdded
	case serviceKindMessage:
		addEventName = MessageAdded
	case serviceKindSignal:
		addEventName = SignalAdded
	case serviceKindSignalType:
		addEventName = SignalTypeAdded
	case serviceKindSignalUnit:
		addEventName = SignalUnitAdded
	case serviceKindSignalEnum:
		addEventName = SignalEnumAdded
	}

	if len(addEventName) > 0 {
		application.Get().EmitEvent(addEventName, s.handler.toResponse(ent))
	}
}

func (s *service[E, R, H]) handleDelete(ent E) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.removeEntity(ent.EntityID().String())
	s.sidebarCtr.sendDelete(ent)
}

func (s *service[E, R, H]) handleClear() {
	s.mux.Lock()
	defer s.mux.Unlock()

	clear(s.entities)
}

func (s *service[E, R, H]) sendHistoryOp(undo, redo func() (E, error)) {
	s.historyCtr.sendOperation(
		s.kind,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			res, err := undo()
			if err != nil {
				return nil, err
			}

			return s.handler.toResponse(res), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			res, err := redo()
			if err != nil {
				return nil, err
			}

			return s.handler.toResponse(res), nil
		},
	)
}

func (s *service[E, R, H]) handle(entityID string, reqDataPtr any, handlerFn func(E, *request, *response[E]) error) (dummyRes R, _ error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	ent, err := s.getEntity(entityID)
	if err != nil {
		return dummyRes, err
	}

	req := newRequest(reqDataPtr)
	res := newResponse[E]()

	if err := handlerFn(ent, req, res); err != nil {
		return dummyRes, err
	}

	if res.changed {
		s.sendHistoryOp(res.undo, res.redo)
	}

	return s.handler.toResponse(ent), nil
}

func (s *service[E, R, H]) addEntity(ent E) {
	s.entities[ent.EntityID()] = ent
}

func (s *service[E, R, H]) getEntity(entityID string) (E, error) {
	ent, ok := s.entities[acmelib.EntityID(entityID)]
	if !ok {
		return ent, errors.New("get entity: not found")
	}

	return ent, nil
}

func (s *service[E, R, H]) removeEntity(entityID string) {
	delete(s.entities, acmelib.EntityID(entityID))
}

func (s *service[E, R, H]) Get(entityID string) (dummyRes R, _ error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	ent, err := s.getEntity(entityID)
	if err != nil {
		return dummyRes, err
	}

	return s.handler.toResponse(ent), nil
}

func (s *service[E, R, H]) GetInvalidNames(entityID string) []string {
	s.mux.RLock()
	defer s.mux.RUnlock()

	names := []string{}
	for _, tmpEnt := range s.entities {
		if tmpEnt.EntityID() == acmelib.EntityID(entityID) {
			continue
		}

		names = append(names, tmpEnt.Name())
	}

	return names
}

func (s *service[E, R, H]) ListBase() []BaseEntity {
	s.mux.RLock()
	defer s.mux.RUnlock()

	var res []BaseEntity
	for _, ent := range s.entities {
		res = append(res, newBaseEntity(ent))
	}

	return res
}

func (s *service[E, R, H]) getController() *serviceController[E] {
	return &serviceController[E]{
		getFn: s.getEntity,

		loadCh:   s.loadCh,
		addCh:    s.addCh,
		deleteCh: s.deleteCh,
		clearCh:  s.clearCh,
	}
}

type serviceController[E entity] struct {
	getFn func(entityID string) (E, error)

	loadCh   chan<- []E
	addCh    chan<- E
	deleteCh chan<- E
	clearCh  chan<- struct{}
}

func (sc *serviceController[E]) get(entityID string) (E, error) {
	return sc.getFn(entityID)
}

func (sc *serviceController[E]) sendLoad(entities []E) {
	sc.loadCh <- entities
}

func (sc *serviceController[E]) sendAdd(ent E) {
	sc.addCh <- ent
}

func (sc *serviceController[E]) sendDelete(ent E) {
	sc.deleteCh <- ent
}

func (sc *serviceController[E]) sendClear() {
	sc.clearCh <- struct{}{}
}

type busController = serviceController[*acmelib.Bus]
type nodeController = serviceController[*acmelib.Node]
type messageController = serviceController[*acmelib.Message]
type signalController = serviceController[acmelib.Signal]
type signalTypeController = serviceController[*acmelib.SignalType]
type signalUnitController = serviceController[*acmelib.SignalUnit]
type signalEnumController = serviceController[*acmelib.SignalEnum]
