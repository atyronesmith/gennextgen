package netconf

type NetConfMacvlan struct {
	NetConf
	Master     string `json:"master"`
	Mode       string `json:"mode"`
	MTU        int    `json:"mtu"`
	Mac        string `json:"mac,omitempty"`
	LinkContNs bool   `json:"linkInContainer,omitempty"`
}

func NewNetConfMacvlan() *NetConfMacvlan {
	nc := &NetConfMacvlan{}

	nc.CNIVersion = "0.3.1"
	nc.Type = "macvlan"

	return nc
}
