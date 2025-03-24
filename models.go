package main

import (
	"slices"
	"time"

	"github.com/squadracorsepolito/acmelib"
)

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

///////////////////
// COMMON MODELS //
///////////////////

type EntityKind string

const (
	EntityKindNetwork      EntityKind = "network"
	EntityKindBus          EntityKind = "bus"
	EntityKindNode         EntityKind = "node"
	EntityKindMessage      EntityKind = "message"
	EntityKindSignal       EntityKind = "signal"
	EntityKindSignalType   EntityKind = "signal-type"
	EntityKindSignalUnit   EntityKind = "signal-unit"
	EntityKindSignalEnum   EntityKind = "signal-enum"
	EntityKindCANIDBuilder EntityKind = "can-id-builder"
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
	case acmelib.EntityKindCANIDBuilder:
		return EntityKindCANIDBuilder
	default:
		return EntityKindNetwork
	}
}

type entity interface {
	EntityKind() acmelib.EntityKind
	EntityID() acmelib.EntityID
	Name() string
	Desc() string
	CreateTime() time.Time
}

type BaseEntity struct {
	EntityKind EntityKind `json:"entityKind"`
	EntityID   string     `json:"entityId"`
	Name       string     `json:"name"`
	Desc       string     `json:"desc"`
	CreateTime time.Time  `json:"createTime"`
}

func newBaseEntity(e entity) BaseEntity {
	return BaseEntity{
		EntityKind: newEntityKind(e.EntityKind()),
		EntityID:   e.EntityID().String(),
		Name:       e.Name(),
		Desc:       e.Desc(),
		CreateTime: e.CreateTime(),
	}
}

////////////////////
// NETWORK MODELS //
////////////////////

type Network struct {
	BaseEntity

	Buses []BusBase `json:"buses"`
}

func newNetwork(net *acmelib.Network) Network {
	if net == nil {
		return Network{}
	}

	res := Network{
		BaseEntity: newBaseEntity(net),

		Buses: []BusBase{},
	}

	for _, bus := range net.Buses() {
		res.Buses = append(res.Buses, newBusBase(bus))
	}

	return res
}

/////////////////
// NODE MODELS //
/////////////////

type BaseNode struct {
	BaseEntity
}

func newBaseNode(node *acmelib.Node) BaseNode {
	if node == nil {
		return BaseNode{}
	}

	return BaseNode{
		BaseEntity: newBaseEntity(node),
	}
}

type NodeInterface struct {
	Number           int          `json:"number"`
	AttachedBus      BaseEntity   `json:"attachedBus"`
	SentMessages     []BaseEntity `json:"sentMessages"`
	ReceivedMessages []BaseEntity `json:"receivedMessages"`
}

func newNodeInterface(nodeInt *acmelib.NodeInterface) NodeInterface {
	sentMessages := []BaseEntity{}
	for _, tmpMsg := range nodeInt.SentMessages() {
		sentMessages = append(sentMessages, newBaseEntity(tmpMsg))
	}

	receivedMessages := []BaseEntity{}
	for _, tmpMsg := range nodeInt.ReceivedMessages() {
		receivedMessages = append(receivedMessages, newBaseEntity(tmpMsg))
	}

	res := NodeInterface{
		Number:           nodeInt.Number(),
		SentMessages:     sentMessages,
		ReceivedMessages: receivedMessages,
	}

	if nodeInt.ParentBus() != nil {
		res.AttachedBus = newBaseEntity(nodeInt.ParentBus())
	}

	return res
}

type Node struct {
	BaseEntity

	ID         uint            `json:"id"`
	Interfaces []NodeInterface `json:"interfaces"`
}

func newNode(node *acmelib.Node) Node {
	if node == nil {
		return Node{}
	}

	res := Node{
		BaseEntity: newBaseEntity(node),

		ID:         uint(node.ID()),
		Interfaces: []NodeInterface{},
	}

	for _, nodeInt := range node.Interfaces() {
		res.Interfaces = append(res.Interfaces, newNodeInterface(nodeInt))
	}

	return res
}

////////////////////
// MESSAGE MODELS //
////////////////////

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

type MessagePriority string

const (
	MessagePriorityVeryHigh MessagePriority = "very-high"
	MessagePriorityHigh     MessagePriority = "high"
	MessagePriorityMedium   MessagePriority = "medium"
	MessagePriorityLow      MessagePriority = "low"
)

