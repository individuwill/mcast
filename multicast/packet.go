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
	"fmt"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

type Packet struct {
	TTL          int
	Port         int
	Address      net.IP
	RouterAlert  bool
	Raw          bool
	IGMPVersion  int // 1, 2, or 3
	Interface    *net.Interface
	Message      []byte
	Protocol     string // 'udp' or 'ip:2'/'ip4:2'
	LocalAddress *net.UDPAddr
	udpConn      *net.UDPConn
	packetConn   *ipv4.PacketConn
	ipConn       net.PacketConn
	rawConn      *ipv4.RawConn
	padding      []byte
	TOS          int
}

func NewPacket() *Packet {
	return &Packet{
		TTL:         50,
		Port:        5050,
		Address:     net.ParseIP("239.1.1.50"),
		IGMPVersion: 2,
		Protocol:    "udp",
		Message:     []byte("test"),
	}
}

func (p *Packet) SetAddress(address string) {
	p.Address = net.ParseIP(address)
}

func (p *Packet) SetMessageText(message string) {
	p.Message = []byte(message)
}

func (p *Packet) ConnectUDP() error {
	udpAddr, err := net.ResolveUDPAddr("udp", p.AddressAndPort())
	if err != nil {
		log.Println("Failed to resolve UDP address")
		return err
	}

	p.udpConn, err = net.DialUDP("udp", p.LocalAddress, udpAddr)
	if err != nil {
		log.Println("Problem dialing UDP")
		return err
	}

	p.packetConn = ipv4.NewPacketConn(p.udpConn)
	return nil
}

func (p *Packet) SendUDP() error {
	if p.udpConn == nil || p.packetConn == nil {
		err := p.ConnectUDP()
		if err != nil {
			log.Println("Unable to iniate UDP connection")
			return err
		}
	}

	var err error
	p.packetConn.SetMulticastTTL(p.TTL)
	p.packetConn.SetTOS(p.TOS)
	if len(p.padding) > 0 {
		copy(p.padding, p.Message)
		_, err = p.udpConn.Write(p.padding)
	} else {
		_, err = p.udpConn.Write(p.Message)
	}
	return err
}

func (p *Packet) AddressAndPort() string {
	return fmt.Sprintf("%v:%v", p.Address, p.Port)
}

func (p *Packet) ConnectRaw() error {
	var err error
	p.ipConn, err = net.ListenPacket(p.Protocol, "0.0.0.0")
	if err != nil {
		log.Println("Failed to dial")
		return err
	}

	p.rawConn, err = ipv4.NewRawConn(p.ipConn)
	if err != nil {
		log.Println("probelem getting raw socket")
		return err
	}
	return nil
}

func (p *Packet) SendRaw() error {
	if p.ipConn == nil || p.rawConn == nil {
		err := p.ConnectRaw()
		if err != nil {
			log.Println("Problem making raw connection")
			return err
		}
	}

	// router alert: https://tools.ietf.org/html/rfc2113
	// option is 0x94
	// order is option flags or val, length, 0, 0
	var options []byte
	if p.RouterAlert {
		options = []byte{0x94, 0x04, 0x0, 0x0}
	}
	hlen := ipv4.HeaderLen + len(options)
	header := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      hlen,
		TOS:      p.TOS,
		TotalLen: hlen + len(p.Message),
		TTL:      p.TTL,
		Protocol: 2,
		Dst:      p.Address,
		Options:  options,
	}
	err := p.rawConn.WriteTo(header, p.Message, nil)
	if err != nil {
		log.Println("problem writing socket")
		return err
	}
	return nil
}

func (p *Packet) Send() error {
	if p.Protocol != "ip:2" && p.Protocol != "ip4:2" {
		return p.SendUDP()
	}
	return p.SendRaw()
}

func (p *Packet) Close() error {
	if p.udpConn != nil {
		return p.udpConn.Close()
	}
	if p.packetConn != nil {
		return p.packetConn.Close()
	}
	if p.ipConn != nil {
		return p.ipConn.Close()
	}
	if p.rawConn != nil {
		return p.rawConn.Close()
	}
	return nil
}
