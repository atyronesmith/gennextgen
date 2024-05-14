package types

import (
	"path/filepath"
	"strings"

	"github.com/atyronesmith/gennextgen/pkg/utils"
)

func ProcessConfigSettings(cdl *ConfigDownload) error {
	for roleIndex, role := range cdl.Roles {
		configPath := filepath.Join("config-download", "overcloud", role.Name, "config_settings.yaml")
		cfgSet, err := utils.YamlToMap(utils.GetFullPath(configPath))
		if err != nil {
			return err
		}
		csm := make(map[string][]TripleoRoleConfigSetting)
		cdl.Roles[roleIndex].ConfigSettings = csm

		for k, v := range cfgSet {
			path := strings.Split(k, "::")
			settingKey := path[len(path)-1]
			cs := TripleoRoleConfigSetting{
				Service: path[0],
				Path:    k,
				Value:   v,
			}
			if len(path) > 1 {
				cs.Section = path[1]
			}
			csm[settingKey] = append(csm[settingKey], cs)
		}
	}

	return nil
}
