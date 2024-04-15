package generate

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/atyronesmith/gennextgen/pkg/types"
	nncp "github.com/atyronesmith/gennextgen/pkg/types/nncp"
	"github.com/atyronesmith/gennextgen/pkg/utils"
)

func GenNNCP(outDir string, cdl *types.ConfigDownload) error {
	cfg := utils.GetConfig()

	nncpList := make([]nncp.NNCP, len(cfg.WorkerNodes.Names))

	for _, worker := range cfg.WorkerNodes.Names {
		n := nncp.NewNNCP()

		// Need to know the name of the OCP worker node where an OSP service will run
		n.ObjectMeta.Name = worker

		for _, network := range cdl.Networks {
			var ni nncp.Interfaces

			ni.Description = fmt.Sprintf("%s vlan interface", network.Name)
			ni.State = nncp.Up
			ni.Mtu = network.Mtu
			if network.VlanId == 0 {
				ni.Name = cfg.WorkerNodes.Interface
				ni.IntfType = nncp.TypeEthernet
			} else {
				ni.Name = fmt.Sprintf("%s.%d", cfg.WorkerNodes.Interface, network.VlanId)
				ni.IntfType = nncp.TypeVLAN
				ni.Vlan = &nncp.Vlan{
					BaseIface: cfg.WorkerNodes.Interface,
					Id:        network.VlanId,
				}
			}
			pl, err := strconv.Atoi(strings.Split(network.Cidr, "/")[1])
			if err != nil {
				return err
			}

			ni.Ipv4 = &nncp.IpAddress{
				Enabled: true,
				Dhcp:    false,
				Address: []nncp.Address{
					{
						Ip:           "0.0.0.0",
						PrefixLength: pl,
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
