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
	SidebarItemKindGroup      SidebarItemKind = "group"
	SidebarItemKindNetwork    SidebarItemKind = "network"
	SidebarItemKindBus        SidebarItemKind = "bus"
	SidebarItemKindNode       SidebarItemKind = "node"
	SidebarItemKindMessage    SidebarItemKind = "message"
	SidebarItemKindSignalType SidebarItemKind = "signal-type"
	SidebarItemKindSignalUnit SidebarItemKind = "signal-unit"
	SidebarItemKindSignalEnum SidebarItemKind = "signal-enum"
)

type SidebarItem struct {
	Kind     SidebarItemKind `json:"kind"`
	Prefix   string          `json:"prefix"`
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Children []SidebarItem   `json:"children"`
}

type Sidebar struct {
	Root SidebarItem `json:"root"`
}

type sidebarItem struct {
	kind     SidebarItemKind
	prefix   string
	entityID acmelib.EntityID
	name     string
	parent   *sidebarItem
	children []*sidebarItem
}

func newSidebarItem(kind SidebarItemKind, entityID acmelib.EntityID, prefix, name string) *sidebarItem {
	return &sidebarItem{
		kind:     kind,
		entityID: entityID,
		prefix:   prefix,
		name:     name,
		children: []*sidebarItem{},
	}
}

func newBusSidebarItem(bus *acmelib.Bus) *sidebarItem {
	return newSidebarItem(SidebarItemKindBus, bus.EntityID(), SidebarBusesPrefix, bus.Name())
}

func newNodeSidebarItem(node *acmelib.Node) *sidebarItem {
	return newSidebarItem(SidebarItemKindNode, node.EntityID(), SidebarNodesPrefix, node.Name())
}

func newMessageSidebarItem(message *acmelib.Message) *sidebarItem {
	return newSidebarItem(SidebarItemKindMessage, message.EntityID(), SidebarMessagesPrefix, message.Name())
}

func newSignalTypeSidebarItem(sigType *acmelib.SignalType) *sidebarItem {
	return newSidebarItem(SidebarItemKindSignalType, sigType.EntityID(), SidebarSignalTypesPrefix, sigType.Name())
}

func newSignalUnitSidebarItem(sigUnit *acmelib.SignalUnit) *sidebarItem {
	return newSidebarItem(SidebarItemKindSignalUnit, sigUnit.EntityID(), SidebarSignalUnitsPrefix, sigUnit.Name())
}

func newSignalEnumSidebarItem(sigEnum *acmelib.SignalEnum) *sidebarItem {
	return newSidebarItem(SidebarItemKindSignalEnum, sigEnum.EntityID(), SidebarSignalEnumsPrefix, sigEnum.Name())
}

func newGroupSidebarItem(id, name string) *sidebarItem {
	return newSidebarItem(SidebarItemKindGroup, acmelib.EntityID(id), "", name)
}

func newMessageBusGroupSidebarItem(bus *acmelib.Bus) *sidebarItem {
	return newSidebarItem(SidebarItemKindGroup, bus.EntityID(), SidebarMessagesPrefix, bus.Name())
}

func newMessageNodeGroupSidebarItem(nodeInt *acmelib.NodeInterface) *sidebarItem {
	busEntID := nodeInt.ParentBus().EntityID()
	nodeEntID := nodeInt.Node().EntityID()

	prefix := fmt.Sprintf("%s:%s", SidebarMessagesPrefix, busEntID)
	name := fmt.Sprintf("%s:%d", nodeInt.Node().Name(), nodeInt.Number())

	return newSidebarItem(SidebarItemKindGroup, nodeEntID, prefix, name)
}

func (si *sidebarItem) addChild(child *sidebarItem) {
	si.children = append(si.children, child)
	child.parent = si
}

func (si *sidebarItem) removeChild(child *sidebarItem) {
	si.children = slices.DeleteFunc(si.children, func(c *sidebarItem) bool {
		return c.kind == child.kind && c.entityID == child.entityID && c.name == child.name
	})
	child.parent = nil
}

func (si *sidebarItem) convert() SidebarItem {
	base := si.convertBase()

	if si.kind != SidebarItemKindNetwork {
		slices.SortFunc(si.children, func(a, b *sidebarItem) int {
			return strings.Compare(a.name, b.name)
		})
	}

	for _, child := range si.children {
		base.Children = append(base.Children, child.convert())
	}

	return base
}

