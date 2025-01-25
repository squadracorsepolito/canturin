package main

import (
	"strings"

	"github.com/squadracorsepolito/acmelib"
)

type SignalKind string

const (
	SignalKindStandard    SignalKind = "standard"
	SignalKindEnum        SignalKind = "enum"
	SignalKindMultiplexed SignalKind = "multiplexed"
)

func newSignalKind(kind acmelib.SignalKind) SignalKind {
	switch kind {
	case acmelib.SignalKindStandard:
		return SignalKindStandard
	case acmelib.SignalKindEnum:
		return SignalKindEnum
	case acmelib.SignalKindMultiplexer:
		return SignalKindMultiplexed
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
	case SignalKindMultiplexed:
		return acmelib.SignalKindMultiplexer
	default:
		return acmelib.SignalKindStandard
	}
}

type StandardSignal struct {
	SignalType SignalTypeBrief `json:"signalType"`
	SignalUnit BaseEntity      `json:"signalUnit"`
}

func newStandardSignal(sig *acmelib.StandardSignal) StandardSignal {
	res := StandardSignal{
		SignalType: newSignalTypeBrief(sig.Type()),
	}

	if sig.Unit() != nil {
		res.SignalUnit = newBaseEntity(sig.Unit())
	}

	return res
}

type Signal struct {
	base

	Paths []EntityPath `json:"paths"`

	ParentMessage BaseEntity `json:"parentMessage"`

	Kind     SignalKind `json:"kind"`
	StartPos int        `json:"startPos"`
	Size     int        `json:"size"`

	Standard StandardSignal `json:"standard"`
}

func newSignal(sig acmelib.Signal) Signal {
	res := Signal{
		base: getBase(sig),

		Paths: newSignalEntityPaths(sig),

		Kind:     newSignalKind(sig.Kind()),
		StartPos: sig.GetStartBit(),
		Size:     sig.GetSize(),
	}

	parMsg := sig.ParentMessage()
	if parMsg != nil {
		res.ParentMessage = newBaseEntity(parMsg)
	}

	switch sig.Kind() {
	case acmelib.SignalKindStandard:
		stdSig, err := sig.ToStandard()
		if err != nil {
			panic(err)
		}
		res.Standard = newStandardSignal(stdSig)
	}

	return res
}

type SignalService struct {
	*service[acmelib.Signal, Signal, *signalHandler]
}

func newSignalService(sidebar *sidebarController, sigTypeSrv *SignalTypeService, sigUnitSrv *SignalUnitService) *SignalService {
	return &SignalService{
		service: newService(serviceKindSignal, newSignalHandler(sidebar, sigTypeSrv, sigUnitSrv), sidebar),
	}
}

func (s *SignalService) GetInvalidNames(entityID string) []string {
	s.mux.RLock()
	defer s.mux.RUnlock()

	names := []string{}

	currSig, err := s.getEntity(entityID)
	if err != nil {
		return names
	}

	parentMsg := currSig.ParentMessage()
	if parentMsg == nil {
		return names
	}

	for _, tmpSig := range parentMsg.Signals() {
		if tmpSig.EntityID() == acmelib.EntityID(entityID) {
			continue
		}

		names = append(names, tmpSig.Name())
	}

	return names
}

func (s *SignalService) UpdateName(entityID string, req UpdateNameReq) (Signal, error) {
	return s.handle(entityID, &req, s.handler.updateName)
}

func (s *SignalService) UpdateDesc(entityID string, req UpdateDescReq) (Signal, error) {
	return s.handle(entityID, &req, s.handler.updateDesc)
}

func (s *SignalService) UpdateSignalType(entityID string, req UpdateSignalTypeReq) (Signal, error) {
	return s.handle(entityID, &req, s.handler.updateSignalType)
}

func (s *SignalService) UpdateSignalUnit(entityID string, req UpdateSignalUnitReq) (Signal, error) {
	return s.handle(entityID, &req, s.handler.updateSignalUnit)
}

type signalRes = response[acmelib.Signal]

type signalHandler struct {
	*commonServiceHandler

	sigTypeSrv *SignalTypeService
	sigUnitSrv *SignalUnitService
}

