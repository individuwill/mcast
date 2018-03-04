package multicast

import (
	"fmt"
	"net"

	"golang.org/x/net/ipv4"
)

func getInterface(interfaceName string) (*net.Interface, error) {
	var localInterface *net.Interface
	var err error
	if interfaceName != "" {
		localInterface, err = net.InterfaceByName(interfaceName)
	}
	return localInterface, err
}

func getUDPConnection(address string, port int, localInterface *net.Interface) (*net.UDPConn, error) {
	var udpConn *net.UDPConn
	var err error
	ip := net.ParseIP(address)
	udpAddr := &net.UDPAddr{IP: ip, Port: port}
	if ip.IsMulticast() {
		udpConn, err = net.ListenMulticastUDP("udp", localInterface, udpAddr)
	} else {
		udpConn, err = net.ListenUDP("udp", udpAddr)
	}
	return udpConn, err
}

type Message struct {
	Data []byte
	CM   *ipv4.ControlMessage
	Src  net.Addr
}

func messagePrinter(messageCh <-chan Message, showData bool) {
	for message := range messageCh {
		if message.CM != nil {
			fmt.Printf("*Received %d bytes on %v with ttl: %v from %v*\n",
				len(message.Data), message.CM.Dst, message.CM.TTL, message.Src)
		} else {
			fmt.Printf("*Received %d bytes from %v*\n", len(message.Data), message.Src)
		}
		if showData {
			fmt.Printf("%s\n", message.Data)
			fmt.Println()
		}
	}
}

func Receive(address string, port int, interfaceName string, showData bool) error {
	localInterface, err := getInterface(interfaceName)
	if err != nil {
		return err
	}

	udpConn, err := getUDPConnection(address, port, localInterface)
	if err != nil {
		return err
	}
	defer udpConn.Close()

	packetConn := ipv4.NewPacketConn(udpConn)
	defer packetConn.Close()
	packetConn.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true)
	buf := make([]byte, 2048)

	messageCh := make(chan Message, 1000)
	go messagePrinter(messageCh, showData)

	for {
		n, cm, src, err := packetConn.ReadFrom(buf)
		if err != nil {
			return err
		}
		data := make([]byte, n)
		copy(data, buf)
		messageCh <- Message{Data: data, CM: cm, Src: src}
	}
}
