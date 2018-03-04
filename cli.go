package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/individuwill/mcast/multicast"
	"golang.org/x/net/ipv4"
)

// igmpv2: https://tools.ietf.org/html/rfc2236

func sendTextTTL(message string) {
	// TODO: Adjust TTL
	udpAddr, err := net.ResolveUDPAddr("udp", "239.1.1.5:5050")
	if err != nil {
		log.Fatal(err)
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	newUDPConn := ipv4.NewPacketConn(udpConn)
	if err != nil {
		log.Fatal(err)
	}
	defer newUDPConn.Close()
	defer udpConn.Close()
	newUDPConn.SetMulticastTTL(20)
	/*
			n, err := newUDPConn.WriteTo([]byte(message), nil, nil)
			if err != nil {
				log.Fatal(err)
			}
		log.Printf("Wrote %d bytes", n)
	*/
	udpConn.Write([]byte(message))
	log.Println(udpConn.LocalAddr())
}

func sendText(message string) {
	// TODO: Adjust TTL
	udpAddr, err := net.ResolveUDPAddr("udp", "239.1.1.5:5050")
	if err != nil {
		log.Fatal(err)
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()
	udpConn.Write([]byte(message))
	log.Println(udpConn.LocalAddr())
}

func receive() {
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

func computeChecksum(buf []byte) uint16 {
	sum := uint32(0)

	for ; len(buf) >= 2; buf = buf[2:] {
		sum += uint32(buf[0])<<8 | uint32(buf[1])
	}
	if len(buf) > 0 {
		sum += uint32(buf[0]) << 8
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	csum := ^uint16(sum)
	/*
	 * From RFC 768:
	 * If the computed checksum is zero, it is transmitted as all ones (the
	 * equivalent in one's complement arithmetic). An all zero transmitted
	 * checksum value means that the transmitter generated no checksum (for
	 * debugging or for higher level protocols that don't care).
	 */
	if csum == 0 {
		csum = 0xffff
	}
	return csum
}

func query() {
	ip, err := net.ResolveIPAddr("ip:2", "224.0.0.1")
	if err != nil {
		log.Println("Failed to resolve")
		log.Fatal(err)
	}
	ipConn, err := net.DialIP("ip4:2", nil, ip)
	if err != nil {
		log.Println("Failed to dial")
		log.Fatal(err)
	}
	defer ipConn.Close()
	seconds := 10
	msg := []byte{0x11, byte(seconds * 10), 0, 0,
		0, 0, 0, 0}
	checksum := computeChecksum(msg)
	log.Printf("Checksum: %x\n", checksum)
	msg[2] = byte(checksum >> 8)   //byte(0xee) // byte(checksum & 0x00FF)
	msg[3] = byte(checksum & 0xFF) //byte(0x9b)          // byte((checksum & 0xFF00) >> 1)
	n, err := ipConn.Write(msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes\n", n)
}

func rawJoin(address string) {
	ip, err := net.ResolveIPAddr("ip:2", address)
	if err != nil {
		log.Println("Failed to resolve")
		log.Fatal(err)
	}
	ipConn, err := net.DialIP("ip4:2", nil, ip)
	if err != nil {
		log.Println("Failed to dial")
		log.Fatal(err)
	}
	defer ipConn.Close()

	seconds := 10
	ip4 := ip.IP.To4()
	msg := []byte{0x12, byte(seconds * 10), 0, 0,
		ip4[0], ip4[1], ip4[2], ip4[3]}
	log.Println(ip.IP.To4()[0])
	checksum := computeChecksum(msg)
	log.Printf("Checksum: %x\n", checksum)
	msg[2] = byte(checksum >> 8)   //byte(0xee) // byte(checksum & 0x00FF)
	msg[3] = byte(checksum & 0xFF) //byte(0x9b)          // byte((checksum & 0xFF00) >> 1)
	n, err := ipConn.Write(msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes\n", n)
}

func rawJoinRouterAlert(address string) {
	ip := net.ParseIP(address)
	//ipConn, err := net.DialIP("ip:2", nil, ip)
	ipConn, err := net.ListenPacket("ip4:2", "0.0.0.0")
	if err != nil {
		log.Println("Failed to dial")
		log.Fatal(err)
	}
	defer ipConn.Close()

	seconds := 10
	ip4 := ip.To4()
	msg := []byte{0x16, byte(seconds * 10), 0, 0,
		ip4[0], ip4[1], ip4[2], ip4[3]}
	checksum := computeChecksum(msg)
	log.Printf("Checksum: %x\n", checksum)
	msg[2] = byte(checksum >> 8)   //byte(0xee) // byte(checksum & 0x00FF)
	msg[3] = byte(checksum & 0xFF) //byte(0x9b)          // byte((checksum & 0xFF00) >> 1)

	rawConn, err := ipv4.NewRawConn(ipConn)
	if err != nil {
		log.Println("probelem getting raw socket")
		log.Fatal(err)
	}
	defer rawConn.Close()
	// router alert: https://tools.ietf.org/html/rfc2113
	// option is 0x94
	// order is option flags or val, length, 0, 0
	options := []byte{0x94, 0x04, 0x0, 0x0}
	hlen := ipv4.HeaderLen + len(options)
	header := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      hlen,
		TOS:      0,
		TotalLen: hlen + len(msg),
		TTL:      1,
		Protocol: 2,
		Dst:      net.ParseIP(address),
		Options:  options,
	}
	err = rawConn.WriteTo(header, msg, nil)
	if err != nil {
		log.Println("problem writing socket")
		log.Fatal(err)
	}
	/*
			n, err := ipConn.Write(msg)
			if err != nil {
				log.Fatal(err)
			}
		log.Printf("Wrote %d bytes\n", n)
	*/
}

func builtinJoin(address string) {
	c, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Println("Problem creating listen")
		log.Fatal(err)
	}
	defer c.Close()
	p := ipv4.NewPacketConn(c)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println("Problem resolving UDP address")
		log.Fatal(err)
	}
	en0, err := net.InterfaceByName("en0")
	if err != nil {
		log.Println("Problem getting interface")
		log.Fatal(err)
	}
	err = p.JoinGroup(en0, addr)
	if err != nil {
		log.Println("Problem joining group")
		log.Fatal(err)
	}
	defer p.LeaveGroup(nil, addr)
}

func rawLeave(address string) {
	ip := net.ParseIP(address)
	//ipConn, err := net.DialIP("ip:2", nil, ip)
	ipConn, err := net.ListenPacket("ip4:2", "0.0.0.0")
	if err != nil {
		log.Println("Failed to dial")
		log.Fatal(err)
	}
	defer ipConn.Close()

	seconds := 10
	ip4 := ip.To4()
	msg := []byte{0x17, byte(seconds * 10), 0, 0,
		ip4[0], ip4[1], ip4[2], ip4[3]}
	checksum := computeChecksum(msg)
	log.Printf("Checksum: %x\n", checksum)
	msg[2] = byte(checksum >> 8)   //byte(0xee) // byte(checksum & 0x00FF)
	msg[3] = byte(checksum & 0xFF) //byte(0x9b)          // byte((checksum & 0xFF00) >> 1)

	rawConn, err := ipv4.NewRawConn(ipConn)
	if err != nil {
		log.Println("probelem getting raw socket")
		log.Fatal(err)
	}
	defer rawConn.Close()
	// router alert: https://tools.ietf.org/html/rfc2113
	// option is 0x94
	// order is option flags or val, length, 0, 0
	options := []byte{0x94, 0x04, 0x0, 0x0}
	hlen := ipv4.HeaderLen + len(options)
	header := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      hlen,
		TOS:      0,
		TotalLen: hlen + len(msg),
		TTL:      1,
		Protocol: 2,
		Dst:      net.ParseIP("224.0.0.2"),
		Options:  options,
	}
	err = rawConn.WriteTo(header, msg, nil)
	if err != nil {
		log.Println("problem writing socket")
		log.Fatal(err)
	}
}

const (
	// subcommand keywords
	sendWord    = "send"
	receiveWord = "receive"
	queryWord   = "query"
	joinWord    = "join"
	leaveWord   = "leave"
	helpWord    = "help"

	// shared default values for subcommands
	defaultSendRecvAddress = "239.1.1.50"
	defaultSendRecvPort    = 5050
	defaultSendTTL         = 50
)

func showHelpMessage() {
	fmt.Printf("Specify a sub command of: %s, %s, %s, %s, %s, %s\n\n",
		sendWord, receiveWord, queryWord, joinWord, leaveWord, helpWord)
	fmt.Println("This program will allow you to test multicast and IGMP functionality.")
	fmt.Println("For help on a specific command, use 'help' followed by that command")
	fmt.Printf("Ex: mcast help %s\n\n", joinWord)
	flag.PrintDefaults()
	os.Exit(3)
}

func processSendCommand(sendGroup *string, sendPort *int, sendInterfaceIP, sendText *string, sendTTL, sendTOS, sendPadding, sendInterval, sendStart, sendMax *int) {
	if strings.Contains(*sendGroup, "/") { // is really a many sender
		network, mask, err := multicast.SplitCIDR(*sendGroup)
		if err != nil {
			fmt.Printf("Couldn't parse the mask")
			os.Exit(1)
		}
		s, err := multicast.NewManySender(network, int(mask), *sendPort, *sendTTL)
		if err != nil {
			fmt.Printf("Couldn't create many sender")
			os.Exit(1)
		}
		s.SetTOS(*sendTOS)
		s.SetMessagePadding(*sendPadding)
		sourceAddress := "host-chosen-address"
		if sendInterfaceIP != nil && *sendInterfaceIP != "" {
			err := s.SetLocalAddress(*sendInterfaceIP)
			if err != nil {
				fmt.Printf("There was a problem with the local address\n%v\n", err)
				os.Exit(1)
			}
			sourceAddress = *sendInterfaceIP
		}
		fmt.Printf("Sending from %v to %v:%d\n", sourceAddress, *sendGroup, *sendPort)
		s.Start(*sendText, *sendInterval, *sendStart, *sendMax)
	} else {
		s := multicast.NewSender(*sendGroup, *sendPort, *sendTTL)
		s.SetTOS(*sendTOS)
		s.SetMessagePadding(*sendPadding)
		sourceAddress := "host-chosen-address"
		if sendInterfaceIP != nil && *sendInterfaceIP != "" {
			err := s.SetLocalAddress(*sendInterfaceIP)
			if err != nil {
				fmt.Printf("There was a problem with the local address\n%v\n", err)
				os.Exit(1)
			}
			sourceAddress = *sendInterfaceIP
		}
		fmt.Printf("Sending from %v to %v:%d\n", sourceAddress, *sendGroup, *sendPort)
		if *sendMax == 1 {
			err := s.One(*sendText)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
		} else if *sendMax > 0 {
			err := s.Max(*sendText, *sendInterval, *sendStart, *sendMax)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
		} else {
			err := s.Forever(*sendText, *sendInterval, *sendStart)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
		}
	}
}

func processReceiveCommand(receiveGroup *string, receivePort *int, receiveInterface *string, receiveShowData *bool) {
	visibleInterface := "host-chosen"
	if *receiveInterface != "" {
		visibleInterface = *receiveInterface
	}
	fmt.Printf("Listening on %v:%d interface: %v\n", *receiveGroup, *receivePort, visibleInterface)
	err := multicast.Receive(*receiveGroup, *receivePort, *receiveInterface, *receiveShowData)
	if err != nil {
		fmt.Println("Problem receiving")
		fmt.Println(err)
	}
}

func processQueryCommand(queryInterface, queryInterfaceIP *string, queryInterval, queryMaxResponseTime *int, queryPlayNice *bool) {

	panic("Not implemented")
}

func processJoinCommand(joinGroup, joinInterface, joinInterfaceIP *string, joinRaw *bool, joinInterval *int, joinRouterAlert *bool, joinIGMPVersion *int) {

	panic("Not implemented")
}

func processLeaveCommand(leaveGroup, leaveInterface, leaveInterfaceIP *string, leaveInterval, leaveMax *int) {

	panic("Not implemented")
}

func processCommands() {
	// sub commands
	sendCommand := flag.NewFlagSet(sendWord, flag.ExitOnError)
	receiveCommand := flag.NewFlagSet(receiveWord, flag.ExitOnError)
	queryCommand := flag.NewFlagSet(queryWord, flag.ExitOnError)
	joinCommand := flag.NewFlagSet(joinWord, flag.ExitOnError)
	leaveCommand := flag.NewFlagSet(leaveWord, flag.ExitOnError)

	// send subcommand
	sendGroup := sendCommand.String("group", defaultSendRecvAddress, "destination multicast group address. Can use CIDR notation to send on multiple addresses.")
	sendPort := sendCommand.Int("port", defaultSendRecvPort, "destination port")
	sendInterfaceIP := sendCommand.String("interface-ip", "", "interface to use defined by IP addrress. default allows system to decide. must be in 0.0.0.0:0000 format")
	sendTTL := sendCommand.Int("ttl", defaultSendTTL, "IP ttl (time to live)")
	sendTOS := sendCommand.Int("tos", 0, "TOS / DSCP to be set. Only works on unicast addresses. 0xB8 for voice")
	sendText := sendCommand.String("text", "This is test number: {c}", "text to send to the receiver. Use '{c}' to access counter")
	sendPadding := sendCommand.Int("padding", 0, "Length to pad the message")
	sendInterval := sendCommand.Int("interval", 1000, "interval between sending messages (milliseconds).")
	sendStart := sendCommand.Int("start-value", 1, "non-negative start value message incrementer")
	sendMax := sendCommand.Int("max", 0, "number of packets to send. '0' for unlimited")

	// recieve subcommand
	receiveGroup := receiveCommand.String("group", defaultSendRecvAddress, "multicast group address to listen on. Can use CIDR notation to listen on multiple addresses.")
	receivePort := receiveCommand.Int("port", defaultSendRecvPort, "port to listen on")
	receiveInterface := receiveCommand.String("interface", "", "interface name use. default allows system to decide")
	receiveShowData := receiveCommand.Bool("show", true, "Print the data received to the console.")

	// query subcommand
	queryInterface := queryCommand.String("interface", "", "interface name use. default allows system to decide")
	queryInterfaceIP := queryCommand.String("interface-ip", "", "interface to use defined by IP addrress. default allows system to decide")
	queryInterval := queryCommand.Int("interval", 5, "interval between queries")
	queryMaxResponseTime := queryCommand.Int("max-response-time", 10, "maximum response time for the queried (seconds)")
	queryPlayNice := queryCommand.Bool("play-nice", false, "be silent if another querier is present")

	// join subcommand
	joinGroup := joinCommand.String("group", defaultSendRecvAddress, "multicast group to join")
	joinInterface := joinCommand.String("interface", "", "interface name use. default allows system to decide")
	joinInterfaceIP := joinCommand.String("interface-ip", "", "interface to use defined by IP addrress. default allows system to decide")
	joinRaw := joinCommand.Bool("raw", true, "send join as raw forged packet")
	joinInterval := joinCommand.Int("interval", 10, "interval between sending IGMP report 'join' (milliseconds)")
	joinRouterAlert := joinCommand.Bool("router-alert", true, "set router alert flag in IP packet")
	joinIGMPVersion := joinCommand.Int("igmp-version", 2, "igmp version to use for join")

	// leave subcommand
	leaveGroup := leaveCommand.String("group", defaultSendRecvAddress, "multicast group to send leave for")
	leaveInterface := leaveCommand.String("interface", "", "interface name use. default allows system to decide")
	leaveInterfaceIP := leaveCommand.String("interface-ip", "", "interface to use defined by IP addrress. default allows system to decide")
	leaveInterval := leaveCommand.Int("interval", 5, "interval between sending IGMP leave (milliseconds)")
	leaveMax := leaveCommand.Int("max", 1, "maximum IGMP leaves to send. 0 for infinite")

	// ensure at least 1 subcommand was specified
	if len(os.Args) < 2 {
		showHelpMessage()
	}

	args := os.Args[2:]
	switch os.Args[1] {
	case sendWord:
		sendCommand.Parse(args)
		processSendCommand(sendGroup, sendPort, sendInterfaceIP, sendText, sendTTL, sendTOS, sendPadding, sendInterval, sendStart, sendMax)
	case receiveWord:
		receiveCommand.Parse(args)
		processReceiveCommand(receiveGroup, receivePort, receiveInterface, receiveShowData)
	case queryWord:
		queryCommand.Parse(args)
		processQueryCommand(queryInterface, queryInterfaceIP, queryInterval, queryMaxResponseTime, queryPlayNice)
	case joinWord:
		joinCommand.Parse(args)
		processJoinCommand(joinGroup, joinInterface, joinInterfaceIP, joinRaw, joinInterval, joinRouterAlert, joinIGMPVersion)
	case leaveWord:
		leaveCommand.Parse(args)
		processLeaveCommand(leaveGroup, leaveInterface, leaveInterfaceIP, leaveInterval, leaveMax)
	case helpWord:
		if len(args) == 1 {
			switch args[0] {
			case sendWord:
				sendCommand.PrintDefaults()
				os.Exit(0)
			case receiveWord:
				receiveCommand.PrintDefaults()
				os.Exit(0)
			case queryWord:
				queryCommand.PrintDefaults()
				os.Exit(0)
			case joinWord:
				joinCommand.PrintDefaults()
				os.Exit(0)
			case leaveWord:
				leaveCommand.PrintDefaults()
				os.Exit(0)
			default:
				fmt.Printf("Use subcommand '%s' followed by a valid subcommand\n", helpWord)
				os.Exit(4)
			}
		}
		fallthrough
	default:
		showHelpMessage()
	}
}

func testBasic() {
	s := multicast.NewSender("239.1.1.9", 5051, 10)
	defer s.Close()
	//s.One("test")
	//s.Forever("test {c}", 1, 1)
	s.SetLocalAddress("172.16.0.178:5050")
	s.Max("test {c}", 1, 12, 3)
}

func testMany() {
	m, err := multicast.NewManySender("239.1.1.0", 32, 5051, 10)
	//m, err := multicast.NewManySender("255.255.255.255", 32, 5051, 10)
	//m, err := multicast.NewManySender("172.16.0.0", 32, 5051, 10)
	if err != nil {
		log.Println("Problem getting new many sender")
		log.Fatal(err)
	}
	m.SetMessagePadding(1400)
	m.SetTOS(0xB8)
	m.Start("This is test message {c}", 1000, 1, 1)
	time.Sleep(5 * time.Second)
}

func testReceive() {
	err := multicast.Receive("0.0.0.0", 5050, "", true)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Println("Starting...")
	//testBasic()
	processCommands()
	//testReceive()
	//fmt.Println(multicast.IPList("239.1.0.0", 23))
	//testMany()
	//sendText("test")
	//sendTextTTL("mytest")
	//receive()
	//query()
	//rawJoin("239.1.1.9")
	//rawJoinRouterAlert("239.1.1.10")
	//builtinJoin("239.1.1.6:5050")
	//rawLeave("239.1.1.7")
	log.Println("Terminating...")
}
