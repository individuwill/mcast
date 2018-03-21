# mcast

## Description
This is a command line utility and library written to test out multicast
traffic flows on the network, and stress test the network and devices.

This utility is capable of entirely disrupting an L2 environment that isn't
robustly configured, so caution is urged.

mcast is a command line utility capable of sending and recieving multicast
or generic UDP traffic. It also allows simulation of IGMP joins, leaves, and
querying.

## License
This repository is licenses under GPLv3. See [LICENSE.md](./LICENSE.md) for details.

Copyright (C) 2018 Will Smith

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

## Usage

mcast is driven by subcommands. You invoke mcast like:

    $ mcast [-help] subcommand -option...

The list of subcommands are:

* help
* send
* receive
* join
* leave
* query

Each subcommand then has a set of options to control its behavior. Many of
the commands share similar options, and the option syntax is the same when
this is the case. Below is a detailed explanation of each subcommand and
its associated options

### help

### send
Will send UDP traffic to an IP address, usually a multicast one. Will send continuously
in a loop at specified interval until the program is terminated or max number of messages
are sent.

The options are:
* -group : IP destination address. Can use CIDR notation to send to multiple address
    * default : 239.1.1.5
* -port : Destination UDP port
    * default : 5050
* -interface-ip : Interface to use defined by IP addrress. Must be in 0.0.0.0:0000 format. Default allows system to decide. 
* -ttl : IP ttl (time to live)
    * default : 50
* -tos : TOS / DSCP to be set. Only works on unicast addresses. 0xB8 for voice.
    * default : 0
* -text : Text / data to send to the receiver. Use '{c}' to access counter
    * default : This is test number: {c}
* -padding : Length to pad the message. Will make message take up specified number of bytes.
    * default : 0
* -interval : Interval between sending messages (milliseconds).
    * default : 1000
* -start-value : Non-negative start value message incrementer / counter
    * default : 1
* -max : Number of packets to send. '0' for continuous send
    * default : 0

### receive

### join

### leave

### query