package main

type request struct {
	data any
}

// newRequest returns a new generic request.
//
// IMPORTANT: the data attribute must be a pointer.
func newRequest(data any) *request {
	return &request{data}
}

/////////////////////
// COMMON REQUESTS //
/////////////////////

type commonCreateReq struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type UpdateNameReq struct {
	Name string `json:"name"`
}

func (r *request) toUpdateName() *UpdateNameReq {
	req, ok := r.data.(*UpdateNameReq)
	if !ok {
		panic("cannot convert to UpdateNameReq")
	}
	return req
}

type UpdateDescReq struct {
	Desc string `json:"desc"`
}

func (r *request) toUpdateDesc() *UpdateDescReq {
	req, ok := r.data.(*UpdateDescReq)
	if !ok {
		panic("cannot convert to UpdateDescReq")
	}
	return req
}

//////////////////
// BUS REQUESTS //
//////////////////

type CreateBusReq struct {
	commonCreateReq

	BusType  BusType `json:"busType"`
	Baudrate int     `json:"baudrate"`
}

func (r *request) toCreateBus() *CreateBusReq {
	req, ok := r.data.(*CreateBusReq)
	if !ok {
		panic("cannot convert to CreateBusReq")
	}
	return req
}

type UpdateBaudrateReq struct {
	Baudrate int `json:"baudrate"`
}

func (r *request) toUpdateBaudrate() *UpdateBaudrateReq {
	req, ok := r.data.(*UpdateBaudrateReq)
	if !ok {
		panic("cannot convert to UpdateBaudrateReq")
	}
	return req
}

type UpdateBusTypeReq struct {
	BusType BusType `json:"busType"`
}

func (r *request) toUpdateBusType() *UpdateBusTypeReq {
	req, ok := r.data.(*UpdateBusTypeReq)
	if !ok {
		panic("cannot convert to UpdateBusTypeReq")
	}
	return req
}

///////////////////
// NODE REQUESTS //
///////////////////

type CreateNodeReq struct {
	commonCreateReq

	NodeID         uint `json:"nodeId"`
	InterfaceCount int  `json:"interfaceCount"`
}

func (r *request) toCreateNode() *CreateNodeReq {
	req, ok := r.data.(*CreateNodeReq)
	if !ok {
		panic("cannot convert to CreateNodeReq")
	}
	return req
}

type UpdateNodeIDReq struct {
	NodeID uint `json:"nodeId"`
}

func (r *request) toUpdateNodeID() *UpdateNodeIDReq {
	req, ok := r.data.(*UpdateNodeIDReq)
	if !ok {
		panic("cannot convert to UpdateNodeIDReq")
	}
	return req
}

type UpdateAttachedBusReq struct {
	BusEntityID     string `json:"busEntityId"`
	InterfaceNumber int    `json:"interfaceNumber"`
}

func (r *request) toUpdateAttachedBus() *UpdateAttachedBusReq {
	req, ok := r.data.(*UpdateAttachedBusReq)
	if !ok {
		panic("cannot convert to UpdateAttachedBusReq")
	}
	return req
}

type RemoveSentMessagesReq struct {
	InterfaceNumber  int      `json:"interfaceNumber"`
	MessageEntityIDs []string `json:"messageEntityIds"`
}

func (r *request) toRemoveSentMessages() *RemoveSentMessagesReq {
	req, ok := r.data.(*RemoveSentMessagesReq)
	if !ok {
		panic("cannot convert to RemoveSentMessagesReq")
	}
	return req
}

type RemoveReceivedMessagesReq struct {
	InterfaceNumber  int      `json:"interfaceNumber"`
	MessageEntityIDs []string `json:"messageEntityIds"`
}

func (r *request) toRemoveReceivedMessages() *RemoveReceivedMessagesReq {
	req, ok := r.data.(*RemoveReceivedMessagesReq)
	if !ok {
		panic("cannot convert to RemoveReceivedMessagesReq")
	}
	return req
}

//////////////////////
// MESSAGE REQUESTS //
//////////////////////

type commondMessageSignalReq struct {
	SignalEntityID string `json:"signalEntityId"`
}

type UpdateMessageIDReq struct {
	MessageID uint `json:"messageId"`
}

func (r *request) toUpdateMessageID() *UpdateMessageIDReq {
	req, ok := r.data.(*UpdateMessageIDReq)
	if !ok {
		panic("cannot convert to UpdateMessageIDReq")
	}
	return req
}

type UpdateStaticCANIDReq struct {
	StaticCANID uint `json:"staticCanId"`
}

func (r *request) toUpdateStaticCANID() *UpdateStaticCANIDReq {
	req, ok := r.data.(*UpdateStaticCANIDReq)
	if !ok {
		panic("cannot convert to UpdateStaticCANIDReq")
	}
	return req
}

