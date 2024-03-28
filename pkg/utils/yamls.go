package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func YamlToMap(filename string) (map[string]interface{}, error) {

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Unable to read file %s,  #%v ", filename, err)

		return nil, err
	}
	obj := make(map[string]interface{})

	err = yaml.Unmarshal(yamlFile, &obj)
	if err != nil {
		fmt.Printf("Error while Unmarshaling %s: %v", filename, err)

		return nil, err
	}

	return obj, nil
}

func YamlList(filename string) ([]map[string]interface{}, error) {

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Unable to read file %s,  #%v ", filename, err)

		return nil, err
	}

	var obj []map[string]interface{}

	err = yaml.Unmarshal(yamlFile, &obj)
	if err != nil {
		fmt.Printf("Error while Unmarshaling %s: %v", filename, err)
		return nil, err
	}

	return obj, nil
}
