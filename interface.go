/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : interface.go

* Purpose :

* Creation Date : 04-13-2017

* Last Modified : Thu 13 Apr 2017 11:57:55 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"net"
)

var myips = selfIPs()

func isMyIP(ip net.IP) bool {
	for _, v := range myips {
		if ip.Equal(v) {
			return true
		}
	}
	return false
}

func selfIPs() []net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return []net.IP{}
	}
	var ips []net.IP
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips
}
