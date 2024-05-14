package nncp

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewNNCP() *NNCP {
	return &NNCP{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "nmstate.io/v1",
			Kind:       "NodeNetworkConfigurationPolicy",
		},
	}
}

// NNCP represents a neural network control protocol.
type NNCP struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata,omitempty"`

	Spec Spec `yaml:"spec"`
}

type Spec struct {
	DesiredState DesiredState      `yaml:"desiredState"`
	NodeSelector map[string]string `yaml:"nodeSelector"`
}

// NetworkManager backend has two set of DNS configurations:

// Global DNS set by D-BUS interface or NetworkManager.conf.
// Interface DNS stored in NetworkManager connection as ipv4.dns or ipv6.dns.
// Nmstate will try to use global DNS via D-BUS interface call, and only use interface level DNS for any of these use case:

// Has IPv6 link-local address as name server: e.g. fe80::deef:1%eth1
// User want static DNS server appended before dynamic one. In this case, user should define auto-dns: true explicitly along with static DNS.
// User want to force DNS server stored in interface for static IP interface. This case, user need to state static DNS config along with static IP config.

type DesiredState struct {
	DnsResolver DnsResolver  `yaml:"dns-resolver,omitempty"`
	Interfaces  []Interfaces `yaml:"interfaces"`
	Routes      Routes       `yaml:"routes,omitempty"`
}

type Metadata struct {
	Labels map[string]string `yaml:"labels,omitempty"`
	Name   string            `yaml:"name"`
}

type Interfaces struct {
	Name        string     `yaml:"name"`
	IntfType    IntfType   `yaml:"type"`
	State       State      `yaml:"state"`
	Description string     `yaml:"description"`
	ProfileName string     `yaml:"profile-name,omitempty"`
	Identifier  Identifier `yaml:"identifier,omitempty"`
	MacAddress  string     `yaml:"mac-address,omitempty"`
	Mtu         int        `yaml:"mtu,omitempty"`
	MinMtu      int        `yaml:"min-mtu,omitempty"`
	MaxMtu      int        `yaml:"max-mtu,omitempty"`

	Ipv4            *IpAddress       `yaml:"ipv4,omitempty"`
	Ipv6            *IpAddress       `yaml:"ipv6,omitempty"`
	Ethernet        *Ethernet        `yaml:"ethernet,omitempty"`
	SrIov           *SrIov           `yaml:"sr-iov,omitempty"`
	LinkAggregation *LinkAggregation `yaml:"link-aggregation,omitempty"`
	Vlan            *Vlan            `yaml:"vlan,omitempty"`
}

type IpAddress struct {
	Enabled           bool      `yaml:"enabled"`
	Dhcp              bool      `yaml:"dhcp"`
	Address           []Address `yaml:"address,omitempty"`
	AutoConf          bool      `yaml:"autoconf,omitempty"`
	AutoDns           bool      `yaml:"auto-dns,omitempty"`
	AutoGateway       bool      `yaml:"auto-gateway,omitempty"`
	AutoRoutes        bool      `yaml:"auto-routes,omitempty"`
	AutoRouteTabledId int       `yaml:"auto-route-table-id,omitempty"`
	AutoRouteMetrics  int       `yaml:"auto-route-metrics,omitempty"`
}

type Address struct {
	Ip           string `yaml:"ip,omitempty"`
	PrefixLength int    `yaml:"prefix-length,omitempty"`
	MptcpFlags   string `yaml:"mptcp-flags,omitempty"`
}

type Ethernet struct {
	Speed           int    `yaml:"speed,omitempty"`
	Duplex          Duplex `yaml:"duplex,omitempty"`
	AutoNegotiation bool   `yaml:"auto-negotiation,omitempty"`
}

type SrIov struct {
	TotalVfs int  `yaml:"total-vfs,omitempty"`
	Vfs      []Vf `yaml:"vfs,omitempty"`
}

