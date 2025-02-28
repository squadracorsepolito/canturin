package main

import (
	"errors"
	"slices"
	"strings"
	"sync"

	"github.com/squadracorsepolito/acmelib"
)

type MessageSendType string

const (
	MessageSendTypeUnset                      MessageSendType = "unset"
	MessageSendTypeCyclic                     MessageSendType = "cyclic"
	MessageSendTypeCyclicIfActive             MessageSendType = "cyclic_if_active"
	MessageSendTypeCyclicAndTriggered         MessageSendType = "cyclic_and_triggered"
	MessageSendTypeCyclicIfActiveAndTriggered MessageSendType = "cyclic_if_active_and_triggered"
)

func newMessageSendType(st acmelib.MessageSendType) MessageSendType {
	switch st {
	case acmelib.MessageSendTypeCyclic:
		return MessageSendTypeCyclic
	case acmelib.MessageSendTypeCyclicIfActive:
		return MessageSendTypeCyclicIfActive
	case acmelib.MessageSendTypeCyclicAndTriggered:
		return MessageSendTypeCyclicAndTriggered
	case acmelib.MessageSendTypeCyclicIfActiveAndTriggered:
		return MessageSendTypeCyclicIfActiveAndTriggered
	default:
		return MessageSendTypeUnset
	}
}

func (st MessageSendType) parse() acmelib.MessageSendType {
	switch st {
	case MessageSendTypeCyclic:
		return acmelib.MessageSendTypeCyclic
	case MessageSendTypeCyclicIfActive:
		return acmelib.MessageSendTypeCyclicIfActive
	case MessageSendTypeCyclicAndTriggered:
		return acmelib.MessageSendTypeCyclicAndTriggered
	case MessageSendTypeCyclicIfActiveAndTriggered:
		return acmelib.MessageSendTypeCyclicIfActiveAndTriggered
	default:
		return acmelib.MessageSendTypeUnset
	}
}

type MessageByteOrder string

const (
	MessageByteOrderLittleEndian MessageByteOrder = "little-endian"
	MessageByteOrderBigEndian    MessageByteOrder = "big-endian"
)

func newMessageByteOrder(bo acmelib.MessageByteOrder) MessageByteOrder {
	switch bo {
	case acmelib.MessageByteOrderLittleEndian:
		return MessageByteOrderLittleEndian
	case acmelib.MessageByteOrderBigEndian:
		return MessageByteOrderBigEndian
	default:
		return MessageByteOrderLittleEndian
	}
}

func (bo MessageByteOrder) parse() acmelib.MessageByteOrder {
	switch bo {
	case MessageByteOrderLittleEndian:
		return acmelib.MessageByteOrderLittleEndian
	case MessageByteOrderBigEndian:
		return acmelib.MessageByteOrderBigEndian
	default:
		return acmelib.MessageByteOrderLittleEndian
	}
}

type Message struct {
	base

	Paths []EntityPath `json:"paths"`

	HasStaticCANID bool `json:"hasStaticCANID"`
	ID             uint `json:"id"`
	CANID          uint `json:"canId"`

	SizeByte               int              `json:"sizeByte"`
	AvailableTrailingBytes int              `json:"availableTrailingBytes"`
	MaxAvailableSpace      int              `json:"maxAvailableSpace"`
	ByteOrder              MessageByteOrder `json:"byteOrder"`

	CycleTime      int             `json:"cycleTime"`
	SendType       MessageSendType `json:"sendType"`
	DelayTime      int             `json:"delayTime"`
	StartDelayTime int             `json:"startDelayTime"`

	Signals []Signal `json:"signals"`

	Receivers []Node0 `json:"receivers"`

	SenderNode BaseEntity `json:"senderNode"`
	ParentBus  BaseEntity `json:"parentBus"`
}

