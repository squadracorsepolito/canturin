package main

import (
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
}

type SignalUnitStub struct {
	entityStub
}

type base struct {
	EntityID   string    `json:"entityId"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	CreateTime time.Time `json:"createTime"`
}

type entity interface {
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

type Network struct {
	base

	Buses []Bus `json:"buses"`
}

type Bus struct {
	base

	Nodes []Node `json:"nodes"`
}

type Node struct {
	base

	SendedMessages []Message `json:"sendedMessages"`
}

type Message struct {
	base

	ID             acmelib.MessageID `json:"id"`
	HasStaticCANID bool              `json:"hasStaticCANID"`
	CANID          acmelib.CANID     `json:"canId"`

	SizeByte  int                      `json:"sizeByte"`
	ByteOrder acmelib.MessageByteOrder `json:"byteOrder"`
	Signals   []Signal                 `json:"signals"`

	Receivers []Node `json:"receivers"`
}

type Signal struct {
	base

	Kind     acmelib.SignalKind `json:"kind"`
	StartPos int                `json:"startPos"`
	Size     int                `json:"size"`
}

type SignalType struct {
	base
}

type withReferences struct {
	References []entityStub `json:"references"`
}

func getWithReferences(refs []entity) withReferences {
	res := withReferences{
		References: []entityStub{},
	}

	for _, tmpEnt := range refs {
		res.References = append(res.References, getEntityStub(tmpEnt))
	}

	return res
}
