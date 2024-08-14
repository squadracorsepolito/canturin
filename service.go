package main

import (
	"errors"

	"github.com/squadracorsepolito/acmelib"
)

type serviceConverterFn[T entity, M any] func(T) M

type service[T entity, M any] struct {
	pool      map[acmelib.EntityID]T
	poolInsCh chan T

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
		case sigUnit := <-s.poolInsCh:
			s.pool[sigUnit.EntityID()] = sigUnit
		}
	}
}

func (s *service[T, M]) Open(entityID string) error {
	tmpEntID := acmelib.EntityID(entityID)
	sigUnit, ok := s.pool[tmpEntID]
	if !ok {
		return errors.New("not found")
	}
	s.opened[tmpEntID] = sigUnit
	return nil
}

func (s *service[T, M]) Close(entityID string) {
	delete(s.opened, acmelib.EntityID(entityID))
}

func (s *service[T, M]) Get(entityID string) (dummyRes M, _ error) {
	item, ok := s.pool[acmelib.EntityID(entityID)]
	if !ok {
		return dummyRes, errors.New("not found")
	}
	return s.converterFn(item), nil
}

func (s *service[T, M]) GetOpen(entityID string) (dummyRes M, _ error) {
	item, ok := s.opened[acmelib.EntityID(entityID)]
	if !ok {
		return dummyRes, errors.New("not found")
	}
	return s.converterFn(item), nil
}
