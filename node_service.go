package main

import (
	"github.com/squadracorsepolito/acmelib"
)

type NodeInterface struct {
	Number           int          `json:"number"`
	AttachedBus      BaseEntity   `json:"attachedBus"`
	SentMessages     []BaseEntity `json:"sentMessages"`
	ReceivedMessages []BaseEntity `json:"receivedMessages"`
}

func getNodeInterface(nodeInt *acmelib.NodeInterface) NodeInterface {
	sentMessages := []BaseEntity{}
	for _, tmpMsg := range nodeInt.SentMessages() {
		sentMessages = append(sentMessages, getBaseEntity(tmpMsg))
	}

	receivedMessages := []BaseEntity{}
	for _, tmpMsg := range nodeInt.ReceivedMessages() {
		receivedMessages = append(receivedMessages, getBaseEntity(tmpMsg))
	}

	res := NodeInterface{
		Number:           nodeInt.Number(),
		SentMessages:     sentMessages,
		ReceivedMessages: receivedMessages,
	}

	if nodeInt.ParentBus() != nil {
		res.AttachedBus = getBaseEntity(nodeInt.ParentBus())
	}

	return res
}

type Node struct {
	base

	ID         uint            `json:"id"`
	Interfaces []NodeInterface `json:"interfaces"`
}

type NodeService struct {
	*service0[*acmelib.Node, Node, *nodeHandler]
}

func newNodeService(sidebar *sidebarController, bus *BusService) *NodeService {
	return &NodeService{
		service0: newService0(serviceKindNode, newNodeHandler(sidebar, bus), sidebar),
	}
}

func (s *NodeService) GetInvalidNodeIDs(entityID string) []uint {
	s.mux.RLock()
	defer s.mux.RUnlock()

	nodeIDs := []uint{}

	for _, tmpNode := range s.entities {
		if tmpNode.EntityID().String() == entityID {
			continue
		}

		nodeIDs = append(nodeIDs, uint(tmpNode.ID()))
	}

	return nodeIDs
}

func (s *NodeService) Create(req CreateNodeReq) (Node, error) {
	node := acmelib.NewNode(req.Name, acmelib.NodeID(req.NodeID), int(req.InterfaceCount))
	node.SetDesc(req.Desc)

	s.mux.Lock()
	defer s.mux.Unlock()

	s.addEntity(node)
	s.sidebar.sendAdd(node)

	s.sendHistoryOp(
		func() (*acmelib.Node, error) {
			s.removeEntity(node.EntityID().String())
			s.sidebar.sendDelete(node)
			return node, nil
		},
		func() (*acmelib.Node, error) {
			s.addEntity(node)
			s.sidebar.sendAdd(node)
			return node, nil
		},
	)

	return s.handler.toResponse(node), nil
}

func (s *NodeService) Delete(entityID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	node, err := s.getEntity(entityID)
	if err != nil {
		return err
	}

	nodeIntBusMap := make(map[int]*acmelib.NodeInterface)

	parBuses := []*acmelib.Bus{}
	parBusIdx := 0

	for _, nodeInt := range node.Interfaces() {
		if nodeInt.ParentBus() != nil {
			parBuses = append(parBuses, nodeInt.ParentBus())
			nodeIntBusMap[parBusIdx] = nodeInt
			parBusIdx++
		}
	}

	for _, tmpBus := range parBuses {
		if err := tmpBus.RemoveNodeInterface(node.EntityID()); err != nil {
			return err
		}
	}

	s.removeEntity(entityID)
	s.sidebar.sendDelete(node)

	s.sendHistoryOp(
		func() (*acmelib.Node, error) {
			for idx, tmpBus := range parBuses {
				if nodeInt, ok := nodeIntBusMap[idx]; ok {
					if err := tmpBus.AddNodeInterface(nodeInt); err != nil {
						return nil, err
					}
				}
			}

			s.addEntity(node)
			s.sidebar.sendAdd(node)

			return node, nil
		},
		func() (*acmelib.Node, error) {
			for _, tmpBus := range parBuses {
				if err := tmpBus.RemoveNodeInterface(node.EntityID()); err != nil {
					return nil, err
				}
			}

			s.removeEntity(node.EntityID().String())
			s.sidebar.sendDelete(node)

			return node, nil
		},
	)

	return nil
}

