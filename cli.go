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

// The mcast cli program for sending, receiving, joining
// and querying UDP and multicast traffic
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/individuwill/mcast/multicast"
)

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

func processJoinCommand(joinGroup *string, joinPort *int, joinInterface *string, joinRaw *bool, joinInterval *int, joinRouterAlert *bool, joinIGMPVersion *int) {
	if *joinRaw {
		err := multicast.JoinRaw(*joinGroup, *joinPort, *joinInterface, *joinInterval, *joinRouterAlert, *joinIGMPVersion)
		if err != nil {
			fmt.Println("Problem with raw join of the group")
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if *joinRouterAlert {
			fmt.Println("Can not enable router alert manually with non-raw join.")
			os.Exit(1)
		}
		if strings.Contains(*joinGroup, "/") {
			fmt.Println("Can not use CIDR notation without raw option")
			os.Exit(1)
		}
		err := multicast.Join(*joinGroup, *joinPort, *joinInterface)
		if err != nil {
			fmt.Println("Problem joining the group")
			fmt.Println(err)
			os.Exit(1)
		}
	}
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
	joinGroup := joinCommand.String("group", defaultSendRecvAddress, "multicast group to join. Can use CIDR notation for multiple joins when raw option is used.")
	joinPort := joinCommand.Int("port", 5050, "Port to use for join.")
	joinInterface := joinCommand.String("interface", "", "interface name use. default allows system to decide")
	joinRaw := joinCommand.Bool("raw", false, "send join as raw forged packet")
	joinInterval := joinCommand.Int("interval", 10000, "interval between sending IGMP report 'join' (milliseconds). Can only be used if raw option is enabled.")
	joinRouterAlert := joinCommand.Bool("router-alert", false, "set router alert flag in IP packet. Can only be used if raw option is enabled.")
	joinIGMPVersion := joinCommand.Int("igmp-version", 2, "igmp version to use for join. Can only be used if raw option is enabled.")

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
		processJoinCommand(joinGroup, joinPort, joinInterface, joinRaw, joinInterval, joinRouterAlert, joinIGMPVersion)
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