func newMessage(msg *acmelib.Message) Message {
	res := Message{
		base: getBase(msg),

		Paths: newMessageEntityPaths(msg),

		HasStaticCANID: msg.HasStaticCANID(),
		ID:             uint(msg.ID()),
		CANID:          uint(msg.GetCANID()),

		SizeByte:               msg.SizeByte(),
		AvailableTrailingBytes: msg.SizeByte(),
		MaxAvailableSpace:      0,
		ByteOrder:              newMessageByteOrder(msg.ByteOrder()),

		CycleTime:      msg.CycleTime(),
		SendType:       newMessageSendType(msg.SendType()),
		DelayTime:      msg.DelayTime(),
		StartDelayTime: msg.StartDelayTime(),

		Signals: []Signal{},

		Receivers: []Node0{},
	}

	if nodeInt := msg.SenderNodeInterface(); nodeInt != nil {
		res.SenderNode = newBaseEntity(nodeInt.Node())

		if bus := nodeInt.ParentBus(); bus != nil {
			res.ParentBus = newBaseEntity(bus)
		}
	}

	signals := msg.Signals()

	if len(signals) > 0 {
		lastSig := signals[len(signals)-1]
		trailingBits := msg.SizeByte()*8 - lastSig.GetStartBit() - lastSig.GetSize()
		res.AvailableTrailingBytes = trailingBits / 8
	}

	holes := []int{}
	currPos := 0
	for _, sig := range signals {
		res.Signals = append(res.Signals, Signal{
			base: getBase(sig),

			Kind:     newSignalKind(sig.Kind()),
			StartPos: sig.GetStartBit(),
			Size:     sig.GetSize(),
		})

		if currPos < sig.GetStartBit() {
			holes = append(holes, sig.GetStartBit()-currPos)
		}

		currPos = sig.GetStartBit() + sig.GetSize()
	}

	if currPos < msg.SizeByte()*8 {
		holes = append(holes, msg.SizeByte()*8-currPos)
	}

	if len(holes) > 0 {
		res.MaxAvailableSpace = slices.Max(holes)
	}

	return res
}

type MessageService struct {
	*service[*acmelib.Message, Message, *messageHandler]
}

func newMessageService(mux *sync.RWMutex, sidebarCtr *sidebarController, signalCtr *signalController) *MessageService {
	return &MessageService{
		service: newService(serviceKindMessage, newMessageHandler(sidebarCtr, signalCtr), mux, sidebarCtr),
	}
}

func (s *MessageService) GetInvalidMessageIDs(entityID string) []uint {
	s.mux.RLock()
	defer s.mux.RUnlock()

	messageIDs := []uint{}

	msg, err := s.getEntity(entityID)
	if err != nil {
		return messageIDs
	}

	nodeInt := msg.SenderNodeInterface()
	if nodeInt == nil {
		return messageIDs
	}

	for _, tmpMsg := range nodeInt.SentMessages() {
		if tmpMsg.EntityID().String() == entityID || tmpMsg.HasStaticCANID() {
			continue
		}

		messageIDs = append(messageIDs, uint(tmpMsg.ID()))
	}

	return messageIDs
}

func (s *MessageService) GetInvalidCANIDs(entityID string, busEntityID string) []uint {
	s.mux.RLock()
	defer s.mux.RUnlock()

	canIDs := []uint{}
	for _, tmpMsg := range s.entities {
		if tmpMsg.EntityID().String() == entityID {
			continue
		}

		nodeInt := tmpMsg.SenderNodeInterface()
		if nodeInt == nil {
			continue
		}

		bus := nodeInt.ParentBus()
		if bus == nil {
			continue
		}

		if bus.EntityID().String() == busEntityID {
			canIDs = append(canIDs, uint(tmpMsg.GetCANID()))
		}
	}

	return canIDs
}

func (s *MessageService) GetSpaceLeft(entityID string) int {
	s.mux.RLock()
	defer s.mux.RUnlock()

	msg, err := s.getEntity(entityID)
	if err != nil {
		return 0
	}

	spaceLeft := msg.SizeByte() * 8
	for _, sig := range msg.Signals() {
		spaceLeft -= sig.GetSize()
	}

	return spaceLeft
}

func (s *MessageService) UpdateName(entityID string, req UpdateNameReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.updateName)
}

func (s *MessageService) UpdateDesc(entityID string, req UpdateDescReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.updateDesc)
}

func (s *MessageService) UpdateMessageID(entityID string, req UpdateMessageIDReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.updateMessageID)
}

func (s *MessageService) UpdateStaticCANID(entityID string, req UpdateStaticCANIDReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.updateStaticCANID)
}

func (s *MessageService) UpdateSizeByte(entityID string, req UpdateSizeByteReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.updateSizeByte)
}

