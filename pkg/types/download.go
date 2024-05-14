package types

import (
	"fmt"
	"os"

	"github.com/atyronesmith/gennextgen/pkg/utils"
)

type DownloadMap map[string]interface{}

func NewDownloadMap() DownloadMap {
	return make(DownloadMap)
}

func (dm *DownloadMap) ProcessDir(baseDir string) {
	fileList, err := utils.SearchFileRegex(baseDir, "*.yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range fileList {
		if err := dm.ParseDownloadFile(file); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}

func (dm *DownloadMap) ParseDownloadFile(path string) error {
	m, err := utils.YamlToMap(path)
	if err != nil {
		return err
	}

	for k, v := range m {
		(*dm)[k] = v
	}

	return nil
}
