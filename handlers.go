package main

import (
	"encoding/json"
	"errors"
	"log"
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

func (ns *NetworkService) GetNetwork() *Network {
	net := &Network{
		base: getBase(ns.network),
	}

	for _, tmpBus := range ns.network.Buses() {
		bus := &Bus{
			base: getBase(tmpBus),
		}

		for _, tmpNodeInt := range tmpBus.NodeInterfaces() {
			node := &Node{
				base: getBase(tmpNodeInt.Node()),
			}

			for _, tmpSendMsg := range tmpNodeInt.Messages() {
				node.SendedMessages = append(node.SendedMessages, &Message{
					base: getBase(tmpSendMsg),
				})
			}

			bus.Nodes = append(bus.Nodes, node)
		}

		net.Buses = append(net.Buses, bus)
	}

	a, err := json.Marshal(net)
	if err != nil {
		panic(err)
	}
	log.Print(len(a))

	return net
}

func (ns *NetworkService) GetMessage(busID, nodeID, msgID string) (res Message, _ error) {
	buses := ns.network.Buses()
	for _, bus := range buses {
		if bus.EntityID() == acmelib.EntityID(busID) {
			for _, nodeInt := range bus.NodeInterfaces() {
				if nodeInt.Node().EntityID() == acmelib.EntityID(nodeID) {
					for _, msg := range nodeInt.Messages() {
						if msg.EntityID() == acmelib.EntityID(msgID) {
							res := Message{
								base: getBase(msg),

								SizeByte: msg.SizeByte(),
								Signals:  []Signal{},
							}

							for _, sig := range msg.Signals() {
								res.Signals = append(res.Signals, Signal{
									base: getBase(sig),

									Kind:     sig.Kind(),
									StartPos: sig.GetStartBit(),
									Size:     sig.GetSize(),
								})
							}

							return res, nil
						}
					}
				}
			}
		}
	}

	return res, errors.New("not found")
}
