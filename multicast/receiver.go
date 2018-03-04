package multicast

import (
	"log"
	"net"
)

func Receive(address string) {
	b := make([]byte, 1024, 1024)
	oob := make([]byte, 1024, 1024)
	udpAddr, err := net.ResolveUDPAddr("udp", "239.1.1.5:5050")
	if err != nil {
		log.Fatal(err)
	}
	udpConn, err := net.ListenMulticastUDP("udp", nil, udpAddr)

	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	for {
		n, oobn, flags, fromAddr, err := udpConn.ReadMsgUDP(b, oob)
		if err != nil {
			log.Fatal(err)
		}
		if n > 0 {
			log.Printf("Read %d bytes, %d oob, with flags %d from %s\n", n, oobn, flags, fromAddr)
			log.Printf("Bytes: %s\n", b[:n])
			log.Printf("OOB: %s\n", oob[:oobn])
		}
	}
}
