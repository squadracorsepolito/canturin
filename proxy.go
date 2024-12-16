package main

import (
	"github.com/squadracorsepolito/acmelib"
)

type appProxy struct {
	network *acmelib.Network

	historyOperationCh chan *operation
}

func newAppProxy() *appProxy {
	return &appProxy{

		historyOperationCh: make(chan *operation),
	}
}

func (p *appProxy) pushHistoryOperation(domain operationDomain, undo, redo operationFunc) {
	p.historyOperationCh <- &operation{
		domain: domain,
		undo:   undo,
		redo:   redo,
	}
}
