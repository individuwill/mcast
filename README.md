# mcast

This is a little utility written to test out multicast on the network.

There are some features that allow for stress testing frame/packet routing
and switch/router state.

This utility is capable of entirely disrupting an L2 environment that isn't
robustly configured, so caution is urged.

mcast is a command line utility capable of sending and recieving multicast
or generic UDP traffic. It also allows simulation of IGMP joins, leaves, and
querying.
