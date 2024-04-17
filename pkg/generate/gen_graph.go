package generate

import (
	"github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/utils"
)

func GenGraph(outDir string, cdl *types.ConfigDownload) error {
	networks, err := utils.ProcessTemplate("network.tmpl", "networks", utils.GetFuncMap(), cdl)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(networks.Bytes(), outDir, "networks.dot")

	return err
}
