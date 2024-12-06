package process

import "github.com/atyronesmith/gennextgen/pkg/types"

type TripleoData struct {
	Environment *types.TripleoOvercloudEnvironment
	Deployment  *types.TripleoOvercloudBaremetalDeployment
	Roles       *types.TripleoOvercloudRolesData
}

func GetTripleoData(environmentFile string, baremetalDeployment string, rolesData string) (*TripleoData, error) {
	environment, err := types.GetTripleoOvercloudEnvironment(environmentFile)
	if err != nil {
		return nil, err
	}

	deployment, err := types.GetTripleoOvercloudBaremetalDeployment(baremetalDeployment)
	if err != nil {
		return nil, err
	}

	roles, err := types.GetTripleoOvercloudRolesData(rolesData)
	if err != nil {
		return nil, err
	}

	return &TripleoData{
		Environment: environment,
		Deployment:  deployment,
		Roles:       roles,
	}, nil
}
