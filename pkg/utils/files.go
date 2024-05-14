package utils

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rootDir string = "./overcloud-deploy"

const (
	TRIPLEO_OVERCLOUD_ROLES_DATA   string = "tripleo-overcloud-roles-data.yaml"
	TRIPLEO_OVERCLOUD_ENVIRONMENT  string = "environment/tripleo-overcloud-environment.yaml"
	TRIPLEO_ANSIBLE_INVENTORY_YAML string = "tripleo-ansible-inventory.yaml"
	OVERCLOUD_EXPORT               string = "overcloud-export.yaml"
	OVERCLOUD_PASSWORDS            string = "overcloud-passwords.yaml"
	GLOBAL_VARS                    string = "config-download/overcloud/global_vars.yaml"
	TRIPLEO_OVERCLOUD_NETWORK_DATA string = "tripleo-overcloud-network-data.yaml"
	BAREMETAL_DEPLOY               string = "tripleo-overcloud-baremetal-deployment.yaml"
	DEPLOY_STEPS_ONE               string = "config-download/overcloud/external_deploy_steps_tasks_step1.yaml"
	// Add more file names as needed
)

func GetFullPath(fileName string) string {
	return rootDir + "/" + fileName
}

func SetRootDir(dir string) {
	rootDir = dir
}

func SearchFile(targetFileName string) (string, error) {
	var foundFilePath string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() && info.Name() == targetFileName {
			foundFilePath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if foundFilePath == "" {
		return "", fmt.Errorf("file not found: %s starting at %s", targetFileName, rootDir)
	}
	return foundFilePath, nil
}

func WalkDir(dirname string) {
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("dir: %v: name: %s\n", info.IsDir(), path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func SearchFileRegex(path string, targetFileRegex string) ([]string, error) {

	targetFileRegex = strings.ReplaceAll(targetFileRegex, "*", ".*")
	targetFileRegex = strings.ReplaceAll(targetFileRegex, "?", ".")

	fileReg, err := regexp.Compile(targetFileRegex)
	if err != nil {
		return nil, err
	}

	fileList := make([]string, 0)

	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && fileReg.Match([]byte(filepath.Base(path))) {
			fullPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			fileList = append(fileList, fullPath)
		}
		return nil
	})

	return fileList, err
}

func WriteByteData(buf []byte, dir string, fileName string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("unable to mkdir: %s", dir)
		}
	}
	dFile, err := os.Create(dir + "/" + fileName)
	if err != nil {
		fmt.Printf("unable to create/open file: %s", fileName)
	}
	defer dFile.Close()

	fmt.Printf("Writing %s...\n", fileName)

	if _, err = dFile.Write(buf); err != nil {
		return fmt.Errorf("unable to write instructions to: %s", dir+"/"+fileName)
	}
	return nil
}

func GetTripleoOvercloudNetworkData() ([]map[string]interface{}, error) {
	// Assume for now a 1:1 mapping between the number of 17.1 OSP controllers and
	// the number of OCP master nodes
	networks, err := YamlToList(GetFullPath(TRIPLEO_OVERCLOUD_NETWORK_DATA))
	if err != nil {
		return nil, err
	}

	return networks, nil
}
