package main

import (
	"context"
	"log"
	"slices"
	"strings"
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

type sidebarAddReq struct {
	kind     SidebarNodeKind
	entityID acmelib.EntityID
	name     string
	parentID acmelib.EntityID
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
	Kind            SidebarNodeKind `json:"kind"`
	Name            string          `json:"name"`
	EntityID        string          `json:"entityId"`
	ParentEntityIDs []string        `json:"parentEntityIds"`
	Children        []SidebarNode   `json:"children"`
}

type sidebarNode struct {
	kind     SidebarNodeKind
	name     string
	entityID acmelib.EntityID
	parent   *sidebarNode
	children []*sidebarNode
}

func (sn *sidebarNode) Convert(parentsIDs ...string) SidebarNode {
	res := SidebarNode{
		Kind:            sn.kind,
		Name:            sn.name,
		EntityID:        sn.entityID.String(),
		ParentEntityIDs: parentsIDs,
	}

	parentsIDs = append(parentsIDs, sn.entityID.String())
	slices.SortFunc(sn.children, func(a, b *sidebarNode) int { return strings.Compare(a.name, b.name) })
	for _, child := range sn.children {
		res.Children = append(res.Children, child.Convert(parentsIDs...))
	}

	return res
}

type SidebarService struct {
	root  *sidebarNode
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

		case req := <-proxy.sidebarAddCh:
			log.Print(req)
			s.add(req)

		case req := <-proxy.sidebarRemoveCh:
			s.remove(req)

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

	s.nodes[req.network.EntityID()] = netNode

	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	for _, bus := range req.network.Buses() {
		busNode := &sidebarNode{
			kind:     SidebarNodeKindBus,
			name:     bus.Name(),
			entityID: bus.EntityID(),
			parent:   netNode,
		}

		for _, nodeInt := range bus.NodeInterfaces() {
			nodeNode := &sidebarNode{
				kind:     SidebarNodeKindNode,
				name:     nodeInt.Node().Name(),
				entityID: nodeInt.Node().EntityID(),
				parent:   busNode,
			}

			for _, msg := range nodeInt.Messages() {
				msgNode := &sidebarNode{
					kind:     SidebarNodeKindMessage,
					name:     msg.Name(),
					entityID: msg.EntityID(),
					parent:   nodeNode,
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
			parent:   netNode,
		}

		netNode.children = append(netNode.children, sigTypeNode)
		s.nodes[sigType.EntityID()] = sigTypeNode
	}

	for _, sigUnit := range sigUnits {
		sigUnitNode := &sidebarNode{
			kind:     SidebarNodeKindSignalUnit,
			name:     sigUnit.Name(),
			entityID: sigUnit.EntityID(),
			parent:   netNode,
		}

		netNode.children = append(netNode.children, sigUnitNode)
		s.nodes[sigUnit.EntityID()] = sigUnitNode
	}

	for _, sigEnum := range sigEnums {
		sigEnumNode := &sidebarNode{
			kind:     SidebarNodeKindSignalEnum,
			name:     sigEnum.Name(),
			entityID: sigEnum.EntityID(),
			parent:   netNode,
		}

		netNode.children = append(netNode.children, sigEnumNode)
		s.nodes[sigEnum.EntityID()] = sigEnumNode
	}

	s.root = netNode

	app.EmitEvent(SidebarLoad)
}

func (s *SidebarService) update(req *sidebarUpdateReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	node, ok := s.nodes[req.entityID]
	if !ok {
		return
	}

	node.name = req.name

	app.EmitEvent(SidebarUpdate, node.Convert())
}

func (s *SidebarService) add(req *sidebarAddReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	parent, ok := s.nodes[req.parentID]
	if !ok {
		return
	}

	parent.children = append(parent.children, &sidebarNode{
		kind:     req.kind,
		name:     req.name,
		entityID: req.entityID,
		parent:   parent,
	})

	s.nodes[req.entityID] = &sidebarNode{
		kind:     req.kind,
		name:     req.name,
		entityID: req.entityID,
		parent:   parent,
	}

	app.EmitEvent(SidebarAdd, parent.Convert())
}

func (s *SidebarService) remove(req *sidebarRemoveReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	node, ok := s.nodes[req.entityID]
	if !ok {
		return
	}

	parent := node.parent
	parent.children = slices.DeleteFunc(parent.children, func(child *sidebarNode) bool {
		return child.entityID == req.entityID
	})

	delete(s.nodes, req.entityID)

	app.EmitEvent(SidebarRemove, parent.Convert())
}

func (s *SidebarService) GetTree() SidebarNode {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.root.Convert()
}