func (s *MessageService) UpdateByteOrder(entityID string, req UpdateByteOrderReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.updateByteOrder)
}

func (s *MessageService) UpdateCycleTime(entityID string, req UpdateCycleTimeReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.updateCycleTime)
}

func (s *MessageService) UpdateSendType(entityID string, req UpdateSendTypeReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.updateSendType)
}

func (s *MessageService) UpdateDelayTime(entityID string, req UpdateDelayTimeReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.updateDelayTime)
}

func (s *MessageService) UpdateStartDelayTime(entityID string, req UpdateStartDelayTimeReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.updateStartDelayTime)
}

func (s *MessageService) AddSignal(entityID string, req AddSignalReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.addSignal)
}

func (s *MessageService) DeleteSignals(entityID string, req DeleteSignalsReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.deleteSignals)
}

func (s *MessageService) CompactSignals(entityID string) (Message, error) {
	return s.handle(entityID, nil, s.handler.compactSignals)
}

func (s *MessageService) ReorderSignal(entityID string, req ReorderSignalReq) (Message, error) {
	return s.handle(entityID, &req, s.handler.reorderSignalHandler)
}

type messageRes = response[*acmelib.Message]

type messageHandler struct {
	*commonServiceHandler

	signalCtr *signalController
}

func newMessageHandler(sidebar *sidebarController, signalCtr *signalController) *messageHandler {
	return &messageHandler{
		commonServiceHandler: newCommonServiceHandler(sidebar),

		signalCtr: signalCtr,
	}
}

func (h *messageHandler) toResponse(msg *acmelib.Message) Message {
	return newMessage(msg)
}

