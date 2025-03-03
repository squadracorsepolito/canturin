package main

import (
	"time"

	"github.com/squadracorsepolito/acmelib"
)

type entityStub struct {
	EntityID string `json:"entityId"`
	Name     string `json:"name"`
}

type NetworkStub struct {
	entityStub

	Buses       []BusStub        `json:"buses"`
	SignalUnits []SignalUnitStub `json:"signalUnits"`
	SignalTypes []SignalTypeStub `json:"signalTypes"`
}

type BusStub struct {
	entityStub

	Nodes []NodeStub `json:"nodes"`
}

type NodeStub struct {
	entityStub

	SendedMessages []MessageStub `json:"sendedMessages"`
}

type MessageStub struct {
	entityStub

	Signals []entityStub `json:"signals"`
}

type SignalUnitStub struct {
	entityStub
}

type SignalTypeStub struct {
	entityStub
}

type base struct {
	EntityID   string    `json:"entityId"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	CreateTime time.Time `json:"createTime"`
}

type entity interface {
	EntityKind() acmelib.EntityKind
	EntityID() acmelib.EntityID
	Name() string
	Desc() string
	CreateTime() time.Time
}

func getBase(e entity) base {
	return base{
		EntityID:   e.EntityID().String(),
		Name:       e.Name(),
		Desc:       e.Desc(),
		CreateTime: e.CreateTime(),
	}
}

type Bus0 struct {
	base

	Nodes []Node0 `json:"nodes"`
}

type Node0 struct {
	base

	SendedMessages []Message `json:"sendedMessages"`
}

type SignalReference struct {
	Bus     entityStub `json:"bus"`
	Node    entityStub `json:"node"`
	Message entityStub `json:"message"`
	Signal  entityStub `json:"signal"`
}

type entityWithRefs[T entity] interface {
	References() []T
}

type EntityKind string

const (
	EntityKindNetwork    EntityKind = "network"
	EntityKindBus        EntityKind = "bus"
	EntityKindNode       EntityKind = "node"
	EntityKindMessage    EntityKind = "message"
	EntityKindSignal     EntityKind = "signal"
	EntityKindSignalType EntityKind = "signal-type"
	EntityKindSignalUnit EntityKind = "signal-unit"
	EntityKindSignalEnum EntityKind = "signal-enum"
)

func newEntityKind(kind acmelib.EntityKind) EntityKind {
	switch kind {
	case acmelib.EntityKindNetwork:
		return EntityKindNetwork
	case acmelib.EntityKindBus:
		return EntityKindBus
	case acmelib.EntityKindNode:
		return EntityKindNode
	case acmelib.EntityKindMessage:
		return EntityKindMessage
	case acmelib.EntityKindSignal:
		return EntityKindSignal
	case acmelib.EntityKindSignalType:
		return EntityKindSignalType
	case acmelib.EntityKindSignalUnit:
		return EntityKindSignalUnit
	case acmelib.EntityKindSignalEnum:
		return EntityKindSignalEnum
	default:
		return EntityKindNetwork
	}
}

type EntityPath struct {
	Kind     EntityKind `json:"kind"`
	EntityID string     `json:"entityId"`
	Name     string     `json:"name"`
}

func newEntityPath(ent entity) EntityPath {
	return EntityPath{
		Kind:     newEntityKind(ent.EntityKind()),
		EntityID: ent.EntityID().String(),
		Name:     ent.Name(),
	}
}

func newBusEntityPaths(bus *acmelib.Bus) []EntityPath {
	res := []EntityPath{}

	parNet := bus.ParentNetwork()
	if parNet != nil {
		res = append(res, newEntityPath(parNet))
	}

	res = append(res, newEntityPath(bus))

	return res
}

func newNodeInterfaceEntityPaths(nodeInt *acmelib.NodeInterface) []EntityPath {
	res := []EntityPath{}

	parBus := nodeInt.ParentBus()
	if parBus != nil {
		res = newBusEntityPaths(parBus)
	}

	res = append(res, newEntityPath(nodeInt.Node()))

	return res
}

func newMessageEntityPaths(msg *acmelib.Message) []EntityPath {
	res := []EntityPath{}

	parNodeInt := msg.SenderNodeInterface()
	if parNodeInt != nil {
		res = newNodeInterfaceEntityPaths(parNodeInt)
	}

	res = append(res, newEntityPath(msg))

	return res
}

func newSignalEntityPaths(sig acmelib.Signal) []EntityPath {
	res := []EntityPath{}

	parMsg := sig.ParentMessage()
	if parMsg != nil {
		res = newMessageEntityPaths(parMsg)
	}

	res = append(res, newEntityPath(sig))

	return res
}
