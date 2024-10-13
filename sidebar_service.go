package main

import (
	"context"
	"log"
	"sync"

	"github.com/squadracorsepolito/acmelib"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type sidebarLoadReq struct {
	network *acmelib.Network
}

type sidebarUpdateReq struct {
	entityID acmelib.EntityID
	name     string
}

type sidebarRemoveReq struct {
	entityID acmelib.EntityID
}

type SidebarNodeKind string

const (
	SidebarNodeKindNetwork    SidebarNodeKind = "network"
	SidebarNodeKindBus        SidebarNodeKind = "bus"
	SidebarNodeKindNode       SidebarNodeKind = "node"
	SidebarNodeKindMessage    SidebarNodeKind = "message"
	SidebarNodeKindSignalType SidebarNodeKind = "signal-type"
	SidebarNodeKindSignalUnit SidebarNodeKind = "signal-unit"
	SidebarNodeKindSignalEnum SidebarNodeKind = "signal-enum"
)

type SidebarNode struct {
	Kind     SidebarNodeKind `json:"kind"`
	Name     string          `json:"name"`
	EntityID string          `json:"entityId"`
	Children []SidebarNode   `json:"children"`
}

type sidebarNode struct {
	kind     SidebarNodeKind
	name     string
	entityID acmelib.EntityID
	children []*sidebarNode
}

func (sn *sidebarNode) Convert() SidebarNode {
	res := SidebarNode{
		Kind:     sn.kind,
		Name:     sn.name,
		EntityID: sn.entityID.String(),
	}

	for _, child := range sn.children {
		res.Children = append(res.Children, child.Convert())
	}

	return res
}

type SidebarService struct {
	tree  *sidebarNode
	nodes map[acmelib.EntityID]*sidebarNode

	mux sync.RWMutex

	stopCh chan struct{}
}

func newSidebarService() *SidebarService {

	return &SidebarService{
		nodes: make(map[acmelib.EntityID]*sidebarNode),

		stopCh: make(chan struct{}),
	}
}

func (s *SidebarService) run() {
	for {
		select {
		case req := <-proxy.sidebarLoadCh:
			s.load(req)

		case req := <-proxy.sidebarUpdateCh:
			s.update(req)

		case req := <-proxy.sidebarRemoveCh:
			log.Print(req.entityID.String())

		case <-s.stopCh:
			return
		}
	}
}

func (s *SidebarService) OnStartup(_ context.Context, _ application.ServiceOptions) error {
	go s.run()
	return nil
}

func (s *SidebarService) OnShutdown() {
	s.stopCh <- struct{}{}
}

func (s *SidebarService) load(req *sidebarLoadReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	clear(s.nodes)

	netNode := &sidebarNode{
		kind:     SidebarNodeKindNetwork,
		name:     req.network.Name(),
		entityID: req.network.EntityID(),
	}

	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	for _, bus := range req.network.Buses() {
		busNode := &sidebarNode{
			kind:     SidebarNodeKindBus,
			name:     bus.Name(),
			entityID: bus.EntityID(),
		}

		for _, nodeInt := range bus.NodeInterfaces() {
			nodeNode := &sidebarNode{
				kind:     SidebarNodeKindNode,
				name:     nodeInt.Node().Name(),
				entityID: nodeInt.Node().EntityID(),
			}

			for _, msg := range nodeInt.Messages() {
				msgNode := &sidebarNode{
					kind:     SidebarNodeKindMessage,
					name:     msg.Name(),
					entityID: msg.EntityID(),
				}

				for _, sig := range msg.Signals() {
					switch sig.Kind() {
					case acmelib.SignalKindStandard:
						stdSig, err := sig.ToStandard()
						if err != nil {
							panic(err)
						}

						sigTypes[stdSig.Type().EntityID()] = stdSig.Type()

						if stdSig.Unit() != nil {
							sigUnits[stdSig.Unit().EntityID()] = stdSig.Unit()
						}

					case acmelib.SignalKindEnum:
						enumSig, err := sig.ToEnum()
						if err != nil {
							panic(err)
						}

						sigEnums[enumSig.Enum().EntityID()] = enumSig.Enum()
					}
				}

				nodeNode.children = append(nodeNode.children, msgNode)
				s.nodes[msg.EntityID()] = msgNode
			}

			busNode.children = append(busNode.children, nodeNode)
			s.nodes[nodeInt.Node().EntityID()] = nodeNode
		}

		netNode.children = append(netNode.children, busNode)
		s.nodes[bus.EntityID()] = busNode
	}

	for _, sigType := range sigTypes {
		sigTypeNode := &sidebarNode{
			kind:     SidebarNodeKindSignalType,
			name:     sigType.Name(),
			entityID: sigType.EntityID(),
		}

		netNode.children = append(netNode.children, sigTypeNode)
		s.nodes[sigType.EntityID()] = sigTypeNode
	}

	for _, sigUnit := range sigUnits {
		sigUnitNode := &sidebarNode{
			kind:     SidebarNodeKindSignalUnit,
			name:     sigUnit.Name(),
			entityID: sigUnit.EntityID(),
		}

		netNode.children = append(netNode.children, sigUnitNode)
		s.nodes[sigUnit.EntityID()] = sigUnitNode
	}

	for _, sigEnum := range sigEnums {
		sigEnumNode := &sidebarNode{
			kind:     SidebarNodeKindSignalEnum,
			name:     sigEnum.Name(),
			entityID: sigEnum.EntityID(),
		}

		netNode.children = append(netNode.children, sigEnumNode)
		s.nodes[sigEnum.EntityID()] = sigEnumNode
	}

	s.tree = netNode

	app.EmitEvent(SidebarLoad)
}

func (s *SidebarService) update(req *sidebarUpdateReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	node := s.nodes[req.entityID]
	node.name = req.name

	app.EmitEvent(SidebarUpdate, node.Convert())
}

func (s *SidebarService) GetTree() SidebarNode {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.tree.Convert()
}
