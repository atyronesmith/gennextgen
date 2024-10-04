package process

import "github.com/atyronesmith/gennextgen/pkg/types"

type TripleoData struct {
	Environment *types.TripleoOvercloudEnvironment
	Deployment  *types.TripleoOvercloudBaremetalDeployment
	Roles       *types.TripleoOvercloudRolesData
}

func GetTripleoData(toeFile string, tobdFile string, tordFile string) (*TripleoData, error) {
	environment, err := types.GetTripleoOvercloudEnvironment(toeFile)
	if err != nil {
		return nil, err
	}

	deployment, err := types.GetTripleoOvercloudBaremetalDeployment(tobdFile)
	if err != nil {
		return nil, err
	}

	roles, err := types.GetTripleoOvercloudRolesData(tordFile)
	if err != nil {
		return nil, err
	}

	return &TripleoData{
		Environment: environment,
		Deployment:  deployment,
		Roles:       roles,
	}, nil
}
