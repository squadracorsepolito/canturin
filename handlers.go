package main

import (
	"os"

	"github.com/squadracorsepolito/acmelib"
)

type NetworkService struct {
	network *acmelib.Network
}

func newNetworkService() (*NetworkService, error) {
	wireFile, err := os.Open("./testdata/SC24.binpb")
	if err != nil {
		return nil, err
	}
	defer wireFile.Close()

	net, err := acmelib.LoadNetwork(wireFile, acmelib.SaveEncodingWire)
	if err != nil {
		return nil, err
	}

	return &NetworkService{
		network: net,
	}, nil
}

func (ns *NetworkService) GetNetworkStub() NetworkStub {
	res := NetworkStub{
		entityStub: getEntityStub(ns.network),

		Buses:       []BusStub{},
		SignalUnits: []SignalUnitStub{},
	}

	sigUnits := make(map[acmelib.EntityID]*acmelib.SignalUnit)

	for _, tmpBus := range ns.network.Buses() {
		bus := BusStub{
			entityStub: getEntityStub(tmpBus),

			Nodes: []NodeStub{},
		}

		for _, tmpNodeInt := range tmpBus.NodeInterfaces() {
			node := NodeStub{
				entityStub: getEntityStub(tmpNodeInt.Node()),
			}

			for _, tmpSendMsg := range tmpNodeInt.Messages() {
				node.SendedMessages = append(node.SendedMessages, MessageStub{
					entityStub: getEntityStub(tmpSendMsg),
				})

				for _, tmpSig := range tmpSendMsg.Signals() {
					if tmpSig.Kind() == acmelib.SignalKindStandard {
						tmpStdSig, err := tmpSig.ToStandard()
						if err != nil {
							panic(err)
						}

						tmpSigUnit := tmpStdSig.Unit()
						if tmpSigUnit != nil {
							sigUnits[tmpSigUnit.EntityID()] = tmpSigUnit
						}
					}
				}
			}

			bus.Nodes = append(bus.Nodes, node)
		}

		res.Buses = append(res.Buses, bus)
	}

	for _, tmpSigUnit := range sigUnits {
		res.SignalUnits = append(res.SignalUnits, SignalUnitStub{entityStub: getEntityStub(tmpSigUnit)})

		sigUnitCh <- tmpSigUnit
	}

	return res
}

func (ns *NetworkService) GetNetwork() Network {
	net := Network{
		base: getBase(ns.network),
	}

	for _, tmpBus := range ns.network.Buses() {
		bus := Bus{
			base: getBase(tmpBus),
		}

		for _, tmpNodeInt := range tmpBus.NodeInterfaces() {
			node := Node{
				base: getBase(tmpNodeInt.Node()),
			}

			for _, tmpSendMsg := range tmpNodeInt.Messages() {
				node.SendedMessages = append(node.SendedMessages, Message{
					base: getBase(tmpSendMsg),
				})
			}

			bus.Nodes = append(bus.Nodes, node)
		}

		net.Buses = append(net.Buses, bus)
	}

	return net
}
