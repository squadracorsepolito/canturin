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
