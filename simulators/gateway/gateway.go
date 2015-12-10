// Copyright © 2015 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

// package gateway offers a dummy representation of a gateway.
//
// The package can be used to create a dummy gateway.
// Its former use is to provide a handy simulator for further testing of the whole network chain.
package gateway

import (
	"errors"
	"fmt"
	"net"
)

type Gateway struct {
	Id      []byte         // Gateway's Identifier
	alti    int            // GPS altitude in RX meters
	ackr    uint           // Number of upstream datagrams that were acknowledged
	dwnb    uint           // Number of downlink datagrams received
	lati    float64        // GPS latitude, North is +
	long    float64        // GPS longitude, East is +
	rxfw    uint           // Number of radio packets forwarded
	rxnb    uint           // Number of radio packets received
	txnb    uint           // Number of packets emitted
	routers []*net.UDPAddr // List of routers addresses
	quit    chan bool      // Communication channel to stop connections
}

// New create a new gateway from a given id and a list of router addresses
func New(id []byte, routers ...string) (*Gateway, error) {
	if len(id) != 8 {
		return nil, errors.New("Invalid gateway id provided")
	}

	if len(routers) == 0 {
		return nil, errors.New("At least one router address should be provided")
	}

	addresses := make([]*net.UDPAddr, 0)
	var err error
	for _, router := range routers {
		var addr *net.UDPAddr
		addr, err = net.ResolveUDPAddr("udp", router)
		if err != nil {
			break
		}
		addresses = append(addresses, addr)
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid router address. %v", err))
	}

	return &Gateway{
		Id:      id,
		alti:    120,
		lati:    53.3702,
		long:    4.8952,
		routers: addresses,
	}, nil
}

type Imitator interface {
	Mimic()
}
