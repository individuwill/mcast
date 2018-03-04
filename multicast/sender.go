package multicast

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type Sender struct {
	Packet
}

func NewSender(address string, port int, ttl int) *Sender {
	s := &Sender{}
	s.TTL = ttl
	s.Port = port
	s.SetAddress(address)
	return s
}

func (s *Sender) SetMessagePadding(paddingSize int) {
	s.padding = make([]byte, paddingSize)
}

func (s *Sender) SetLocalAddress(address string) error {
	if !strings.Contains(address, ":") {
		return errors.New("Local address must contain a ':' with source port")
	}
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	s.LocalAddress = addr
	return nil
}

func (s *Sender) SetTOS(tos int) {
	s.TOS = tos
}

func (s *Sender) One(message string) error {
	s.SetMessageText(message)
	return s.Send()
}

func (s *Sender) Max(message string, interval int, startValue int, numberOfMessages int) error {
	var text func(int) string
	if strings.Contains(message, "{c}") {
		subStr := strings.Replace(message, "{c}", "%d", 1)
		text = func(x int) string {
			return fmt.Sprintf(subStr, x)
		}
	} else {
		text = func(x int) string { return message }
	}

	d := time.Duration(interval) * time.Millisecond

	for i := 0; numberOfMessages == 0 || i < numberOfMessages; i++ {
		err := s.One(text(startValue))
		if err != nil {
			return err
		}
		startValue++
		time.Sleep(d)
	}
	return nil
}

func (s *Sender) Forever(message string, interval int, startValue int) error {
	return s.Max(message, interval, startValue, 0)
}