func (si *sidebarItem) convertBase() SidebarItem {
	id := ""
	if si.kind == SidebarItemKindGroup {
		id = fmt.Sprintf("%s:%s", SidebarGroupPrefix, si.entityID)
	} else {
		id = si.entityID.String()
	}

	return SidebarItem{
		Kind:     si.kind,
		Prefix:   si.prefix,
		ID:       id,
		Name:     si.name,
		Children: []SidebarItem{},
	}
}

type SidebarService0 struct {
	items map[string]*sidebarItem
	root  *sidebarItem

	mux sync.RWMutex

	loadCh       chan *sidebarLoadReq
	updateNameCh chan *sidebarUpdateNameReq
	addCh        chan *sidebarAddReq0
	deleteCh     chan *sidebarDeleteReq

	stopCh chan struct{}
}

func newSidebarService0() *SidebarService0 {
	return &SidebarService0{
		items: make(map[string]*sidebarItem),

		loadCh:       make(chan *sidebarLoadReq),
		updateNameCh: make(chan *sidebarUpdateNameReq),
		addCh:        make(chan *sidebarAddReq0),
		deleteCh:     make(chan *sidebarDeleteReq),

		stopCh: make(chan struct{}),
	}
}

func (s *SidebarService0) run() {
	for {
		select {
		case req := <-s.loadCh:
			s.load(req)

		case req := <-s.updateNameCh:
			s.updateName(req)

		case req := <-s.addCh:
			s.add(req)

		case req := <-s.deleteCh:
			s.delete(req)

		case <-s.stopCh:
			return
		}
	}
}

func (s *SidebarService0) OnStartup(_ context.Context, _ application.ServiceOptions) error {
	go s.run()
	return nil
}

func (s *SidebarService0) OnShutdown() {
	s.stopCh <- struct{}{}
}

func (s *SidebarService0) addItem(item *sidebarItem) {
	s.items[item.entityID.String()] = item
}

func (s *SidebarService0) getMessageBusGroupKey(bus *acmelib.Bus) string {
	return fmt.Sprintf("%s:message", bus.EntityID())
}

func (s *SidebarService0) getMessageNodeGroupKey(nodeInt *acmelib.NodeInterface) string {
	return fmt.Sprintf("%s:%d:message", nodeInt.Node().EntityID(), nodeInt.Number())
}

func (s *SidebarService0) getMessageNodeGroupName(nodeInt *acmelib.NodeInterface) string {
	return fmt.Sprintf("%s:%d", nodeInt.Node().Name(), nodeInt.Number())
}

func newSidebarLoadReq(network *acmelib.Network) *sidebarLoadReq {
	return &sidebarLoadReq{network: network}
}

func (s *SidebarService0) sendLoad(req *sidebarLoadReq) {
	s.loadCh <- req
}

