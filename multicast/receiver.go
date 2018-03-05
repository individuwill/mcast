package multicast

import (
	"fmt"
	"net"
	"strings"

	"golang.org/x/net/ipv4"
)

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

type message struct {
	Data []byte
	CM   *ipv4.ControlMessage
	Src  net.Addr
}

func messagePrinter(messageCh <-chan message, showData bool) {
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

func receive(address string, port int, interfaceName string, showData bool, messageCh chan message) error {
	localInterface, err := GetInterface(interfaceName)
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

	for {
		n, cm, src, err := packetConn.ReadFrom(buf)
		if err != nil {
			return err
		}
		data := make([]byte, n)
		copy(data, buf)
		messageCh <- message{Data: data, CM: cm, Src: src}
	}
}

// Receive will listen on the port and addresses provided for incoming UDP messages.
// Uponn receipt of a UDP message, a message will be printed to the console.
// If showData is true, the data contained in that message will also be printed
// with the assumption that the data is a string.
// address can be in CIDR notation, in which case all of the addresses falling
// within that network will be listened on
func Receive(address string, port int, interfaceName string, showData bool) error {
	messageCh := make(chan message, 1000)
	go messagePrinter(messageCh, showData)

	if strings.Contains(address, "/") {
		network, mask, err := SplitCIDR(address)
		if err != nil {
			return err
		}
		ips, err := IPList(network, mask)
		if err != nil {
			return err
		}
		for _, ip := range ips {
			go receive(ip.String(), port, interfaceName, showData, messageCh)
		}
		for { // block this method call
		}
	}

	return receive(address, port, interfaceName, showData, messageCh)
}
