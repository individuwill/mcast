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
	"strings"
	"sync"

	"golang.org/x/net/ipv4"
)

func joinRaw(address string, port int, interfaceName string, interval int, routerAlert bool, igmpVersion int) error {
	/*

		d := time.Duration(interval) * time.Millisecond
		time.Sleep(d)
	*/

	return nil
}

func JoinRaw(address string, port int, interfaceName string, interval int, routerAlert bool, igmpVersion int) error {
	if strings.Contains(address, "/") {
		ips, err := IPListCIDR(address)
		if err != nil {
			return err
		}
		wg := new(sync.WaitGroup)
		for _, ip := range ips {
			wg.Add(1)
			func(ipAddress string) {
				err := joinRaw(ipAddress, port, interfaceName, interval, routerAlert, igmpVersion)
				if err != nil {
					log.Printf("Problem with join of %v:%d\n", ipAddress, port)
					log.Println(err)
				}
				wg.Done()
			}(ip.String())
		}
		wg.Wait()
		return nil
	}
	return joinRaw(address, port, interfaceName, interval, routerAlert, igmpVersion)
}

// Join will use the system built in IGMP group join mechanisms to join a group.
// You may not see any IGMP requests sent if the system isn't ready to send them
// (group is currently joined, and timers are good). Join does not explicitly
// send a leave request.
func Join(joinIPString string, joinPort int, interfaceName string) error {
	address := fmt.Sprintf("%v:%d", joinIPString, joinPort)
	c, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Println("Problem with initial listen")
		return err
	}
	defer c.Close()

	p := ipv4.NewPacketConn(c)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println("Problem resolving address")
		return err
	}
	joinInterface, err := GetInterface(interfaceName)
	if err != nil {
		log.Println("Problem gettin interface")
		return err
	}

	err = p.JoinGroup(joinInterface, addr)
	if err != nil {
		log.Println("Problem joining")
		return err
	}
	return nil
}
