package main

import (
	"log"
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

	// Push the network data to the frontend
	proxy.pushSidebarLoad(net)
	proxy.network = net

	// Log the network data
    log.Println("Network data loaded and pushed to frontend")

	// Push the signal types, units, and enums to the frontend
	sigTypes := make(map[acmelib.EntityID]*acmelib.SignalType)
	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)
	sigEnums := make(map[acmelib.EntityID]*acmelib.SignalEnum)

	// Iterate over the buses, node interfaces and messages
	for _, bus := range net.Buses() {
		
		for _, nodeInt := range bus.NodeInterfaces() {
			
			for _, msg := range nodeInt.Messages() {
				// Push the message to the frontend
				proxy.pushMessage(msg)
				log.Println("Message pushed to frontend:", msg)

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
	// Log the signal types, units, and enums
    log.Println("Signal types, units, and enums processed and pushed to frontend")

	// Push the signal types, units, and enums to the frontend
	for _, sigType := range sigTypes {
		proxy.pushSignalType(sigType)
	}

	for _, sigUnit := range sigUnits {
		proxy.pushSignalUnit(sigUnit)
	}

	for _, sigEnum := range sigEnums {
		proxy.pushSignalEnum(sigEnum)
	}
}
