package main

import (
	"sync"

	"github.com/squadracorsepolito/acmelib"
)

type CANIDBuilderService struct {
	*service[*acmelib.CANIDBuilder, CANIDBuilder, *canIDBuilderHandler]
}

func newCANIDBuilderService(mux *sync.RWMutex, sidebarCtr *sidebarController) *CANIDBuilderService {
	handler := newCANIDBuilderHandler(sidebarCtr)

	return &CANIDBuilderService{
		service: newService(serviceKindCANIDBuilder, handler, mux, sidebarCtr),
	}
}

func (s *CANIDBuilderService) UpdateName(entityID string, req UpdateNameReq) (CANIDBuilder, error) {
	return s.handle(entityID, &req, s.handler.updateName)
}

func (s *CANIDBuilderService) UpdateDesc(entityID string, req UpdateDescReq) (CANIDBuilder, error) {
	return s.handle(entityID, &req, s.handler.updateDesc)
}

func (s *CANIDBuilderService) CalculateCANIDs(entityID string, req CalculateCANIDReq) ([]uint, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	builder, err := s.getEntity(entityID)
	if err != nil {
		return nil, err
	}

	msgPriority := req.MessagePrioriry.parse()
	msgID := acmelib.MessageID(req.MessageID)
	nodeID := acmelib.NodeID(req.NodeID)

	canIDs := builder.CalculatePartials(msgPriority, msgID, nodeID)

	res := make([]uint, len(canIDs))
	for idx, canID := range canIDs {
		res[idx] = uint(canID)
	}

	return res, nil
}

func (s *CANIDBuilderService) InsertOperation(entityID string, req InsertOperationReq) (CANIDBuilder, error) {
	return s.handle(entityID, &req, s.handler.insertOperation)
}

func (s *CANIDBuilderService) DeleteOperation(entityID string, req DeleteOperationReq) (CANIDBuilder, error) {
	return s.handle(entityID, &req, s.handler.deleteOperation)
}

func (s *CANIDBuilderService) UpdateOperationFrom(entityID string, req UpdateOperationFromReq) (CANIDBuilder, error) {
	return s.handle(entityID, &req, s.handler.updateOperationFrom)
}

func (s *CANIDBuilderService) UpdateOperationLen(entityID string, req UpdateOperationLenReq) (CANIDBuilder, error) {
	return s.handle(entityID, &req, s.handler.updateOperationLen)
}

type canIDBuilderRes = response[*acmelib.CANIDBuilder]

type canIDBuilderHandler struct {
	sidebarCtr *sidebarController
}

func newCANIDBuilderHandler(sidebarCtr *sidebarController) *canIDBuilderHandler {
	return &canIDBuilderHandler{sidebarCtr: sidebarCtr}
}

func (h *canIDBuilderHandler) toResponse(canIDBuilder *acmelib.CANIDBuilder) CANIDBuilder {
	return newCANIDBuilder(canIDBuilder)
}

func (h *canIDBuilderHandler) updateName(canIDBuilder *acmelib.CANIDBuilder, req *request, res *canIDBuilderRes) error {
	parsedReq := req.toUpdateName()

	name := parsedReq.Name

	oldName := canIDBuilder.Name()
	if name == oldName {
		return nil
	}

	if err := canIDBuilder.UpdateName(name); err != nil {
		return err
	}

	h.sidebarCtr.sendUpdateName(canIDBuilder)

	res.setUndo(
		func() (*acmelib.CANIDBuilder, error) {
			if err := canIDBuilder.UpdateName(oldName); err != nil {
				return nil, err
			}

			h.sidebarCtr.sendUpdateName(canIDBuilder)

			return canIDBuilder, nil
		},
	)

	res.setRedo(
		func() (*acmelib.CANIDBuilder, error) {
			if err := canIDBuilder.UpdateName(name); err != nil {
				return nil, err
			}

			h.sidebarCtr.sendUpdateName(canIDBuilder)

			return canIDBuilder, nil
		},
	)

	return nil
}

func (h *canIDBuilderHandler) updateDesc(canIDBuilder *acmelib.CANIDBuilder, req *request, res *canIDBuilderRes) error {
	parsedReq := req.toUpdateDesc()

	desc := parsedReq.Desc

	oldDesc := canIDBuilder.Desc()
	if desc == oldDesc {
		return nil
	}

	canIDBuilder.SetDesc(desc)

	res.setUndo(
		func() (*acmelib.CANIDBuilder, error) {
			canIDBuilder.SetDesc(oldDesc)
			return canIDBuilder, nil
		},
	)

	res.setRedo(
		func() (*acmelib.CANIDBuilder, error) {
			canIDBuilder.SetDesc(desc)
			return canIDBuilder, nil
		},
	)

	return nil
}

