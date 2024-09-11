package types

import (
	"path/filepath"
	"strings"

	"github.com/atyronesmith/gennextgen/pkg/utils"
)

func ProcessConfigSettings(cdl *ConfigDownload) error {
	for _, role := range cdl.Roles {
		configPath := filepath.Join("config-download", "overcloud", role.Name, "config_settings.yaml")
		cfgSet, err := utils.YamlToMap(utils.GetFullPath(configPath))
		if err != nil {
			return err
		}

		for k, v := range cfgSet {
			varName := strings.ReplaceAll(k, "::", ".")
			err := SetConfigSetting(role, varName, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
