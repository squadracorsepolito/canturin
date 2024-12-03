package main

import (
	"github.com/squadracorsepolito/acmelib"
)

type NodeBase struct {
	base
}

func getNodeBase(node *acmelib.Node) NodeBase {
	if node == nil {
		return NodeBase{}
	}

	return NodeBase{
		base: getBase(node),
	}
}

type NodeMessage struct {
	base
}

func getNodeMessage(msg *acmelib.Message) NodeMessage {
	return NodeMessage{
		base: getBase(msg),
	}
}

type NodeInterface struct {
	Number           int           `json:"number"`
	AttachedBus      BusBase       `json:"attachedBus"`
	SentMessages     []NodeMessage `json:"sentMessages"`
	ReceivedMessages []NodeMessage `json:"receivedMessages"`
}

func getNodeInterface(nodeInt *acmelib.NodeInterface) NodeInterface {
	sentMessages := []NodeMessage{}
	for _, tmpMsg := range nodeInt.SentMessages() {
		sentMessages = append(sentMessages, getNodeMessage(tmpMsg))
	}

	receivedMessages := []NodeMessage{}
	for _, tmpMsg := range nodeInt.ReceivedMessages() {
		receivedMessages = append(receivedMessages, getNodeMessage(tmpMsg))
	}

	return NodeInterface{
		AttachedBus:      getBusBase(nodeInt.ParentBus()),
		Number:           nodeInt.Number(),
		SentMessages:     sentMessages,
		ReceivedMessages: receivedMessages,
	}
}

type Node struct {
	base

	ID         uint            `json:"id"`
	Interfaces []NodeInterface `json:"interfaces"`
}

type NodeService struct {
	*service[*acmelib.Node, Node]
}

func getNode(node *acmelib.Node) Node {
	interfaces := []NodeInterface{}

	for _, tmpInt := range node.Interfaces() {
		interfaces = append(interfaces, getNodeInterface(tmpInt))
	}

	return Node{
		base: getBase(node),

		ID:         uint(node.ID()),
		Interfaces: interfaces,
	}
}

func newNodeService() *NodeService {
	return &NodeService{
		service: newService(proxy.nodeCh, getNode),
	}
}

func (s *NodeService) sendSidebarUpdateName(node *acmelib.Node) {
	manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(node.EntityID().String(), node.Name()))

	for _, nodeInt := range node.Interfaces() {
		nodeIntKey := manager.sidebar.getMessageNodeGroupKey(nodeInt)
		nodeIntName := manager.sidebar.getMessageNodeGroupName(nodeInt)

		manager.sidebar.sendUpdateName(newSidebarUpdateNameReq(nodeIntKey, nodeIntName))
	}
}

func (s *NodeService) GetInvalidNames(entityID string) []string {
	s.mux.Lock()
	defer s.mux.Unlock()

	names := []string{}
	for _, tmpNode := range s.pool {
		if tmpNode.EntityID() != acmelib.EntityID(entityID) {
			names = append(names, tmpNode.Name())
		}
	}

	return names
}

func (s *NodeService) GetInvalidIDs(entityID string) []uint {
	s.mux.Lock()
	defer s.mux.Unlock()

	ids := []uint{}
	for _, tmpNode := range s.pool {
		if tmpNode.EntityID() != acmelib.EntityID(entityID) {
			ids = append(ids, uint(tmpNode.ID()))
		}
	}

	return ids
}

