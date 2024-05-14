package generate

import (
	"embed"
	_ "embed"

	"github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/utils"
	v1 "k8s.io/api/core/v1"
)

func GenSecrets(mappingYaml string, templateDir embed.FS, outDir string, cdl *types.ConfigDownload) error {

	parameterDefaults, err := utils.YamlToMap(utils.GetFullPath(utils.OVERCLOUD_PASSWORDS))
	if err != nil {
		return err
	}
	passwords := parameterDefaults["parameter_defaults"].(map[string]interface{})

	err = genKeystoneSecret(templateDir, outDir, passwords)
	if err != nil {
		return err
	}

	err = genOpenStackSecret(outDir, cdl)
	if err != nil {
		return err
	}
	return nil
}

func genKeystoneSecret(templateDir embed.FS, outDir string, passwords map[string]interface{}) error {

	secret, err := utils.ProcessTemplate(templateDir, "keystone-secret.tmpl", "keystone", utils.GetFuncMap(), passwords)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(secret.Bytes(), outDir, "keystone-secret.yaml")

	return err
}

func genOpenStackSecret(outDir string, cdl *types.ConfigDownload) error {

	secret := v1.Secret{}
	secret.Type = "Opaque"
	secret.Namespace = "openstack"
	secret.Name = "osp-secret"
	secret.Data = make(map[string][]byte)

	for index, s := range cdl.Passwords {
		secret.Data[index] = []byte(s)
	}

	y, err := utils.StructToYamlK8s(secret)
	if err != nil {
		return err
	}

	err = utils.WriteByteData(y, outDir, "osp-secret.yaml")

	return err
}
