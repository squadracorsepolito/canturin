package main

import (
	"errors"

	"github.com/squadracorsepolito/acmelib"
)

type MessageService struct {
	ns *NetworkService

	message *acmelib.Message
}

func (ms *MessageService) toFrontend() Message {
	msg := Message{
		base: getBase(ms.message),

		ID:             ms.message.ID(),
		HasStaticCANID: ms.message.HasStaticCANID(),
		CANID:          ms.message.GetCANID(),

		SizeByte: ms.message.SizeByte(),
		Signals:  []Signal{},

		Receivers: []Node{},
	}

	for _, sig := range ms.message.Signals() {
		msg.Signals = append(msg.Signals, Signal{
			base: getBase(sig),

			Kind:     sig.Kind(),
			StartPos: sig.GetStartBit(),
			Size:     sig.GetSize(),
		})
	}

	for _, rec := range ms.message.Receivers() {
		msg.Receivers = append(msg.Receivers, Node{
			base: getBase(rec.Node()),
		})
	}

	return msg
}

func (ms *MessageService) Register(busEntID, nodeEntID, msgEntID string) error {
	for _, bus := range ms.ns.network.Buses() {
		if bus.EntityID() == acmelib.EntityID(busEntID) {
			for _, nodeInt := range bus.NodeInterfaces() {
				if nodeInt.Node().EntityID() == acmelib.EntityID(nodeEntID) {
					for _, msg := range nodeInt.Messages() {
						if msg.EntityID() == acmelib.EntityID(msgEntID) {
							ms.message = msg
							return nil
						}
					}
				}
			}
		}
	}

	return errors.New("not found")
}

func (ms *MessageService) Get() Message {
	return ms.toFrontend()
}

func (ms *MessageService) SetDesc(desc string) Message {
	ms.message.SetDesc(desc)
	return ms.toFrontend()
}

func (ms *MessageService) UpdateName(name string) (Message, error) {
	if err := ms.message.UpdateName(name); err != nil {
		return Message{}, err
	}
	return ms.toFrontend(), nil
}
