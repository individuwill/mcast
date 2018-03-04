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
