package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func MarshalArray[T any](mapInput interface{}, things *[]T) {
	marshalledInput, err := yaml.Marshal(mapInput)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(marshalledInput, &things)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func MarshalMap[T any](mapInput interface{}, things *T) {
	marshalledInput, err := yaml.Marshal(mapInput)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(marshalledInput, &things)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
