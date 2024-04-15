package utils

import (
	"fmt"
	"net"
)

func IpOffset(cidr string, offset int) (net.IP, error) {

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	ip = ip.Mask(ipnet.Mask)
	for i := 0; i < offset; inc(ip) {
		if !ipnet.Contains(ip) {
			return nil, fmt.Errorf("Offset %d out of range for %s", offset, cidr)
		}
		i++
	}

	return ip, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
