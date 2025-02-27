package main

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/squadracorsepolito/acmelib"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type SidebarItemKind string

const (
	SidebarItemKindGroup         SidebarItemKind = "group"
	SidebarItemKindNetwork       SidebarItemKind = "network"
	SidebarItemKindBus           SidebarItemKind = "bus"
	SidebarItemKindNode          SidebarItemKind = "node"
	SidebarItemKindNodeInterface SidebarItemKind = "node-interface"
	SidebarItemKindMessage       SidebarItemKind = "message"
	SidebarItemKindSignal        SidebarItemKind = "signal"
	SidebarItemKindSignalType    SidebarItemKind = "signal-type"
	SidebarItemKindSignalUnit    SidebarItemKind = "signal-unit"
	SidebarItemKindSignalEnum    SidebarItemKind = "signal-enum"
)

type SidebarItem struct {
	Kind     SidebarItemKind `json:"kind"`
	ID       string          `json:"id"`
	Path     string          `json:"path"`
	Name     string          `json:"name"`
	Children []SidebarItem   `json:"children"`
}

type Sidebar struct {
	Root SidebarItem `json:"root"`
}

type sidebarItem struct {
	kind     SidebarItemKind
	id       string
	path     string
	name     string
	parent   *sidebarItem
	children []*sidebarItem
}

func newNodeIntSidebarItemID(nodeInt *acmelib.NodeInterface) string {
	return fmt.Sprintf("%s:%d", nodeInt.Node().EntityID(), nodeInt.Number())
}

func newNodeIntSidebarItemName(nodeInt *acmelib.NodeInterface) string {
	return fmt.Sprintf("%s:%d", nodeInt.Node().Name(), nodeInt.Number())
}

func newSidebarItem(kind SidebarItemKind, id, name string) *sidebarItem {
	return &sidebarItem{
		kind:     kind,
		id:       id,
		path:     id,
		name:     name,
		children: []*sidebarItem{},
	}
}

func (si *sidebarItem) getKey() string {
	splPath := strings.Split(si.path, "/")
	return splPath[len(splPath)-1]
}

func (si *sidebarItem) addChild(child *sidebarItem) {
	si.children = append(si.children, child)
	child.parent = si
	child.path = fmt.Sprintf("%s/%s", si.path, child.getKey())
}

func (si *sidebarItem) removeChild(child *sidebarItem) {
	si.children = slices.DeleteFunc(si.children, func(c *sidebarItem) bool {
		return c.kind == child.kind && c.id == child.id && c.name == child.name
	})
	child.parent = nil
	child.path = child.getKey()
}

func (si *sidebarItem) convertBase() SidebarItem {
	return SidebarItem{
		Kind:     si.kind,
		ID:       si.id,
		Path:     si.path,
		Name:     si.name,
		Children: []SidebarItem{},
	}
}

func (si *sidebarItem) convert() SidebarItem {
	res := si.convertBase()

	if si.kind != SidebarItemKindNetwork {
		slices.SortFunc(si.children, func(a, b *sidebarItem) int {
			return strings.Compare(a.name, b.name)
		})
	}

	for _, child := range si.children {
		res.Children = append(res.Children, child.convert())
	}

	return res
}

type sidebarLoadReq struct {
	network *acmelib.Network
}

func newSidebarLoadReq(network *acmelib.Network) *sidebarLoadReq {
	return &sidebarLoadReq{network: network}
}

type sidebarUpdateNameReq struct {
	itemKey string
	name    string
}

func newSidebarUpdateNameReq(itemKey string, name string) *sidebarUpdateNameReq {
	return &sidebarUpdateNameReq{itemKey: itemKey, name: name}
}

type sidebarAddReq struct {
	item          *sidebarItem
	parentItemKey string
}

func newSidebarAddReq(item *sidebarItem, parentItemKey string) *sidebarAddReq {
	return &sidebarAddReq{
		item:          item,
		parentItemKey: parentItemKey,
	}
}

type sidebarDeleteReq struct {
	itemKey string
}

func newSidebarDeleteReq(itemKey string) *sidebarDeleteReq {
	return &sidebarDeleteReq{itemKey: itemKey}
}

type SidebarService struct {
	items map[string]*sidebarItem
	root  *sidebarItem

	mux sync.RWMutex

	loadCh       chan *sidebarLoadReq
	updateNameCh chan *sidebarUpdateNameReq
	addCh        chan *sidebarAddReq
	deleteCh     chan *sidebarDeleteReq
}

func newSidebarService() *SidebarService {
	return &SidebarService{
		items: make(map[string]*sidebarItem),

		loadCh:       make(chan *sidebarLoadReq),
		updateNameCh: make(chan *sidebarUpdateNameReq),
		addCh:        make(chan *sidebarAddReq),
		deleteCh:     make(chan *sidebarDeleteReq),
	}
}

func (s *SidebarService) run(ctx context.Context) {
	for {
		select {
		case req := <-s.loadCh:
			s.handleLoad(req)

		case req := <-s.updateNameCh:
			s.handleUpdateName(req)

		case req := <-s.addCh:
			s.handleAdd(req)

		case req := <-s.deleteCh:
			s.handleDelete(req)

		case <-ctx.Done():
			return
		}
	}
}

