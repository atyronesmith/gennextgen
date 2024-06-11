package types

import (
	"fmt"
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

		for k, v := range cfgSet {
			varName := strings.ReplaceAll(k, "::", ".")
			if _, ok := cdl.Roles[roleIndex].ConfigSettings[varName]; ok {
				return fmt.Errorf("ProcessConfigSettings: <%s> already exists for role %s", varName, role.Name)
			}
			cdl.Roles[roleIndex].ConfigSettings[varName] = v
		}
	}

	return nil
}
