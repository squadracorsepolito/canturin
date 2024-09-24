package main

import (
	"github.com/squadracorsepolito/acmelib"
)

type MessageService struct {
	*service[*acmelib.Message, Message]
}

func convertMessage(m *acmelib.Message) Message {
	msg := Message{
		base: getBase(m),

		ID:             m.ID(),
		HasStaticCANID: m.HasStaticCANID(),
		CANID:          m.GetCANID(),

		SizeByte: m.SizeByte(),
		Signals:  []Signal{},

		Receivers: []Node{},
	}

	for _, sig := range m.Signals() {
		msg.Signals = append(msg.Signals, Signal{
			base: getBase(sig),

			Kind:     sig.Kind(),
			StartPos: sig.GetStartBit(),
			Size:     sig.GetSize(),
		})
	}

	return msg
}

func newMessageService(messageCh chan *acmelib.Message) *MessageService {
	return &MessageService{
		newService(messageCh, convertMessage),
	}
}
