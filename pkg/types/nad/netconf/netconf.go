package netconf

type NetConf struct {
	CNIVersion string `json:"cniVersion,omitempty"`

	Name         string          `json:"name,omitempty"`
	Type         string          `json:"type,omitempty"`
	Capabilities map[string]bool `json:"capabilities,omitempty"`
	IPAM         interface{}     `json:"ipam,omitempty"`
	DNS          string          `json:"dns,omitempty"`
}
