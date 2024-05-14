package generate

import (
	"bytes"
	"fmt"

	"github.com/atyronesmith/gennextgen/pkg/types"
	nncp "github.com/atyronesmith/gennextgen/pkg/types/nncp"
	"github.com/atyronesmith/gennextgen/pkg/utils"
)

// Generate an NodeNetworkConfigurationPolicy for each worker node
// that will host control plane services.  The NNCP will contain
// information regarding the network interfaces, DNS resolver, and
// routes for each worker node.  The NNCP will be used to configure
// the network interfaces on the worker nodes.
// The NNCP will be generated in YAML format and written to the
// output directory.
// The NNCP will be generated based on the configuration download
// file that contains information about the networks that are
// available to the worker nodes.

func GenNNCP(outDir string, cdl *types.ConfigDownload) error {
	cfg := utils.GetConfig()

	nncpList := make([]nncp.NNCP, len(cfg.WorkerNodes))

	// TODO -- Initially, map 1:1 between OSP 17.1 Controllers and Worker nodes

	// Take the name of the worker node from the config file
	for _, worker := range cfg.WorkerNodes {
		n := nncp.NewNNCP()

		// The name of the policy
		n.ObjectMeta.Name = fmt.Sprintf("%s-%s", worker, cfg.Interface)

		n.ObjectMeta.Labels = map[string]string{
			"osp/interface": cfg.Interface,
		}

		n.Spec.NodeSelector = map[string]string{
			"kubernetes.io/hostname":         worker.Name,
			"node-role.kubernetes.io/worker": "",
		}

		for _, network := range cdl.Networks {
			var ni nncp.Interfaces

			ni.Description = fmt.Sprintf("%s vlan interface", network.Name)
			ni.State = nncp.Up
			ni.Mtu = network.Mtu
			if network.VlanId == 0 {
				ni.Name = cfg.Interface
				ni.IntfType = nncp.TypeEthernet
			} else {
				ni.Name = fmt.Sprintf("%s.%d", cfg.Interface, network.VlanId)
				ni.IntfType = nncp.TypeVLAN
				ni.Vlan = &nncp.Vlan{
					BaseIface: cfg.Interface,
					Id:        network.VlanId,
				}
			}

			ni.Ipv4 = &nncp.IpAddress{
				Enabled: true,
				Dhcp:    false,
				Address: []nncp.Address{
					{
						Ip:           "0.0.0.0",
						PrefixLength: network.PrefixLen,
					},
				},
			}
			ni.Ipv6 = &nncp.IpAddress{
				Enabled: false,
				Dhcp:    false,
			}
			n.Spec.DesiredState.Interfaces = append(n.Spec.DesiredState.Interfaces, ni)
		}

		n.Spec.DesiredState.DnsResolver = nncp.DnsResolver{
			Config: nncp.DnsConfig{
				Search:  []string{},
				Servers: []string{"10.1.1.1", "10.2.1.1"},
			},
		}

		n.Spec.DesiredState.Routes = nncp.Routes{
			Config: []nncp.RouteEntry{
				{
					Destination:      "192.168.122.100",
					NextHopInterface: "enp6s0.20",
					NextHopAddress:   "192.168.122.1",
				},
			},
		}

		n.Spec.NodeSelector = map[string]string{
			"kubernetes.io/hostname":         "worker1",
			"node-role.kubernetes.io/worker": "",
		}

		nncpList = append(nncpList, *n)
	}

	yaml, err := utils.StructToYaml(nncpList)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(yaml, outDir, "nncp.yaml")
	if err != nil {
		return err
	}

	//	fmt.Printf("%s", yaml)

	return nil
}

func GenNNCPFile(root string, nncpv []nncp.NNCP) error {
	var yamlOut bytes.Buffer

	for _, nncp := range nncpv {
		yaml, err := utils.StructToYaml(nncp)
		if err != nil {
			return err
		}
		yamlOut.WriteString("---\n")
		yamlOut.Write(yaml)
	}
	fmt.Printf("%s", yamlOut.String())

	return nil
}
