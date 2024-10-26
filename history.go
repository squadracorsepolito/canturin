package main

// import "github.com/squadracorsepolito/acmelib"

// type operationKind int

// const (
// 	operationKindSetDesc operationKind = iota
// 	operationKindUpdateName
// )

// type operationDomain int

// const (
// 	operationDomainNetwork operationDomain = iota
// 	operationDomainBus
// 	operationDomainNode
// 	operationDomainMessage
// )

// type messageOperationHandler struct{}

// func (h *messageOperationHandler) apply() any {
// 	return nil
// }

// type operation interface {
// 	kind() operationKind

// 	entityID() acmelib.EntityID
// 	before() any
// 	after() any
// }

// type history struct {
// 	operations []operation
// 	currIdx    int
// }

// func newHistory() *history {
// 	return &history{
// 		operations: []operation{},
// 		currIdx:    -1,
// 	}
// }

// func (h *history) push(op operation) {
// 	if len(h.operations) > h.currIdx+1 {
// 		h.operations = h.operations[:h.currIdx+1]
// 	}

// 	h.operations = append(h.operations, op)
// 	h.currIdx++
// }

// func (h *history) undo() {
// 	if h.currIdx == 0 {
// 		return
// 	}

// 	h.currIdx--
// }

// func (h *history) redo() {
// 	if h.currIdx == len(h.operations)-1 {
// 		return
// 	}

// 	h.currIdx++
// }