func (s *NodeService) UpdateName(entityID string, req UpdateNameReq) (Node, error) {
	return s.handle(entityID, &req, s.handler.updateName)
}

func (s *NodeService) UpdateDesc(entityID string, req UpdateDescReq) (Node, error) {
	return s.handle(entityID, &req, s.handler.updateDesc)
}

func (s *NodeService) UpdateNodeID(entityID string, req UpdateNodeIDReq) (Node, error) {
	return s.handle(entityID, &req, s.handler.updateNodeID)
}

func (s *NodeService) UpdateAttachedBus(entityID string, req UpdateAttachedBusReq) (Node, error) {
	return s.handle(entityID, &req, s.handler.updateAttachedBus)
}

func (s *NodeService) RemoveSentMessages(entityID string, req RemoveSentMessagesReq) (Node, error) {
	return s.handle(entityID, &req, s.handler.removeSentMessages)
}

func (s *NodeService) RemoveReceivedMessages(entityID string, req RemoveReceivedMessagesReq) (Node, error) {
	return s.handle(entityID, &req, s.handler.removeReceivedMessages)
}

type nodeRes = response[*acmelib.Node]

type nodeHandler struct {
	*commonServiceHandler

	bus *BusService
}

func newNodeHandler(sidebar *sidebarController, bus *BusService) *nodeHandler {
	return &nodeHandler{
		commonServiceHandler: newCommonServiceHandler(sidebar),

		bus: bus,
	}
}

func (h *nodeHandler) toResponse(node *acmelib.Node) Node {
	res := Node{
		base: getBase(node),

		ID:         uint(node.ID()),
		Interfaces: []NodeInterface{},
	}

	for _, nodeInt := range node.Interfaces() {
		res.Interfaces = append(res.Interfaces, getNodeInterface(nodeInt))
	}

	return res
}

func (h *nodeHandler) updateName(node *acmelib.Node, req *request, res *nodeRes) error {
	parsedReq := req.toUpdateName()

	name := parsedReq.Name

	oldName := node.Name()
	if name == oldName {
		return nil
	}

	if err := node.UpdateName(name); err != nil {
		return err
	}

	h.sidebar.sendUpdateName(node)

	res.setUndo(
		func() (*acmelib.Node, error) {
			if err := node.UpdateName(oldName); err != nil {
				return nil, err
			}

			h.sidebar.sendUpdateName(node)

			return node, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Node, error) {
			if err := node.UpdateName(name); err != nil {
				return nil, err
			}

			h.sidebar.sendUpdateName(node)

			return node, nil
		},
	)

	return nil
}

func (h *nodeHandler) updateDesc(node *acmelib.Node, req *request, res *nodeRes) error {
	parsedRes := req.toUpdateDesc()

	desc := parsedRes.Desc

	oldDesc := node.Desc()
	if desc == oldDesc {
		return nil
	}

	node.SetDesc(desc)

	res.setUndo(
		func() (*acmelib.Node, error) {
			node.SetDesc(oldDesc)
			return node, nil
		},
	)

	res.setUndo(
		func() (*acmelib.Node, error) {
			node.SetDesc(desc)
			return node, nil
		},
	)

	return nil
}

func (h *nodeHandler) updateNodeID(node *acmelib.Node, req *request, res *nodeRes) error {
	parsedRes := req.toUpdateNodeID()

	nodeID := acmelib.NodeID(parsedRes.NodeID)

	oldNodeID := node.ID()
	if nodeID == oldNodeID {
		return nil
	}

	if err := node.UpdateID(nodeID); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.Node, error) {
			if err := node.UpdateID(oldNodeID); err != nil {
				return nil, err
			}

			return node, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Node, error) {
			if err := node.UpdateID(nodeID); err != nil {
				return nil, err
			}

			return node, nil
		},
	)

	return nil
}

