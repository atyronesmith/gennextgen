package types

// NNCP represents a neural network control protocol.
type NNCP struct {
	DnsResolver DnsResolver  `yaml:"dns-resolver,omitempty"`
	Interfaces  []Interfaces `yaml:"interfaces"`
	Routes      Routes       `yaml:"routes,omitempty"`
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

	Ipv4            IpAddress       `yaml:"ipv4,omitempty"`
	Ipv6            IpAddress       `yaml:"ipv6,omitempty"`
	Ethernet        Ethernet        `yaml:"ethernet,omitempty"`
	SrIov           SrIov           `yaml:"sr-iov,omitempty"`
	LinkAggregation LinkAggregation `yaml:"link-aggregation,omitempty"`
	Vlan            Vlan            `yaml:"vlan,omitempty"`
}

type IpAddress struct {
	Enabled           bool      `yaml:"enabled"`
	Dhcp              bool      `yaml:"dhcp"`
	Address           []Address `yaml:"address"`
	AutoConf          bool      `yaml:"autoconf"`
	AutoDns           bool      `yaml:"auto-dns"`
	AutoGateway       bool      `yaml:"auto-gateway"`
	AutoRoutes        bool      `yaml:"auto-routes"`
	AutoRouteTabledId int       `yaml:"auto-route-table-id"`
	AutoRouteMetrics  int       `yaml:"auto-route-metrics"`
}

type Address struct {
	Ip           string `yaml:"ip"`
	PrefixLength int    `yaml:"prefix-length"`
	MptcpFlags   string `yaml:"mptcp-flags"`
}

type Ethernet struct {
	Speed           int    `yaml:"speed"`
	Duplex          Duplex `yaml:"duplex"`
	AutoNegotiation bool   `yaml:"auto-negotiation"`
}

type SrIov struct {
	TotalVfs int
	Vfs      []Vf
}

type Vf struct {
	MacAddress string
	SpoofCheck bool
	Trust      bool
	MinTxRate  int
	MaxTxRate  int
	VlandId    int
	Qos        int
	VlanProto  VlanProto
}

type LinkAggregation struct {
	Mode    Mode
	Options Options
	Port    []string
}

type Options struct {
	AllSlavesActive string
	ArpAllTargets   bool
	ArpInterval     int
	ArpValidate     bool
	DownDelay       int
	LpInterval      int
	Miimon          int
	MinLinks        int
	PacketsPerSlave int
	PrimaryReselect bool
	ResendIGMP      int
	UpDelay         int
	UseCarrier      bool
}

type Vlan struct {
	BaseIface string
	Id        int
	Protocol  VlanProto
}

type Routes struct {
	Config []RouteEntry
}

type RouteEntry struct {
	State            RouteState
	Destination      string `yaml:"destination"`
	NextHopInterface string `yaml:"next-hop-interface"`
	NextHopAddress   string `yaml:"next-hop-address,omitempty"`
	Metric           int    `yaml:"metric,omitempty"`
	TableId          int    `yaml:"table-id,omitempty"`
	Cwnd             int    `yaml:"cwnd,omitempty"`
}

type DnsResolver struct {
	Config struct {
		Server []string `yaml:"server"`
		Search []string `yaml:"search"`
	} `yaml:"config"`
}

type RouteState int

const (
	_ RouteState = iota
	StatePresent
	StateAbsent
)

type VlanProto int

const (
	_ VlanProto = iota
	Dot1Q
	Dot1AD
)

type Mode int

const (
	ActiveBackup Mode = iota
	BalanceRR
	BalanceXOR
	Broadcast
	_8023ad
	BalanceTLB
	BalanceALB
)

func (m Mode) String() string {
	switch m {
	case ActiveBackup:
		return "active-backup"
	case BalanceRR:
		return "balance-rr"
	case BalanceXOR:
		return "balance-xor"
	case Broadcast:
		return "broadcast"
	case _8023ad:
		return "802.3ad"
	case BalanceTLB:
		return "balance-tlb"
	case BalanceALB:
		return "balance-alb"
	default:
		return "Unknown"
	}
}

type State int

const (
	Up State = iota
	Down
	Ignore
	Absent
)

func (s State) String() string {
	switch s {
	case Up:
		return "up"
	case Down:
		return "down"
	case Ignore:
		return "ignore"
	case Absent:
		return "absent"
	default:
		return "Unknown"
	}
}

type Identifier int

const (
	Name Identifier = iota
	MacAddress
)

func (i Identifier) String() string {
	switch i {
	case Name:
		return "name"
	case MacAddress:
		return "mac-address"
	default:
		return "Unknown"
	}
}

type IntfType int

const (
	TypeEthernet IntfType = iota
	TypeVLAN
	TypeLinuxBridge
	TypeBond
)

func (i IntfType) String() string {
	switch i {
	case TypeEthernet:
		return "ethernet"
	case TypeVLAN:
		return "vlan"
	case TypeLinuxBridge:
		return "linux-bridge"
	case TypeBond:
		return "bond"
	default:
		return "Unknown"
	}
}

type Duplex int

const (
	Half Duplex = iota
	Full
)

func (d Duplex) String() string {
	switch d {
	case Half:
		return "half"
	case Full:
		return "full"
	default:
		return "Unknown"
	}
}
