package main

import (
	"context"
	"errors"
	"sync"

	"github.com/squadracorsepolito/acmelib"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type serviceConverterFn[T entity, M any] func(T) M

type service[T entity, M any] struct {
	pool      map[acmelib.EntityID]T
	poolInsCh chan T

	mux sync.RWMutex
	// opened map[acmelib.EntityID]T

	converterFn serviceConverterFn[T, M]

	stopCh chan struct{}
}

func newService[T entity, M any](poolInsCh chan T, converterFn serviceConverterFn[T, M]) *service[T, M] {
	return &service[T, M]{
		pool:      make(map[acmelib.EntityID]T),
		poolInsCh: poolInsCh,

		// opened: make(map[acmelib.EntityID]T),

		converterFn: converterFn,

		stopCh: make(chan struct{}),
	}
}

func (s *service[T, M]) run() {
	for {
		select {
		case item := <-s.poolInsCh:
			s.mux.Lock()
			s.pool[item.EntityID()] = item
			s.mux.Unlock()

		case <-s.stopCh:
			return
		}
	}
}

func (s *service[T, M]) OnStartup(_ context.Context, _ application.ServiceOptions) error {
	go s.run()

	return nil
}

func (s *service[T, M]) OnShutdown() {
	s.stopCh <- struct{}{}
	close(s.poolInsCh)
}

// func (s *service[T, M]) Open(entityID string) error {
// 	tmpEntID := acmelib.EntityID(entityID)

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	item, ok := s.pool[tmpEntID]
// 	if !ok {
// 		return errors.New("cannot open: not found")
// 	}
// 	s.opened[tmpEntID] = item
// 	return nil
// }

// func (s *service[T, M]) Close(entityID string) {
// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	delete(s.opened, acmelib.EntityID(entityID))
// }

func (s *service[T, M]) getEntity(entityID string) (T, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	item, ok := s.pool[acmelib.EntityID(entityID)]
	if !ok {
		return item, errors.New("get: not found")
	}

	return item, nil
}

func (s *service[T, M]) Get(entityID string) (dummyRes M, _ error) {
	item, err := s.getEntity(entityID)
	if err != nil {
		return dummyRes, err
	}
	return s.converterFn(item), nil
}

// func (s *service[T, M]) GetOpen(entityID string) (dummyRes M, _ error) {
// 	s.mux.RLock()
// 	defer s.mux.RUnlock()

// 	item, ok := s.opened[acmelib.EntityID(entityID)]
// 	if !ok {
// 		return dummyRes, errors.New("get open: not found")
// 	}
// 	return s.converterFn(item), nil
// }
