package types

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/atyronesmith/gennextgen/pkg/utils"
)

type ArgsType interface {
	int | bool | string
}

type ArgDefs map[string]ArgDef

type ArgDef struct {
	Name         string      `yaml:"name"`
	Role         string      `yaml:"role,omitempty"`
	Type         string      `yaml:"type"`
	Required     bool        `yaml:"required,omitempty"`
	Default      interface{} `yaml:"default,omitempty"`
	Description  string      `yaml:"description,omitempty"`
	Elements     string      `yaml:"elements,omitempty"`
	VersionAdded string      `yaml:"version_added,omitempty"`
}

func (args *ArgDefs) ParseArgDefFile(path string) error {
	m, err := utils.YamlToMap(path)
	if err != nil {
		return err
	}

	opts, ok := m["argument_specs"].(map[string]interface{})["main"].(map[string]interface{})["options"].(map[string]interface{})
	if !ok {
		return err
	}
	var re = regexp.MustCompile(`\/roles\/([^\/]+)`)
	for k, v := range opts {
		ac := ArgDef{}
		if err := ac.ParseArgDef(v.(map[string]interface{})); err != nil {
			err = fmt.Errorf("%w; Definition: %v", err, k)
			return err
		}
		if m := re.FindAllStringSubmatch(path, 1); len(m) > 0 {
			ac.Role = m[0][1]
		}
		ac.Name = k
		(*args)[k] = ac
	}

	return nil
}

func (ac *ArgDef) ParseArgDef(obj map[string]interface{}) error {
	if name, ok := obj["name"].(string); ok {
		ac.Name = name
	}

	if required, ok := obj["required"].(string); ok {
		requiredVal, err := strconv.ParseBool(required)
		if err != nil {
			return fmt.Errorf("Required value must be a boolean")
		}
		ac.Required = requiredVal
	} else {
		ac.Required = false
	}

	if desc, ok := obj["description"].(string); ok {
		ac.Description = desc
	}

	var argType string

	if t, ok := obj["type"].(string); ok {
		argType = t
	} else {
		argType = "str"
	}

	ac.Type = argType

	// Required is only required if True
	if def, ok := obj["default"]; ok {
		switch argType {
		case "path":
			fallthrough
		case "str":
			if def == nil {
				ac.Default = ""
			} else {
				defStr, ok := def.(string)
				if !ok {
					return fmt.Errorf("Default value must be a string with a str type: %v", def)
				}
				ac.Default = defStr
			}
		case "int":
			defInt, ok := def.(int)
			if !ok {
				return fmt.Errorf("Default value must be an integer with an int type")
			}
			ac.Default = defInt
		case "bool":
			defBool, ok := def.(bool)
			if !ok {
				return fmt.Errorf("Default value must be a boolean with a bool type")
			}
			ac.Default = defBool
		case "float":
			defFloat, ok := def.(float64)
			if !ok {
				return fmt.Errorf("Default value must be a float with a float type")
			}
			ac.Default = defFloat
		case "list":
			ele, ok := obj["element"].(string)
			if !ok {
				ele = "str"
			}
			ac.Elements = ele
			ac.Default = def
		case "dict":
			ac.Default = def
		default:
			return fmt.Errorf("Unknown type: %s", argType)
		}
	}

	return nil
}
