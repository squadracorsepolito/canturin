package main

import (
	"slices"
	"strings"
	"time"

	"github.com/squadracorsepolito/acmelib"
)

type entityStub struct {
	EntityID string `json:"entityId"`
	Name     string `json:"name"`
}

func getEntityStub(e entity) entityStub {
	return entityStub{
		EntityID: e.EntityID().String(),
		Name:     e.Name(),
	}
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

func getSignalReferences(refs entityWithRefs[*acmelib.StandardSignal]) []SignalReference {
	res := []SignalReference{}

	for _, tmpStdSig := range refs.References() {
		tmpMsg := tmpStdSig.ParentMessage()
		tmpNode := tmpMsg.SenderNodeInterface().Node()
		tmpBus := tmpMsg.SenderNodeInterface().ParentBus()

		res = append(res, SignalReference{
			Bus:     getEntityStub(tmpBus),
			Node:    getEntityStub(tmpNode),
			Message: getEntityStub(tmpMsg),
			Signal:  getEntityStub(tmpStdSig),
		})
	}

	slices.SortFunc(res, func(a, b SignalReference) int {
		busCmp := strings.Compare(a.Bus.Name, b.Bus.Name)
		if busCmp != 0 {
			return busCmp
		}

		nodeCmp := strings.Compare(a.Node.Name, b.Node.Name)
		if nodeCmp != 0 {
			return nodeCmp
		}

		msgCmp := strings.Compare(a.Message.Name, b.Message.Name)
		if msgCmp != 0 {
			return msgCmp
		}

		return strings.Compare(a.Signal.Name, b.Signal.Name)
	})

	return res
}
