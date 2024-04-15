package netconf

import (
	"net"

	whereabouts "github.com/k8snetworkplumbingwg/whereabouts/pkg/types"
)

func WAAddRange(config *whereabouts.IPAMConfig, ipRange string, ipRangeStart string, ipRangeEnd string) {
	config.Range = string(net.ParseIP(ipRange))
	config.RangeStart = net.ParseIP(ipRangeStart)
	config.RangeEnd = net.ParseIP(ipRangeEnd)
}

// WAConfig returns a JSON byte array of the Whereabouts IPAM configuration
func WAConfig() *whereabouts.IPAMConfig {
	wa := whereabouts.IPAMConfig{}

	wa.Type = "whereabouts"

	return &wa
}
