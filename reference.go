package main

import "github.com/squadracorsepolito/acmelib"

type ReferenceKind string

const (
	ReferenceKindBus     ReferenceKind = "bus"
	ReferenceKindNode    ReferenceKind = "node"
	ReferenceKindMessage ReferenceKind = "message"
	ReferenceKindSignal  ReferenceKind = "signal"
)

func newReferenceKind(kind acmelib.EntityKind) ReferenceKind {
	switch kind {
	case acmelib.EntityKindBus:
		return ReferenceKindBus
	case acmelib.EntityKindNode:
		return ReferenceKindNode
	case acmelib.EntityKindMessage:
		return ReferenceKindMessage
	case acmelib.EntityKindSignal:
		return ReferenceKindSignal
	default:
		return ReferenceKindBus
	}
}

type Reference struct {
	Kind     ReferenceKind `json:"kind"`
	EntityID string        `json:"entityId"`
	Name     string        `json:"name"`
	Children []Reference   `json:"children"`
}

type reference struct {
	kind        ReferenceKind
	entityID    acmelib.EntityID
	name        string
	children    []*reference
	childrenMap map[acmelib.EntityID]struct{}
}

func newReference(ent entity) *reference {
	return &reference{
		kind:        newReferenceKind(ent.EntityKind()),
		entityID:    ent.EntityID(),
		name:        ent.Name(),
		childrenMap: map[acmelib.EntityID]struct{}{},
	}
}

func (r *reference) addChild(child *reference) {
	if _, ok := r.childrenMap[child.entityID]; ok {
		return
	}

	r.children = append(r.children, child)
	r.childrenMap[child.entityID] = struct{}{}
}

func (r *reference) toResponse() Reference {
	res := Reference{
		Kind:     r.kind,
		EntityID: r.entityID.String(),
		Name:     r.name,
		Children: []Reference{},
	}

	for _, child := range r.children {
		res.Children = append(res.Children, child.toResponse())
	}

	return res
}

/*
func getReferencesFromSignals(signals []acmelib.Signal) []Reference {
	rootRefs := []*reference{}
	refs := make(map[acmelib.EntityID]*reference)
	for _, sig := range signals {
		sigRef := newReference(sig)
		refs[sig.EntityID()] = sigRef

		var msgRef *reference
		msg := sig.ParentMessage()
		sigRef.entityID = msg.EntityID()
		msgRef, ok := refs[msg.EntityID()]
		if !ok {
			msgRef = newReference(msg)
			refs[msg.EntityID()] = msgRef
		}
		msgRef.addChild(sigRef)

		var nodeRef *reference
		node := msg.SenderNodeInterface().Node()
		nodeRef, ok = refs[node.EntityID()]
		if !ok {
			nodeRef = newReference(node)
			refs[node.EntityID()] = nodeRef
		}
		nodeRef.addChild(msgRef)

		var busRef *reference
		bus := msg.SenderNodeInterface().ParentBus()
		busRef, ok = refs[bus.EntityID()]
		if !ok {
			busRef = newReference(bus)
			refs[bus.EntityID()] = busRef
			rootRefs = append(rootRefs, busRef)
		}
		busRef.addChild(nodeRef)
	}

	res := []Reference{}

	for _, tmpRef := range rootRefs {
		res = append(res, tmpRef.toResponse())
	}

	return res
}
*/