func (h *nodeHandler) updateAttachedBus(node *acmelib.Node, req *request, res *nodeRes) error {
	h.bus.mux.Lock()
	defer h.bus.mux.Unlock()

	parsedRes := req.toUpdateAttachedBus()

	busEntID := parsedRes.BusEntityID
	intNum := parsedRes.InterfaceNumber
	nodeEntID := node.EntityID()

	nodeInt, err := node.GetInterface(intNum)
	if err != nil {
		return err
	}

	oldBus := nodeInt.ParentBus()
	if oldBus != nil {
		if oldBus.EntityID().String() == busEntID {
			return nil
		}

		if err := oldBus.RemoveNodeInterface(nodeEntID); err != nil {
			return err
		}
	}

	bus, err := h.bus.getEntity(busEntID)
	if err != nil {
		return err
	}

	if err := bus.AddNodeInterface(nodeInt); err != nil {
		return err
	}

	res.setUndo(
		func() (*acmelib.Node, error) {
			h.bus.mux.Lock()
			defer h.bus.mux.Unlock()

			if err := bus.RemoveNodeInterface(nodeEntID); err != nil {
				return nil, err
			}

			if oldBus == nil {
				return node, nil
			}

			if err := oldBus.AddNodeInterface(nodeInt); err != nil {
				return nil, err
			}

			return node, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Node, error) {
			h.bus.mux.Lock()
			defer h.bus.mux.Unlock()

			if oldBus != nil {
				if err := oldBus.RemoveNodeInterface(nodeInt.Node().EntityID()); err != nil {
					return nil, err
				}
			}

			if err := bus.AddNodeInterface(nodeInt); err != nil {
				return nil, err
			}

			return node, nil
		},
	)

	return nil
}

func (h *nodeHandler) removeSentMessages(node *acmelib.Node, req *request, res *nodeRes) error {
	parsedReq := req.toRemoveSentMessages()

	intNum := parsedReq.InterfaceNumber
	msgEntIDs := parsedReq.MessageEntityIDs

	if len(msgEntIDs) == 0 {
		return nil
	}

	nodeInt, err := node.GetInterface(intNum)
	if err != nil {
		return err
	}

	parsedMsgEntIDs := make(map[string]struct{})
	for _, msgEntID := range msgEntIDs {
		parsedMsgEntIDs[msgEntID] = struct{}{}
	}

	msgToRemove := []*acmelib.Message{}
	for _, sentMsg := range nodeInt.SentMessages() {
		if _, ok := parsedMsgEntIDs[sentMsg.EntityID().String()]; ok {
			msgToRemove = append(msgToRemove, sentMsg)
		}
	}

	for _, tmpMsg := range msgToRemove {
		if err := nodeInt.RemoveSentMessage(tmpMsg.EntityID()); err != nil {
			return err
		}

		h.sidebar.sendDelete(tmpMsg)
	}

	res.setUndo(
		func() (*acmelib.Node, error) {
			for _, tmpMsg := range msgToRemove {
				if err := nodeInt.AddSentMessage(tmpMsg); err != nil {
					return nil, err
				}

				h.sidebar.sendAdd(tmpMsg)
			}

			return node, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Node, error) {
			for _, tmpMsg := range msgToRemove {
				if err := nodeInt.RemoveSentMessage(tmpMsg.EntityID()); err != nil {
					return nil, err
				}

				h.sidebar.sendDelete(tmpMsg)
			}

			return node, nil
		},
	)

	return nil
}

