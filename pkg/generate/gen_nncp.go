package generate

import (
	"github.com/atyronesmith/gennextgen/pkg/types"
)

func GenNNCP(root string) error {
	var nncp types.NNCP

	nncp.Interfaces = append(nncp.Interfaces, types.Interfaces{
		Name:     "enp6s0.20",
		State:    types.Up,
		IntfType: types.TypeVLAN,
		Vlan: types.Vlan{
			BaseIface: "eth0",
			Id:        10,
		},
		Ipv4: types.IpAddress{
			Enabled: true,
			Dhcp:    false,
			Address: []types.Address{{
				Ip:           "192.168.111.1",
				PrefixLength: 24,
			}},
		},
		Ipv6: types.IpAddress{
			Enabled: false,
		},
	})

	return nil
}
