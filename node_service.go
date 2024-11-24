package main

import "github.com/squadracorsepolito/acmelib"

type NodeMessage struct {
	base
}

func getNodeMessage(msg *acmelib.Message) NodeMessage {
	return NodeMessage{
		base: getBase(msg),
	}
}

type NodeInterfaceBus struct {
	base
}

func getNodeInterfaceBus(bus *acmelib.Bus) NodeInterfaceBus {
	if bus == nil {
		return NodeInterfaceBus{}
	}

	return NodeInterfaceBus{
		base: getBase(bus),
	}
}

type NodeInterface struct {
	base

	Number           int              `json:"number"`
	Bus              NodeInterfaceBus `json:"bus"`
	SentMessages     []NodeMessage    `json:"sentMessages"`
	ReceivedMessages []NodeMessage    `json:"receivedMessages"`
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
		base: getBase(nodeInt.Node()),

		Bus:              getNodeInterfaceBus(nodeInt.ParentBus()),
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
		Interfaces: make([]NodeInterface, len(node.Interfaces())),
	}
}

func newNodeService() *NodeService {
	return &NodeService{
		service: newService(proxy.nodeCh, getNode),
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

	proxy.pushSidebarUpdate(node.EntityID(), name)

	proxy.pushHistoryOperation(
		operationDomainNode,
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := node.UpdateName(oldName); err != nil {
				return Node{}, err
			}

			proxy.pushSidebarUpdate(node.EntityID(), oldName)

			return s.converterFn(node), nil
		},
		func() (any, error) {
			s.mux.Lock()
			defer s.mux.Unlock()

			if err := node.UpdateName(name); err != nil {
				return Node{}, err
			}

			proxy.pushSidebarUpdate(node.EntityID(), name)

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
