package main

import (
	"os"

	"github.com/squadracorsepolito/acmelib"
)

/*
 * Load the network data from the given file path
 */
func loadNetwork(path string) {
	// Open the file
	wireFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer wireFile.Close()

	// Load the network data
	net, err := acmelib.LoadNetwork(wireFile, acmelib.SaveEncodingWire)
	if err != nil {
		panic(err)
	}

	manager.sidebar.sendLoad(newSidebarLoadReq(net))

	proxy.network = net

	nodes := make(map[acmelib.EntityID]*acmelib.Node)

	// Push the signal types, units, and enums to the frontend
	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	// Iterate over the buses, node interfaces and messages
	for _, bus := range net.Buses() {
		manager.bus.sendLoad(bus)

		for _, nodeInt := range bus.NodeInterfaces() {
			tmpNode := nodeInt.Node()
			nodes[tmpNode.EntityID()] = tmpNode

			for _, msg := range nodeInt.SentMessages() {
				manager.message.sendLoad(msg)

				// Iterate over the signals
				for _, sig := range msg.Signals() {
					switch sig.Kind() {
					// Standard signal
					case acmelib.SignalKindStandard:
						stdSig, err := sig.ToStandard()
						if err != nil {
							panic(err)
						}
						sigTypes[stdSig.Type().EntityID()] = stdSig.Type()

						if stdSig.Unit() != nil {
							sigUnits[stdSig.Unit().EntityID()] = stdSig.Unit()
						}

					// Enum signal
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

	for _, node := range nodes {
		manager.node.sendLoad(node)
	}

	// Push the signal types, units, and enums to the frontend
	for _, sigType := range sigTypes {
		manager.signalType.sendLoad(sigType)
	}

	for _, sigUnit := range sigUnits {
		manager.signalUnit.sendLoad(sigUnit)
	}

	for _, sigEnum := range sigEnums {
		manager.signalEnum.sendLoad(sigEnum)
	}
}
