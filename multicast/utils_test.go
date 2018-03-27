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
	"testing"
)

func TestSplitCIDRNoSlash(t *testing.T) {
	a, b, err := SplitCIDR("239.1.1.1")
	if err != nil {
		t.Error("Received error for ip without a slash", err)
	}
	if b != 32 {
		t.Error("Incorrect mask. Expecting 32, but got ", b)
	}
	if a != "239.1.1.1" {
		t.Error("Did not receive expected IP address. Expected 239.1.1.1 instead got ", a)
	}
}

func TestIP4ToIntAndIntToIP4(t *testing.T) {
	testList := []uint32{0, 2147483648, 4294967295}
	//for i := uint32(0); i <= 10000; i++ { // 4294967295
	for _, i := range testList {
		ip := IntToIP4(i)
		c := IP4ToInt(ip)
		if c != i {
			t.Errorf("IP to int conversion inconsistent. Expected %d, found %d for %v", i, c, ip)
		}
	}
}

func TestSplitCIDRWithSlash(t *testing.T) {
	a, b, err := SplitCIDR("239.1.1.1/24")
	if err != nil {
		t.Error("Received error for ip without a slash", err)
	}
	if b != 24 {
		t.Error("Incorrect mask. Expecting 32, but got ", b)
	}
	if a != "239.1.1.1" {
		t.Error("Did not receive expected IP address. Expected 239.1.1.1 instead got ", a)
	}
}

func TestIPListCIDRSingle(t *testing.T) {
	a, err := IPListCIDR("239.1.1.1")
	if err != nil {
		t.Error("Received error for basic IP conversion", err)
	}
	if len(a) != 1 {
		t.Error("Got more than 1 IP when expecting only 1", len(a))
	}
}

func TestIPListCIDR32(t *testing.T) {
	a, err := IPListCIDR("239.1.1.1/32")
	if err != nil {
		t.Error("Received error for IP with /32", err)
	}
	if len(a) != 1 {
		t.Errorf("Expected %d address, instead got %d", 1, len(a))
	}
}

func TestIPListCIDR24(t *testing.T) {
	a, err := IPListCIDR("239.1.1.0/24")
	if err != nil {
		t.Error("Received error for IP with /24", err)
	}
	if len(a) != 256 {
		t.Errorf("Expected %d addresses, instead got %d", 256, len(a))
	}
}

func TestIPListCIDRWrongNetwork(t *testing.T) {
	a, err := IPListCIDR("239.1.30.1/21")
	if err != nil {
		t.Error("Received error for  IP with /21", err)
	}
	if len(a) != 2048 {
		t.Errorf("Expected %d addresses, instead got %d", 2048, len(a))
	}
}
