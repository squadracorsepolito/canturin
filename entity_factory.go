package main

import (
	"fmt"

	"github.com/squadracorsepolito/acmelib"
)

type entityFactory struct{}

func newEntityFactory() *entityFactory {
	return &entityFactory{}
}

func (f *entityFactory) getNewName(entityName string, takenNames map[string]struct{}) string {
	res := ""
	count := 0

	for {
		res = fmt.Sprintf("new_%s_%d", entityName, count)
		if _, ok := takenNames[res]; !ok {
			break
		}
		count++
	}

	return res
}

func (f *entityFactory) createBus(net *acmelib.Network) (*acmelib.Bus, error) {
	takenNames := make(map[string]struct{})
	for _, tmpBus := range net.Buses() {
		takenNames[tmpBus.Name()] = struct{}{}
	}

	bus := acmelib.NewBus(f.getNewName("bus", takenNames))
	if err := net.AddBus(bus); err != nil {
		return nil, err
	}

	return bus, nil
}

func (f *entityFactory) createNode(net *acmelib.Network) (*acmelib.Node, error) {
	takenNames := make(map[string]struct{})
	takenNodeIDs := make(map[acmelib.NodeID]struct{})

	for _, tmpBus := range net.Buses() {
		for _, tmpNodeInt := range tmpBus.NodeInterfaces() {
			tmpNode := tmpNodeInt.Node()
			takenNames[tmpNode.Name()] = struct{}{}
			takenNodeIDs[tmpNode.ID()] = struct{}{}
		}
	}

	nodeID := acmelib.NodeID(1)
	for {
		if _, ok := takenNodeIDs[nodeID]; !ok {
			break
		}
		nodeID++
	}

	node := acmelib.NewNode(f.getNewName("node", takenNames), nodeID, 1)

	return node, nil
}
