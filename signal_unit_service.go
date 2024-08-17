package main

import (
	"slices"
	"strings"

	"github.com/squadracorsepolito/acmelib"
)

type SignalUnit struct {
	base

	Symbol     string            `json:"symbol"`
	References []SignalReference `json:"references"`
}

type SignalUnitService struct {
	*service[*acmelib.SignalUnit, SignalUnit]
}

func newSignalUnitService(sigUnitCh chan *acmelib.SignalUnit) *SignalUnitService {
	return &SignalUnitService{
		service: newService(sigUnitCh, func(su *acmelib.SignalUnit) SignalUnit {
			res := SignalUnit{
				base: getBase(su),

				Symbol:     su.Symbol(),
				References: []SignalReference{},
			}

			for _, tmpStdSig := range su.References() {
				tmpMsg := tmpStdSig.ParentMessage()
				tmpNode := tmpMsg.SenderNodeInterface().Node()
				tmpBus := tmpMsg.SenderNodeInterface().ParentBus()

				res.References = append(res.References, SignalReference{
					Bus:     getEntityStub(tmpBus),
					Node:    getEntityStub(tmpNode),
					Message: getEntityStub(tmpMsg),
					Signal:  getEntityStub(tmpStdSig),
				})
			}

			slices.SortFunc(res.References, func(a, b SignalReference) int {
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
		}),
	}
}

type SignalReference struct {
	Bus     entityStub `json:"bus"`
	Node    entityStub `json:"node"`
	Message entityStub `json:"message"`
	Signal  entityStub `json:"signal"`
}
