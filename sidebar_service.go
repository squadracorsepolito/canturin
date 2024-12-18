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

type sidebarItemFactory struct{}

func (f *sidebarItemFactory) newItem(kind SidebarItemKind, entityID acmelib.EntityID, prefix, name string) *sidebarItem {
	return &sidebarItem{
		kind:     kind,
		entityID: entityID,
		prefix:   prefix,
		name:     name,
		children: []*sidebarItem{},
	}
}

func (f *sidebarItemFactory) newGroup(prefix, name string) (string, *sidebarItem) {
	return prefix, f.newItem(SidebarItemKindGroup, acmelib.EntityID(prefix), "", name)
}

func (f *sidebarItemFactory) newNetwork(net entity) (string, *sidebarItem) {
	return net.EntityID().String(), f.newItem(SidebarItemKindNetwork, net.EntityID(), "", net.Name())
}

func (f *sidebarItemFactory) newBus(bus entity) (string, *sidebarItem) {
	return bus.EntityID().String(), f.newItem(SidebarItemKindBus, bus.EntityID(), SidebarBusesPrefix, bus.Name())
}

func (f *sidebarItemFactory) newMessageBus(bus entity) (string, *sidebarItem) {
	key := f.getMessageBusKey(bus)
	item := f.newItem(SidebarItemKindBus, bus.EntityID(), SidebarMessagesPrefix, bus.Name())
	return key, item
}

func (f *sidebarItemFactory) newNode(node entity) (string, *sidebarItem) {
	return node.EntityID().String(), f.newItem(SidebarItemKindNode, node.EntityID(), SidebarNodesPrefix, node.Name())
}

func (f *sidebarItemFactory) newMessageNode(nodeInt *acmelib.NodeInterface) (string, *sidebarItem) {
	parBus := nodeInt.ParentBus()
	prefix := ""
	if parBus != nil {
		prefix = fmt.Sprintf("%s:%s", SidebarMessagesPrefix, parBus.EntityID())
	}

	name := f.getMessageNodeName(nodeInt)

	key := f.getMessageNodeKey(nodeInt)
	item := f.newItem(SidebarItemKindNode, nodeInt.Node().EntityID(), prefix, name)

	return key, item
}

func (f *sidebarItemFactory) newMessage(msg entity) (string, *sidebarItem) {
	return msg.EntityID().String(), f.newItem(SidebarItemKindMessage, msg.EntityID(), SidebarMessagesPrefix, msg.Name())
}

func (f *sidebarItemFactory) newSignalType(sigType entity) (string, *sidebarItem) {
	return sigType.EntityID().String(), f.newItem(SidebarItemKindSignalType, sigType.EntityID(), SidebarSignalTypesPrefix, sigType.Name())
}

func (f *sidebarItemFactory) newSignalUnit(sigUnit entity) (string, *sidebarItem) {
	return sigUnit.EntityID().String(), f.newItem(SidebarItemKindSignalUnit, sigUnit.EntityID(), SidebarSignalUnitsPrefix, sigUnit.Name())
}

func (f *sidebarItemFactory) newSignalEnum(sigEnum entity) (string, *sidebarItem) {
	return sigEnum.EntityID().String(), f.newItem(SidebarItemKindSignalEnum, sigEnum.EntityID(), SidebarSignalEnumsPrefix, sigEnum.Name())
}

func (f *sidebarItemFactory) getMessageBusKey(bus entity) string {
	return fmt.Sprintf("%s:%s", bus.EntityID(), SidebarMessagesPrefix)
}

func (f *sidebarItemFactory) getMessageNodeKey(nodeInt *acmelib.NodeInterface) string {
	return fmt.Sprintf("%s:%d:%s", nodeInt.Node().EntityID(), nodeInt.Number(), SidebarMessagesPrefix)
}

func (f *sidebarItemFactory) getMessageNodeName(nodeInt *acmelib.NodeInterface) string {
	return fmt.Sprintf("%s:%d", nodeInt.Node().Name(), nodeInt.Number())
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
	itemKey       string
	parentItemKey string
}