func newSignalHandler(sidebar *sidebarController, sigTypeSrv *SignalTypeService, sigUnitSrv *SignalUnitService) *signalHandler {
	return &signalHandler{
		commonServiceHandler: newCommonServiceHandler(sidebar),

		sigTypeSrv: sigTypeSrv,
		sigUnitSrv: sigUnitSrv,
	}
}

func (h *signalHandler) toResponse(sig acmelib.Signal) Signal {
	return newSignal(sig)
}

func (h *signalHandler) updateName(sig acmelib.Signal, req *request, res *signalRes) error {
	parsedReq := req.toUpdateName()

	name := strings.TrimSpace(parsedReq.Name)

	oldName := sig.Name()
	if name == oldName {
		return nil
	}

	if err := sig.UpdateName(name); err != nil {
		return err
	}

	h.sidebar.sendUpdateName(sig)

	res.setUndo(
		func() (acmelib.Signal, error) {
			if err := sig.UpdateName(oldName); err != nil {
				return nil, err
			}
			h.sidebar.sendUpdateName(sig)
			return sig, nil
		},
	)

	res.setRedo(
		func() (acmelib.Signal, error) {
			if err := sig.UpdateName(name); err != nil {
				return nil, err
			}
			h.sidebar.sendUpdateName(sig)
			return sig, nil
		},
	)

	return nil
}

func (h *signalHandler) updateDesc(sig acmelib.Signal, req *request, res *signalRes) error {
	parsedReq := req.toUpdateDesc()

	desc := parsedReq.Desc

	oldDesc := sig.Desc()
	if desc == oldDesc {
		return nil
	}

	sig.SetDesc(desc)

	res.setUndo(
		func() (acmelib.Signal, error) {
			sig.SetDesc(oldDesc)
			return sig, nil
		},
	)

	res.setRedo(
		func() (acmelib.Signal, error) {
			sig.SetDesc(desc)
			return sig, nil
		},
	)

	return nil
}

func (h *signalHandler) updateSignalType(sig acmelib.Signal, req *request, res *signalRes) error {
	stdSig, err := sig.ToStandard()
	if err != nil {
		return err
	}

	parsedReq := req.toUpdateSignalType()
	sigTypeEntID := parsedReq.SignalTypeEntityID

	oldSigType := stdSig.Type()
	if sigTypeEntID == oldSigType.EntityID().String() {
		return nil
	}

	h.sigTypeSrv.mux.Lock()
	defer h.sigTypeSrv.mux.Unlock()

	sigType, err := h.sigTypeSrv.getEntity(sigTypeEntID)
	if err != nil {
		return err
	}

	if err := stdSig.SetType(sigType); err != nil {
		return err
	}

	res.setUndo(
		func() (acmelib.Signal, error) {
			if err := stdSig.SetType(oldSigType); err != nil {
				return nil, err
			}
			return sig, nil
		},
	)

	res.setRedo(
		func() (acmelib.Signal, error) {
			if err := stdSig.SetType(sigType); err != nil {
				return nil, err
			}
			return sig, nil
		},
	)

	return nil
}

func (h *signalHandler) updateSignalUnit(sig acmelib.Signal, req *request, res *signalRes) error {
	stdSig, err := sig.ToStandard()
	if err != nil {
		return err
	}

	parsedReq := req.toUpdateSignalUnit()
	sigUnitEntID := parsedReq.SignalUnitEntityID

	isClearing := len(sigUnitEntID) == 0

	oldSigUnit := stdSig.Unit()

	if oldSigUnit == nil && isClearing {
		return nil
	}

	if oldSigUnit != nil && sigUnitEntID == oldSigUnit.EntityID().String() {
		return nil
	}

	h.sigUnitSrv.mux.Lock()
	defer h.sigUnitSrv.mux.Unlock()

	var sigUnit *acmelib.SignalUnit
	sigUnit = nil

	if !isClearing {
		sigUnit, err = h.sigUnitSrv.getEntity(sigUnitEntID)
		if err != nil {
			return err
		}
	}

	stdSig.SetUnit(sigUnit)

	res.setUndo(
		func() (acmelib.Signal, error) {
			stdSig.SetUnit(oldSigUnit)
			return sig, nil
		},
	)

	res.setRedo(
		func() (acmelib.Signal, error) {
			stdSig.SetUnit(sigUnit)
			return sig, nil
		},
	)

	return nil
}
