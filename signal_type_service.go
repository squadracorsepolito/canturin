package main

import (
	"fmt"
	"strings"

	"github.com/squadracorsepolito/acmelib"
)

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

type SignalType struct {
	base

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

type SignalTypeService struct {
	*service[*acmelib.SignalType, SignalType, *signalTypeHandler]
}

func newSignalTypeService(sidebar *sidebarController) *SignalTypeService {
	return &SignalTypeService{
		service: newService(serviceKindSignalType, newSignalTypeHandler(sidebar), sidebar),
	}
}

func (s *SignalTypeService) Create(req CreateSignalTypeReq) (SignalType, error) {
	name := req.Name
	desc := req.Desc
	size := req.Size
	kind := req.Kind
	signed := req.Signed
	min := req.Min
	max := req.Max
	scale := req.Scale
	offset := req.Offset

	sigType := &acmelib.SignalType{}
	switch kind {
	case SignalTypeKindCustom:
		tmpSigType, err := acmelib.NewCustomSignalType(name, size, signed, min, max, scale, offset)
		if err != nil {
			return SignalType{}, err
		}
		sigType = tmpSigType

	case SignalTypeKindFlag:
		sigType = acmelib.NewFlagSignalType(name)

	case SignalTypeKindInteger:
		tmpSigType, err := acmelib.NewIntegerSignalType(name, size, signed)
		if err != nil {
			return SignalType{}, err
		}
		sigType = tmpSigType

	case SignalTypeKindDecimal:
		tmpSigType, err := acmelib.NewDecimalSignalType(name, size, signed)
		if err != nil {
			return SignalType{}, err
		}
		sigType = tmpSigType
	}

	if len(desc) > 0 {
		sigType.SetDesc(desc)
	}

	sigType.SetMin(min)
	sigType.SetMax(max)
	sigType.SetScale(scale)
	sigType.SetOffset(offset)

	s.mux.Lock()
	defer s.mux.Unlock()

	s.addEntity(sigType)
	s.sidebar.sendAdd(sigType)

	s.sendHistoryOp(
		func() (*acmelib.SignalType, error) {
			s.removeEntity(sigType.EntityID().String())
			s.sidebar.sendDelete(sigType)
			return sigType, nil
		},
		func() (*acmelib.SignalType, error) {
			s.addEntity(sigType)
			s.sidebar.sendAdd(sigType)
			return sigType, nil
		},
	)

	return s.handler.toResponse(sigType), nil
}

func (s *SignalTypeService) Delete(entityID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	sigType, err := s.getEntity(entityID)
	if err != nil {
		return err
	}

	if sigType.ReferenceCount() > 0 {
		return fmt.Errorf("signal type %s is referenced %d times", sigType.Name(), sigType.ReferenceCount())
	}

	s.removeEntity(entityID)
	s.sidebar.sendDelete(sigType)

	s.sendHistoryOp(
		func() (*acmelib.SignalType, error) {
			s.addEntity(sigType)
			s.sidebar.sendAdd(sigType)
			return sigType, nil
		},
		func() (*acmelib.SignalType, error) {
			s.removeEntity(sigType.EntityID().String())
			s.sidebar.sendDelete(sigType)
			return sigType, nil
		},
	)

	return nil
}

func (s *SignalTypeService) UpdateName(entityID string, req UpdateNameReq) (SignalType, error) {
	return s.handle(entityID, &req, s.handler.updateName)
}

func (s *SignalTypeService) UpdateDesc(entityID string, req UpdateDescReq) (SignalType, error) {
	return s.handle(entityID, &req, s.handler.updateDesc)
}

func (s *SignalTypeService) UpdateMin(entityID string, req UpdateMinReq) (SignalType, error) {
	return s.handle(entityID, &req, s.handler.updateMin)
}

func (s *SignalTypeService) UpdateMax(entityID string, req UpdateMaxReq) (SignalType, error) {
	return s.handle(entityID, &req, s.handler.updateMax)
}

func (s *SignalTypeService) UpdateScale(entityID string, req UpdateScaleReq) (SignalType, error) {
	return s.handle(entityID, &req, s.handler.updateScale)
}

func (s *SignalTypeService) UpdateOffset(entityID string, req UpdateOffsetReq) (SignalType, error) {
	return s.handle(entityID, &req, s.handler.updateOffset)
}

type signalTypeRes = response[*acmelib.SignalType]

type signalTypeHandler struct {
	*commonServiceHandler
}

func newSignalTypeHandler(sidebar *sidebarController) *signalTypeHandler {
	return &signalTypeHandler{
		commonServiceHandler: newCommonServiceHandler(sidebar),
	}
}

func (h *signalTypeHandler) toResponse(sigType *acmelib.SignalType) SignalType {
	refCount := sigType.ReferenceCount()

	res := SignalType{
		base: getBase(sigType),

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
		sigRef.entityID = msg.EntityID()
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

func (h *signalTypeHandler) updateName(sigType *acmelib.SignalType, req *request, res *signalTypeRes) error {
	parsedReq := req.toUpdateName()

	name := strings.TrimSpace(parsedReq.Name)

	oldName := sigType.Name()
	if name == oldName {
		return nil
	}

	sigType.SetName(name)
	h.sidebar.sendUpdateName(sigType)

	res.setUndo(
		func() (*acmelib.SignalType, error) {
			sigType.SetName(oldName)
			h.sidebar.sendUpdateName(sigType)
			return sigType, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalType, error) {
			sigType.SetName(name)
			h.sidebar.sendUpdateName(sigType)
			return sigType, nil
		},
	)

	return nil
}

func (h *signalTypeHandler) updateDesc(sigType *acmelib.SignalType, req *request, res *signalTypeRes) error {
	parsedReq := req.toUpdateDesc()

	desc := parsedReq.Desc

	oldDesc := sigType.Desc()
	if oldDesc == desc {
		return nil
	}

	sigType.SetDesc(desc)

	res.setUndo(
		func() (*acmelib.SignalType, error) {
			sigType.SetDesc(oldDesc)
			return sigType, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalType, error) {
			sigType.SetDesc(desc)
			return sigType, nil
		},
	)

	return nil
}

func (h *signalTypeHandler) updateMin(sigType *acmelib.SignalType, req *request, res *signalTypeRes) error {
	parsedReq := req.toUpdateMin()

	min := parsedReq.Min

	oldMin := sigType.Min()
	if min == oldMin {
		return nil
	}

	sigType.SetMin(min)

	res.setUndo(
		func() (*acmelib.SignalType, error) {
			sigType.SetMin(oldMin)
			return sigType, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalType, error) {
			sigType.SetMin(min)
			return sigType, nil
		},
	)

	return nil
}

func (h *signalTypeHandler) updateMax(sigType *acmelib.SignalType, req *request, res *signalTypeRes) error {
	parsedReq := req.toUpdateMax()

	max := parsedReq.Max

	oldMax := sigType.Max()
	if oldMax == max {
		return nil
	}

	sigType.SetMax(max)

	res.setUndo(
		func() (*acmelib.SignalType, error) {
			sigType.SetMax(oldMax)
			return sigType, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalType, error) {
			sigType.SetMax(max)
			return sigType, nil
		},
	)

	return nil
}

func (h *signalTypeHandler) updateScale(sigType *acmelib.SignalType, req *request, res *signalTypeRes) error {
	parsedReq := req.toUpdateScale()

	scale := parsedReq.Scale

	oldScale := sigType.Scale()
	if oldScale == scale {
		return nil
	}

	sigType.SetScale(scale)

	res.setUndo(
		func() (*acmelib.SignalType, error) {
			sigType.SetScale(oldScale)
			return sigType, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalType, error) {
			sigType.SetScale(scale)
			return sigType, nil
		},
	)

	return nil
}

func (h *signalTypeHandler) updateOffset(sigType *acmelib.SignalType, req *request, res *signalTypeRes) error {
	parsedReq := req.toUpdateOffset()

	offset := parsedReq.Offset

	oldOffset := sigType.Offset()
	if oldOffset == offset {
		return nil
	}

	sigType.SetOffset(offset)

	res.setUndo(
		func() (*acmelib.SignalType, error) {
			sigType.SetOffset(oldOffset)
			return sigType, nil
		},
	)

	res.setRedo(
		func() (*acmelib.SignalType, error) {
			sigType.SetOffset(offset)
			return sigType, nil
		},
	)

	return nil
}
