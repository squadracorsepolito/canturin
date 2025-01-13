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

type Signal struct {
	base

	Kind     SignalKind `json:"kind"`
	StartPos int        `json:"startPos"`
	Size     int        `json:"size"`
}

func newSignal(sig acmelib.Signal) Signal {
	return Signal{
		base: getBase(sig),

		Kind:     newSignalKind(sig.Kind()),
		StartPos: sig.GetStartBit(),
		Size:     sig.GetSize(),
	}
}

type SignalService struct {
	*service[acmelib.Signal, Signal, *signalHandler]
}

func newSignalService(sidebar *sidebarController) *SignalService {
	return &SignalService{
		service: newService(serviceKindSignal, newSignalHandler(sidebar), sidebar),
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

type signalRes = response[acmelib.Signal]

type signalHandler struct {
	*commonServiceHandler
}

func newSignalHandler(sidebar *sidebarController) *signalHandler {
	return &signalHandler{
		commonServiceHandler: newCommonServiceHandler(sidebar),
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

	res.setUndo(
		func() (acmelib.Signal, error) {
			if err := sig.UpdateName(oldName); err != nil {
				return nil, err
			}
			return sig, nil
		},
	)

	res.setRedo(
		func() (acmelib.Signal, error) {
			if err := sig.UpdateName(name); err != nil {
				return nil, err
			}
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