type UpdateSizeByteReq struct {
	SizeByte int `json:"sizeByte"`
}

func (r *request) toUpdateSizeByte() *UpdateSizeByteReq {
	req, ok := r.data.(*UpdateSizeByteReq)
	if !ok {
		panic("cannot convert to UpdateSizeByteReq")
	}
	return req
}

type UpdateByteOrderReq struct {
	ByteOrder MessageByteOrder `json:"byteOrder"`
}

func (r *request) toUpdateByteOrder() *UpdateByteOrderReq {
	req, ok := r.data.(*UpdateByteOrderReq)
	if !ok {
		panic("cannot convert to UpdateByteOrderReq")
	}
	return req
}

type UpdateSendTypeReq struct {
	SendType MessageSendType `json:"sendType"`
}

func (r *request) toUpdateSendType() *UpdateSendTypeReq {
	req, ok := r.data.(*UpdateSendTypeReq)
	if !ok {
		panic("cannot convert to UpdateSendTypeReq")
	}
	return req
}

type UpdateCycleTimeReq struct {
	CycleTime int `json:"cycleTime"`
}

func (r *request) toUpdateCycleTime() *UpdateCycleTimeReq {
	req, ok := r.data.(*UpdateCycleTimeReq)
	if !ok {
		panic("cannot convert to UpdateCycleTimeReq")
	}
	return req
}

type UpdateDelayTimeReq struct {
	DelayTime int `json:"delayTime"`
}

func (r *request) toUpdateDelayTime() *UpdateDelayTimeReq {
	req, ok := r.data.(*UpdateDelayTimeReq)
	if !ok {
		panic("cannot convert to UpdateDelayTimeReq")
	}
	return req
}

type UpdateStartDelayTimeReq struct {
	StartDelayTime int `json:"startDelayTime"`
}

func (r *request) toUpdateStartDelayTime() *UpdateStartDelayTimeReq {
	req, ok := r.data.(*UpdateStartDelayTimeReq)
	if !ok {
		panic("cannot convert to UpdateStartDelayTimeReq")
	}
	return req
}

type AddSignalReq struct {
	SignalKind SignalKind `json:"signalKind"`
}

func (r *request) toAddSignal() *AddSignalReq {
	req, ok := r.data.(*AddSignalReq)
	if !ok {
		panic("cannot convert to AddSignalReq")
	}
	return req
}

type DeleteSignalsReq struct {
	SignalEntityIDs []string `json:"signalEntityIds"`
}

func (r *request) toDeleteSignals() *DeleteSignalsReq {
	req, ok := r.data.(*DeleteSignalsReq)
	if !ok {
		panic("cannot convert to DeleteSignalsReq")
	}
	return req
}

type ReorderSignalReq struct {
	commondMessageSignalReq

	From int `json:"from"`
	To   int `json:"to"`
}

func (r *request) toReorderSignal() *ReorderSignalReq {
	req, ok := r.data.(*ReorderSignalReq)
	if !ok {
		panic("cannot convert to ReorderSignalReq")
	}
	return req
}

/////////////////////
// SIGNAL REQUESTS //
/////////////////////

type UpdateSignalTypeReq struct {
	SignalTypeEntityID string `json:"signalTypeEntityId"`
}

func (r *request) toUpdateSignalType() *UpdateSignalTypeReq {
	req, ok := r.data.(*UpdateSignalTypeReq)
	if !ok {
		panic("cannot convert to UpdateSignalTypeReq")
	}
	return req
}

type UpdateSignalUnitReq struct {
	SignalUnitEntityID string `json:"signalUnitEntityId"`
}

func (r *request) toUpdateSignalUnit() *UpdateSignalUnitReq {
	req, ok := r.data.(*UpdateSignalUnitReq)
	if !ok {
		panic("cannot convert to UpdateSignalUnitReq")
	}
	return req
}

type UpdateSignalEnumReq struct {
	SignalEnumEntityID string `json:"signalEnumEntityId"`
}

func (r *request) toUpdateSignalEnum() *UpdateSignalEnumReq {
	req, ok := r.data.(*UpdateSignalEnumReq)
	if !ok {
		panic("cannot convert to UpdateSignalEnumReq")
	}
	return req
}

//////////////////////////
// SIGNAL TYPE REQUESTS //
//////////////////////////

type CreateSignalTypeReq struct {
	commonCreateReq

	Kind   SignalTypeKind `json:"kind"`
	Size   int            `json:"size"`
	Signed bool           `json:"signed"`
	Min    float64        `json:"min"`
	Max    float64        `json:"max"`
	Scale  float64        `json:"scale"`
	Offset float64        `json:"offset"`
}

