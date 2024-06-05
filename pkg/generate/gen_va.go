package generate

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/atyronesmith/gennextgen/pkg/utils"
	"github.com/go-git/go-git/v5"
)

func CreateVADirs(root string, templateDir embed.FS) error {
	// Create the directory structure for the VA
	// default
	// ├── edpm
	// │   ├── deployment
	// │   |   ├── .gitignore
	// │   |   ├── kustomization.yaml
	// │   |   ├── values
	// │   |
	// │   ├── nodeset
	// │       ├── .gitignore
	// │       ├── kustomization.yaml
	// │       ├── values
	// │
	// ├── nncp
	// │   ├── kustomization.yaml
	// │   ├── values
	// │
	// ├── kustomization.yaml
	// ├── service-values.yaml

	// Create the root directory
	err := os.MkdirAll(root, os.ModePerm)
	if err != nil {
		return err
	}

	err = cloneVaLib(root)
	if err != nil {
		return err
	}

	// Create the subdirectories
	subdirs := []string{
		filepath.Join(root, "edpm", "deployment"),
		filepath.Join(root, "edpm", "nodeset"),
		filepath.Join(root, "nncp"),
	}
	for _, subdir := range subdirs {
		err := os.MkdirAll(subdir, os.ModePerm)
		if err != nil {
			// Handle the error
			panic(err)
		}
	}

	// Create the files
	files := []string{
		filepath.Join(root, "edpm", "deployment", ".gitignore"),
		filepath.Join(root, "edpm", "deployment", "kustomization.yaml"),
		filepath.Join(root, "edpm", "deployment", "values"),
		filepath.Join(root, "edpm", "nodeset", ".gitignore"),
		filepath.Join(root, "edpm", "nodeset", "kustomization.yaml"),
		filepath.Join(root, "edpm", "nodeset", "values"),
		filepath.Join(root, "nncp", "kustomization.yaml"),
		filepath.Join(root, "nncp", "values"),
		filepath.Join(root, "kustomization.yaml"),
		filepath.Join(root, "service-values.yaml"),
	}
	for _, file := range files {
		_, err := os.Create(file)
		if err != nil {
			return err
		}
	}

	err = createDeploymentValues(root)
	if err != nil {
		return err
	}

	err = createDeploymentKustomization(root, templateDir)
	if err != nil {
		return err
	}
	err = createNodesetKustomization(root, templateDir)
	if err != nil {
		return err
	}

	return nil
}

type EDPMDeploymentConfigMap struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string                 `yaml:"name"`
		Annotations map[string]interface{} `yaml:"annotations"`
	} `yaml:"metadata"`
	Data interface{} `yaml:"data"`
}

func createDeploymentValues(outDir string) error {
	var edcm EDPMDeploymentConfigMap

	edcm.APIVersion = "v1"
	edcm.Kind = "ConfigMap"
	edcm.Metadata.Name = "edpm-deployment-values"
	edcm.Metadata.Annotations = map[string]interface{}{
		"config.kubernetes.io/local-config": "true",
	}
	edcm.Data = make(map[string]interface{})

	// Convert the struct to YAML
	yaml, err := utils.StructToYamlK8s(edcm)
	if err != nil {
		return err
	}

	file := filepath.Join(outDir, "va", "edpm", "deployment")

	return utils.WriteByteData(yaml, file, "values.yaml")
}

func createDeploymentKustomization(outDir string, templateDir embed.FS) error {
	secret, err := utils.ProcessTemplate(templateDir, "edpm-deployment-kustomization.tmpl", "deployment", utils.GetFuncMap(), nil)
	if err != nil {
		return err
	}
	file := filepath.Join(outDir, "va", "edpm", "deployment")
	return utils.WriteByteData(secret.Bytes(), file, "kustomization.yaml")
}

func createNodesetKustomization(outDir string, templateDir embed.FS) error {
	nodeset, err := utils.ProcessTemplate(templateDir, "edpm-nodeset-kustomization.tmpl", "nodeset", utils.GetFuncMap(), nil)
	if err != nil {
		return err
	}
	file := filepath.Join(outDir, "va", "edpm", "nodeset")
	return utils.WriteByteData(nodeset.Bytes(), file, "kustomization.yaml")
}

func cloneVaLib(outDir string) error {
	// Clone the VA library
	// git clone

	// Clone the given repository to the given directory

	_, err := git.PlainClone(outDir, false, &git.CloneOptions{
		URL:      "https://github.com/openstack-k8s-operators/architecture.git",
		Progress: os.Stdout,
	})
	if err == git.ErrRepositoryAlreadyExists {
		fmt.Printf("Skipping clone...  Repository %s already exists.\n", outDir)
	} else if err != nil {
		fmt.Printf("Error cloning the repository: %s\n", err)
		return err
	}
	return nil
}