func (h *nodeHandler) removeReceivedMessages(node *acmelib.Node, req *request, res *nodeRes) error {
	parsedRes := req.toRemoveReceivedMessages()

	intNum := parsedRes.InterfaceNumber
	msgEntIDs := parsedRes.MessageEntityIDs

	if len(msgEntIDs) == 0 {
		return nil
	}

	nodeInt, err := node.GetInterface(intNum)
	if err != nil {
		return err
	}

	parsedMsgEntIDs := make(map[string]struct{})
	for _, msgEntID := range msgEntIDs {
		parsedMsgEntIDs[msgEntID] = struct{}{}
	}

	msgToRemove := []*acmelib.Message{}
	for _, receivedMsg := range nodeInt.ReceivedMessages() {
		if _, ok := parsedMsgEntIDs[receivedMsg.EntityID().String()]; ok {
			msgToRemove = append(msgToRemove, receivedMsg)
		}
	}

	for _, tmpMsg := range msgToRemove {
		if err := nodeInt.RemoveReceivedMessage(tmpMsg.EntityID()); err != nil {
			return err
		}
	}

	res.setUndo(
		func() (*acmelib.Node, error) {
			for _, tmpMsg := range msgToRemove {
				if err := nodeInt.AddReceivedMessage(tmpMsg); err != nil {
					return nil, err
				}
			}
			return node, nil
		},
	)

	res.setRedo(
		func() (*acmelib.Node, error) {
			for _, tmpMsg := range msgToRemove {
				if err := nodeInt.RemoveReceivedMessage(tmpMsg.EntityID()); err != nil {
					return nil, err
				}
			}
			return node, nil
		},
	)

	return nil
}

//
//
//

// type NodeService0 struct {
// 	*service[*acmelib.Node, Node]
// }

// func getNode(node *acmelib.Node) Node {
// 	interfaces := []NodeInterface{}

// 	// for _, tmpInt := range node.Interfaces() {
// 	// 	interfaces = append(interfaces, getNodeInterface(tmpInt))
// 	// }

// 	return Node{
// 		base: getBase(node),

// 		ID:         uint(node.ID()),
// 		Interfaces: interfaces,
// 	}
// }

// func newNodeService() *NodeService0 {
// 	return &NodeService0{
// 		service: newService(getNode),
// 	}
// }

// func (s *NodeService0) sendSidebarAdd(node *acmelib.Node) {
// 	item := newNodeSidebarItem(node)
// 	manager.sidebar.sendAdd(newSidebarAddReq(item, SidebarNodesPrefix))

// 	for _, nodeInt := range node.Interfaces() {
// 		if nodeInt.ParentBus() == nil {
// 			continue
// 		}

// 		tmpItem := newMessageNodeGroupSidebarItem(nodeInt)
// 		manager.sidebar.sendAdd(newSidebarAddReq(tmpItem, nodeInt.ParentBus().String()))
// 	}
// }

// func (s *NodeService0) sendSidebarUpdateName(node *acmelib.Node) {
// 	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(node.EntityID().String(), node.Name()))

// 	for _, nodeInt := range node.Interfaces() {
// 		if nodeInt.ParentBus() == nil {
// 			continue
// 		}

// 		nodeIntKey := manager.sidebar.getMessageNodeGroupKey(nodeInt)
// 		nodeIntName := manager.sidebar.getMessageNodeGroupName(nodeInt)

// 		manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(nodeIntKey, nodeIntName))
// 	}
// }

// func (s *NodeService0) sendSidebarDelete(node *acmelib.Node) {
// 	manager.sidebar.sendDelete(newSidebarDeleteReq(node.EntityID().String()))

// 	for _, nodeInt := range node.Interfaces() {
// 		if nodeInt.ParentBus() == nil {
// 			continue
// 		}

// 		nodeIntKey := manager.sidebar.getMessageNodeGroupKey(nodeInt)
// 		manager.sidebar.sendDelete(newSidebarDeleteReq(nodeIntKey))
// 	}
// }

// func (s *NodeService0) GetInvalidIDs(entityID string) []uint {
// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	ids := []uint{}
// 	for _, tmpNode := range s.pool {
// 		if tmpNode.EntityID() != acmelib.EntityID(entityID) {
// 			ids = append(ids, uint(tmpNode.ID()))
// 		}
// 	}

// 	return ids
// }

// func (s *NodeService0) Create(name, desc string, id, interfaceCount int) (Node, error) {
// 	node := acmelib.NewNode(name, acmelib.NodeID(id), interfaceCount)
// 	node.SetDesc(desc)

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	s.sendSidebarAdd(node)

// 	return s.converterFn(node), nil
// }

// func (s *NodeService0) UpdateName(entityID, name string) (Node, error) {
// 	node, err := s.getEntity(entityID)
// 	if err != nil {
// 		return Node{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	oldName := node.Name()
// 	if name == oldName {
// 		return s.converterFn(node), nil
// 	}