func (h *canIDBuilderHandler) insertOperation(canIDBuilder *acmelib.CANIDBuilder, req *request, res *canIDBuilderRes) error {
	parsedReq := req.toInsertOperation()

	kind := parsedReq.OpKind.parse()
	from := parsedReq.OpFrom
	len := parsedReq.OpLen
	opIdx := parsedReq.OpIndex

	if err := canIDBuilder.InsertOperation(kind, from, len, opIdx); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.CANIDBuilder, error) {
			if err := canIDBuilder.RemoveOperation(opIdx); err != nil {
				return nil, err
			}
			return canIDBuilder, nil
		},
	)

	res.setRedo(
		func() (*acmelib.CANIDBuilder, error) {
			if err := canIDBuilder.InsertOperation(kind, from, len, opIdx); err != nil {
				return nil, err
			}
			return canIDBuilder, nil
		},
	)

	return nil
}

func (h *canIDBuilderHandler) deleteOperation(canIDBuilder *acmelib.CANIDBuilder, req *request, res *canIDBuilderRes) error {
	parsedReq := req.toDeleteOperation()

	idx := parsedReq.OpIndex

	operations := canIDBuilder.Operations()
	if idx < 0 || idx >= len(operations) {
		return nil
	}

	op := operations[idx]

	if err := canIDBuilder.RemoveOperation(idx); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.CANIDBuilder, error) {
			if err := canIDBuilder.InsertOperation(op.Kind(), op.From(), op.Len(), idx); err != nil {
				return nil, err
			}

			return canIDBuilder, nil
		},
	)

	res.setRedo(
		func() (*acmelib.CANIDBuilder, error) {
			if err := canIDBuilder.RemoveOperation(idx); err != nil {
				return nil, err
			}

			return canIDBuilder, nil
		},
	)

	return nil
}

func (h *canIDBuilderHandler) updateOperationFrom(canIDBuilder *acmelib.CANIDBuilder, req *request, res *canIDBuilderRes) error {
	parsedReq := req.toUpdateOperationFrom()

	from := parsedReq.OpFrom
	opIdx := parsedReq.OpIndex

	operations := canIDBuilder.Operations()
	if opIdx < 0 || opIdx >= len(operations) {
		return nil
	}

	op := operations[opIdx]
	oldFrom := op.From()
	if from == oldFrom {
		return nil
	}

	if err := canIDBuilder.RemoveOperation(opIdx); err != nil {
		return err
	}

	if err := canIDBuilder.InsertOperation(op.Kind(), from, op.Len(), opIdx); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.CANIDBuilder, error) {
			if err := canIDBuilder.RemoveOperation(opIdx); err != nil {
				return nil, err
			}

			if err := canIDBuilder.InsertOperation(op.Kind(), oldFrom, op.Len(), opIdx); err != nil {
				return nil, err
			}

			return canIDBuilder, nil
		},
	)

	res.setRedo(
		func() (*acmelib.CANIDBuilder, error) {
			if err := canIDBuilder.RemoveOperation(opIdx); err != nil {
				return nil, err
			}

			if err := canIDBuilder.InsertOperation(op.Kind(), from, op.Len(), opIdx); err != nil {
				return nil, err
			}

			return canIDBuilder, nil
		},
	)

	return nil
}

func (h *canIDBuilderHandler) updateOperationLen(canIDBuilder *acmelib.CANIDBuilder, req *request, res *canIDBuilderRes) error {
	parsedReq := req.toUpdateOperationLen()

	opLen := parsedReq.OpLen
	opIdx := parsedReq.OpIndex

	operations := canIDBuilder.Operations()
	if opIdx < 0 || opIdx >= len(operations) {
		return nil
	}

	op := operations[opIdx]
	oldLen := op.Len()
	if opLen == oldLen {
		return nil
	}

	if err := canIDBuilder.RemoveOperation(opIdx); err != nil {
		return err
	}

	if err := canIDBuilder.InsertOperation(op.Kind(), op.From(), opLen, opIdx); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.CANIDBuilder, error) {
			if err := canIDBuilder.RemoveOperation(opIdx); err != nil {
				return nil, err
			}

			if err := canIDBuilder.InsertOperation(op.Kind(), op.From(), oldLen, opIdx); err != nil {
				return nil, err
			}

			return canIDBuilder, nil
		},
	)

	res.setRedo(
		func() (*acmelib.CANIDBuilder, error) {
			if err := canIDBuilder.RemoveOperation(opIdx); err != nil {
				return nil, err
			}

			if err := canIDBuilder.InsertOperation(op.Kind(), op.From(), opLen, opIdx); err != nil {
				return nil, err
			}

			return canIDBuilder, nil
		},
	)

	return nil
}