func (r *request) toCreateSignalType() *CreateSignalTypeReq {
	req, ok := r.data.(*CreateSignalTypeReq)
	if !ok {
		panic("cannot convert to CreateSignalTypeReq")
	}
	return req
}

type UpdateMinReq struct {
	Min float64 `json:"min"`
}

func (r *request) toUpdateMin() *UpdateMinReq {
	req, ok := r.data.(*UpdateMinReq)
	if !ok {
		panic("cannot convert to UpdateMinReq")
	}
	return req
}

type UpdateMaxReq struct {
	Max float64 `json:"max"`
}

func (r *request) toUpdateMax() *UpdateMaxReq {
	req, ok := r.data.(*UpdateMaxReq)
	if !ok {
		panic("cannot convert to UpdateMaxReq")
	}
	return req
}

type UpdateScaleReq struct {
	Scale float64 `json:"scale"`
}

func (r *request) toUpdateScale() *UpdateScaleReq {
	req, ok := r.data.(*UpdateScaleReq)
	if !ok {
		panic("cannot convert to UpdateScaleReq")
	}
	return req
}

type UpdateOffsetReq struct {
	Offset float64 `json:"offset"`
}

func (r *request) toUpdateOffset() *UpdateOffsetReq {
	req, ok := r.data.(*UpdateOffsetReq)
	if !ok {
		panic("cannot convert to UpdateOffsetReq")
	}
	return req
}

//////////////////////////
// SIGNAL UNIT REQUESTS //
//////////////////////////

type CreateSignalUnitReq struct {
	commonCreateReq

	Kind   SignalUnitKind `json:"kind"`
	Symbol string         `json:"symbol"`
}

func (r *request) toCreateSignalUnit() *CreateSignalUnitReq {
	req, ok := r.data.(*CreateSignalUnitReq)
	if !ok {
		panic("cannot convert to CreateSignalUnitReq")
	}
	return req
}

type UpdateSignalUnitKindReq struct {
	Kind SignalUnitKind `json:"kind"`
}

func (r *request) toUpdateSignalUnitKind() *UpdateSignalUnitKindReq {
	req, ok := r.data.(*UpdateSignalUnitKindReq)
	if !ok {
		panic("cannot convert to UpdateSignalUnitKindReq")
	}
	return req
}

type UpdateSymbolReq struct {
	Symbol string `json:"symbol"`
}

func (r *request) toUpdateSymbol() *UpdateSymbolReq {
	req, ok := r.data.(*UpdateSymbolReq)
	if !ok {
		panic("cannot convert to UpdateSymbolReq")
	}
	return req
}

//////////////////////////
// SIGNAL ENUM REQUESTS //
//////////////////////////

type CreateSignalEnumReq struct {
	commonCreateReq

	MinSize int `json:"minSize"`
}

func (r *request) toCreateSignalEnum() *CreateSignalEnumReq {
	req, ok := r.data.(*CreateSignalEnumReq)
	if !ok {
		panic("cannot convert to CreateSignalEnumReq")
	}
	return req
}

type commonSignalEnumValueReq struct {
	ValueEntityID string `json:"valueEntityId"`
}

type RemoveValuesReq struct {
	ValueEntityIDs []string `json:"valueEntityIds"`
}

func (r *request) toRemoveValues() *RemoveValuesReq {
	req, ok := r.data.(*RemoveValuesReq)
	if !ok {
		panic("cannot convert to RemoveValuesReq")
	}
	return req
}

type ReorderValueReq struct {
	commonSignalEnumValueReq

	From int `json:"from"`
	To   int `json:"to"`
}

func (r *request) toReorderValue() *ReorderValueReq {
	req, ok := r.data.(*ReorderValueReq)
	if !ok {
		panic("cannot convert to ReorderValueReq")
	}
	return req
}

type UpdateValueNameReq struct {
	commonSignalEnumValueReq
	UpdateNameReq
}

func (r *request) toUpdateValueName() *UpdateValueNameReq {
	req, ok := r.data.(*UpdateValueNameReq)
	if !ok {
		panic("cannot convert to UpdateValueNameReq")
	}
	return req
}

type UpdateValueDescReq struct {
	commonSignalEnumValueReq
	UpdateDescReq
}

func (r *request) toUpdateValueDesc() *UpdateValueDescReq {
	req, ok := r.data.(*UpdateValueDescReq)
	if !ok {
		panic("cannot convert to UpdateValueDescReq")
	}
	return req
}

type UpdateValueIndexReq struct {
	commonSignalEnumValueReq

	Index int `json:"index"`
}

func (r *request) toUpdateValueIndex() *UpdateValueIndexReq {
	req, ok := r.data.(*UpdateValueIndexReq)
	if !ok {
		panic("cannot convert to UpdateValueIndexReq")
	}
	return req
}
