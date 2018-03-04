package multicast

import (
	"log"
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

func Receive(address string, port int, interfaceName string) error {
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

	for {
		n, cm, src, err := packetConn.ReadFrom(buf)
		if err != nil {
			return err
		}
		if cm != nil {
			log.Printf("Received %d bytes on %v with ttl: %v from %v\n", n, cm.Dst, cm.TTL, src)
		} else {
			log.Printf("Received %d bytes from %v\n", n, src)
		}
	}
	return nil
}
