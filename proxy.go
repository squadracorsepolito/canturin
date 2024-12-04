package main

import (
	"github.com/squadracorsepolito/acmelib"
)

type appProxy struct {
	network *acmelib.Network

	historyOperationCh chan *operation

	busCh     chan *acmelib.Bus
	nodeCh    chan *acmelib.Node
	messageCh chan *acmelib.Message

	sigTypeCh    chan *acmelib.SignalType
	sigUnitCh    chan *acmelib.SignalUnit
	signalEnumCh chan *acmelib.SignalEnum
}

func newAppProxy() *appProxy {
	return &appProxy{

		historyOperationCh: make(chan *operation),

		busCh:     make(chan *acmelib.Bus),
		nodeCh:    make(chan *acmelib.Node),
		messageCh: make(chan *acmelib.Message),

		sigTypeCh:    make(chan *acmelib.SignalType),
		sigUnitCh:    make(chan *acmelib.SignalUnit),
		signalEnumCh: make(chan *acmelib.SignalEnum),
	}
}

func (p *appProxy) pushHistoryOperation(domain operationDomain, undo, redo operationFunc) {
	p.historyOperationCh <- &operation{
		domain: domain,
		undo:   undo,
		redo:   redo,
	}
}

func (p *appProxy) pushBus(bus *acmelib.Bus) {
	p.busCh <- bus
}

func (p *appProxy) pushNode(node *acmelib.Node) {
	p.nodeCh <- node
}

func (p *appProxy) pushMessage(msg *acmelib.Message) {
	p.messageCh <- msg
}

func (p *appProxy) pushSignalType(sigType *acmelib.SignalType) {
	p.sigTypeCh <- sigType
}

func (p *appProxy) pushSignalUnit(sigUnit *acmelib.SignalUnit) {
	p.sigUnitCh <- sigUnit
}

func (p *appProxy) pushSignalEnum(sigEnum *acmelib.SignalEnum) {
	p.signalEnumCh <- sigEnum
}