func (h *messageHandler) updateName(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toUpdateName()

	name := strings.TrimSpace(parsedReq.Name)

	oldName := msg.Name()
	if name == oldName {
		return nil
	}

	if err := msg.UpdateName(name); err != nil {
		return err
	}
	h.sidebarCtr.sendUpdateName(msg)

	res.setUndo(
		func() (*acmelib.Message, error) {
			if err := msg.UpdateName(oldName); err != nil {
				return nil, err
			}
			h.sidebarCtr.sendUpdateName(msg)
			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			if err := msg.UpdateName(name); err != nil {
				return nil, err
			}
			h.sidebarCtr.sendUpdateName(msg)
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) updateDesc(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toUpdateDesc()

	desc := parsedReq.Desc

	oldDesc := msg.Desc()
	if desc == oldDesc {
		return nil
	}

	msg.SetDesc(desc)

	res.setUndo(
		func() (*acmelib.Message, error) {
			msg.SetDesc(oldDesc)
			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			msg.SetDesc(desc)
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) updateMessageID(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toUpdateMessageID()

	msgID := acmelib.MessageID(parsedReq.MessageID)

	oldMsgID := msg.ID()
	if msgID == oldMsgID {
		return nil
	}

	if err := msg.UpdateID(msgID); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.Message, error) {
			if err := msg.UpdateID(oldMsgID); err != nil {
				return nil, err
			}
			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			if err := msg.UpdateID(msgID); err != nil {
				return nil, err
			}
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) updateStaticCANID(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toUpdateStaticCANID()

	staticCANID := acmelib.CANID(parsedReq.StaticCANID)

	oldStaticCANID := msg.GetCANID()
	if staticCANID == oldStaticCANID {
		return nil
	}

	var oldID acmelib.MessageID
	wasStatic := msg.HasStaticCANID()
	if !wasStatic {
		oldID = msg.ID()
	}

	if err := msg.SetStaticCANID(staticCANID); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.Message, error) {
			if wasStatic {
				if err := msg.SetStaticCANID(oldStaticCANID); err != nil {
					return nil, err
				}
			} else {
				if err := msg.UpdateID(oldID); err != nil {
					return nil, err
				}
			}

			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			if err := msg.SetStaticCANID(staticCANID); err != nil {
				return nil, err
			}
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) updateSizeByte(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toUpdateSizeByte()

	sizeByte := parsedReq.SizeByte

	oldSizeByte := msg.SizeByte()
	if sizeByte == oldSizeByte {
		return nil
	}

	if err := msg.UpdateSizeByte(sizeByte); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.Message, error) {
			if err := msg.UpdateSizeByte(oldSizeByte); err != nil {
				return nil, err
			}
			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			if err := msg.UpdateSizeByte(sizeByte); err != nil {
				return nil, err
			}
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) updateByteOrder(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toUpdateByteOrder()

	byteOrder := parsedReq.ByteOrder.parse()

	oldByteOrder := msg.ByteOrder()
	if byteOrder == oldByteOrder {
		return nil
	}

	msg.SetByteOrder(byteOrder)

	res.setUndo(
		func() (*acmelib.Message, error) {
			msg.SetByteOrder(oldByteOrder)
			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			msg.SetByteOrder(byteOrder)
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) updateCycleTime(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toUpdateCycleTime()

	cycleTime := parsedReq.CycleTime

	oldCycleTime := msg.CycleTime()
	if cycleTime == oldCycleTime {
		return nil
	}

	msg.SetCycleTime(cycleTime)

	res.setUndo(
		func() (*acmelib.Message, error) {
			msg.SetCycleTime(oldCycleTime)
			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			msg.SetCycleTime(cycleTime)
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) updateSendType(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toUpdateSendType()

	sendType := parsedReq.SendType.parse()

	oldSendType := msg.SendType()
	if sendType == oldSendType {
		return nil
	}

	msg.SetSendType(sendType)

	res.setUndo(
		func() (*acmelib.Message, error) {
			msg.SetSendType(oldSendType)
			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			msg.SetSendType(sendType)
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) updateDelayTime(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toUpdateDelayTime()

	delayTime := parsedReq.DelayTime

	oldDelayTime := msg.DelayTime()
	if delayTime == oldDelayTime {
		return nil
	}

	msg.SetDelayTime(delayTime)

	res.setUndo(
		func() (*acmelib.Message, error) {
			msg.SetDelayTime(oldDelayTime)
			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			msg.SetDelayTime(delayTime)
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) updateStartDelayTime(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toUpdateStartDelayTime()

	startDelatTime := parsedReq.StartDelayTime

	oldStartDelayTime := msg.StartDelayTime()
	if startDelatTime == oldStartDelayTime {
		return nil
	}

	msg.SetStartDelayTime(startDelatTime)

	res.setUndo(
		func() (*acmelib.Message, error) {
			msg.SetStartDelayTime(oldStartDelayTime)
			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			msg.SetStartDelayTime(startDelatTime)
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) addSignal(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toAddSignal()

	sigKind := parsedReq.SignalKind.parse()

	otherSignals := []entity{}

	type payloadHole struct {
		startPos int
		size     int
	}

	payloadHoles := []payloadHole{}
	currPos := 0

	for _, sig := range msg.Signals() {
		otherSignals = append(otherSignals, sig)

		if currPos < sig.GetStartBit() {
			payloadHoles = append(payloadHoles, payloadHole{startPos: currPos, size: sig.GetStartBit() - currPos})
		}

		currPos = sig.GetStartBit() + sig.GetSize()
	}
	sigName := getNextNewName("signal", otherSignals)

	if currPos < msg.SizeByte()*8 {
		payloadHoles = append(payloadHoles, payloadHole{startPos: currPos, size: msg.SizeByte()*8 - currPos})
	}

	if len(payloadHoles) == 0 {
		return errors.New("payload is full")
	}

	startPos := slices.MaxFunc(payloadHoles, func(a, b payloadHole) int {
		return a.size - b.size
	}).startPos

	var sig acmelib.Signal
	switch sigKind {
	case acmelib.SignalKindStandard:
		tmpSig, err := acmelib.NewStandardSignal(sigName, defaultSignalType)
		if err != nil {
			return err
		}
		sig = tmpSig

	case acmelib.SignalKindEnum:
		tmpSig, err := acmelib.NewEnumSignal(sigName, defaultSignalEnum)
		if err != nil {
			return err
		}
		sig = tmpSig

	case acmelib.SignalKindMultiplexer:
	}

	if err := msg.InsertSignal(sig, startPos); err != nil {
		return err
	}

	h.signalCtr.sendAdd(sig)

	res.setUndo(
		func() (*acmelib.Message, error) {
			if err := msg.RemoveSignal(sig.EntityID()); err != nil {
				return nil, err
			}

			h.signalCtr.sendDelete(sig)

			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			if err := msg.InsertSignal(sig, startPos); err != nil {
				return nil, err
			}

			h.signalCtr.sendAdd(sig)

			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) deleteSignals(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toDeleteSignals()

	if len(parsedReq.SignalEntityIDs) == 0 {
		return nil
	}

	remSigIDs := make(map[string]struct{})
	for _, sigID := range parsedReq.SignalEntityIDs {
		remSigIDs[sigID] = struct{}{}
	}

	remSignals := []acmelib.Signal{}
	remStartPos := make(map[string]int)
	for _, sig := range msg.Signals() {
		tmpID := sig.EntityID().String()

		if _, ok := remSigIDs[tmpID]; ok {
			remSignals = append(remSignals, sig)
			remStartPos[tmpID] = sig.GetStartBit()
		}
	}

	for _, sig := range remSignals {
		if err := msg.RemoveSignal(sig.EntityID()); err != nil {
			return err
		}
	}

	for _, sig := range remSignals {
		h.signalCtr.sendDelete(sig)
	}

	res.setUndo(
		func() (*acmelib.Message, error) {
			for _, sig := range remSignals {
				startPos, ok := remStartPos[sig.EntityID().String()]
				if !ok {
					continue
				}

				if err := msg.InsertSignal(sig, startPos); err != nil {
					return nil, err
				}
			}

			for _, sig := range remSignals {
				h.signalCtr.sendAdd(sig)
			}

			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			for _, sig := range remSignals {
				if err := msg.RemoveSignal(sig.EntityID()); err != nil {
					return nil, err
				}
			}

			for _, sig := range remSignals {
				h.signalCtr.sendDelete(sig)
			}

			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) compactSignals(msg *acmelib.Message, _ *request, res *messageRes) error {
	signals := []acmelib.Signal{}
	startPos := make(map[acmelib.EntityID]int)

	for _, sig := range msg.Signals() {
		signals = append(signals, sig)
		startPos[sig.EntityID()] = sig.GetStartBit()
	}

	msg.CompactSignals()

	res.setUndo(
		func() (*acmelib.Message, error) {
			msg.RemoveAllSignals()

			for _, sig := range signals {
				startPos, ok := startPos[sig.EntityID()]
				if !ok {
					continue
				}

				if err := msg.InsertSignal(sig, startPos); err != nil {
					return nil, err
				}
			}

			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			msg.CompactSignals()
			return msg, nil
		},
	)

	return nil
}

func (h *messageHandler) reorderSignal(msg *acmelib.Message, sig acmelib.Signal, from, to int) error {
	if err := msg.RemoveSignal(sig.EntityID()); err != nil {
		return err
	}

	otherSignals := msg.Signals()

	newStartPos := 0
	if to < from {
		// move up
		nearestSig := otherSignals[from-1]
		targetSig := otherSignals[to]

		newStartPos = targetSig.GetStartBit()
		offset := sig.GetSize() + sig.GetStartBit() - nearestSig.GetSize() - nearestSig.GetStartBit()

		for i := from - 1; i >= to; i = i - 1 {
			msg.ShiftSignalRight(otherSignals[i].EntityID(), offset)
		}
	} else {
		// move down
		nearestSig := otherSignals[from]
		targetSig := otherSignals[to-1]

		newStartPos = targetSig.GetStartBit() + targetSig.GetSize() - sig.GetSize()
		offset := nearestSig.GetStartBit() - sig.GetStartBit()

		for i := from; i < to; i++ {
			msg.ShiftSignalLeft(otherSignals[i].EntityID(), offset)
		}
	}

	if err := msg.InsertSignal(sig, newStartPos); err != nil {
		return err
	}

	return nil
}

func (h *messageHandler) reorderSignalHandler(msg *acmelib.Message, req *request, res *messageRes) error {
	parsedReq := req.toReorderSignal()

	from := parsedReq.From
	to := parsedReq.To

	if from == to {
		return nil
	}

	currSig, err := msg.GetSignal(acmelib.EntityID(parsedReq.SignalEntityID))
	if err != nil {
		return err
	}

	if err := h.reorderSignal(msg, currSig, from, to); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.Message, error) {
			if err := h.reorderSignal(msg, currSig, to, from); err != nil {
				return nil, err
			}
			return msg, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Message, error) {
			if err := h.reorderSignal(msg, currSig, from, to); err != nil {
				return nil, err
			}
			return msg, nil
		},
	)

	return nil
}