func (s *SidebarService0) load(req *sidebarLoadReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	clear(s.items)

	net := req.network

	nodes := make(map[acmelib.EntityID]*acmelib.Node)
	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	// the network is the root item
	netItem := newSidebarItem(SidebarItemKindNetwork, net.EntityID(), "", net.Name())
	s.addItem(netItem)
	s.root = netItem

	// it groups all the buses
	busGroupItem := newGroupSidebarItem(SidebarBusesPrefix, "Buses")
	s.addItem(busGroupItem)
	netItem.addChild(busGroupItem)

	// it groups all the nodes
	nodeGroupItem := newGroupSidebarItem(SidebarNodesPrefix, "Nodes")
	s.addItem(nodeGroupItem)
	netItem.addChild(nodeGroupItem)

	// it groups all the messages
	msgGroupItem := newGroupSidebarItem(SidebarMessagesPrefix, "Messages")
	s.addItem(msgGroupItem)
	netItem.addChild(msgGroupItem)

	// it groups all the signal types
	sigTypeGroupItem := newGroupSidebarItem(SidebarSignalTypesPrefix, "Signal Types")
	s.addItem(sigTypeGroupItem)
	netItem.addChild(sigTypeGroupItem)

	// it groups all the signal units
	sigUnitGroupItem := newGroupSidebarItem(SidebarSignalUnitsPrefix, "Signal Units")
	s.addItem(sigUnitGroupItem)
	netItem.addChild(sigUnitGroupItem)

	// it groups all the signal enums
	sigEnumGroupItem := newGroupSidebarItem(SidebarSignalEnumsPrefix, "Signal Enums")
	s.addItem(sigEnumGroupItem)
	netItem.addChild(sigEnumGroupItem)

	// add buses and nodes
	for _, bus := range net.Buses() {
		busItem := newBusSidebarItem(bus)
		s.addItem(busItem)
		busGroupItem.addChild(busItem)

		for _, nodeInt := range bus.NodeInterfaces() {
			node := nodeInt.Node()
			nodeItem := newNodeSidebarItem(node)
			s.addItem(nodeItem)
			nodeGroupItem.addChild(nodeItem)

			nodes[node.EntityID()] = node
		}

		// add bus group for messages
		msgBusGroupItem := newMessageBusGroupSidebarItem(bus)
		s.items[s.getMessageBusGroupKey(bus)] = msgBusGroupItem
		msgGroupItem.addChild(msgBusGroupItem)
	}

	// add messages
	for _, node := range nodes {
		for _, nodeInt := range node.Interfaces() {
			// add node group for messages
			msgNodeGroupItem := newMessageNodeGroupSidebarItem(nodeInt)
			s.items[s.getMessageNodeGroupKey(nodeInt)] = msgNodeGroupItem

			if nodeInt.ParentBus() == nil {
				continue
			}

			// add node group into bus group for messages
			msgBusGroupItem := s.items[s.getMessageBusGroupKey(nodeInt.ParentBus())]
			msgBusGroupItem.addChild(msgNodeGroupItem)

			for _, msg := range nodeInt.SentMessages() {
				msgItem := newMessageSidebarItem(msg)
				s.addItem(msgItem)
				msgNodeGroupItem.addChild(msgItem)

				// selecting signal types/units/enums
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
			}
		}
	}

	// add signal types
	for _, sigType := range sigTypes {
		sigTypeItem := newSignalTypeSidebarItem(sigType)
		s.addItem(sigTypeItem)
		sigTypeGroupItem.addChild(sigTypeItem)
	}

	// add signal units
	for _, sigUnit := range sigUnits {
		sigUnitItem := newSignalUnitSidebarItem(sigUnit)
		s.addItem(sigUnitItem)
		sigUnitGroupItem.addChild(sigUnitItem)
	}

	// add signal enums
	for _, sigEnum := range sigEnums {
		sigEnumItem := newSignalEnumSidebarItem(sigEnum)
		s.addItem(sigEnumItem)
		sigEnumGroupItem.addChild(sigEnumItem)
	}
}

type sidebarUpdateNameReq struct {
	itemKey string
	name    string
}

func newSidebarUpdateNameReq(itemKey string, name string) *sidebarUpdateNameReq {
	return &sidebarUpdateNameReq{itemKey: itemKey, name: name}
}

func (s *SidebarService0) sendUpdateName(req *sidebarUpdateNameReq) {
	s.updateNameCh <- req
}

func (s *SidebarService0) updateName(req *sidebarUpdateNameReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	item, ok := s.items[req.itemKey]
	if !ok {
		return
	}

	item.name = req.name

	app.EmitEvent(SidebarUpdateName, item.convertBase())
}

type sidebarAddReq0 struct {
	item          *sidebarItem
	parentItemKey string
}

func newSidebarAddReq(item *sidebarItem, parentItemKey string) *sidebarAddReq0 {
	return &sidebarAddReq0{
		item:          item,
		parentItemKey: parentItemKey,
	}
}

func (s *SidebarService0) sendAdd(req *sidebarAddReq0) {
	s.addCh <- req
}

func (s *SidebarService0) add(req *sidebarAddReq0) {
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

type sidebarDeleteReq struct {
	itemKey string
}

func newSidebarDeleteReq(itemKey string) *sidebarDeleteReq {
	return &sidebarDeleteReq{itemKey: itemKey}
}

func (s *SidebarService0) sendDelete(req *sidebarDeleteReq) {
	s.deleteCh <- req
}

func (s *SidebarService0) delete(req *sidebarDeleteReq) {
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

func (s *SidebarService0) Get() Sidebar {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if s.root == nil {
		return Sidebar{Root: SidebarItem{}}
	}

	return Sidebar{Root: s.root.convert()}
}
