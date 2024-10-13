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

	converterFn serviceConverterFn[T, M]

	stopCh chan struct{}
}

func newService[T entity, M any](poolInsCh chan T, converterFn serviceConverterFn[T, M]) *service[T, M] {
	return &service[T, M]{
		pool:      make(map[acmelib.EntityID]T),
		poolInsCh: poolInsCh,

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

func (s *service[T, M]) GetNames() []string {
	s.mux.RLock()
	defer s.mux.RUnlock()

	names := []string{}
	for _, item := range s.pool {
		names = append(names, item.Name())
	}
	return names
}
