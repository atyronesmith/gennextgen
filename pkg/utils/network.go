package utils

import (
	"fmt"
	"net/netip"
)

func IpOffset(cidr netip.Prefix, offset int) (netip.Addr, error) {

	zero := netip.Addr{}

	start := cidr.Addr().Next()
	for i := 0; i < offset; i++ {
		if start == zero {
			return zero, fmt.Errorf("Error getting IP offset")
		}
		start = start.Next()
	}
	return start, nil
}