// 	if err := node.UpdateName(name); err != nil {
// 		return Node{}, err
// 	}

// 	s.sendSidebarUpdateName(node)

// 	proxy.pushHistoryOperation(
// 		operationDomainNode,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := node.UpdateName(oldName); err != nil {
// 				return Node{}, err
// 			}

// 			s.sendSidebarUpdateName(node)

// 			return s.converterFn(node), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := node.UpdateName(name); err != nil {
// 				return Node{}, err
// 			}

// 			s.sendSidebarUpdateName(node)

// 			return s.converterFn(node), nil
// 		},
// 	)

// 	return s.converterFn(node), nil
// }

// func (s *NodeService0) UpdateDesc(entityID, desc string) (Node, error) {
// 	node, err := s.getEntity(entityID)
// 	if err != nil {
// 		return Node{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	oldDesc := node.Desc()
// 	if desc == oldDesc {
// 		return s.converterFn(node), nil
// 	}

// 	node.SetDesc(desc)

// 	proxy.pushHistoryOperation(
// 		operationDomainNode,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			node.SetDesc(oldDesc)

// 			return s.converterFn(node), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			node.SetDesc(desc)

// 			return s.converterFn(node), nil
// 		},
// 	)

// 	return s.converterFn(node), nil
// }

// func (s *NodeService0) UpdateID(entityID string, id uint) (Node, error) {
// 	node, err := s.getEntity(entityID)
// 	if err != nil {
// 		return Node{}, err
// 	}

// 	nodeID := acmelib.NodeID(id)

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	oldNodeID := node.ID()
// 	if nodeID == oldNodeID {
// 		return s.converterFn(node), nil
// 	}

// 	if err := node.UpdateID(nodeID); err != nil {
// 		return Node{}, err
// 	}

// 	proxy.pushHistoryOperation(
// 		operationDomainNode,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := node.UpdateID(oldNodeID); err != nil {
// 				return Node{}, err
// 			}

// 			return s.converterFn(node), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := node.UpdateID(nodeID); err != nil {
// 				return Node{}, err
// 			}

// 			return s.converterFn(node), nil
// 		},
// 	)

// 	return s.converterFn(node), nil
// }

// func (s *NodeService0) AttachBus(nodeEntID string, intNumber int, busEntID string) (Node, error) {
// 	node, err := s.getEntity(nodeEntID)
// 	if err != nil {
// 		return Node{}, err
// 	}

// 	nodeInt, err := node.GetInterface(intNumber)
// 	if err != nil {
// 		return Node{}, err
// 	}

// 	bus, err := manager.bus.getEntity(busEntID)
// 	if err != nil {
// 		return Node{}, err
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	oldBus := nodeInt.ParentBus()
// 	if oldBus.EntityID() == bus.EntityID() {
// 		return s.converterFn(node), nil
// 	}

// 	if err := bus.AddNodeInterface(nodeInt); err != nil {
// 		return Node{}, err
// 	}

// 	proxy.pushHistoryOperation(
// 		operationDomainNode,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := bus.RemoveNodeInterface(nodeInt.Node().EntityID()); err != nil {
// 				return Node{}, err
// 			}

// 			if oldBus != nil {
// 				if err := oldBus.AddNodeInterface(nodeInt); err != nil {
// 					return Node{}, err
// 				}
// 			}

// 			return s.converterFn(node), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			if err := bus.AddNodeInterface(nodeInt); err != nil {
// 				return Node{}, err
// 			}

// 			return s.converterFn(node), nil
// 		},
// 	)

// 	return s.converterFn(node), nil
// }

// func (s *NodeService0) RemoveSentMessages(nodeEntID string, intNumber int, messageEntIDs ...string) (Node, error) {
// 	node, err := s.getEntity(nodeEntID)
// 	if err != nil {
// 		return Node{}, err
// 	}

// 	if len(messageEntIDs) == 0 {
// 		return s.converterFn(node), nil
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	nodeInt, err := node.GetInterface(intNumber)
// 	if err != nil {
// 		return Node{}, err
// 	}

