package main

import (
	"github.com/squadracorsepolito/acmelib"
)

type MessageService struct {
	*service[*acmelib.Message, Message, *messageHandler]
}

func convertMessage(m *acmelib.Message) Message {
	msg := Message{
		base: getBase(m),

		ID:             m.ID(),
		HasStaticCANID: m.HasStaticCANID(),
		CANID:          m.GetCANID(),

		SizeByte: m.SizeByte(),
		Signals:  []Signal{},

		Receivers: []Node0{},
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

func newMessageService(sidebar *sidebarController) *MessageService {
	return &MessageService{
		service: newService(serviceKindMessage, newMessageHandler(sidebar), sidebar),
	}
}

type messageHandler struct {
	*commonServiceHandler
}

func newMessageHandler(sidebar *sidebarController) *messageHandler {
	return &messageHandler{
		commonServiceHandler: newCommonServiceHandler(sidebar),
	}
}

func (h *messageHandler) toResponse(msg *acmelib.Message) Message {
	res := Message{
		base: getBase(msg),

		ID:             msg.ID(),
		HasStaticCANID: msg.HasStaticCANID(),
		CANID:          msg.GetCANID(),

		SizeByte: msg.SizeByte(),
		Signals:  []Signal{},

		Receivers: []Node0{},
	}

	for _, sig := range msg.Signals() {
		res.Signals = append(res.Signals, Signal{
			base: getBase(sig),

			Kind:     sig.Kind(),
			StartPos: sig.GetStartBit(),
			Size:     sig.GetSize(),
		})
	}

	return res
}
