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

	proxy.pushSidebarLoad(net)
	proxy.network = net

	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	for _, bus := range net.Buses() {

		for _, nodeInt := range bus.NodeInterfaces() {

			for _, msg := range nodeInt.Messages() {
				proxy.pushMessage(msg)

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

	for _, sigType := range sigTypes {
		proxy.pushSignalType(sigType)
	}

	for _, sigUnit := range sigUnits {
		proxy.pushSignalUnit(sigUnit)
	}

	for _, sigEnum := range sigEnums {
		proxy.pushLoadSignalEnum(sigEnum)
	}
}
