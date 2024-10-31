package main

import (
	"github.com/squadracorsepolito/acmelib"
)

type appProxy struct {
	network *acmelib.Network

	sidebarLoadCh   chan *sidebarLoadReq
	sidebarUpdateCh chan *sidebarUpdateReq
	sidebarAddCh    chan *sidebarAddReq
	sidebarRemoveCh chan *sidebarRemoveReq

	historyOperationCh chan *operation

	messageCh chan *acmelib.Message

	sigTypeCh    chan *acmelib.SignalType
	sigUnitCh    chan *acmelib.SignalUnit
	signalEnumCh chan *acmelib.SignalEnum
}

func newAppProxy() *appProxy {
	return &appProxy{
		sidebarLoadCh:   make(chan *sidebarLoadReq),
		sidebarUpdateCh: make(chan *sidebarUpdateReq),
		sidebarAddCh:    make(chan *sidebarAddReq),
		sidebarRemoveCh: make(chan *sidebarRemoveReq),

		historyOperationCh: make(chan *operation),

		messageCh: make(chan *acmelib.Message),

		sigTypeCh:    make(chan *acmelib.SignalType),
		sigUnitCh:    make(chan *acmelib.SignalUnit),
		signalEnumCh: make(chan *acmelib.SignalEnum),
	}
}

func (p *appProxy) pushSidebarLoad(network *acmelib.Network) {
	p.sidebarLoadCh <- &sidebarLoadReq{
		network: network,
	}
}

func (p *appProxy) pushSidebarUpdate(entID acmelib.EntityID, name string) {
	p.sidebarUpdateCh <- &sidebarUpdateReq{
		entityID: entID,
		name:     name,
	}
}

func (p *appProxy) pushSidebarAdd(kind SidebarNodeKind, entID, parentID acmelib.EntityID, name string) {
	p.sidebarAddCh <- &sidebarAddReq{
		kind:     kind,
		entityID: entID,
		name:     name,
		parentID: parentID,
	}
}

func (p *appProxy) pushSidebarRemove(entID acmelib.EntityID) {
	p.sidebarRemoveCh <- &sidebarRemoveReq{
		entityID: entID,
	}
}

func (p *appProxy) pushHistoryOperation(domain operationDomain, undo, redo operationFunc) {
	p.historyOperationCh <- &operation{
		domain: domain,
		undo:   undo,
		redo:   redo,
	}
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
