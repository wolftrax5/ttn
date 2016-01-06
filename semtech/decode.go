// Copyright © 2015 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package semtech

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Unmarshal parses a raw response from a server and turn in into a packet.
// Will return an error if the response fields are incorrect.
func Unmarshal(raw []byte) (*Packet, error) {
	size := len(raw)

	if size < 4 {
		return nil, errors.New("Invalid raw data format")
	}

	packet := &Packet{
		Version:    raw[0],
		Token:      raw[1:3],
		Identifier: raw[3],
		GatewayId:  nil,
		Payload:    nil,
	}

	if packet.Version != VERSION {
		return nil, errors.New("Unreckognized protocol version")
	}

	if packet.Identifier > PULL_ACK {
		return nil, errors.New("Unreckognized protocol identifier")
	}

	cursor := 4
	if packet.Identifier == PULL_DATA || packet.Identifier == PUSH_DATA {
		if size < 12 {
			return nil, errors.New("Invalid gateway identifier")
		}
		packet.GatewayId = raw[cursor:12]
		cursor = 12
	}

	var err error
	if size > cursor && (packet.Identifier == PUSH_DATA || packet.Identifier == PULL_RESP) {
		packet.Payload = new(Payload)
		packet.Payload.Raw = raw[cursor:]
		err = json.Unmarshal(raw[cursor:], packet.Payload)
	}

	return packet, err
}

// UnmarshalJSON implements the Unmarshaler interface from encoding/json
func (t *Timeparser) UnmarshalJSON(raw []byte) error {
	str := strings.Trim(string(raw), `"`)
	v, err := time.Parse("2006-01-02 15:04:05 GMT", str)
	if err != nil {
		v, err = time.Parse(time.RFC3339, str)
	}
	if err != nil {
		v, err = time.Parse(time.RFC3339Nano, str)
	}
	if err != nil {
		return errors.New("Unkown date format. Unable to parse time")
	}
	t.Value = &v
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface from encoding/json
func (d *Datrparser) UnmarshalJSON(raw []byte) error {
	v := strings.Trim(string(raw), `"`)

	if v == "" {
		return errors.New("Invalid datr format")
	}

	d.Value = &v
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface from encoding/json
func (p *Payload) UnmarshalJSON(raw []byte) error {
	proxy := payloadProxy{
		ProxStat: &statProxy{},
		ProxTXPK: &txpkProxy{},
	}

	if err := json.Unmarshal(raw, &proxy); err != nil {
		return err
	}

	if proxy.ProxStat.Stat != nil {
		if proxy.ProxStat.Time != nil {
			proxy.ProxStat.Stat.Time = proxy.ProxStat.Time.Value
		}
		p.Stat = proxy.ProxStat.Stat
	}

	if proxy.ProxTXPK.TXPK != nil {
		if proxy.ProxTXPK.Time != nil {
			proxy.ProxTXPK.TXPK.Time = proxy.ProxTXPK.Time.Value
		}
		if proxy.ProxTXPK.Datr != nil {
			proxy.ProxTXPK.TXPK.Datr = proxy.ProxTXPK.Datr.Value
		}
		p.TXPK = proxy.ProxTXPK.TXPK
	}

	for _, rxpk := range proxy.ProxRXPK {
		if rxpk.Time != nil {
			rxpk.RXPK.Time = rxpk.Time.Value
		}
		if rxpk.Datr != nil {
			rxpk.RXPK.Datr = rxpk.Datr.Value
		}
		p.RXPK = append(p.RXPK, *rxpk.RXPK)
	}

	return nil
}