func (s *SidebarService) clear() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.root = nil
	clear(s.items)
}

func (s *SidebarService) OnStartup(ctx context.Context, _ application.ServiceOptions) error {
	go s.run(ctx)
	return nil
}

func (s *SidebarService) OnShutdown() {}

func (s *SidebarService) addItem(item *sidebarItem) {
	s.items[item.getKey()] = item
}

func (s *SidebarService) sendLoad(req *sidebarLoadReq) {
	s.loadCh <- req
}

func (s *SidebarService) handleLoad(req *sidebarLoadReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	clear(s.items)

	nodes := make(map[acmelib.EntityID]*acmelib.Node)
	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	net := req.network

	// the network is the root item
	netItem := newSidebarItem(SidebarItemKindNetwork, net.EntityID().String(), net.Name())
	s.addItem(netItem)
	s.root = netItem

	// it groups all the nodes
	nodeGroupItem := newSidebarItem(SidebarItemKindGroup, SidebarNodeGroupID, "Nodes")
	s.addItem(nodeGroupItem)
	netItem.addChild(nodeGroupItem)

	// it groups all the signal types
	sigTypeGroupItem := newSidebarItem(SidebarItemKindGroup, SidebarSignalTypeGroupID, "Signal Types")
	s.addItem(sigTypeGroupItem)
	netItem.addChild(sigTypeGroupItem)

	// it groups all the signal units
	sigUnitGroupItem := newSidebarItem(SidebarItemKindGroup, SidebarSignalUnitGroupID, "Signal Units")
	s.addItem(sigUnitGroupItem)
	netItem.addChild(sigUnitGroupItem)

	// it groups all the signal enums
	sigEnumGroupItem := newSidebarItem(SidebarItemKindGroup, SidebarSignalEnumGroupID, "Signal Enums")
	s.addItem(sigEnumGroupItem)
	netItem.addChild(sigEnumGroupItem)

	for _, bus := range net.Buses() {
		// add the bus
		busItem := newSidebarItem(SidebarItemKindBus, bus.EntityID().String(), bus.Name())
		s.addItem(busItem)
		netItem.addChild(busItem)

		for _, nodeInt := range bus.NodeInterfaces() {
			node := nodeInt.Node()
			nodes[node.EntityID()] = node

			// add the node interface to the bus
			nodeIntItem := newSidebarItem(SidebarItemKindNodeInterface, newNodeIntSidebarItemID(nodeInt), newNodeIntSidebarItemName(nodeInt))
			s.addItem(nodeIntItem)
			busItem.addChild(nodeIntItem)

			for _, msg := range nodeInt.SentMessages() {
				// add the message to the node interface
				msgItem := newSidebarItem(SidebarItemKindMessage, msg.EntityID().String(), msg.Name())
				s.addItem(msgItem)
				nodeIntItem.addChild(msgItem)

				for _, sig := range msg.Signals() {
					// add the signal to the message
					sigItem := newSidebarItem(SidebarItemKindSignal, sig.EntityID().String(), sig.Name())
					s.addItem(sigItem)
					msgItem.addChild(sigItem)

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
			}
		}
	}

	// add nodes
	for _, node := range nodes {
		nodeItem := newSidebarItem(SidebarItemKindNode, node.EntityID().String(), node.Name())
		s.addItem(nodeItem)
		nodeGroupItem.addChild(nodeItem)
	}

	// add signal types
	for _, sigType := range sigTypes {
		sigTypeItem := newSidebarItem(SidebarItemKindSignalType, sigType.EntityID().String(), sigType.Name())
		s.addItem(sigTypeItem)
		sigTypeGroupItem.addChild(sigTypeItem)
	}

	// add signal units
	for _, sigUnit := range sigUnits {
		sigUnitItem := newSidebarItem(SidebarItemKindSignalUnit, sigUnit.EntityID().String(), sigUnit.Name())
		s.addItem(sigUnitItem)
		sigUnitGroupItem.addChild(sigUnitItem)
	}

	// add signal enums
	for _, sigEnum := range sigEnums {
		sigEnumItem := newSidebarItem(SidebarItemKindSignalEnum, sigEnum.EntityID().String(), sigEnum.Name())
		s.addItem(sigEnumItem)
		sigEnumGroupItem.addChild(sigEnumItem)
	}

	app.EmitEvent(SidebarLoad)
}

func (s *SidebarService) handleUpdateName(req *sidebarUpdateNameReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	item, ok := s.items[req.itemKey]
	if !ok {
		return
	}

	item.name = req.name

	app.EmitEvent(SidebarUpdateName, item.convertBase())
}

func (s *SidebarService) handleAdd(req *sidebarAddReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	parent, ok := s.items[req.parentItemKey]
	if !ok {
		return
	}

	s.addItem(req.item)
	parent.addChild(req.item)

	app.EmitEvent(SidebarAdd, parent.convert())
}

func (s *SidebarService) handleDelete(req *sidebarDeleteReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	item, ok := s.items[req.itemKey]
	if !ok {
		return
	}

	parent := item.parent
	parent.removeChild(item)

	delete(s.items, req.itemKey)

	app.EmitEvent(SidebarRemove, parent.convert())
}

func (s *SidebarService) getController() *sidebarController {
	return &sidebarController{
		updateNameCh: s.updateNameCh,
		addCh:        s.addCh,
		deleteCh:     s.deleteCh,
	}
}

func (s *SidebarService) Get() Sidebar {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if s.root == nil {
		return Sidebar{Root: SidebarItem{}}
	}

	return Sidebar{Root: s.root.convert()}
}

type sidebarController struct {
	updateNameCh chan<- *sidebarUpdateNameReq
	addCh        chan<- *sidebarAddReq
	deleteCh     chan<- *sidebarDeleteReq
}

func (s *sidebarController) sendUpdateName(ent entity) {
	itemKey := ent.EntityID().String()
	name := ent.Name()

	s.updateNameCh <- newSidebarUpdateNameReq(itemKey, name)

	switch ent.EntityKind() {
	case acmelib.EntityKindNode:
		node, ok := ent.(*acmelib.Node)
		if !ok {
			panic("entity is not an acmelib.Node")
		}

		for _, nodeInt := range node.Interfaces() {
			if nodeInt.ParentBus() == nil {
				continue
			}

			s.updateNameCh <- newSidebarUpdateNameReq(newNodeIntSidebarItemID(nodeInt), newNodeIntSidebarItemName(nodeInt))
		}
	}
}

func (s *sidebarController) sendAdd(ent entity) {
	switch ent.EntityKind() {
	case acmelib.EntityKindBus:
		bus, ok := ent.(*acmelib.Bus)
		if !ok {
			panic("entity is not an acmelib.Bus")
		}

		parNet := bus.ParentNetwork()
		if parNet == nil {
			return
		}

		busItem := newSidebarItem(SidebarItemKindBus, bus.EntityID().String(), bus.Name())
		s.addCh <- newSidebarAddReq(busItem, parNet.EntityID().String())

	case acmelib.EntityKindNode:
		node, ok := ent.(*acmelib.Node)
		if !ok {
			panic("entity is not an acmelib.Node")
		}

		nodeItem := newSidebarItem(SidebarItemKindNode, node.EntityID().String(), node.Name())
		s.addCh <- newSidebarAddReq(nodeItem, node.EntityID().String())

		for _, nodeInt := range node.Interfaces() {
			parBus := nodeInt.ParentBus()
			if parBus == nil {
				continue
			}

			nodeIntItem := newSidebarItem(SidebarItemKindNodeInterface, newNodeIntSidebarItemID(nodeInt), newNodeIntSidebarItemName(nodeInt))
			s.addCh <- newSidebarAddReq(nodeIntItem, parBus.EntityID().String())
		}

	case acmelib.EntityKindMessage:
		msg, ok := ent.(*acmelib.Message)
		if !ok {
			panic("entity is not an acmelib.Message")
		}

		parNodeInt := msg.SenderNodeInterface()
		if parNodeInt == nil {
			return
		}

		msgItem := newSidebarItem(SidebarItemKindMessage, msg.EntityID().String(), msg.Name())
		s.addCh <- newSidebarAddReq(msgItem, newNodeIntSidebarItemID(parNodeInt))

	case acmelib.EntityKindSignal:
		sig, ok := ent.(acmelib.Signal)
		if !ok {
			panic("entity is not an acmelib.Signal")
		}

		parMsg := sig.ParentMessage()
		if parMsg == nil {
			return
		}

		sigItem := newSidebarItem(SidebarItemKindSignal, sig.EntityID().String(), sig.Name())
		s.addCh <- newSidebarAddReq(sigItem, parMsg.EntityID().String())

	case acmelib.EntityKindSignalType:
		sigTypeItem := newSidebarItem(SidebarItemKindSignalType, ent.EntityID().String(), ent.Name())
		s.addCh <- newSidebarAddReq(sigTypeItem, SidebarSignalTypeGroupID)

	case acmelib.EntityKindSignalUnit:
		sigUnitItem := newSidebarItem(SidebarItemKindSignalUnit, ent.EntityID().String(), ent.Name())
		s.addCh <- newSidebarAddReq(sigUnitItem, SidebarSignalUnitGroupID)

	case acmelib.EntityKindSignalEnum:
		sigEnumItem := newSidebarItem(SidebarItemKindSignalEnum, ent.EntityID().String(), ent.Name())
		s.addCh <- newSidebarAddReq(sigEnumItem, SidebarSignalEnumGroupID)
	}
}

func (s *sidebarController) sendDelete(ent entity) {
	s.deleteCh <- newSidebarDeleteReq(ent.EntityID().String())

	switch ent.EntityKind() {
	case acmelib.EntityKindNode:
		node, ok := ent.(*acmelib.Node)
		if !ok {
			panic("entity is not an acmelib.Node")
		}

		for _, nodeInt := range node.Interfaces() {
			if nodeInt.ParentBus() == nil {
				continue
			}

			s.deleteCh <- newSidebarDeleteReq(newNodeIntSidebarItemID(nodeInt))
		}
	}
}
