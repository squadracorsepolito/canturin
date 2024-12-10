package main

import (
	"os"

	"github.com/squadracorsepolito/acmelib"
)

func loadNetwork(path string) {
	wireFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer wireFile.Close()

	net, err := acmelib.LoadNetwork(wireFile, acmelib.SaveEncodingWire)
	if err != nil {
		panic(err)
	}

	manager.sidebar.sendLoad(newSidebarLoadReq(net))

	proxy.network = net

	nodes := make(map[acmelib.EntityID]*acmelib.Node)

	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	for _, bus := range net.Buses() {
		manager.bus.sendLoad(bus)

		for _, nodeInt := range bus.NodeInterfaces() {
			tmpNode := nodeInt.Node()
			nodes[tmpNode.EntityID()] = tmpNode

			for _, msg := range nodeInt.SentMessages() {
				manager.message.sendLoad(msg)

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

	for _, node := range nodes {
		manager.node.sendLoad(node)
	}

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
