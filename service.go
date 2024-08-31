package main

import (
	"errors"
	"sync"

	"github.com/squadracorsepolito/acmelib"
)

type serviceConverterFn[T entity, M any] func(T) M

type service[T entity, M any] struct {
	pool      map[acmelib.EntityID]T
	poolInsCh chan T

	mux    sync.RWMutex
	opened map[acmelib.EntityID]T

	converterFn serviceConverterFn[T, M]
}

func newService[T entity, M any](poolInsCh chan T, converterFn serviceConverterFn[T, M]) *service[T, M] {
	return &service[T, M]{
		pool:      make(map[acmelib.EntityID]T),
		poolInsCh: poolInsCh,

		opened: make(map[acmelib.EntityID]T),

		converterFn: converterFn,
	}
}

func (s *service[T, M]) run() {
	for {
		select {
		case item := <-s.poolInsCh:
			s.mux.Lock()
			s.pool[item.EntityID()] = item
			s.mux.Unlock()
		}
	}
}

func (s *service[T, M]) Open(entityID string) error {
	tmpEntID := acmelib.EntityID(entityID)

	s.mux.Lock()
	defer s.mux.Unlock()

	item, ok := s.pool[tmpEntID]
	if !ok {
		return errors.New("cannot open: not found")
	}
	s.opened[tmpEntID] = item
	return nil
}

func (s *service[T, M]) Close(entityID string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.opened, acmelib.EntityID(entityID))
}

func (s *service[T, M]) Get(entityID string) (dummyRes M, _ error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	item, ok := s.pool[acmelib.EntityID(entityID)]
	if !ok {
		return dummyRes, errors.New("get: not found")
	}
	return s.converterFn(item), nil
}

func (s *service[T, M]) GetOpen(entityID string) (dummyRes M, _ error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	item, ok := s.opened[acmelib.EntityID(entityID)]
	if !ok {
		return dummyRes, errors.New("get open: not found")
	}
	return s.converterFn(item), nil
}
