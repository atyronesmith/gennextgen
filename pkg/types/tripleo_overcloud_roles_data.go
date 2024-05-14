package types

import "github.com/atyronesmith/gennextgen/pkg/utils"

type TripleoOvercloudRolesData []TripleoOvercloudRolesDataEntry

type TripleoOvercloudRolesDataEntry struct {
	HostnameFormatDefault   *string                `json:"HostnameFormatDefault,omitempty"`
	RoleParametersDefault   map[string]interface{} `json:"RoleParametersDefault"`
	ServicesDefault         []string               `json:"ServicesDefault"`
	DefaultRouteNetworks    []string               `json:"default_route_networks,omitempty"`
	DeprecatedNICConfigName *string                `json:"deprecated_nic_config_name,omitempty"`
	DeprecatedParamFlavor   *string                `json:"deprecated_param_flavor,omitempty"`
	DeprecatedParamImage    *string                `json:"deprecated_param_image,omitempty"`
	Description             string                 `json:"description"`
	Name                    string                 `json:"name"`
	Networks                map[string]interface{} `json:"networks"`
	Tags                    []string               `json:"tags"`
	UpdateSerial            int64                  `json:"update_serial"`
	UsesDeprecatedParams    *bool                  `json:"uses_deprecated_params,omitempty"`
}

var tripleoOvercloudRolesData TripleoOvercloudRolesData

func GetTripleoOvercloudRolesData() (*TripleoOvercloudRolesData, error) {
	if len(tripleoOvercloudRolesData) == 0 {
		tripleoOvercloudRolesData = TripleoOvercloudRolesData{}
		err := utils.YamlToStruct(utils.GetFullPath(utils.TRIPLEO_OVERCLOUD_ROLES_DATA), &tripleoOvercloudRolesData)
		if err != nil {
			return nil, err
		}
	}
	return &tripleoOvercloudRolesData, nil

}