type Vf struct {
	MacAddress string    `yaml:"mac-address,omitempty"`
	SpoofCheck bool      `yaml:"spoof-check,omitempty"`
	Trust      bool      `yaml:"trust,omitempty"`
	MinTxRate  int       `yaml:"min-tx-rate,omitempty"`
	MaxTxRate  int       `yaml:"max-tx-rate,omitempty"`
	VlandId    int       `yaml:"vlan-id,omitempty"`
	Qos        int       `yaml:"qos,omitempty"`
	VlanProto  VlanProto `yaml:"vlan-protocol,omitempty"`
}

type LinkAggregation struct {
	Mode    Mode     `yaml:"mode,omitempty"`
	Options Options  `yaml:"options,omitempty"`
	Port    []string `yaml:"port,omitempty"`
}

type Options struct {
	AllSlavesActive string `yaml:"all-slaves-active,omitempty"`
	ArpAllTargets   bool   `yaml:"arp-all-targets,omitempty"`
	ArpInterval     int    `yaml:"arp-interval,omitempty"`
	ArpValidate     bool   `yaml:"arp-validate,omitempty"`
	DownDelay       int    `yaml:"down-delay,omitempty"`
	LpInterval      int    `yaml:"lp-interval,omitempty"`
	Miimon          int    `yaml:"miimon,omitempty"`
	MinLinks        int    `yaml:"min-links,omitempty"`
	PacketsPerSlave int    `yaml:"packets-per-slave,omitempty"`
	PrimaryReselect bool   `yaml:"primary-reselect,omitempty"`
	ResendIGMP      int    `yaml:"resend-igmp,omitempty"`
	UpDelay         int    `yaml:"up-delay,omitempty"`
	UseCarrier      bool   `yaml:"use-carrier,omitempty"`
}

type Vlan struct {
	BaseIface string    `yaml:"base-interface"`
	Id        int       `yaml:"id"`
	Protocol  VlanProto `yaml:"protocol,omitempty"`
}

type Routes struct {
	Config []RouteEntry `yaml:"config,omitempty"`
}

type RouteEntry struct {
	State            RouteState `yaml:"state,omitempty"`
	Destination      string     `yaml:"destination"`
	NextHopInterface string     `yaml:"next-hop-interface"`
	NextHopAddress   string     `yaml:"next-hop-address,omitempty"`
	Metric           int        `yaml:"metric,omitempty"`
	TableId          int        `yaml:"table-id,omitempty"`
	Cwnd             int        `yaml:"cwnd,omitempty"`
}

type DnsResolver struct {
	Config DnsConfig `yaml:"config"`
}

type DnsConfig struct {
	Servers []string `yaml:"server"`
	Search  []string `yaml:"search"`
}

type RouteState string

const (
	RSabsent RouteState = "absent"
)

type VlanProto string

const (
	Dot1Q  VlanProto = "802.1q"
	Dot1AD VlanProto = "802.1ad"
)

type Mode string

const (
	ActiveBackup Mode = "active-backup"
	BalanceRR    Mode = "balance-rr"
	BalanceXOR   Mode = "balance-xor"
	Broadcast    Mode = "broadcast"
	_8023ad      Mode = "802.3ad"
	BalanceTLB   Mode = "balance-tlb"
	BalanceALB   Mode = "balance-alb"
)

type State string

const (
	Up     State = "up"
	Down   State = "down"
	Ignore State = "ignore"
	Absent State = "absent"
)

type Identifier string

const (
	Name       Identifier = "name"
	MacAddress Identifier = "mac-address"
)

type IntfType string

const (
	TypeEthernet    IntfType = "ethernet"
	TypeVLAN        IntfType = "vlan"
	TypeLinuxBridge IntfType = "linux-bridge"
	TypeBond        IntfType = "bond"
)

type Duplex string

const (
	Half Duplex = "half"
	Full Duplex = "full"
)