// 	targetIDs := make(map[string]struct{})
// 	for _, messageEntID := range messageEntIDs {
// 		targetIDs[messageEntID] = struct{}{}
// 	}

// 	sentMessages := []*acmelib.Message{}
// 	for _, tmpMsg := range nodeInt.SentMessages() {
// 		tmpEntID := tmpMsg.EntityID()

// 		_, ok := targetIDs[tmpEntID.String()]
// 		if !ok {
// 			continue
// 		}

// 		sentMessages = append(sentMessages, tmpMsg)
// 	}

// 	for _, tmpMsg := range sentMessages {
// 		tmpEntID := tmpMsg.EntityID()

// 		if err := nodeInt.RemoveSentMessage(tmpEntID); err != nil {
// 			return Node{}, err
// 		}

// 		manager.sidebar.sendDelete(newSidebarDeleteReq(tmpEntID.String()))
// 	}

// 	proxy.pushHistoryOperation(
// 		operationDomainNode,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			for _, tmpMsg := range sentMessages {
// 				if err := nodeInt.AddSentMessage(tmpMsg); err != nil {
// 					return Node{}, err
// 				}

// 				msgNodeKey := manager.sidebar.getMessageNodeGroupKey(nodeInt)
// 				manager.sidebar.sendAdd(newSidebarAddReq(newMessageSidebarItem(tmpMsg), msgNodeKey))
// 			}

// 			return s.converterFn(node), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			for _, tmpMsg := range sentMessages {
// 				tmpEntID := tmpMsg.EntityID()

// 				if err := nodeInt.RemoveSentMessage(tmpMsg.EntityID()); err != nil {
// 					return Node{}, err
// 				}

// 				manager.sidebar.sendDelete(newSidebarDeleteReq(tmpEntID.String()))
// 			}

// 			return s.converterFn(node), nil
// 		},
// 	)

// 	return s.converterFn(node), nil
// }

// func (s *NodeService0) RemoveReceivedMessages(nodeEntID string, intNumber int, messageEntIDs ...string) (Node, error) {
// 	node, err := s.getEntity(nodeEntID)
// 	if err != nil {
// 		return Node{}, err
// 	}

// 	if len(messageEntIDs) == 0 {
// 		return s.converterFn(node), nil
// 	}

// 	s.mux.Lock()
// 	defer s.mux.Unlock()

// 	nodeInt, err := node.GetInterface(intNumber)
// 	if err != nil {
// 		return Node{}, err
// 	}

// 	targetIDs := make(map[string]struct{})
// 	for _, messageEntID := range messageEntIDs {
// 		targetIDs[messageEntID] = struct{}{}
// 	}

// 	receivedMessages := []*acmelib.Message{}
// 	for _, tmpMsg := range nodeInt.ReceivedMessages() {
// 		tmpEntID := tmpMsg.EntityID()

// 		_, ok := targetIDs[tmpEntID.String()]
// 		if !ok {
// 			continue
// 		}

// 		receivedMessages = append(receivedMessages, tmpMsg)
// 	}

// 	for _, tmpMsg := range receivedMessages {
// 		if err := nodeInt.RemoveReceivedMessage(tmpMsg.EntityID()); err != nil {
// 			return Node{}, err
// 		}
// 	}

// 	proxy.pushHistoryOperation(
// 		operationDomainNode,
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			for _, tmpMsg := range receivedMessages {
// 				if err := nodeInt.AddReceivedMessage(tmpMsg); err != nil {
// 					return Node{}, err
// 				}
// 			}

// 			return s.converterFn(node), nil
// 		},
// 		func() (any, error) {
// 			s.mux.Lock()
// 			defer s.mux.Unlock()

// 			for _, tmpMsg := range receivedMessages {
// 				if err := nodeInt.RemoveReceivedMessage(tmpMsg.EntityID()); err != nil {
// 					return Node{}, err
// 				}
// 			}

// 			return s.converterFn(node), nil
// 		},
// 	)

// 	return s.converterFn(node), nil
// }
