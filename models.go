package main

import (
	"time"

	"github.com/squadracorsepolito/acmelib"
)

type base struct {
	ID         string    `json:"id"`
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
		ID:         e.EntityID().String(),
		Name:       e.Name(),
		Desc:       e.Desc(),
		CreateTime: e.CreateTime(),
	}
}

type Network struct {
	base

	Buses []*Bus `json:"buses"`
}

type Bus struct {
	base

	Nodes []*Node `json:"nodes"`
}

type Node struct {
	base

	SendedMessages []*Message `json:"sendedMessages"`
}

type Message struct {
	base

	SizeByte int      `json:"sizeByte"`
	Signals  []Signal `json:"signals"`
}

type Signal struct {
	base

	Kind     acmelib.SignalKind `json:"kind"`
	StartPos int                `json:"startPos"`
	Size     int                `json:"size"`
}