func newMessagePriority(priority acmelib.MessagePriority) MessagePriority {
	switch priority {
	case acmelib.MessagePriorityVeryHigh:
		return MessagePriorityVeryHigh
	case acmelib.MessagePriorityHigh:
		return MessagePriorityHigh
	case acmelib.MessagePriorityMedium:
		return MessagePriorityMedium
	case acmelib.MessagePriorityLow:
		return MessagePriorityLow
	default:
		return MessagePriorityVeryHigh
	}
}

func (mp MessagePriority) parse() acmelib.MessagePriority {
	switch mp {
	case MessagePriorityVeryHigh:
		return acmelib.MessagePriorityVeryHigh
	case MessagePriorityHigh:
		return acmelib.MessagePriorityHigh
	case MessagePriorityMedium:
		return acmelib.MessagePriorityMedium
	case MessagePriorityLow:
		return acmelib.MessagePriorityLow
	default:
		return acmelib.MessagePriorityVeryHigh
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

type BaseMessage struct {
	BaseEntity

	SizeByte int `json:"sizeByte"`
}

func newBaseMessage(msg *acmelib.Message) BaseMessage {
	if msg == nil {
		return BaseMessage{}
	}

	return BaseMessage{
		BaseEntity: newBaseEntity(msg),

		SizeByte: msg.SizeByte(),
	}
}

type Message struct {
	BaseMessage

	Paths []EntityPath `json:"paths"`

	HasStaticCANID bool `json:"hasStaticCANID"`
	ID             uint `json:"id"`
	CANID          uint `json:"canId"`

	AvailableTrailingBytes int              `json:"availableTrailingBytes"`
	MaxAvailableSpace      int              `json:"maxAvailableSpace"`
	ByteOrder              MessageByteOrder `json:"byteOrder"`

	Priority       MessagePriority `json:"priority"`
	CycleTime      int             `json:"cycleTime"`
	SendType       MessageSendType `json:"sendType"`
	DelayTime      int             `json:"delayTime"`
	StartDelayTime int             `json:"startDelayTime"`

	Signals []BaseSignal `json:"signals"`

	Receivers []BaseNode `json:"receivers"`

	SenderNode BaseEntity `json:"senderNode"`
	ParentBus  BaseEntity `json:"parentBus"`
}

func newMessage(msg *acmelib.Message) Message {
	if msg == nil {
		return Message{}
	}

	res := Message{
		BaseMessage: newBaseMessage(msg),

		Paths: newMessageEntityPaths(msg),

		HasStaticCANID: msg.HasStaticCANID(),
		ID:             uint(msg.ID()),
		CANID:          uint(msg.GetCANID()),

		AvailableTrailingBytes: msg.SizeByte(),
		MaxAvailableSpace:      0,
		ByteOrder:              newMessageByteOrder(msg.ByteOrder()),

		Priority:       newMessagePriority(msg.Priority()),
		CycleTime:      msg.CycleTime(),
		SendType:       newMessageSendType(msg.SendType()),
		DelayTime:      msg.DelayTime(),
		StartDelayTime: msg.StartDelayTime(),

		Signals: []BaseSignal{},

		Receivers: []BaseNode{},
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
		res.Signals = append(res.Signals, newBaseSignal(sig))

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

	for _, receiver := range msg.Receivers() {
		res.Receivers = append(res.Receivers, newBaseNode(receiver.Node()))
	}

	return res
}

///////////////////
// SIGNAL MODELS //
///////////////////

type SignalKind string

const (
	SignalKindStandard    SignalKind = "standard"
	SignalKindEnum        SignalKind = "enum"
	SignalKindMultiplexer SignalKind = "multiplexer"
)

func newSignalKind(kind acmelib.SignalKind) SignalKind {
	switch kind {
	case acmelib.SignalKindStandard:
		return SignalKindStandard
	case acmelib.SignalKindEnum:
		return SignalKindEnum
	case acmelib.SignalKindMultiplexer:
		return SignalKindMultiplexer
	default:
		return SignalKindStandard
	}
}

func (sk SignalKind) parse() acmelib.SignalKind {
	switch sk {
	case SignalKindStandard:
		return acmelib.SignalKindStandard
	case SignalKindEnum:
		return acmelib.SignalKindEnum
	case SignalKindMultiplexer:
		return acmelib.SignalKindMultiplexer
	default:
		return acmelib.SignalKindStandard
	}
}

type BaseSignal struct {
	BaseEntity

	Kind             SignalKind `json:"kind"`
	StartPos         int        `json:"startPos"`
	RelativeStartPos int        `json:"relativeStartPos"`
	Size             int        `json:"size"`
}

func newBaseSignal(sig acmelib.Signal) BaseSignal {
	if sig == nil {
		return BaseSignal{}
	}

	return BaseSignal{
		BaseEntity: newBaseEntity(sig),

		Kind:             newSignalKind(sig.Kind()),
		StartPos:         sig.GetStartBit(),
		RelativeStartPos: sig.GetRelativeStartPos(),
		Size:             sig.GetSize(),
	}
}

type StandardSignal struct {
	SignalType SignalTypeBrief `json:"signalType"`
	SignalUnit BaseEntity      `json:"signalUnit"`
}

func newStandardSignal(stdSig *acmelib.StandardSignal) StandardSignal {
	res := StandardSignal{
		SignalType: newSignalTypeBrief(stdSig.Type()),
	}

	if stdSig.Unit() != nil {
		res.SignalUnit = newBaseEntity(stdSig.Unit())
	}

	return res
}

type EnumSignal struct {
	SignalEnum SignalEnumBrief `json:"signalEnum"`
}

func newEnumSignal(enumSig *acmelib.EnumSignal) EnumSignal {
	return EnumSignal{
		SignalEnum: newSignalEnumBrief(enumSig.Enum()),
	}
}

type MultiplexerSignalGroup struct {
	ID      int          `json:"id"`
	Signals []BaseSignal `json:"signals"`
}

func newMultiplexerSignalGroup(groupID int, group []acmelib.Signal) MultiplexerSignalGroup {
	res := MultiplexerSignalGroup{
		ID:      groupID,
		Signals: []BaseSignal{},
	}

	for _, sig := range group {
		res.Signals = append(res.Signals, newBaseSignal(sig))
	}

	return res
}

type MultiplexerSignal struct {
	GroupCount     int                      `json:"groupCount"`
	GroupCountSize int                      `json:"groupCountSize"`
	GroupSize      int                      `json:"groupSize"`
	Groups         []MultiplexerSignalGroup `json:"groups"`
	EmptyGroups    []int                    `json:"emptyGroups"`
}

func newMultiplexedSignal(muxSig *acmelib.MultiplexerSignal) MultiplexerSignal {
	if muxSig == nil {
		return MultiplexerSignal{}
	}

	res := MultiplexerSignal{
		GroupCount:     muxSig.GroupCount(),
		GroupCountSize: muxSig.GetGroupCountSize(),
		GroupSize:      muxSig.GroupSize(),
		Groups:         []MultiplexerSignalGroup{},
		EmptyGroups:    []int{},
	}

	for groupID, group := range muxSig.GetSignalGroups() {
		if len(group) == 0 {
			res.EmptyGroups = append(res.EmptyGroups, groupID)
			continue
		}

		res.Groups = append(res.Groups, newMultiplexerSignalGroup(groupID, group))
	}

	return res
}

type Signal struct {
	BaseSignal

	Paths []EntityPath `json:"paths"`

	ParentMessage BaseMessage `json:"parentMessage"`

	Standard    StandardSignal    `json:"standard"`
	Enum        EnumSignal        `json:"enum"`
	Multiplexer MultiplexerSignal `json:"multiplexer"`
}

func newSignal(sig acmelib.Signal) Signal {
	res := Signal{
		BaseSignal: newBaseSignal(sig),

		Paths: newSignalEntityPaths(sig),
	}

	parMsg := sig.ParentMessage()
	if parMsg != nil {
		res.ParentMessage = newBaseMessage(parMsg)
	}

	switch sig.Kind() {
	case acmelib.SignalKindStandard:
		stdSig, err := sig.ToStandard()
		if err != nil {
			panic(err)
		}
		res.Standard = newStandardSignal(stdSig)

	case acmelib.SignalKindEnum:
		enumSig, err := sig.ToEnum()
		if err != nil {
			panic(err)
		}
		res.Enum = newEnumSignal(enumSig)

	case acmelib.SignalKindMultiplexer:
		muxSig, err := sig.ToMultiplexer()
		if err != nil {
			panic(err)
		}
		res.Multiplexer = newMultiplexedSignal(muxSig)
	}

	return res
}

////////////////////////
// SIGNAL TYPE MODELS //
////////////////////////

type SignalTypeKind string

const (
	SignalTypeKindCustom  SignalTypeKind = "custom"
	SignalTypeKindFlag    SignalTypeKind = "flag"
	SignalTypeKindInteger SignalTypeKind = "integer"
	SignalTypeKindDecimal SignalTypeKind = "decimal"
)

func newSignalTypeKind(kind acmelib.SignalTypeKind) SignalTypeKind {
	return SignalTypeKind(kind.String())
}

type SignalTypeBrief struct {
	BaseEntity

	Kind SignalTypeKind `json:"kind"`
	Size int            `json:"size"`
}

func newSignalTypeBrief(sigType *acmelib.SignalType) SignalTypeBrief {
	return SignalTypeBrief{
		BaseEntity: newBaseEntity(sigType),

		Kind: newSignalTypeKind(sigType.Kind()),
		Size: sigType.Size(),
	}
}

type SignalType struct {
	BaseEntity

	Kind   SignalTypeKind `json:"kind"`
	Size   int            `json:"size"`
	Signed bool           `json:"signed"`
	Min    float64        `json:"min"`
	Max    float64        `json:"max"`
	Scale  float64        `json:"scale"`
	Offset float64        `json:"offset"`

	ReferenceCount int         `json:"referenceCount"`
	References     []Reference `json:"references"`
}

func newSignalType(sigType *acmelib.SignalType) SignalType {
	if sigType == nil {
		return SignalType{}
	}

	refCount := sigType.ReferenceCount()

	res := SignalType{
		BaseEntity: newBaseEntity(sigType),

		Kind:   newSignalTypeKind(sigType.Kind()),
		Size:   int(sigType.Size()),
		Signed: sigType.Signed(),
		Min:    sigType.Min(),
		Max:    sigType.Max(),
		Scale:  sigType.Scale(),
		Offset: sigType.Offset(),

		ReferenceCount: refCount,
	}

	if refCount == 0 {
		return res
	}

	rootRefs := []*reference{}
	refs := make(map[acmelib.EntityID]*reference)
	for _, sig := range sigType.References() {
		sigRef := newReference(sig)
		refs[sig.EntityID()] = sigRef

		var msgRef *reference
		msg := sig.ParentMessage()
		msgRef, ok := refs[msg.EntityID()]
		if !ok {
			msgRef = newReference(msg)
			refs[msg.EntityID()] = msgRef
		}
		msgRef.addChild(sigRef)

		var nodeRef *reference
		node := msg.SenderNodeInterface().Node()
		nodeRef, ok = refs[node.EntityID()]
		if !ok {
			nodeRef = newReference(node)
			refs[node.EntityID()] = nodeRef
		}
		nodeRef.addChild(msgRef)

		var busRef *reference
		bus := msg.SenderNodeInterface().ParentBus()
		busRef, ok = refs[bus.EntityID()]
		if !ok {
			busRef = newReference(bus)
			refs[bus.EntityID()] = busRef
			rootRefs = append(rootRefs, busRef)
		}
		busRef.addChild(nodeRef)
	}

	for _, tmpRef := range rootRefs {
		res.References = append(res.References, tmpRef.toResponse())
	}

	return res
}

////////////////////////
// SIGNAL UNIT MODELS //
////////////////////////

type SignalUnitKind string

const (
	SignalUnitKindCustom      SignalUnitKind = "custom"
	SignalUnitKindTemperature SignalUnitKind = "temperature"
	SignalUnitKindElectrical  SignalUnitKind = "electrical"
	SignalUnitKindPower       SignalUnitKind = "power"
)

func newSignalUnitKind(kind acmelib.SignalUnitKind) SignalUnitKind {
	switch kind {
	case acmelib.SignalUnitKindCustom:
		return SignalUnitKindCustom
	case acmelib.SignalUnitKindTemperature:
		return SignalUnitKindTemperature
	case acmelib.SignalUnitKindElectrical:
		return SignalUnitKindElectrical
	case acmelib.SignalUnitKindPower:
		return SignalUnitKindPower
	default:
		return SignalUnitKindCustom
	}
}

func (k SignalUnitKind) parse() acmelib.SignalUnitKind {
	switch k {
	case SignalUnitKindCustom:
		return acmelib.SignalUnitKindCustom
	case SignalUnitKindTemperature:
		return acmelib.SignalUnitKindTemperature
	case SignalUnitKindElectrical:
		return acmelib.SignalUnitKindElectrical
	case SignalUnitKindPower:
		return acmelib.SignalUnitKindPower
	default:
		return acmelib.SignalUnitKindCustom
	}
}

type SignalUnitBrief struct {
	BaseEntity

	Kind SignalUnitKind `json:"kind"`
}

func newSignalUnitBrief(sigUnit *acmelib.SignalUnit) SignalUnitBrief {
	return SignalUnitBrief{
		BaseEntity: newBaseEntity(sigUnit),

		Kind: newSignalUnitKind(sigUnit.Kind()),
	}
}

type SignalUnit struct {
	BaseEntity

	Kind   SignalUnitKind `json:"kind"`
	Symbol string         `json:"symbol"`

	ReferenceCount int         `json:"referenceCount"`
	References     []Reference `json:"references"`
}

func newSignalUnit(sigUnit *acmelib.SignalUnit) SignalUnit {
	refCount := sigUnit.ReferenceCount()

	res := SignalUnit{
		BaseEntity: newBaseEntity(sigUnit),

		Kind:   newSignalUnitKind(sigUnit.Kind()),
		Symbol: sigUnit.Symbol(),

		ReferenceCount: refCount,
	}

	if refCount == 0 {
		return res
	}

	rootRefs := []*reference{}
	refs := make(map[acmelib.EntityID]*reference)
	for _, stdSig := range sigUnit.References() {
		sigRef := newReference(stdSig)
		refs[stdSig.EntityID()] = sigRef

		var msgRef *reference
		msg := stdSig.ParentMessage()
		msgRef, ok := refs[msg.EntityID()]
		if !ok {
			msgRef = newReference(msg)
			refs[msg.EntityID()] = msgRef
		}
		msgRef.addChild(sigRef)

		var nodeRef *reference
		node := msg.SenderNodeInterface().Node()
		nodeRef, ok = refs[node.EntityID()]
		if !ok {
			nodeRef = newReference(node)
			refs[node.EntityID()] = nodeRef
		}
		nodeRef.addChild(msgRef)

		var busRef *reference
		bus := msg.SenderNodeInterface().ParentBus()
		busRef, ok = refs[bus.EntityID()]
		if !ok {
			busRef = newReference(bus)
			refs[bus.EntityID()] = busRef
			rootRefs = append(rootRefs, busRef)
		}
		busRef.addChild(nodeRef)
	}

	for _, tmpRef := range rootRefs {
		res.References = append(res.References, tmpRef.toResponse())
	}

	return res
}

////////////////////////
// SIGNAL ENUM MODELS //
////////////////////////

type SignalEnumValue struct {
	BaseEntity

	Index int `json:"index"`
}

func newSignalEnumValue(sigEnumValue *acmelib.SignalEnumValue) SignalEnumValue {
	if sigEnumValue == nil {
		return SignalEnumValue{}
	}

	return SignalEnumValue{
		BaseEntity: newBaseEntity(sigEnumValue),

		Index: sigEnumValue.Index(),
	}
}

type SignalEnumBrief struct {
	BaseEntity

	Size int `json:"size"`
}

func newSignalEnumBrief(sigEnum *acmelib.SignalEnum) SignalEnumBrief {
	return SignalEnumBrief{
		BaseEntity: newBaseEntity(sigEnum),

		Size: sigEnum.GetSize(),
	}
}

type SignalEnum struct {
	BaseEntity

	Size     int               `json:"size"`
	MinSize  int               `json:"minSize"`
	MaxIndex int               `json:"maxIndex"`
	Values   []SignalEnumValue `json:"values"`

	References []Reference `json:"references"`
}

func newSignalEnum(sigEnum *acmelib.SignalEnum) SignalEnum {
	if sigEnum == nil {
		return SignalEnum{}
	}

	values := []SignalEnumValue{}
	for _, val := range sigEnum.Values() {
		values = append(values, newSignalEnumValue(val))
	}

	res := SignalEnum{
		BaseEntity: newBaseEntity(sigEnum),

		Size:     sigEnum.GetSize(),
		MinSize:  sigEnum.MinSize(),
		MaxIndex: sigEnum.MaxIndex(),
		Values:   values,
	}

	if len(sigEnum.References()) == 0 {
		return res
	}

	rootRefs := []*reference{}
	refs := make(map[acmelib.EntityID]*reference)
	for _, sig := range sigEnum.References() {
		sigRef := newReference(sig)
		refs[sig.EntityID()] = sigRef

		var msgRef *reference
		msg := sig.ParentMessage()
		msgRef, ok := refs[msg.EntityID()]
		if !ok {
			msgRef = newReference(msg)
			refs[msg.EntityID()] = msgRef
		}
		msgRef.addChild(sigRef)

		var nodeRef *reference
		node := msg.SenderNodeInterface().Node()
		nodeRef, ok = refs[node.EntityID()]
		if !ok {
			nodeRef = newReference(node)
			refs[node.EntityID()] = nodeRef
		}
		nodeRef.addChild(msgRef)

		var busRef *reference
		bus := msg.SenderNodeInterface().ParentBus()
		busRef, ok = refs[bus.EntityID()]
		if !ok {
			busRef = newReference(bus)
			refs[bus.EntityID()] = busRef
			rootRefs = append(rootRefs, busRef)
		}
		busRef.addChild(nodeRef)
	}

	for _, tmpRef := range rootRefs {
		res.References = append(res.References, tmpRef.toResponse())
	}

	return res
}

///////////////////////////
// CAN ID BUILDER MODELS //
///////////////////////////

type CANIDBuilderOpKind string

const (
	CANIDBuilderOpKindMessagePriority CANIDBuilderOpKind = "message-priority"
	CANIDBuilderOpKindMessageID       CANIDBuilderOpKind = "message-id"
	CANIDBuilderOpKindNodeID          CANIDBuilderOpKind = "node-id"
	CANIDBuilderOpKindBitMask         CANIDBuilderOpKind = "bit-mask"
)

func newCANIDBuilderOpKind(kind acmelib.CANIDBuilderOpKind) CANIDBuilderOpKind {
	switch kind {
	case acmelib.CANIDBuilderOpKindMessagePriority:
		return CANIDBuilderOpKindMessagePriority
	case acmelib.CANIDBuilderOpKindMessageID:
		return CANIDBuilderOpKindMessageID
	case acmelib.CANIDBuilderOpKindNodeID:
		return CANIDBuilderOpKindNodeID
	case acmelib.CANIDBuilderOpKindBitMask:
		return CANIDBuilderOpKindBitMask
	default:
		return CANIDBuilderOpKindMessageID
	}
}

func (bok CANIDBuilderOpKind) parse() acmelib.CANIDBuilderOpKind {
	switch bok {
	case CANIDBuilderOpKindMessagePriority:
		return acmelib.CANIDBuilderOpKindMessagePriority
	case CANIDBuilderOpKindMessageID:
		return acmelib.CANIDBuilderOpKindMessageID
	case CANIDBuilderOpKindNodeID:
		return acmelib.CANIDBuilderOpKindNodeID
	case CANIDBuilderOpKindBitMask:
		return acmelib.CANIDBuilderOpKindBitMask
	default:
		return acmelib.CANIDBuilderOpKindMessageID
	}
}

type CANIDBuilderOp struct {
	Kind CANIDBuilderOpKind `json:"kind"`
	From int                `json:"from"`
	Len  int                `json:"len"`
}

func newCANIDBuilderOp(op *acmelib.CANIDBuilderOp) CANIDBuilderOp {
	if op == nil {
		return CANIDBuilderOp{}
	}

	return CANIDBuilderOp{
		Kind: newCANIDBuilderOpKind(op.Kind()),
		From: op.From(),
		Len:  op.Len(),
	}
}

type CANIDBuilder struct {
	BaseEntity

	Operations []CANIDBuilderOp `json:"operations"`
}

func newCANIDBuilder(builder *acmelib.CANIDBuilder) CANIDBuilder {
	if builder == nil {
		return CANIDBuilder{}
	}

	res := CANIDBuilder{
		BaseEntity: newBaseEntity(builder),

		Operations: []CANIDBuilderOp{},
	}

	for _, op := range builder.Operations() {
		res.Operations = append(res.Operations, newCANIDBuilderOp(op))
	}

	return res
}