func newSidebarAddReq(item *sidebarItem, itemKey, parentItemKey string) *sidebarAddReq {
	return &sidebarAddReq{
		item:          item,
		itemKey:       itemKey,
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
	f *sidebarItemFactory

	items map[string]*sidebarItem
	root  *sidebarItem

	mux sync.RWMutex

	loadCh       chan *sidebarLoadReq
	updateNameCh chan *sidebarUpdateNameReq
	addCh        chan *sidebarAddReq
	deleteCh     chan *sidebarDeleteReq

	stopCh chan struct{}
}

func newSidebarService() *SidebarService {
	return &SidebarService{
		f: &sidebarItemFactory{},

		items: make(map[string]*sidebarItem),

		loadCh:       make(chan *sidebarLoadReq),
		updateNameCh: make(chan *sidebarUpdateNameReq),
		addCh:        make(chan *sidebarAddReq),
		deleteCh:     make(chan *sidebarDeleteReq),

		stopCh: make(chan struct{}),
	}
}

func (s *SidebarService) run() {
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

func (s *SidebarService) OnStartup(_ context.Context, _ application.ServiceOptions) error {
	go s.run()
	return nil
}

func (s *SidebarService) OnShutdown() {
	s.stopCh <- struct{}{}
}

func (s *SidebarService) addItem(itemKey string, item *sidebarItem) {
	s.items[itemKey] = item
}

func (s *SidebarService) sendLoad(req *sidebarLoadReq) {
	s.loadCh <- req
}

func (s *SidebarService) load(req *sidebarLoadReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	clear(s.items)

	net := req.network

	nodes := make(map[acmelib.EntityID]*acmelib.Node)
	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	// the network is the root item
	netKey, netItem := s.f.newNetwork(net)
	s.addItem(netKey, netItem)
	s.root = netItem

	// it groups all the buses
	busGroupKey, busGroupItem := s.f.newGroup(SidebarBusesPrefix, "Buses")
	s.addItem(busGroupKey, busGroupItem)
	netItem.addChild(busGroupItem)

	// it groups all the nodes
	nodeGroupKey, nodeGroupItem := s.f.newGroup(SidebarNodesPrefix, "Nodes")
	s.addItem(nodeGroupKey, nodeGroupItem)
	netItem.addChild(nodeGroupItem)

	// it groups all the messages
	msgGroupKey, msgGroupItem := s.f.newGroup(SidebarMessagesPrefix, "Messages")
	s.addItem(msgGroupKey, msgGroupItem)
	netItem.addChild(msgGroupItem)

	// it groups all the signal types
	sigTypeGroupKey, sigTypeGroupItem := s.f.newGroup(SidebarSignalTypesPrefix, "Signal Types")
	s.addItem(sigTypeGroupKey, sigTypeGroupItem)
	netItem.addChild(sigTypeGroupItem)

	// it groups all the signal units
	sigUnitGroupKey, sigUnitGroupItem := s.f.newGroup(SidebarSignalUnitsPrefix, "Signal Units")
	s.addItem(sigUnitGroupKey, sigUnitGroupItem)
	netItem.addChild(sigUnitGroupItem)

	// it groups all the signal enums
	sigEnumGroupKey, sigEnumGroupItem := s.f.newGroup(SidebarSignalEnumsPrefix, "Signal Enums")
	s.addItem(sigEnumGroupKey, sigEnumGroupItem)
	netItem.addChild(sigEnumGroupItem)

	// add buses and nodes
	for _, bus := range net.Buses() {
		busKey, busItem := s.f.newBus(bus)
		s.addItem(busKey, busItem)
		busGroupItem.addChild(busItem)

		for _, nodeInt := range bus.NodeInterfaces() {
			node := nodeInt.Node()
			nodeKey, nodeItem := s.f.newNode(node)
			s.addItem(nodeKey, nodeItem)
			nodeGroupItem.addChild(nodeItem)

			nodes[node.EntityID()] = node
		}

		// add bus group for messages
		msgBusKey, msgBusItem := s.f.newMessageBus(bus)
		s.addItem(msgBusKey, msgBusItem)
		msgGroupItem.addChild(msgBusItem)
	}

	// add messages
	for _, node := range nodes {
		for _, nodeInt := range node.Interfaces() {
			parBus := nodeInt.ParentBus()
			if parBus == nil {
				continue
			}

			// add node group for messages
			msgNodeKey, msgNodeItem := s.f.newMessageNode(nodeInt)
			s.addItem(msgNodeKey, msgNodeItem)

			// add node group into bus group for messages
			msgBusItem := s.items[s.f.getMessageBusKey(parBus)]
			msgBusItem.addChild(msgNodeItem)

			for _, msg := range nodeInt.SentMessages() {
				msgKey, msgItem := s.f.newMessage(msg)
				s.addItem(msgKey, msgItem)
				msgNodeItem.addChild(msgItem)

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
		sigTypeKey, sigTypeItem := s.f.newSignalType(sigType)
		s.addItem(sigTypeKey, sigTypeItem)
		sigTypeGroupItem.addChild(sigTypeItem)
	}

	// add signal units
	for _, sigUnit := range sigUnits {
		sigUnitKey, sigUnitItem := s.f.newSignalUnit(sigUnit)
		s.addItem(sigUnitKey, sigUnitItem)
		sigUnitGroupItem.addChild(sigUnitItem)
	}

	// add signal enums
	for _, sigEnum := range sigEnums {
		sigEnumKey, sigEnumItem := s.f.newSignalEnum(sigEnum)
		s.addItem(sigEnumKey, sigEnumItem)
		sigEnumGroupItem.addChild(sigEnumItem)
	}

	app.EmitEvent(SidebarLoad)
}

func (s *SidebarService) sendUpdateName(req *sidebarUpdateNameReq) {
	s.updateNameCh <- req
}

func (s *SidebarService) updateName(req *sidebarUpdateNameReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	item, ok := s.items[req.itemKey]
	if !ok {
		return
	}

	item.name = req.name

	app.EmitEvent(SidebarUpdateName, item.convertBase())
}

func (s *SidebarService) sendAdd(req *sidebarAddReq) {
	s.addCh <- req
}

func (s *SidebarService) add(req *sidebarAddReq) {
	s.mux.Lock()
	defer s.mux.Unlock()

	parent, ok := s.items[req.parentItemKey]
	if !ok {
		return
	}

	// TODO! remove this if
	itemKey := req.item.entityID.String()
	if req.itemKey != "" {
		itemKey = req.itemKey
	}

	s.addItem(itemKey, req.item)
	parent.addChild(req.item)

	app.EmitEvent(SidebarAdd, parent.convert())
}

func (s *SidebarService) sendDelete(req *sidebarDeleteReq) {
	s.deleteCh <- req
}

func (s *SidebarService) delete(req *sidebarDeleteReq) {
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
	f *sidebarItemFactory

	updateNameCh chan<- *sidebarUpdateNameReq
	addCh        chan<- *sidebarAddReq
	deleteCh     chan<- *sidebarDeleteReq
}

func (s *sidebarController) sendUpdateName(ent entity) {
	itemKey := ent.EntityID().String()
	name := ent.Name()

	s.updateNameCh <- newSidebarUpdateNameReq(itemKey, name)

	switch ent.EntityKind() {
	case acmelib.EntityKindBus:
		s.updateNameCh <- newSidebarUpdateNameReq(s.f.getMessageBusKey(ent), name)

	case acmelib.EntityKindNode:
		node, ok := ent.(*acmelib.Node)
		if !ok {
			panic("entity is not an acmelib.Node")
		}

		for _, nodeInt := range node.Interfaces() {
			if nodeInt.ParentBus() == nil {
				continue
			}

			s.updateNameCh <- newSidebarUpdateNameReq(s.f.getMessageNodeKey(nodeInt), s.f.getMessageNodeName(nodeInt))
		}
	}
}

func (s *sidebarController) sendAdd(ent entity) {
	switch ent.EntityKind() {
	case acmelib.EntityKindBus:
		busKey, busItem := s.f.newBus(ent)
		s.addCh <- newSidebarAddReq(busItem, busKey, SidebarBusesPrefix)

		msgBusKey, msgBusItem := s.f.newMessageBus(ent)
		s.addCh <- newSidebarAddReq(msgBusItem, msgBusKey, SidebarMessagesPrefix)

	case acmelib.EntityKindNode:
		nodeKey, nodeItem := s.f.newNode(ent)
		s.addCh <- newSidebarAddReq(nodeItem, nodeKey, SidebarNodesPrefix)

		node, ok := ent.(*acmelib.Node)
		if !ok {
			panic("entity is not an acmelib.Node")
		}

		for _, nodeInt := range node.Interfaces() {
			parBus := nodeInt.ParentBus()
			if parBus == nil {
				continue
			}

			msgNodeKey, msgGroupItem := s.f.newMessageNode(nodeInt)
			s.addCh <- newSidebarAddReq(msgGroupItem, msgNodeKey, s.f.getMessageBusKey(parBus))
		}

	case acmelib.EntityKindMessage:
		msgKey, msgItem := s.f.newMessage(ent)
		s.addCh <- newSidebarAddReq(msgItem, msgKey, SidebarMessagesPrefix)

	case acmelib.EntityKindSignalType:
		sigTypeKey, sigTypeItem := s.f.newSignalType(ent)
		s.addCh <- newSidebarAddReq(sigTypeItem, sigTypeKey, SidebarSignalTypesPrefix)

	case acmelib.EntityKindSignalUnit:
		sigUnitKey, sigUnitItem := s.f.newSignalUnit(ent)
		s.addCh <- newSidebarAddReq(sigUnitItem, sigUnitKey, SidebarSignalUnitsPrefix)

	case acmelib.EntityKindSignalEnum:
		sigEnumKey, sigEnumItem := s.f.newSignalEnum(ent)
		s.addCh <- newSidebarAddReq(sigEnumItem, sigEnumKey, SidebarSignalEnumsPrefix)
	}
}

func (s *sidebarController) sendDelete(ent entity) {
	s.deleteCh <- newSidebarDeleteReq(ent.EntityID().String())

	switch ent.EntityKind() {
	case acmelib.EntityKindBus:
		s.deleteCh <- newSidebarDeleteReq(s.f.getMessageBusKey(ent))

	case acmelib.EntityKindNode:
		node, ok := ent.(*acmelib.Node)
		if !ok {
			panic("entity is not an acmelib.Node")
		}

		for _, nodeInt := range node.Interfaces() {
			if nodeInt.ParentBus() == nil {
				continue
			}

			s.deleteCh <- newSidebarDeleteReq(s.f.getMessageNodeKey(nodeInt))
		}
	}
}
