package main

import (
	"strings"
	"sync"

	"github.com/squadracorsepolito/acmelib"
)

type SignalService struct {
	*service[acmelib.Signal, Signal, *signalHandler]
}

func newSignalService(mux *sync.RWMutex, sidebar *sidebarController, sigTypeCtr *signalTypeController, sigUnitCtr *signalUnitController, sigEnumCtr *signalEnumController) *SignalService {
	return &SignalService{
		service: newService(serviceKindSignal, newSignalHandler(sidebar, sigTypeCtr, sigUnitCtr, sigEnumCtr), mux, sidebar),
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

func (s *SignalService) UpdateSignalEnum(entityID string, req UpdateSignalEnumReq) (Signal, error) {
	return s.handle(entityID, &req, s.handler.updateSignalEnum)
}

func (s *SignalService) DeleteMultiplexedSignals(entityID string, req DeleteMultiplexedSignalsReq) (Signal, error) {
	return s.handle(entityID, &req, s.handler.deleteMultiplexedSignals)
}

type signalRes = response[acmelib.Signal]

type signalHandler struct {
	*commonServiceHandler

	sigTypeCtr *signalTypeController
	sigUnitCtr *signalUnitController
	sigEnumCtr *signalEnumController
}

func newSignalHandler(sidebar *sidebarController, sigTypeCtr *signalTypeController, sigUnitCtr *signalUnitController, sigEnumCtr *signalEnumController) *signalHandler {
	return &signalHandler{
		commonServiceHandler: newCommonServiceHandler(sidebar),

		sigTypeCtr: sigTypeCtr,
		sigUnitCtr: sigUnitCtr,
		sigEnumCtr: sigEnumCtr,
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

	h.sidebarCtr.sendUpdateName(sig)

	res.setUndo(
		func() (acmelib.Signal, error) {
			if err := sig.UpdateName(oldName); err != nil {
				return nil, err
			}
			h.sidebarCtr.sendUpdateName(sig)
			return sig, nil
		},
	)

	res.setRedo(
		func() (acmelib.Signal, error) {
			if err := sig.UpdateName(name); err != nil {
				return nil, err
			}
			h.sidebarCtr.sendUpdateName(sig)
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

	sigType, err := h.sigTypeCtr.get(sigTypeEntID)
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
			return stdSig, nil
		},
	)

	res.setRedo(
		func() (acmelib.Signal, error) {
			if err := stdSig.SetType(sigType); err != nil {
				return nil, err
			}
			return stdSig, nil
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

	var sigUnit *acmelib.SignalUnit
	sigUnit = nil

	if !isClearing {
		sigUnit, err = h.sigUnitCtr.get(sigUnitEntID)
		if err != nil {
			return err
		}
	}

	stdSig.SetUnit(sigUnit)

	res.setUndo(
		func() (acmelib.Signal, error) {
			stdSig.SetUnit(oldSigUnit)
			return stdSig, nil
		},
	)

	res.setRedo(
		func() (acmelib.Signal, error) {
			stdSig.SetUnit(sigUnit)
			return stdSig, nil
		},
	)

	return nil
}

func (h *signalHandler) updateSignalEnum(sig acmelib.Signal, req *request, res *signalRes) error {
	enumSig, err := sig.ToEnum()
	if err != nil {
		return err
	}

	parsedReq := req.toUpdateSignalEnum()
	sigEnumEntID := parsedReq.SignalEnumEntityID

	oldSigEnum := enumSig.Enum()
	if sigEnumEntID == oldSigEnum.EntityID().String() {
		return nil
	}

	sigEnum, err := h.sigEnumCtr.get(sigEnumEntID)
	if err != nil {
		return err
	}

	if err := enumSig.SetEnum(sigEnum); err != nil {
		return err
	}

	res.setUndo(
		func() (acmelib.Signal, error) {
			if err := enumSig.SetEnum(oldSigEnum); err != nil {
				return nil, err
			}
			return enumSig, nil
		},
	)

	res.setRedo(
		func() (acmelib.Signal, error) {
			if err := enumSig.SetEnum(sigEnum); err != nil {
				return nil, err
			}
			return enumSig, nil
		},
	)

	return nil
}

func (h *signalHandler) deleteMultiplexedSignals(sig acmelib.Signal, req *request, res *signalRes) error {
	parsedReq := req.toDeleteMultiplexedSignals()

	if len(parsedReq.SignalEntityIDs) == 0 {
		return nil
	}

	groupID := parsedReq.GroupID

	muxSig, err := sig.ToMultiplexer()
	if err != nil {
		return err
	}

	remSigIDs := make(map[string]struct{})
	for _, sigID := range parsedReq.SignalEntityIDs {
		remSigIDs[sigID] = struct{}{}
	}

	remSignals := []acmelib.Signal{}
	remStartPos := make(map[string]int)
	for _, sig := range muxSig.GetSignalGroup(groupID) {
		tmpID := sig.EntityID().String()

		if _, ok := remSigIDs[tmpID]; ok {
			remSignals = append(remSignals, sig)
			remStartPos[tmpID] = sig.GetRelativeStartPos()
		}
	}

	for _, sig := range remSignals {
		if err := muxSig.RemoveSignal(sig.EntityID()); err != nil {
			return err
		}
	}

	res.setUndo(
		func() (acmelib.Signal, error) {
			for _, sig := range remSignals {
				startPos, ok := remStartPos[sig.EntityID().String()]
				if !ok {
					continue
				}

				if err := muxSig.InsertSignal(sig, startPos, groupID); err != nil {
					return nil, err
				}
			}

			return muxSig, nil
		},
	)

	res.setRedo(
		func() (acmelib.Signal, error) {
			for _, sig := range remSignals {
				if err := muxSig.RemoveSignal(sig.EntityID()); err != nil {
					return nil, err
				}
			}

			return muxSig, nil
		},
	)

	return nil
}
