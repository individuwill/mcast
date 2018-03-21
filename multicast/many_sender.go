/*
*    mcast - Command line tool and library for testing multicast traffic
*    flows and stress testing networks and devices.
*    Copyright (C) 2018 Will Smith
*
*    This program is free software: you can redistribute it and/or modify
*    it under the terms of the GNU General Public License as published by
*    the Free Software Foundation, either version 3 of the License, or
*    (at your option) any later version.
*
*    This program is distributed in the hope that it will be useful,
*    but WITHOUT ANY WARRANTY; without even the implied warranty of
*    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*    GNU General Public License for more details.
*
*    You should have received a copy of the GNU General Public License
*    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package multicast

import (
	"errors"
	"log"
	"net"
	"strings"
	"sync"
)

type ManySender struct {
	Addresses      *[]net.IP
	Port           int
	TTL            int
	MessagePadding int
	TOS            int
	LocalAddress   *net.UDPAddr
}

func NewManySender(network string, mask int, port int, ttl int) (*ManySender, error) {
	m := ManySender{TTL: ttl, Port: port}
	addresses, err := IPList(network, mask)
	if err != nil {
		return nil, err
	}
	m.Addresses = &addresses

	return &m, nil
}

func (m *ManySender) startSender(address, message string, interval, startValue, numberOfMessages int) {
	s := NewSender(address, m.Port, m.TTL)
	s.SetMessagePadding(m.MessagePadding)
	s.SetTOS(m.TOS)
	s.LocalAddress = m.LocalAddress
	err := s.Max(message, interval, startValue, numberOfMessages)
	if err != nil {
		log.Println("Problem sending max messages")
		log.Println(err)
	}
}

func (m *ManySender) Start(message string, interval int, startValue int, numberOfMessages int) {
	wg := new(sync.WaitGroup)
	for _, address := range *m.Addresses {
		wg.Add(1)
		go func(address string) {
			m.startSender(address, message, interval, startValue, numberOfMessages)
			wg.Done()
		}(address.String())
	}
	wg.Wait()
}

func (m *ManySender) SetMessagePadding(paddingSize int) {
	m.MessagePadding = paddingSize
}

func (m *ManySender) SetTOS(tos int) {
	m.TOS = tos
}

func (m *ManySender) SetLocalAddress(address string) error {
	if !strings.Contains(address, ":") {
		return errors.New("Local address must contain a ':' with source port")
	}
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	m.LocalAddress = addr
	return nil
}
