package generate

import (
	"testing"

	"github.com/atyronesmith/gennextgen/pkg/types/nncp"
)

func TestGenNNCPFile(t *testing.T) {
	type args struct {
		root  string
		nncpv []nncp.NNCP
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestGenNNCPFile", args: args{root: "/tmp",
			nncpv: genNNCP()}, wantErr: false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenNNCPFile(tt.args.root, tt.args.nncpv); (err != nil) != tt.wantErr {
				t.Errorf("GenNNCPFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func genNNCP() []nncp.NNCP {
	var nncpv []nncp.NNCP

	worker1 := nncp.NewNNCP()

	worker1.Name = "osp-enp6s0-worker-1"

	worker1.Spec.DesiredState.Interfaces = append(worker1.Spec.DesiredState.Interfaces, nncp.Interfaces{
		Description: "internalapi vlan interface",
		Name:        "enp6s0.20",
		State:       nncp.Up,
		IntfType:    nncp.TypeVLAN,
		Vlan: &nncp.Vlan{
			BaseIface: "enp6s0",
			Id:        20,
		},
		Ipv4: &nncp.IpAddress{
			Enabled: true,
			Dhcp:    false,
			Address: []nncp.Address{{
				Ip:           "172.17.0.10",
				PrefixLength: 24,
			}},
		},
		Ipv6: &nncp.IpAddress{
			Enabled: false,
		},
	})

	worker1.Spec.DesiredState.Interfaces = append(worker1.Spec.DesiredState.Interfaces, nncp.Interfaces{
		Description: "storage vlan interface",
		Name:        "enp6s0.21",
		State:       nncp.Up,
		IntfType:    nncp.TypeVLAN,
		Vlan: &nncp.Vlan{
			BaseIface: "enp6s0",
			Id:        21,
		},
		Ipv4: &nncp.IpAddress{
			Enabled: true,
			Dhcp:    false,
			Address: []nncp.Address{{
				Ip:           "172.18.0.10",
				PrefixLength: 24,
			}},
		},
		Ipv6: &nncp.IpAddress{
			Enabled: false,
		},
	})

	worker1.Spec.DesiredState.Interfaces = append(worker1.Spec.DesiredState.Interfaces, nncp.Interfaces{
		Description: "Configuring enp6s0",
		Name:        "enp6s0",
		State:       nncp.Up,
		IntfType:    nncp.TypeEthernet,
		Mtu:         1500,
		Vlan: &nncp.Vlan{
			BaseIface: "enp6s0",
			Id:        21,
		},
		Ipv4: &nncp.IpAddress{
			Enabled: true,
			Dhcp:    false,
			Address: []nncp.Address{{
				Ip:           "172.18.0.10",
				PrefixLength: 24,
			}},
		},
		Ipv6: &nncp.IpAddress{
			Enabled: false,
		},
	})

	worker1.Spec.DesiredState.DnsResolver = nncp.DnsResolver{
		Config: nncp.DnsConfig{
			Search:  []string{},
			Servers: []string{"10.1.1.1", "10.2.1.1"},
		},
	}

	worker1.Spec.DesiredState.Routes = nncp.Routes{
		Config: []nncp.RouteEntry{
			{
				Destination:      "192.168.122.100",
				NextHopInterface: "enp6s0.20",
				NextHopAddress:   "192.168.122.1",
			},
		},
	}

	worker1.Spec.NodeSelector = map[string]string{
		"kubernetes.io/hostname":         "worker10",
		"node-role.kubernetes.io/worker": "",
	}

	nncpv = append(nncpv, *worker1)

	return nncpv
}