func (s *NodeService) UpdateName(entityID, name string) (Node, error) {
	node, err := s.getEntity(entityID)
	if err != nil {
		return Node{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldName := node.Name()
	if name == oldName {
		return s.converterFn(node), nil
	}

	if err := node.UpdateName(name); err != nil {
		return Node{}, err
	}

	s.sendSidebarUpdateName(node)

	proxy.pushHistoryOperation(
		operationDomainNode,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := node.UpdateName(oldName); err != nil {
				return Node{}, err
			}

			s.sendSidebarUpdateName(node)

			return s.converterFn(node), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := node.UpdateName(name); err != nil {
				return Node{}, err
			}

			s.sendSidebarUpdateName(node)

			return s.converterFn(node), nil
		},
	)

	return s.converterFn(node), nil
}

func (s *NodeService) UpdateDesc(entityID, desc string) (Node, error) {
	node, err := s.getEntity(entityID)
	if err != nil {
		return Node{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldDesc := node.Desc()
	if desc == oldDesc {
		return s.converterFn(node), nil
	}

	node.SetDesc(desc)

	proxy.pushHistoryOperation(
		operationDomainNode,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			node.SetDesc(oldDesc)

			return s.converterFn(node), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			node.SetDesc(desc)

			return s.converterFn(node), nil
		},
	)

	return s.converterFn(node), nil
}

func (s *NodeService) UpdateID(entityID string, id uint) (Node, error) {
	node, err := s.getEntity(entityID)
	if err != nil {
		return Node{}, err
	}

	nodeID := acmelib.NodeID(id)

	s.mux.Lock()
	defer s.mux.Unlock()

	oldNodeID := node.ID()
	if nodeID == oldNodeID {
		return s.converterFn(node), nil
	}

	if err := node.UpdateID(nodeID); err != nil {
		return Node{}, err
	}

	proxy.pushHistoryOperation(
		operationDomainNode,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := node.UpdateID(oldNodeID); err != nil {
				return Node{}, err
			}

			return s.converterFn(node), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := node.UpdateID(nodeID); err != nil {
				return Node{}, err
			}

			return s.converterFn(node), nil
		},
	)

	return s.converterFn(node), nil
}

func (s *NodeService) AttachBus(nodeEntID string, intNumber int, busEntID string) (Node, error) {
	node, err := s.getEntity(nodeEntID)
	if err != nil {
		return Node{}, err
	}

	nodeInt, err := node.GetInterface(intNumber)
	if err != nil {
		return Node{}, err
	}

	bus, err := manager.busService.getEntity(busEntID)
	if err != nil {
		return Node{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldBus := nodeInt.ParentBus()
	if oldBus.EntityID() == bus.EntityID() {
		return s.converterFn(node), nil
	}

	if err := bus.AddNodeInterface(nodeInt); err != nil {
		return Node{}, err
	}

	proxy.pushHistoryOperation(
		operationDomainNode,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := bus.RemoveNodeInterface(nodeInt.Node().EntityID()); err != nil {
				return Node{}, err
			}

			if oldBus != nil {
				if err := oldBus.AddNodeInterface(nodeInt); err != nil {
					return Node{}, err
				}
			}

			return s.converterFn(node), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := bus.AddNodeInterface(nodeInt); err != nil {
				return Node{}, err
			}

			return s.converterFn(node), nil
		},
	)

	return s.converterFn(node), nil
}

func (s *NodeService) RemoveSentMessages(nodeEntID string, intNumber int, messageEntIDs ...string) (Node, error) {
	node, err := s.getEntity(nodeEntID)
	if err != nil {
		return Node{}, err
	}

	if len(messageEntIDs) == 0 {
		return s.converterFn(node), nil
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	nodeInt, err := node.GetInterface(intNumber)
	if err != nil {
		return Node{}, err
	}

	targetIDs := make(map[string]struct{})
	for _, messageEntID := range messageEntIDs {
		targetIDs[messageEntID] = struct{}{}
	}

	sentMessages := []*acmelib.Message{}
	for _, tmpMsg := range nodeInt.SentMessages() {
		tmpEntID := tmpMsg.EntityID()

		_, ok := targetIDs[tmpEntID.String()]
		if !ok {
			continue
		}

		sentMessages = append(sentMessages, tmpMsg)
	}

	for _, tmpMsg := range sentMessages {
		tmpEntID := tmpMsg.EntityID()

		if err := nodeInt.RemoveSentMessage(tmpEntID); err != nil {
			return Node{}, err
		}

		proxy.pushSidebarRemove(tmpEntID)
	}

	proxy.pushHistoryOperation(
		operationDomainNode,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			for _, tmpMsg := range sentMessages {
				if err := nodeInt.AddSentMessage(tmpMsg); err != nil {
					return Node{}, err
				}

				proxy.pushSidebarAdd(SidebarNodeKindMessage, tmpMsg.EntityID(), node.EntityID(), tmpMsg.Name())
			}

			return s.converterFn(node), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			for _, tmpMsg := range sentMessages {
				tmpEntID := tmpMsg.EntityID()

				if err := nodeInt.RemoveSentMessage(tmpMsg.EntityID()); err != nil {
					return Node{}, err
				}

				proxy.pushSidebarRemove(tmpEntID)
			}

			return s.converterFn(node), nil
		},
	)

	return s.converterFn(node), nil
}

func (s *NodeService) RemoveReceivedMessages(nodeEntID string, intNumber int, messageEntIDs ...string) (Node, error) {
	node, err := s.getEntity(nodeEntID)
	if err != nil {
		return Node{}, err
	}

	if len(messageEntIDs) == 0 {
		return s.converterFn(node), nil
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	nodeInt, err := node.GetInterface(intNumber)
	if err != nil {
		return Node{}, err
	}

	targetIDs := make(map[string]struct{})
	for _, messageEntID := range messageEntIDs {
		targetIDs[messageEntID] = struct{}{}
	}

	receivedMessages := []*acmelib.Message{}
	for _, tmpMsg := range nodeInt.ReceivedMessages() {
		tmpEntID := tmpMsg.EntityID()

		_, ok := targetIDs[tmpEntID.String()]
		if !ok {
			continue
		}

		receivedMessages = append(receivedMessages, tmpMsg)
	}

	for _, tmpMsg := range receivedMessages {
		if err := nodeInt.RemoveReceivedMessage(tmpMsg.EntityID()); err != nil {
			return Node{}, err
		}
	}

	proxy.pushHistoryOperation(
		operationDomainNode,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			for _, tmpMsg := range receivedMessages {
				if err := nodeInt.AddReceivedMessage(tmpMsg); err != nil {
					return Node{}, err
				}
			}

			return s.converterFn(node), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			for _, tmpMsg := range receivedMessages {
				if err := nodeInt.RemoveReceivedMessage(tmpMsg.EntityID()); err != nil {
					return Node{}, err
				}
			}

			return s.converterFn(node), nil
		},
	)

	return s.converterFn(node), nil
}
