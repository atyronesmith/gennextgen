package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	k8Yaml "sigs.k8s.io/yaml"

	"gopkg.in/yaml.v3"
)

func YamlToMap(filename string) (map[string]interface{}, error) {

	yamlFile, err := os.ReadFile(string(filename))
	if err != nil {
		return nil, err
	}
	obj := make(map[string]interface{})

	err = yaml.Unmarshal(yamlFile, &obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func YamlToList(filename string) ([]map[string]interface{}, error) {

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var obj []map[string]interface{}

	err = yaml.Unmarshal(yamlFile, &obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func YamlToStruct[T any](filename string, obj *T) error {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, obj)
	if err != nil {
		return err
	}
	return nil
}

func StructToYaml(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("Error while Marshaling %T: %v", obj, err)
	}
	return b.Bytes(), nil
}

func StructToYamlK8s(obj interface{}) ([]byte, error) {
	yaml, err := k8Yaml.Marshal(obj)
	if err != nil {
		fmt.Printf("Error while Marshaling %T: %v", obj, err)
		return nil, err
	}
	return yaml, nil
}

func StructToJson(obj interface{}) ([]byte, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error while Marshaling %T: %v\n", obj, err)
		return nil, err
	}
	return b, nil
}
