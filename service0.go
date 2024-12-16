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

func getBaseEntity(e entity) BaseEntity {
	return BaseEntity{
		EntityID:   e.EntityID().String(),
		Name:       e.Name(),
		Desc:       e.Desc(),
		CreateTime: e.CreateTime(),
	}
}

type serviceKind int

const (
	serviceKindBus serviceKind = iota
	serviceKindNode
	serviceKindMessage
	serviceKindSignalType
	serviceKindSignalUnit
	serviceKindSignalEnum
)

type serviceHandler[E entity, R any] interface {
	toResponse(entity E) R
}

type commonServiceHandler struct {
	sidebar *sidebarController
}

func newCommonServiceHandler(sidebar *sidebarController) *commonServiceHandler {
	return &commonServiceHandler{
		sidebar: sidebar,
	}
}

type service0[E entity, R any, H serviceHandler[E, R]] struct {
	kind serviceKind

	handler H

	mux      sync.RWMutex
	entities map[acmelib.EntityID]E

	loadCh chan E
	stopCh chan struct{}

	sidebar *sidebarController
}

func newService0[E entity, R any, H serviceHandler[E, R]](kind serviceKind, handlers H, sidebar *sidebarController) *service0[E, R, H] {
	return &service0[E, R, H]{
		kind: kind,

		handler: handlers,

		entities: make(map[acmelib.EntityID]E),

		loadCh: make(chan E),
		stopCh: make(chan struct{}),

		sidebar: sidebar,
	}
}

func (s *service0[E, R, H]) sendLoad(ent E) {
	s.loadCh <- ent
}

func (s *service0[E, R, H]) run() {
	for {
		select {
		case ent := <-s.loadCh:
			s.mux.Lock()
			s.addEntity(ent)
			s.mux.Unlock()

		case <-s.stopCh:
			return
		}
	}
}

func (s *service0[E, R, H]) OnStartup(_ context.Context, _ application.ServiceOptions) error {
	go s.run()
	return nil
}

func (s *service0[E, R, H]) OnShutdown() {
	s.stopCh <- struct{}{}
	close(s.loadCh)
}

// func (s *service0[E, R, H]) sendSidebarAdd(ent E) {
// 	var item *sidebarItem
// 	var parItemKey string

// 	switch s.kind {
// 	case serviceKindBus:
// 		item = newSidebarItem(SidebarItemKindBus, ent.EntityID(), SidebarBusesPrefix, ent.Name())
// 		parItemKey = SidebarBusesPrefix

// 	case serviceKindNode:
// 		item = newSidebarItem(SidebarItemKindNode, ent.EntityID(), SidebarNodesPrefix, ent.Name())
// 		parItemKey = SidebarNodesPrefix

// 	case serviceKindMessage:
// 		item = newSidebarItem(SidebarItemKindMessage, ent.EntityID(), SidebarMessagesPrefix, ent.Name())
// 		parItemKey = SidebarMessagesPrefix

// 	case serviceKindSignalType:
// 		item = newSidebarItem(SidebarItemKindSignalType, ent.EntityID(), SidebarSignalTypesPrefix, ent.Name())
// 		parItemKey = SidebarSignalTypesPrefix

// 	case serviceKindSignalUnit:
// 		item = newSidebarItem(SidebarItemKindSignalUnit, ent.EntityID(), SidebarSignalUnitsPrefix, ent.Name())
// 		parItemKey = SidebarSignalUnitsPrefix

// 	case serviceKindSignalEnum:
// 		item = newSidebarItem(SidebarItemKindSignalEnum, ent.EntityID(), SidebarSignalEnumsPrefix, ent.Name())
// 		parItemKey = SidebarSignalEnumsPrefix
// 	}

// 	manager.sidebar.sendAdd(newSidebarAddReq(item, parItemKey))
// }

// func (s *service0[E, R, H]) sendSidebarUpdateName(ent E) {
// 	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(ent.EntityID().String(), ent.Name()))
// }

// func (s *service0[E, R, H]) sendSidebarRemove(ent E) {
// 	manager.sidebar.sendDelete(newSidebarDeleteReq(ent.EntityID().String()))
// }

func (s *service0[E, R, H]) sendHistoryOp(undo, redo func() (E, error)) {
	var opDomain operationDomain

	switch s.kind {
	case serviceKindBus:
		opDomain = operationDomainBus

	case serviceKindNode:
		opDomain = operationDomainNode

	case serviceKindMessage:
		opDomain = operationDomainMessage

	case serviceKindSignalType:
		opDomain = operationDomainSignalType

	case serviceKindSignalUnit:
		opDomain = operationDomainSignalUnit

	case serviceKindSignalEnum:
		opDomain = operationDomainSignalEnum
	}

	proxy.pushHistoryOperation(opDomain,
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

func (s *service0[E, R, H]) handle(entityID string, reqDataPtr any, handlerFn func(E, *request, *response[E]) error) (dummyRes R, _ error) {
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

func (s *service0[E, R, H]) crossHandle(entityID string, handlerFn func(E) error) (E, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	ent, err := s.getEntity(entityID)
	if err != nil {
		return ent, err
	}

	if err := handlerFn(ent); err != nil {
		return ent, err
	}

	return ent, nil
}

func (s *service0[E, R, H]) addEntity(ent E) {
	s.entities[ent.EntityID()] = ent
}

func (s *service0[E, R, H]) getEntity(entityID string) (E, error) {
	ent, ok := s.entities[acmelib.EntityID(entityID)]
	if !ok {
		return ent, errors.New("get entity: not found")
	}

	return ent, nil
}

func (s *service0[E, R, H]) removeEntity(entityID string) {
	delete(s.entities, acmelib.EntityID(entityID))
}

func (s *service0[E, R, H]) Get(entityID string) (dummyRes R, _ error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	ent, err := s.getEntity(entityID)
	if err != nil {
		return dummyRes, err
	}

	return s.handler.toResponse(ent), nil
}

func (s *service0[E, R, H]) GetInvalidNames(entityID string) []string {
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

func (s *service0[E, R, H]) ListBase() []BaseEntity {
	s.mux.RLock()
	defer s.mux.RUnlock()

	var res []BaseEntity
	for _, ent := range s.entities {
		res = append(res, getBaseEntity(ent))
	}

	return res
}
